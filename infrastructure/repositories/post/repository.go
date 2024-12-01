package post

import (
	"context"
	"errors"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/vorotilkin/twitter-posts/domain/models"
	"github.com/vorotilkin/twitter-posts/pkg/database"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/model"
	"github.com/vorotilkin/twitter-posts/schema/gen/posts/public/table"
)

type Repository struct {
	conn *database.Database
}

func (r *Repository) PostByID(ctx context.Context, id int32) (models.Post, error) {
	query, args := table.Post.
		SELECT(
			table.Post.ID,
			table.Post.Body,
			table.Post.UserID,
			table.Post.CreatedAt,
			table.Post.UpdatedAt,
		).
		WHERE(table.Post.ID.EQ(postgres.Int(int64(id)))).
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Post{}, err
	}

	count, err := r.likeCount(ctx, id)
	if err != nil {
		return models.Post{}, err
	}

	comments, err := r.commentsByPostID(ctx, id)
	if err != nil {
		return models.Post{}, err
	}

	return domainPost(post, comments, count), nil
}

func (r *Repository) likeCount(ctx context.Context, postID int32) (int32, error) {
	query, args := table.Like.
		SELECT(
			postgres.COUNT(table.Like.ID),
		).
		WHERE(table.Like.PostID.EQ(postgres.Int(int64(postID))).AND(table.Like.IsLike.IS_TRUE())).
		Sql()

	row := r.conn.QueryRow(ctx, query, args...)
	var count int32

	err := row.Scan(
		&count,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	return count, nil
}

func (r *Repository) commentsByPostID(ctx context.Context, postID int32) ([]models.Comment, error) {
	query, args := table.Comment.
		SELECT(
			table.Comment.ID,
			table.Comment.Body,
			table.Comment.UserID,
			table.Comment.PostID,
			table.Comment.CreatedAt,
			table.Comment.UpdatedAt,
		).
		WHERE(table.Comment.PostID.EQ(postgres.Int(int64(postID)))).
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

	return comments, err
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

	return domainPost(post, nil, 0), nil
}

func NewRepository(conn *database.Database) *Repository {
	return &Repository{conn: conn}
}
