-- Migration: Enhance roles structure with separated entities

-- ====== COURSE CATEGORIES ======
CREATE TABLE IF NOT EXISTS course_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(255),
    "order" INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_course_categories_name ON course_categories(name);
CREATE INDEX idx_course_categories_active ON course_categories(is_active);
CREATE INDEX idx_course_categories_order ON course_categories("order");

-- ====== ENHANCED TEACHERS ======
-- Drop old teachers table constraints if needed
ALTER TABLE teachers 
    DROP CONSTRAINT IF EXISTS teachers_pkey,
    DROP CONSTRAINT IF EXISTS teachers_user_id_fkey;

-- Rename old teachers table
ALTER TABLE IF EXISTS teachers RENAME TO teachers_old;

-- Create new enhanced teachers table
CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_id VARCHAR(50) UNIQUE,
    first_name TEXT,
    last_name TEXT,
    department VARCHAR(255),
    specialization VARCHAR(255),
    qualifications VARCHAR(255),
    bio TEXT,
    office_hours JSONB,
    phone VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_teachers_user_id ON teachers(user_id);
CREATE INDEX idx_teachers_employee_id ON teachers(employee_id);
CREATE INDEX idx_teachers_name ON teachers(first_name, last_name);
CREATE INDEX idx_teachers_department ON teachers(department);
CREATE INDEX idx_teachers_active ON teachers(is_active);

-- Migrate data from old teachers table
INSERT INTO teachers (user_id, first_name, last_name, department, is_active, created_at, updated_at)
SELECT user_id, COALESCE(first_name, ''), COALESCE(last_name, ''), department, true, created_at, NOW()
FROM teachers_old
ON CONFLICT (user_id) DO NOTHING;

-- Drop old table
DROP TABLE IF EXISTS teachers_old;

-- ====== ADMINS ======
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_id VARCHAR(50) UNIQUE,
    department VARCHAR(255),
    access_level VARCHAR(50) NOT NULL CHECK (access_level IN ('super_admin', 'admin')) DEFAULT 'admin',
    permissions JSONB DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_admins_user_id ON admins(user_id);
CREATE INDEX idx_admins_employee_id ON admins(employee_id);
CREATE INDEX idx_admins_access_level ON admins(access_level);
CREATE INDEX idx_admins_active ON admins(is_active);

-- ====== MANAGERS ======
-- Enhance existing managers table
ALTER TABLE managers 
    ADD COLUMN IF NOT EXISTS id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    ADD COLUMN IF NOT EXISTS employee_id VARCHAR(50) UNIQUE,
    ADD COLUMN IF NOT EXISTS department VARCHAR(255),
    ADD COLUMN IF NOT EXISTS manages_categories JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS manages_teachers JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();

ALTER TABLE managers
    DROP CONSTRAINT IF EXISTS managers_pkey CASCADE;

-- Recreate managers table properly
CREATE TABLE IF NOT EXISTS managers_new (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_id VARCHAR(50) UNIQUE,
    first_name TEXT,
    last_name TEXT,
    department VARCHAR(255),
    manages_categories JSONB DEFAULT '[]'::jsonb,
    manages_teachers JSONB DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_managers_user_id ON managers_new(user_id);
CREATE INDEX idx_managers_employee_id ON managers_new(employee_id);
CREATE INDEX idx_managers_department ON managers_new(department);
CREATE INDEX idx_managers_active ON managers_new(is_active);

-- Migrate data if old managers table exists
INSERT INTO managers_new (user_id, first_name, last_name, created_at)
SELECT user_id, COALESCE(first_name, ''), COALESCE(last_name, ''), created_at
FROM managers WHERE user_id IS NOT NULL
ON CONFLICT (user_id) DO NOTHING;

-- Drop old and rename new
DROP TABLE IF EXISTS managers CASCADE;
ALTER TABLE managers_new RENAME TO managers;

-- ====== CATEGORY MANAGERS ======
CREATE TABLE IF NOT EXISTS category_managers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES course_categories(id) ON DELETE CASCADE,
    permission_level VARCHAR(50) NOT NULL CHECK (permission_level IN ('view', 'edit', 'admin')) DEFAULT 'edit',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, category_id)
);

CREATE INDEX idx_category_managers_user_id ON category_managers(user_id);
CREATE INDEX idx_category_managers_category_id ON category_managers(category_id);
CREATE INDEX idx_category_managers_user_category ON category_managers(user_id, category_id);
CREATE INDEX idx_category_managers_active ON category_managers(is_active);

-- ====== UPDATE COURSES TABLE ======
ALTER TABLE courses
    ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES course_categories(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS teacher_id UUID REFERENCES teachers(id) ON DELETE SET NULL;

-- Migrate existing teacher_id from owner_teacher_id
UPDATE courses 
SET teacher_id = (SELECT id FROM teachers WHERE user_id = owner_teacher_id)
WHERE owner_teacher_id IS NOT NULL AND teacher_id IS NULL;

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_courses_category_id ON courses(category_id);
CREATE INDEX IF NOT EXISTS idx_courses_teacher_id ON courses(teacher_id);

-- Update users table to support new roles
UPDATE users SET role = 'teacher' WHERE role = 'teacher';
UPDATE users SET role = 'admin' WHERE role = 'admin';
UPDATE users SET role = 'manager' WHERE role = 'manager';
ALTER TABLE users 
    ADD CONSTRAINT check_role CHECK (role IN ('student', 'teacher', 'admin', 'manager', 'category_manager'));
