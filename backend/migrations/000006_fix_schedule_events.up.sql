-- Add missing columns to schedule_events table
ALTER TABLE schedule_events
ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS event_type TEXT NOT NULL DEFAULT 'class';

-- Update index
CREATE INDEX IF NOT EXISTS idx_schedule_events_created_by ON schedule_events(created_by);
