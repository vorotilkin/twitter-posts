package comment

import (
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/model"
)

func domainComment(comment model.Comment) models.Comment {
	return models.Comment{
		ID:        comment.ID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
	}
}
