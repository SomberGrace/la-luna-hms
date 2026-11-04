-- Persist docket destinations, kitchen progress, and settlement state.
ALTER TABLE restaurant_dockets
  ADD COLUMN IF NOT EXISTS table_number VARCHAR(20),
  ADD COLUMN IF NOT EXISTS notes TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS kitchen_status TEXT NOT NULL DEFAULT 'Pending',
  ADD COLUMN IF NOT EXISTS settled_at TIMESTAMP WITH TIME ZONE;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
      FROM pg_constraint
     WHERE conname = 'restaurant_dockets_kitchen_status_check'
       AND conrelid = 'restaurant_dockets'::regclass
  ) THEN
    ALTER TABLE restaurant_dockets
      ADD CONSTRAINT restaurant_dockets_kitchen_status_check
      CHECK (kitchen_status IN ('Pending', 'Cooking', 'Ready', 'Served'));
  END IF;
END
$$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
      FROM pg_constraint
     WHERE conname = 'restaurant_dockets_settlement_requires_served'
       AND conrelid = 'restaurant_dockets'::regclass
  ) THEN
    ALTER TABLE restaurant_dockets
      ADD CONSTRAINT restaurant_dockets_settlement_requires_served
      CHECK (settled_at IS NULL OR kitchen_status = 'Served');
  END IF;
END
$$;
