-- Create "post" table
CREATE TABLE "post" ("id" serial NOT NULL, "body" text NOT NULL, "user_id" integer NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Create index "idx_post_id" to table: "post"
CREATE INDEX "idx_post_id" ON "post" ("id");
-- Create index "idx_post_user_id" to table: "post"
CREATE INDEX "idx_post_user_id" ON "post" ("user_id");
-- Create "comment" table
CREATE TABLE "comment" ("id" serial NOT NULL, "body" text NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "user_id" integer NOT NULL, "post_id" integer NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_comment_post_id" FOREIGN KEY ("post_id") REFERENCES "post" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_comment_post_id" to table: "comment"
CREATE INDEX "idx_comment_post_id" ON "comment" ("post_id");
-- Create index "idx_comment_user_id" to table: "comment"
CREATE INDEX "idx_comment_user_id" ON "comment" ("user_id");
-- Create "likes" table
CREATE TABLE "likes" ("id" serial NOT NULL, "user_id" integer NOT NULL, "post_id" integer NOT NULL, "is_like" boolean NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_like_post_id" FOREIGN KEY ("post_id") REFERENCES "post" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_like_post_id" to table: "likes"
CREATE INDEX "idx_like_post_id" ON "likes" ("post_id");
-- Create index "idx_like_user_id" to table: "likes"
CREATE INDEX "idx_like_user_id" ON "likes" ("user_id");
