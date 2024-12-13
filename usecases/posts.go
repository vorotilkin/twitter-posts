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
	LikesByPostIDs(ctx context.Context, postIDs []int32) (map[int32][]models.Like, error)
	Like(ctx context.Context, userID, postID int32) (bool, error)
	Dislike(ctx context.Context, userID, postID int32) (bool, error)
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

	likesByPostID, err := s.likeRepository.LikesByPostIDs(ctx, []int32{postID})
	if err != nil {
		return nil, errors.Wrap(err, "like count repo error")
	}

	likes := likesByPostID[postID]

	userID := request.GetUserId()

	post.LikeCount = int32(len(likes))
	post.IsCurrentUserLike = lo.ContainsBy(likes, func(like models.Like) bool {
		return like.UserID == userID
	})

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

	if userIDs := protoFilter.GetFilterUsers().GetUserIds(); len(userIDs) > 0 {
		filter.UserIDs = mo.Some(userIDs)
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

	if len(posts) == 0 {
		return &proto.PostsResponse{
			Posts: []*proto.Post{},
		}, nil
	}

	postIDs := lo.Map(posts, func(post models.Post, _ int) int32 {
		return post.ID
	})

	likesByPostID, err := s.likeRepository.LikesByPostIDs(ctx, postIDs)
	if err != nil {
		return nil, errors.Wrap(err, "like count repo error")
	}

	commentsByPostID, err := s.commentsRepository.CommentsByPostID(ctx, postIDs)
	if err != nil {
		return nil, errors.Wrap(err, "get comments by post ids")
	}

	userID := request.GetCurrentUserId()

	for i, post := range posts {
		likes := likesByPostID[post.ID]

		post.LikeCount = int32(len(likes))
		post.IsCurrentUserLike = lo.ContainsBy(likes, func(like models.Like) bool {
			return like.UserID == userID
		})
		post.Comments = commentsByPostID[post.ID]

		posts[i] = post
	}

	return &proto.PostsResponse{
		Posts: hydrators.ProtoPosts(posts),
	}, nil
}

func (s *PostsServer) CommentsByPostID(ctx context.Context, request *proto.CommentsByPostIDRequest) (*proto.CommentsByPostIDResponse, error) {
	postID := request.GetPostId()
	if postID == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	commentsByPostID, err := s.commentsRepository.CommentsByPostID(ctx, []int32{postID})
	if err != nil {
		return nil, errors.Wrap(err, "get comments by post id")
	}

	return &proto.CommentsByPostIDResponse{Comments: hydrators.ProtoComments(commentsByPostID[postID])}, nil
}

func (s *PostsServer) Like(ctx context.Context, request *proto.LikeRequest) (*proto.LikeResponse, error) {
	userID := request.GetUserId()
	if userID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	postID := request.GetPostId()
	if postID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid post id")
	}

	var (
		ok  bool
		err error
	)

	switch request.GetOperationType() {
	case proto.LikeRequest_OPERATION_TYPE_DISLIKE:
		ok, err = s.likeRepository.Dislike(ctx, userID, postID)
	default:
		ok, err = s.likeRepository.Like(ctx, userID, postID)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LikeResponse{Ok: ok}, nil
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
