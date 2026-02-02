-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ====== USERS (roles: student/teacher/manager/admin) ======
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student','teacher','manager','admin')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);

-- ====== STUDENTS / TEACHERS / MANAGERS (profiles) ======
CREATE TABLE IF NOT EXISTS students (
                                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    group_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS teachers (
                                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    department TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS managers (
                                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- ====== COURSES ======
CREATE TABLE IF NOT EXISTS courses (
                                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    owner_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_courses_owner ON courses(owner_teacher_id);
CREATE INDEX IF NOT EXISTS idx_courses_active ON courses(is_active);

-- ====== ENROLLMENTS ======
CREATE TABLE IF NOT EXISTS enrollments (
                                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','pending','rejected','dropped')),
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(course_id, student_id)
    );

CREATE INDEX IF NOT EXISTS idx_enrollments_course ON enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON enrollments(student_id);

-- ====== ASSIGNMENTS (deadlines) ======
CREATE TABLE IF NOT EXISTS assignments (
                                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    due_at TIMESTAMPTZ,
    max_points INT NOT NULL DEFAULT 100 CHECK (max_points >= 0),
    created_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_assignments_course ON assignments(course_id);
CREATE INDEX IF NOT EXISTS idx_assignments_due ON assignments(due_at);

-- ====== SUBMISSIONS ======
CREATE TABLE IF NOT EXISTS submissions (
                                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
    content_text TEXT,
    file_url TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'submitted' CHECK (status IN ('draft','submitted','late','resubmitted')),
    UNIQUE(assignment_id, student_id)
    );

CREATE INDEX IF NOT EXISTS idx_submissions_assignment ON submissions(assignment_id);
CREATE INDEX IF NOT EXISTS idx_submissions_student ON submissions(student_id);

-- ====== GRADES ======
CREATE TABLE IF NOT EXISTS grades (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL UNIQUE REFERENCES submissions(id) ON DELETE CASCADE,
    graded_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    score INT NOT NULL CHECK (score >= 0),
    feedback TEXT,
    graded_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_grades_teacher ON grades(graded_by_teacher_id);

-- ====== ATTENDANCE ======
CREATE TABLE IF NOT EXISTS attendance_sessions (
                                                   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    created_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL
    );

CREATE INDEX IF NOT EXISTS idx_attendance_sessions_course ON attendance_sessions(course_id);

CREATE TABLE IF NOT EXISTS attendance_marks (
                                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('present','absent','late','excused')),
    marked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(session_id, student_id)
    );

CREATE INDEX IF NOT EXISTS idx_attendance_marks_student ON attendance_marks(student_id);

-- ====== CHAT (rooms + messages) ======
CREATE TABLE IF NOT EXISTS chat_rooms (
                                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    room_type TEXT NOT NULL DEFAULT 'course' CHECK (room_type IN ('course','direct','group')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_chat_rooms_course ON chat_rooms(course_id);

CREATE TABLE IF NOT EXISTS chat_room_members (
                                                 room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(room_id, user_id)
    );

CREATE TABLE IF NOT EXISTS chat_messages (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    sender_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_chat_messages_room ON chat_messages(room_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created ON chat_messages(created_at);

-- ====== NOTIFICATIONS ======
CREATE TABLE IF NOT EXISTS notifications (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- e.g. "deadline", "grade", "message"
    title TEXT NOT NULL,
    body TEXT,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(is_read);

-- ====== SCHEDULE (simple) ======
CREATE TABLE IF NOT EXISTS schedule_events (
                                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_schedule_events_course ON schedule_events(course_id);
CREATE INDEX IF NOT EXISTS idx_schedule_events_starts ON schedule_events(starts_at);

-- ====== ANALYTICS (events log) ======
CREATE TABLE IF NOT EXISTS analytics_events (
                                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event_name TEXT NOT NULL,
    meta JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_name ON analytics_events(event_name);

-- ====== PLAGIARISM (placeholder) ======
CREATE TABLE IF NOT EXISTS plagiarism_reports (
                                                  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    score NUMERIC(5,2) NOT NULL DEFAULT 0,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_plagiarism_reports_submission ON plagiarism_reports(submission_id);

-- ====== SESSIONS (if you use DB sessions; Redis later) ======
CREATE TABLE IF NOT EXISTS sessions (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
