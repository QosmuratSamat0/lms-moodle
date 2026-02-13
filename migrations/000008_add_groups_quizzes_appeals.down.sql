-- Rollback: Remove groups, announcements, quizzes, and grade appeals tables

-- Drop triggers
DROP TRIGGER IF EXISTS update_quizzes_updated_at ON quizzes;
DROP TRIGGER IF EXISTS update_announcements_updated_at ON announcements;
DROP TRIGGER IF EXISTS update_groups_updated_at ON groups;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS quiz_answers;
DROP TABLE IF EXISTS quiz_attempts;
DROP TABLE IF EXISTS quiz_questions;
DROP TABLE IF EXISTS quizzes;
DROP TABLE IF EXISTS grade_appeals;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;

-- Remove teacher_id from courses if exists
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'courses' AND column_name = 'teacher_id') THEN
        DROP INDEX IF EXISTS idx_courses_teacher_id;
        ALTER TABLE courses DROP COLUMN teacher_id;
    END IF;
END $$;
