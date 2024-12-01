//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Comment = newCommentTable("public", "comment", "")

type commentTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	Body      postgres.ColumnString
	CreatedAt postgres.ColumnTimestamp
	UpdatedAt postgres.ColumnTimestamp
	UserID    postgres.ColumnInteger
	PostID    postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CommentTable struct {
	commentTable

	EXCLUDED commentTable
}

// AS creates new CommentTable with assigned alias
func (a CommentTable) AS(alias string) *CommentTable {
	return newCommentTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CommentTable with assigned schema name
func (a CommentTable) FromSchema(schemaName string) *CommentTable {
	return newCommentTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CommentTable with assigned table prefix
func (a CommentTable) WithPrefix(prefix string) *CommentTable {
	return newCommentTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CommentTable with assigned table suffix
func (a CommentTable) WithSuffix(suffix string) *CommentTable {
	return newCommentTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCommentTable(schemaName, tableName, alias string) *CommentTable {
	return &CommentTable{
		commentTable: newCommentTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newCommentTableImpl("", "excluded", ""),
	}
}

func newCommentTableImpl(schemaName, tableName, alias string) commentTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		BodyColumn      = postgres.StringColumn("body")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		UpdatedAtColumn = postgres.TimestampColumn("updated_at")
		UserIDColumn    = postgres.IntegerColumn("user_id")
		PostIDColumn    = postgres.IntegerColumn("post_id")
		allColumns      = postgres.ColumnList{IDColumn, BodyColumn, CreatedAtColumn, UpdatedAtColumn, UserIDColumn, PostIDColumn}
		mutableColumns  = postgres.ColumnList{BodyColumn, CreatedAtColumn, UpdatedAtColumn, UserIDColumn, PostIDColumn}
	)

	return commentTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Body:      BodyColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,
		UserID:    UserIDColumn,
		PostID:    PostIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}