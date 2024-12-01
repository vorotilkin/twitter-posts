package usecases

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/proto"
	"github.com/vorotilkin/twitter-posts/usecases/hydrators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostsRepository interface {
	Create(ctx context.Context, userID int32, body string) (models.Post, error)
	Posts(ctx context.Context, filter models.PostFilter) ([]models.Post, error)
}

type LikeRepository interface {
	LikeCountByPostID(ctx context.Context, postIDs []int32) (map[int32]int32, error)
}

type CommentRepository interface {
	CommentsByPostID(ctx context.Context, postIDs []int32) (map[int32][]models.Comment, error)
}

type PostsServer struct {
	proto.UnimplementedPostsServer
	postsRepository    PostsRepository
	likeRepository     LikeRepository
	commentsRepository CommentRepository
}

func (s *PostsServer) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	if len(request.GetBody()) == 0 || request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "body and user_id is required")
	}

	post, err := s.postsRepository.Create(ctx, request.GetUserId(), request.GetBody())
	if err != nil {
		return nil, err
	}

	return &proto.CreateResponse{
		Post: hydrators.ProtoPost(post),
	}, nil
}

func (s *PostsServer) PostByID(ctx context.Context, request *proto.PostByIDRequest) (*proto.PostByIDResponse, error) {
	postID := request.GetId()
	if postID == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	posts, err := s.postsRepository.Posts(ctx, models.PostFilter{PostIDs: mo.Some([]int32{postID})})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch posts")
	}

	post, exist := lo.Find(posts, func(post models.Post) bool {
		return post.ID == postID
	})
	if !exist {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	likesByPostID, err := s.likeRepository.LikeCountByPostID(ctx, []int32{postID})
	if err != nil {
		return nil, errors.Wrap(err, "like count repo error")
	}

	post.LikeCount = likesByPostID[postID]

	commentsByPostID, err := s.commentsRepository.CommentsByPostID(ctx, []int32{postID})
	if err != nil {
		return nil, errors.Wrap(err, "get comments by post id")
	}

	post.Comments = commentsByPostID[postID]

	return &proto.PostByIDResponse{
		Post: hydrators.ProtoPost(post),
	}, nil
}

func (s *PostsServer) Posts(ctx context.Context, request *proto.PostsRequest) (*proto.PostsResponse, error) {
	protoFilter := request.GetFilters()
	if protoFilter == nil {
		return nil, status.Error(codes.InvalidArgument, "filters is required")
	}

	var filter models.PostFilter

	if userID := protoFilter.GetFilterUser().GetUserId(); userID > 0 {
		filter.UserIDs = mo.Some([]int32{userID})
	}

	if perPage := protoFilter.GetPagination().GetPerPage(); perPage > 0 {
		filter.Pagination = mo.Some(models.Pagination{PerPage: perPage})
	}

	if sort := protoFilter.GetSort().GetSort(); sort > 0 {
		filter.Sort = mo.Some(models.Sort(sort))
	}

	posts, err := s.postsRepository.Posts(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch posts")
	}

	postIDs := lo.Map(posts, func(post models.Post, _ int) int32 {
		return post.ID
	})

	likesByPostID, err := s.likeRepository.LikeCountByPostID(ctx, postIDs)
	if err != nil {
		return nil, errors.Wrap(err, "like count repo error")
	}

	commentsByPostID, err := s.commentsRepository.CommentsByPostID(ctx, postIDs)
	if err != nil {
		return nil, errors.Wrap(err, "get comments by post ids")
	}

	for i, post := range posts {
		posts[i].LikeCount = likesByPostID[post.ID]
		posts[i].Comments = commentsByPostID[post.ID]
	}

	return &proto.PostsResponse{
		Posts: hydrators.ProtoPosts(posts),
	}, nil
}

func NewPostsServer(
	postsRepo PostsRepository,
	likeRepo LikeRepository,
	commentRepo CommentRepository,
) *PostsServer {
	return &PostsServer{
		postsRepository:    postsRepo,
		likeRepository:     likeRepo,
		commentsRepository: commentRepo,
	}
}
