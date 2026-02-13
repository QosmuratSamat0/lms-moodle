DROP INDEX IF EXISTS idx_assignments_category;
ALTER TABLE assignments DROP COLUMN IF EXISTS weight_percentage;
ALTER TABLE assignments DROP COLUMN IF EXISTS grading_category;
