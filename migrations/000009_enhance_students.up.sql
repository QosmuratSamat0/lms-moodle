-- Add new columns to students table for enrollment, group, and major tracking
ALTER TABLE students
ADD COLUMN IF NOT EXISTS id UUID DEFAULT uuid_generate_v4(),
ADD COLUMN IF NOT EXISTS student_code TEXT UNIQUE,
ADD COLUMN IF NOT EXISTS major TEXT,
ADD COLUMN IF NOT EXISTS year INT DEFAULT 1,
ADD COLUMN IF NOT EXISTS gpa DECIMAL(3,2) DEFAULT 0.00,
ADD COLUMN IF NOT EXISTS admitted_at TIMESTAMPTZ DEFAULT now(),
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_students_id ON students(id);
CREATE INDEX IF NOT EXISTS idx_students_student_code ON students(student_code);
CREATE INDEX IF NOT EXISTS idx_students_major ON students(major);
CREATE INDEX IF NOT EXISTS idx_students_year ON students(year);
CREATE INDEX IF NOT EXISTS idx_students_gpa ON students(gpa);

-- Function to update updated_at on students
CREATE OR REPLACE FUNCTION update_students_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for auto-updating updated_at
DROP TRIGGER IF EXISTS trigger_students_updated_at ON students;
CREATE TRIGGER trigger_students_updated_at
BEFORE UPDATE ON students
FOR EACH ROW
EXECUTE FUNCTION update_students_updated_at();
