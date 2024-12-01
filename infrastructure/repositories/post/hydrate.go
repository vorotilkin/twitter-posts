package post

import (
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/model"
)

func domainPost(post model.Post) models.Post {
	return models.Post{
		ID:        post.ID,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		UserID:    post.UserID,
	}
}
