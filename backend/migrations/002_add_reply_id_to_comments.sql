-- Add reply_id column to comments table
-- reply_id is nullable and references the id column of the same table (self-referencing foreign key)
ALTER TABLE comments
ADD COLUMN reply_id INT NULL REFERENCES comments(id) ON DELETE SET NULL;

-- Add an index on reply_id for better query performance
CREATE INDEX idx_comments_reply_id ON comments(reply_id);
