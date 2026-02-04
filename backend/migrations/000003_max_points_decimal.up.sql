-- First drop the old constraint
ALTER TABLE assignments
    DROP CONSTRAINT IF EXISTS assignments_max_points_check;

-- Update existing rows with values > 100 to a default percentage (7.5%)
UPDATE assignments SET max_points = 7.5 WHERE max_points > 100;

-- Change max_points from INT to DECIMAL to support percentage values like 7.5%
ALTER TABLE assignments
    ALTER COLUMN max_points TYPE DECIMAL(5,2) USING max_points::DECIMAL(5,2);

-- Add new constraint to allow decimal values between 0 and 100
ALTER TABLE assignments
    ADD CONSTRAINT assignments_max_points_check CHECK (max_points >= 0 AND max_points <= 100);

-- Set default value for new assignments
ALTER TABLE assignments
    ALTER COLUMN max_points SET DEFAULT 7.5;
