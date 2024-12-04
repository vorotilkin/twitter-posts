package like

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

func (r *Repository) LikesByPostIDs(ctx context.Context, postIDs []int32) (map[int32][]models.Like, error) {
	dbIDs := lo.Map(postIDs, func(id int32, _ int) postgres.Expression {
		return postgres.Int(int64(id))
	})

	query, args := table.Likes.
		SELECT(
			table.Likes.PostID,
			table.Likes.UserID,
		).
		WHERE(table.Likes.PostID.IN(dbIDs...)).
		Sql()

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	likes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Like, error) {
		var result models.Like

		err := row.Scan(
			&result.PostID,
			&result.UserID,
		)
		if err != nil {
			return result, err
		}

		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return lo.GroupBy(likes, func(like models.Like) int32 {
		return like.PostID
	}), nil
}

func (r *Repository) Like(ctx context.Context, userID, postID int32) (bool, error) {
	query, args := table.Likes.
		INSERT(table.Likes.UserID, table.Likes.PostID).
		MODEL(model.Likes{
			UserID: userID,
			PostID: postID,
		}).
		Sql()

	commandTag, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (r *Repository) Dislike(ctx context.Context, userID, postID int32) (bool, error) {
	query, args := table.Likes.
		DELETE().WHERE(
		table.Likes.UserID.EQ(postgres.Int(int64(userID))).
			AND(
				table.Likes.PostID.EQ(postgres.Int(int64(postID))),
			),
	).
		Sql()

	commandTag, err := r.conn.Exec(ctx, query, args...)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
