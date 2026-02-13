-- Remove trigger and function
DROP TRIGGER IF EXISTS trigger_students_updated_at ON students;
DROP FUNCTION IF EXISTS update_students_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_students_id;
DROP INDEX IF EXISTS idx_students_student_code;
DROP INDEX IF EXISTS idx_students_major;
DROP INDEX IF EXISTS idx_students_year;
DROP INDEX IF EXISTS idx_students_gpa;

-- Remove columns added in up migration
ALTER TABLE students
DROP COLUMN IF EXISTS id,
DROP COLUMN IF EXISTS student_code,
DROP COLUMN IF EXISTS major,
DROP COLUMN IF EXISTS year,
DROP COLUMN IF EXISTS gpa,
DROP COLUMN IF EXISTS admitted_at,
DROP COLUMN IF EXISTS updated_at;
