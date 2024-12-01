package comment

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/pkg/database"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/model"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) CommentsByPostID(ctx context.Context, postIDs []int32) (map[int32][]models.Comment, error) {
	dbIDs := lo.Map(postIDs, func(id int32, _ int) postgres.Expression {
		return postgres.Int(int64(id))
	})

	query, args := table.Comment.
		SELECT(
			table.Comment.ID,
			table.Comment.Body,
			table.Comment.UserID,
			table.Comment.PostID,
			table.Comment.CreatedAt,
			table.Comment.UpdatedAt,
		).
		WHERE(table.Comment.PostID.IN(dbIDs...)).
		Sql()

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	comments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Comment, error) {
		var comment model.Comment

		err := row.Scan(
			&comment.ID,
			&comment.Body,
			&comment.UserID,
			&comment.PostID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return models.Comment{}, err
		}

		return domainComment(comment), nil
	})

	return lo.GroupBy(comments, func(comment models.Comment) int32 {
		return comment.PostID
	}), nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
