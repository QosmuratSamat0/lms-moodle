-- Fix missing student profiles
-- This migration ensures all users with role='student' have a corresponding record in students table

INSERT INTO students (user_id, first_name, last_name, group_name, created_at)
SELECT
    u.id,
    '', -- first_name (empty, can be updated later)
    '', -- last_name (empty, can be updated later)
    '', -- group_name (empty, can be updated later)
    u.created_at
FROM users u
WHERE u.role = 'student'
AND NOT EXISTS (SELECT 1 FROM students s WHERE s.user_id = u.id);

-- Similarly for teachers
INSERT INTO teachers (user_id, first_name, last_name, department, created_at)
SELECT
    u.id,
    '',
    '',
    '',
    u.created_at
FROM users u
WHERE u.role = 'teacher'
AND NOT EXISTS (SELECT 1 FROM teachers t WHERE t.user_id = u.id);

-- And for managers
INSERT INTO managers (user_id, first_name, last_name, created_at)
SELECT
    u.id,
    '',
    '',
    u.created_at
FROM users u
WHERE u.role = 'manager'
AND NOT EXISTS (SELECT 1 FROM managers m WHERE m.user_id = u.id);
