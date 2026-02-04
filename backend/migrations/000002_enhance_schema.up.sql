-- Add groups table for student grouping (e.g., SE-2430)
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code TEXT NOT NULL UNIQUE, -- e.g., 'SE-2430'
    name TEXT,
    description TEXT,
    year_of_admission INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_groups_code ON groups(code);

-- Add teacher_course_groups mapping table
CREATE TABLE IF NOT EXISTS teacher_course_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID NOT NULL REFERENCES teachers(user_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(teacher_id, course_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_teacher_course_groups_teacher ON teacher_course_groups(teacher_id);
CREATE INDEX IF NOT EXISTS idx_teacher_course_groups_course ON teacher_course_groups(course_id);
CREATE INDEX IF NOT EXISTS idx_teacher_course_groups_group ON teacher_course_groups(group_id);

-- Enhance students table
ALTER TABLE students 
    ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS major TEXT,
    ADD COLUMN IF NOT EXISTS year_of_study INT,
    ADD COLUMN IF NOT EXISTS gpa NUMERIC(3,2),
    ADD COLUMN IF NOT EXISTS enrollment_status TEXT DEFAULT 'enrolled' CHECK (enrollment_status IN ('enrolled', 'on_leave', 'graduated', 'dropped'));

CREATE INDEX IF NOT EXISTS idx_students_group ON students(group_id);
CREATE INDEX IF NOT EXISTS idx_students_enrollment_status ON students(enrollment_status);

-- Enhance assignments table with allow_late flag
ALTER TABLE assignments 
    ADD COLUMN IF NOT EXISTS allow_late BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES groups(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_assignments_group ON assignments(group_id);

-- Add attachments table for file uploads
CREATE TABLE IF NOT EXISTS attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('submission', 'assignment', 'course', 'profile', 'chat')),
    reference_id UUID, -- Generic reference to related entity
    public_id TEXT NOT NULL, -- Cloudinary public ID
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    original_name TEXT NOT NULL,
    format TEXT,
    resource_type TEXT NOT NULL, -- 'image', 'raw', 'video', etc.
    size BIGINT NOT NULL,
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_attachments_uploaded_by ON attachments(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_attachments_type ON attachments(type);
CREATE INDEX IF NOT EXISTS idx_attachments_reference ON attachments(type, reference_id);
CREATE INDEX IF NOT EXISTS idx_attachments_public_id ON attachments(public_id);

-- Enhance attendance_sessions to link to groups
ALTER TABLE attendance_sessions 
    ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS title TEXT,
    ADD COLUMN IF NOT EXISTS session_date DATE;

CREATE INDEX IF NOT EXISTS idx_attendance_sessions_group ON attendance_sessions(group_id);

-- Enhance attendance_marks with notes and marked_by
ALTER TABLE attendance_marks 
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS marked_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Update submissions to store multiple attachments (reference via attachments table)
-- The existing file_url can remain for backwards compatibility

-- Add submission_attachments junction table for multiple files per submission
CREATE TABLE IF NOT EXISTS submission_attachments (
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    PRIMARY KEY(submission_id, attachment_id)
);

-- Add assignment_attachments junction table for assignment materials
CREATE TABLE IF NOT EXISTS assignment_attachments (
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    PRIMARY KEY(assignment_id, attachment_id)
);
