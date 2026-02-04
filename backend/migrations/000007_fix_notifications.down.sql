-- Revert notifications table to original schema
ALTER TABLE notifications RENAME COLUMN message TO body;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS is_read BOOLEAN NOT NULL DEFAULT false;

UPDATE notifications SET is_read = true WHERE read_at IS NOT NULL;

ALTER TABLE notifications
DROP COLUMN IF EXISTS channel,
DROP COLUMN IF EXISTS read_at,
DROP COLUMN IF EXISTS sent_at;
