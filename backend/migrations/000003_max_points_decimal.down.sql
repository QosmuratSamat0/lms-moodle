-- Revert max_points back to INT
ALTER TABLE assignments
    DROP CONSTRAINT IF EXISTS assignments_max_points_check;

ALTER TABLE assignments
    ALTER COLUMN max_points TYPE INT USING max_points::INT;

ALTER TABLE assignments
    ADD CONSTRAINT assignments_max_points_check CHECK (max_points >= 0);
