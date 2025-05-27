-- Drop views
DROP VIEW IF EXISTS article_stats;

-- Drop triggers
DROP TRIGGER IF EXISTS update_storage_stats ON articles;
DROP TRIGGER IF EXISTS update_articles_updated_at ON articles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS update_user_storage_stats();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_token_hash;
DROP INDEX IF EXISTS idx_sessions_user_id;

DROP INDEX IF EXISTS idx_articles_content_text;
DROP INDEX IF EXISTS idx_articles_category;
DROP INDEX IF EXISTS idx_articles_tags;
DROP INDEX IF EXISTS idx_articles_updated_at;
DROP INDEX IF EXISTS idx_articles_created_at;
DROP INDEX IF EXISTS idx_articles_is_archived;
DROP INDEX IF EXISTS idx_articles_is_favorite;
DROP INDEX IF EXISTS idx_articles_is_read;
DROP INDEX IF EXISTS idx_articles_status;
DROP INDEX IF EXISTS idx_articles_url;
DROP INDEX IF EXISTS idx_articles_user_id;

DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp"; 