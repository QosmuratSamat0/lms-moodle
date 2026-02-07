-- Rollback: Enhance roles structure with separated entities

-- Remove constraints added to users
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS check_role;

-- Update courses table
ALTER TABLE courses
    DROP INDEX IF EXISTS idx_courses_category_id,
    DROP INDEX IF EXISTS idx_courses_teacher_id,
    DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS teacher_id;

-- Drop category managers table
DROP TABLE IF EXISTS category_managers CASCADE;

-- Drop category managers table
DROP TABLE IF EXISTS course_categories CASCADE;

-- Recreate old managers table
CREATE TABLE IF NOT EXISTS managers_old (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO managers_old (user_id, first_name, last_name, created_at)
SELECT user_id, first_name, last_name, created_at FROM managers;

DROP TABLE IF EXISTS managers CASCADE;
ALTER TABLE managers_old RENAME TO managers;

-- Recreate old teachers table
CREATE TABLE IF NOT EXISTS teachers_old (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    department TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO teachers_old (user_id, first_name, last_name, department, created_at)
SELECT user_id, first_name, last_name, department, created_at FROM teachers;

DROP TABLE IF EXISTS teachers CASCADE;
ALTER TABLE teachers_old RENAME TO teachers;

-- Drop admins table
DROP TABLE IF EXISTS admins CASCADE;
