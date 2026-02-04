-- Revert grade score back to INT
ALTER TABLE grades DROP CONSTRAINT IF EXISTS grades_score_check;

ALTER TABLE grades
    ALTER COLUMN score TYPE INT USING score::INT;

ALTER TABLE grades
    ADD CONSTRAINT grades_score_check CHECK (score >= 0);
