-- Drop junction tables
DROP TABLE IF EXISTS assignment_attachments;
DROP TABLE IF EXISTS submission_attachments;

-- Drop attachments table
DROP TABLE IF EXISTS attachments;

-- Remove added columns from attendance_marks
ALTER TABLE attendance_marks 
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS marked_by;

-- Remove added columns from attendance_sessions
ALTER TABLE attendance_sessions 
    DROP COLUMN IF EXISTS group_id,
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS session_date;

-- Remove added columns from assignments
ALTER TABLE assignments 
    DROP COLUMN IF EXISTS allow_late,
    DROP COLUMN IF EXISTS group_id;

-- Remove added columns from students
ALTER TABLE students 
    DROP COLUMN IF EXISTS group_id,
    DROP COLUMN IF EXISTS major,
    DROP COLUMN IF EXISTS year_of_study,
    DROP COLUMN IF EXISTS gpa,
    DROP COLUMN IF EXISTS enrollment_status;

-- Drop teacher_course_groups table
DROP TABLE IF EXISTS teacher_course_groups;

-- Drop groups table
DROP TABLE IF EXISTS groups;
