-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    is_email_verified BOOLEAN DEFAULT false,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Cloud backup preferences
    enable_cloud_backup BOOLEAN DEFAULT false,
    encryption_key VARCHAR(255),
    
    -- Storage statistics
    storage_used BIGINT DEFAULT 0,
    max_storage_limit BIGINT DEFAULT 10737418240, -- 10GB in bytes
    
    -- Preferences
    theme VARCHAR(20) DEFAULT 'auto',
    language VARCHAR(2) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC'
);

-- Create articles table
CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    content TEXT,
    content_text TEXT, -- Plain text version for search
    author VARCHAR(255),
    site_name VARCHAR(255),
    image_url TEXT,
    published_at TIMESTAMP,
    word_count INTEGER DEFAULT 0,
    reading_time INTEGER DEFAULT 0, -- in minutes
    language VARCHAR(10),
    
    -- Storage and encryption
    is_encrypted BOOLEAN DEFAULT false,
    storage_size BIGINT DEFAULT 0,
    local_path TEXT,
    cloud_path TEXT,
    
    -- Status and metadata
    status VARCHAR(20) DEFAULT 'pending', -- pending, processed, failed
    is_read BOOLEAN DEFAULT false,
    is_favorite BOOLEAN DEFAULT false,
    is_archived BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    -- Tags and categories (using JSONB for better performance)
    tags JSONB DEFAULT '[]',
    category VARCHAR(100),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Capture metadata
    capture_method VARCHAR(20) DEFAULT 'api', -- extension, api, manual
    user_agent TEXT,
    captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create sessions table for user sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_active ON users(is_active);

CREATE INDEX idx_articles_user_id ON articles(user_id);
CREATE INDEX idx_articles_url ON articles(url);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_is_read ON articles(is_read);
CREATE INDEX idx_articles_is_favorite ON articles(is_favorite);
CREATE INDEX idx_articles_is_archived ON articles(is_archived);
CREATE INDEX idx_articles_created_at ON articles(created_at);
CREATE INDEX idx_articles_updated_at ON articles(updated_at);
CREATE INDEX idx_articles_tags ON articles USING GIN(tags);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_content_text ON articles USING GIN(to_tsvector('english', content_text));

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to update user storage statistics
CREATE OR REPLACE FUNCTION update_user_storage_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users 
        SET storage_used = storage_used + NEW.storage_size
        WHERE id = NEW.user_id;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE users 
        SET storage_used = storage_used - OLD.storage_size + NEW.storage_size
        WHERE id = NEW.user_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users 
        SET storage_used = storage_used - OLD.storage_size
        WHERE id = OLD.user_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update storage statistics
CREATE TRIGGER update_storage_stats AFTER INSERT OR UPDATE OR DELETE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_user_storage_stats();

-- Create view for article statistics
CREATE VIEW article_stats AS
SELECT 
    user_id,
    COUNT(*) as total_articles,
    COUNT(*) FILTER (WHERE is_read = true) as read_articles,
    COUNT(*) FILTER (WHERE is_read = false) as unread_articles,
    COUNT(*) FILTER (WHERE is_favorite = true) as favorite_articles,
    COUNT(*) FILTER (WHERE is_archived = true) as archived_articles,
    COALESCE(SUM(word_count), 0) as total_word_count,
    COALESCE(SUM(reading_time), 0) as total_reading_time,
    COALESCE(SUM(storage_size), 0) as storage_used
FROM articles
GROUP BY user_id; 