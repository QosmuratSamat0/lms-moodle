-- ====== ASSIGNMENTS ======
CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    max_points DECIMAL(10,2) NOT NULL DEFAULT 100 CHECK (max_points > 0),
    due_at TIMESTAMPTZ,
    allow_late BOOLEAN NOT NULL DEFAULT false,
    created_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_assignments_course ON assignments(course_id);
CREATE INDEX idx_assignments_due ON assignments(due_at);
CREATE INDEX idx_assignments_teacher ON assignments(created_by_teacher_id);

-- ====== SUBMISSIONS ======
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(user_id) ON DELETE CASCADE,
    content_text TEXT,
    file_url TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'submitted' CHECK (status IN ('draft','submitted','late','resubmitted')),
    UNIQUE(assignment_id, student_id)
);

CREATE INDEX idx_submissions_assignment ON submissions(assignment_id);
CREATE INDEX idx_submissions_student ON submissions(student_id);
CREATE INDEX idx_submissions_status ON submissions(status);
