package models

import (
	"github.com/samber/mo"
	"time"
)

type Post struct {
	ID        int32
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int32
	LikeCount int32
	Comments  []Comment
}

type PostFilter struct {
	Sort       mo.Option[Sort]
	Pagination mo.Option[Pagination]
	UserIDs    mo.Option[[]int32]
	PostIDs    mo.Option[[]int32]
}
