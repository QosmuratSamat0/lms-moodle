-- ====== GRADES ======
CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL UNIQUE REFERENCES submissions(id) ON DELETE CASCADE,
    graded_by_teacher_id UUID REFERENCES teachers(user_id) ON DELETE SET NULL,
    score DECIMAL(10,2) NOT NULL CHECK (score >= 0),
    feedback TEXT,
    graded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_grades_submission ON grades(submission_id);
CREATE INDEX idx_grades_teacher ON grades(graded_by_teacher_id);
