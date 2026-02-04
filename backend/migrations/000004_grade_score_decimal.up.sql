-- Change grade score from INT to DECIMAL to support percentage values like 7.5%

-- First drop any existing constraints
ALTER TABLE grades DROP CONSTRAINT IF EXISTS grades_score_check;

-- Update existing scores (if any are over 100, convert them)
-- Keep scores as-is since they should be the actual grade
UPDATE grades SET score = LEAST(score, 100) WHERE score > 100;

-- Change score from INT to DECIMAL
ALTER TABLE grades
    ALTER COLUMN score TYPE DECIMAL(5,2) USING score::DECIMAL(5,2);

-- Add constraint to ensure scores are between 0 and 100
ALTER TABLE grades
    ADD CONSTRAINT grades_score_check CHECK (score >= 0 AND score <= 100);
