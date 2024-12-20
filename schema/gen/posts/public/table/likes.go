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

var Likes = newLikesTable("public", "likes", "")

type likesTable struct {
	postgres.Table

	// Columns
	UserID postgres.ColumnInteger
	PostID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type LikesTable struct {
	likesTable

	EXCLUDED likesTable
}

// AS creates new LikesTable with assigned alias
func (a LikesTable) AS(alias string) *LikesTable {
	return newLikesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new LikesTable with assigned schema name
func (a LikesTable) FromSchema(schemaName string) *LikesTable {
	return newLikesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new LikesTable with assigned table prefix
func (a LikesTable) WithPrefix(prefix string) *LikesTable {
	return newLikesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new LikesTable with assigned table suffix
func (a LikesTable) WithSuffix(suffix string) *LikesTable {
	return newLikesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newLikesTable(schemaName, tableName, alias string) *LikesTable {
	return &LikesTable{
		likesTable: newLikesTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newLikesTableImpl("", "excluded", ""),
	}
}

func newLikesTableImpl(schemaName, tableName, alias string) likesTable {
	var (
		UserIDColumn   = postgres.IntegerColumn("user_id")
		PostIDColumn   = postgres.IntegerColumn("post_id")
		allColumns     = postgres.ColumnList{UserIDColumn, PostIDColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return likesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID: UserIDColumn,
		PostID: PostIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
