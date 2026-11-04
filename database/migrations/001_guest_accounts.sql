-- 001_guest_accounts.sql
-- Adds guest self-registration support for the booking portal.

-- Add GUEST role to user_role enum  
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'GUEST';

-- Add room description/image columns to rooms table for the booking portal
ALTER TABLE rooms ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';
ALTER TABLE rooms ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT '';

-- Update existing rooms with descriptions
UPDATE rooms SET description = CASE room_type::text
  WHEN 'Standard' THEN 'Comfortable room with modern amenities, en-suite bathroom, and garden views. Perfect for business or leisure stays.'
  WHEN 'Executive' THEN 'Spacious executive room with premium furnishings, work desk, mini-bar, and panoramic views of the grounds.'
  WHEN 'Apartment' THEN 'Full apartment-style suite with separate living area, kitchenette, dining space, and premium amenities for extended stays.'
  WHEN 'Sector Suite' THEN 'Our finest accommodation with luxury furnishings, private balcony, jacuzzi bath, and dedicated concierge service.'
  ELSE 'Comfortable room with standard amenities.'
END
WHERE description = '' OR description IS NULL;

UPDATE rooms SET image_url = CASE room_type::text
  WHEN 'Standard' THEN '/images/rooms/standard.jpg'
  WHEN 'Executive' THEN '/images/rooms/executive.jpg'
  WHEN 'Apartment' THEN '/images/rooms/apartment.jpg'
  WHEN 'Sector Suite' THEN '/images/rooms/suite.jpg'
  ELSE '/images/rooms/standard.jpg'
END
WHERE image_url = '' OR image_url IS NULL;
