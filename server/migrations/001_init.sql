-- Migration: 001_init.sql
-- Description: Initial database schema for post-comments service
-- Created: 2024

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT,
    display_name TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create posts table
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create comments table with nested comment support
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    path UUID[] NOT NULL,
    thread_id UUID NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);


-- Users table indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Posts table indexes
CREATE INDEX idx_posts_created_by ON posts(created_by);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);
CREATE INDEX idx_posts_title ON posts USING gin(to_tsvector('english', title));
CREATE INDEX idx_posts_content ON posts USING gin(to_tsvector('english', content));

-- Comments table indexes
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_thread_id ON comments(thread_id);
CREATE INDEX idx_comments_created_by ON comments(created_by);
CREATE INDEX idx_comments_created_at ON comments(created_at);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);
CREATE INDEX idx_comments_path_gin ON comments USING GIN (path);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_posts_updated_at 
    BEFORE UPDATE ON posts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to handle comment path and thread_id
CREATE OR REPLACE FUNCTION set_comment_path_and_thread()
RETURNS TRIGGER AS $$
BEGIN
    -- If this is a top-level comment (no parent)
    IF NEW.parent_id IS NULL THEN
        NEW.path = ARRAY[NEW.id];
        NEW.thread_id = NEW.id;
    ELSE
        -- This is a reply to another comment
        DECLARE
            parent_path UUID[];
            parent_thread_id UUID;
        BEGIN
            -- Get parent comment's path and thread_id
            SELECT path, thread_id INTO parent_path, parent_thread_id
            FROM comments 
            WHERE id = NEW.parent_id;
            
            -- Set the path as parent's path + this comment's id
            NEW.path = parent_path || NEW.id;
            NEW.thread_id = parent_thread_id;
        END;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for comment path and thread_id
CREATE TRIGGER set_comment_path_and_thread_trigger
    BEFORE INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION set_comment_path_and_thread();

-- Insert some sample data for testing (optional)
-- Uncomment the following lines if you want sample data

-- INSERT INTO users (username, email, display_name) VALUES 
-- ('john_doe', 'john@example.com', 'John Doe'),
-- ('jane_smith', 'jane@example.com', 'Jane Smith'),
-- ('bob_wilson', 'bob@example.com', 'Bob Wilson');

-- INSERT INTO posts (title, content, created_by) VALUES 
-- ('Welcome to our platform', 'This is the first post on our platform. Welcome everyone!', 
--  (SELECT id FROM users WHERE username = 'john_doe')),
-- ('How to get started', 'Here are some tips to get started with our service...', 
--  (SELECT id FROM users WHERE username = 'jane_smith'));

-- INSERT INTO comments (content, post_id, created_by) VALUES 
-- ('Great post! Thanks for sharing.', 
--  (SELECT id FROM posts WHERE title = 'Welcome to our platform'), 
--  (SELECT id FROM users WHERE username = 'jane_smith')),
-- ('Looking forward to using this platform.', 
--  (SELECT id FROM posts WHERE title = 'Welcome to our platform'), 
--  (SELECT id FROM users WHERE username = 'bob_wilson')); 