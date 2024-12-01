package post

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

func (r *Repository) Create(ctx context.Context, userID int32, body string) (models.Post, error) {
	query, args := table.Post.
		INSERT(table.Post.UserID, table.Post.Body).
		MODEL(model.Post{
			Body:   body,
			UserID: userID,
		}).
		RETURNING(
			table.Post.ID,
			table.Post.Body,
			table.Post.UserID,
			table.Post.CreatedAt,
			table.Post.UpdatedAt,
		).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	post := model.Post{}

	err := row.Scan(
		&post.ID,
		&post.Body,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return models.Post{}, err
	}

	return domainPost(post), nil
}

func (r *Repository) Posts(ctx context.Context, filter models.PostFilter) ([]models.Post, error) {
	query := table.Post.
		SELECT(
			table.Post.ID,
			table.Post.Body,
			table.Post.UserID,
			table.Post.CreatedAt,
			table.Post.UpdatedAt,
		)

	filter.Sort.ForEach(func(sort models.Sort) {
		if sort.IsAsc() {
			query.ORDER_BY(table.Post.CreatedAt.ASC())
		}
		if sort.IsDesc() {
			query.ORDER_BY(table.Post.UpdatedAt.DESC())
		}
	})

	filter.PostIDs.ForEach(func(postIDs []int32) {
		ids := lo.Map(postIDs, func(postID int32, _ int) postgres.Expression {
			return postgres.Int(int64(postID))
		})
		query.WHERE(table.Post.ID.IN(ids...))
	})

	filter.UserIDs.ForEach(func(userIDs []int32) {
		ids := lo.Map(userIDs, func(userID int32, _ int) postgres.Expression {
			return postgres.Int(int64(userID))
		})
		query.WHERE(table.Post.UserID.IN(ids...))
	})

	filter.Pagination.ForEach(func(pagination models.Pagination) {
		query.LIMIT(int64(pagination.PerPage))
	})

	sql, args := query.Sql()

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	posts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Post, error) {
		var post model.Post

		err := row.Scan(
			&post.ID,
			&post.Body,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return models.Post{}, err
		}

		return domainPost(post), nil
	})

	return posts, err
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
