-- ====== ATTENDANCE ======
CREATE TABLE attendance_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ,
    created_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_attendance_sessions_course ON attendance_sessions(course_id);
CREATE INDEX idx_attendance_sessions_starts ON attendance_sessions(starts_at);

CREATE TABLE attendance_marks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('present','absent','late','excused')),
    marked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(session_id, student_id)
);

CREATE INDEX idx_attendance_marks_session ON attendance_marks(session_id);
CREATE INDEX idx_attendance_marks_student ON attendance_marks(student_id);
