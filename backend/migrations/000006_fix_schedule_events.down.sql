-- Remove added columns from schedule_events table
ALTER TABLE schedule_events
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS event_type;

DROP INDEX IF EXISTS idx_schedule_events_created_by;
