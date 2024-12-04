-- Modify "likes" table
ALTER TABLE "likes" DROP CONSTRAINT "likes_pkey", DROP COLUMN "id", DROP COLUMN "is_like", ADD PRIMARY KEY ("user_id", "post_id");
