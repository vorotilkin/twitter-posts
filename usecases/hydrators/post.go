package hydrators

import (
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoPosts(posts []models.Post) []*proto.Post {
	return lo.Map(posts, func(post models.Post, _ int) *proto.Post {
		return ProtoPost(post)
	})
}

func ProtoPost(post models.Post) *proto.Post {
	return &proto.Post{
		Id:                post.ID,
		Body:              post.Body,
		CreatedAt:         timestamppb.New(post.CreatedAt),
		UpdatedAt:         timestamppb.New(post.UpdatedAt),
		UserId:            post.UserID,
		LikeCounter:       post.LikeCount,
		IsCurrentUserLike: post.IsCurrentUserLike,
		Comments:          ProtoComments(post.Comments),
	}
}

func ProtoComments(comments []models.Comment) []*proto.Comment {
	return lo.Map(comments, func(comment models.Comment, _ int) *proto.Comment {
		return &proto.Comment{
			Id:        comment.ID,
			Body:      comment.Body,
			UserId:    comment.UserID,
			PostId:    comment.PostID,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
		}
	})
}
