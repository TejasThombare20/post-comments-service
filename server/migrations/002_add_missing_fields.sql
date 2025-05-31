-- Migration: 002_add_missing_fields.sql
-- Description: Add missing fields to comments table
-- Created: 2024

-- Add updated_at field to comments table
ALTER TABLE comments ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();

-- Add replies_count field to comments table
ALTER TABLE comments ADD COLUMN replies_count INTEGER DEFAULT 0;

-- Create trigger to automatically update updated_at for comments
CREATE TRIGGER update_comments_updated_at 
    BEFORE UPDATE ON comments 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to update replies_count when a reply is added
CREATE OR REPLACE FUNCTION update_replies_count()
RETURNS TRIGGER AS $$
BEGIN
    -- If this is a reply (has parent_id), increment parent's replies_count
    IF NEW.parent_id IS NOT NULL THEN
        UPDATE comments 
        SET replies_count = replies_count + 1 
        WHERE id = NEW.parent_id;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to update replies_count when a reply is added
CREATE TRIGGER update_replies_count_trigger
    AFTER INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_replies_count();

-- Create function to decrement replies_count when a reply is deleted
CREATE OR REPLACE FUNCTION decrement_replies_count()
RETURNS TRIGGER AS $$
BEGIN
    -- If this was a reply (has parent_id), decrement parent's replies_count
    IF OLD.parent_id IS NOT NULL THEN
        UPDATE comments 
        SET replies_count = GREATEST(replies_count - 1, 0)
        WHERE id = OLD.parent_id;
    END IF;
    
    RETURN OLD;
END;
$$ language 'plpgsql';

-- Create trigger to decrement replies_count when a reply is deleted
CREATE TRIGGER decrement_replies_count_trigger
    AFTER UPDATE OF deleted_at ON comments
    FOR EACH ROW
    WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
    EXECUTE FUNCTION decrement_replies_count();

-- Update existing comments to have updated_at = created_at
UPDATE comments SET updated_at = created_at WHERE updated_at IS NULL;

-- Create index for replies_count for better performance
CREATE INDEX idx_comments_replies_count ON comments(replies_count);
CREATE INDEX idx_comments_updated_at ON comments(updated_at); 