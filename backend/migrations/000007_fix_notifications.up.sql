-- Fix notifications table schema to match repository expectations
-- Rename body to message
ALTER TABLE notifications RENAME COLUMN body TO message;

-- Add missing columns
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS channel TEXT NOT NULL DEFAULT 'in_app',
ADD COLUMN IF NOT EXISTS read_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS sent_at TIMESTAMPTZ;

-- Update is_read based on read_at
UPDATE notifications SET read_at = created_at WHERE is_read = true;

-- Drop old is_read column (now we use read_at IS NOT NULL to check if read)
ALTER TABLE notifications DROP COLUMN IF EXISTS is_read;
