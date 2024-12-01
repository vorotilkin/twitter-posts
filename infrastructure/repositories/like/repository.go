package like

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/vorotilkin/twitter-posts/pkg/database"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) LikeCountByPostID(ctx context.Context, postIDs []int32) (map[int32]int32, error) {
	dbIDs := lo.Map(postIDs, func(id int32, _ int) postgres.Expression {
		return postgres.Int(int64(id))
	})

	query, args := table.Likes.
		SELECT(
			table.Likes.PostID,
			postgres.COUNT(table.Likes.ID),
		).
		WHERE(table.Likes.PostID.IN(dbIDs...).AND(table.Likes.IsLike.IS_TRUE())).
		GROUP_BY(table.Likes.PostID).
		Sql()

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	type dbResult struct {
		PostID    int32
		LikeCount int32
	}

	likes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (dbResult, error) {
		var result dbResult

		err := row.Scan(
			&result.PostID,
			&result.LikeCount,
		)
		if err != nil {
			return dbResult{}, err
		}

		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return lo.SliceToMap(likes, func(like dbResult) (int32, int32) {
		return like.PostID, like.LikeCount
	}), nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
