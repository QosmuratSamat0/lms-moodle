-- Add grading category and weight percentage columns to assignments
ALTER TABLE assignments ADD COLUMN IF NOT EXISTS grading_category TEXT DEFAULT 'register_midterm'
    CHECK (grading_category IN ('register_midterm', 'register_endterm', 'final', 'bonus'));

ALTER TABLE assignments ADD COLUMN IF NOT EXISTS weight_percentage DECIMAL(5,2) DEFAULT 0
    CHECK (weight_percentage >= 0 AND weight_percentage <= 100);

CREATE INDEX IF NOT EXISTS idx_assignments_category ON assignments(grading_category);
