table "post" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "body" {
    null = false
    type = text
  }
  column "user_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_post_id" {
    columns = [column.id]
  }
  index "idx_post_user_id" {
    columns = [column.user_id]
  }
}
table "likes" {
  schema = schema.public

  column "user_id" {
    null = false
    type = integer
  }

  column "post_id" {
    null = false
    type = integer
  }

  primary_key {
    columns = [column.user_id, column.post_id]
  }

  index "idx_like_user_id" {
    columns = [column.user_id]
  }

  index "idx_like_post_id" {
    columns = [column.post_id]
  }

  foreign_key "fk_like_post_id" {
    columns    = [column.post_id]
    ref_columns = [table.post.column.id]
    on_delete = CASCADE
  }
}

table "comment" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "body" {
    null = false
    type = text
  }

  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "user_id" {
    null = false
    type = integer
  }

  column "post_id" {
    null = false
    type = integer
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_comment_user_id" {
    columns = [column.user_id]
  }

  index "idx_comment_post_id" {
    columns = [column.post_id]
  }

  foreign_key "fk_comment_post_id" {
      columns    = [column.post_id]
      ref_columns = [table.post.column.id]
      on_delete = CASCADE
  }
}
schema "public" {
  comment = "standard public schema"
}
