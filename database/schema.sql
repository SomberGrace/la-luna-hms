-- =============================================================================
-- La Luna Hotel Management System - Comprehensive PostgreSQL Schema & DDL
-- Database Engine: PostgreSQL 14+
-- Stack: Golang Backend + Svelte/Tailwind Frontend
-- =============================================================================

-- Enable UUID extension for standard GUID primary key generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clean up existing tables and custom types (Reverse Dependency Order)
DROP TABLE IF EXISTS report_documents CASCADE;
DROP TABLE IF EXISTS consolidated_reports CASCADE;
DROP TABLE IF EXISTS executive_summaries CASCADE;
DROP TABLE IF EXISTS restaurant_reports CASCADE;
DROP TABLE IF EXISTS kitchen_reports CASCADE;
DROP TABLE IF EXISTS bar_reports CASCADE;
DROP TABLE IF EXISTS docket_items CASCADE;
DROP TABLE IF EXISTS restaurant_dockets CASCADE;
DROP TABLE IF EXISTS menu_item_ingredients CASCADE;
DROP TABLE IF EXISTS menu_items CASCADE;
DROP TABLE IF EXISTS kitchen_stock CASCADE;
DROP TABLE IF EXISTS bar_ledgers CASCADE;
DROP TABLE IF EXISTS bar_items CASCADE;
DROP TABLE IF EXISTS folios CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS guests CASCADE;
DROP TABLE IF EXISTS user_accounts CASCADE;

DROP TYPE IF EXISTS user_role CASCADE;
DROP TYPE IF EXISTS room_type_enum CASCADE;
DROP TYPE IF EXISTS bed_type_enum CASCADE;
DROP TYPE IF EXISTS room_status_enum CASCADE;
DROP TYPE IF EXISTS payment_mode_enum CASCADE;
DROP TYPE IF EXISTS folio_status_enum CASCADE;
DROP TYPE IF EXISTS menu_category_enum CASCADE;
DROP TYPE IF EXISTS report_status_enum CASCADE;

-- =============================================================================
-- 1. ENUMERATED TYPES (Domain Validation Rules)
-- =============================================================================

CREATE TYPE user_role AS ENUM (
  'ADMIN',
  'MANAGER',
  'FRONT_DESK',
  'CHEF',
  'BAR_ATTENDANT',
  'WAITER'
);

CREATE TYPE room_type_enum AS ENUM (
  'Standard',
  'Executive',
  'Apartment',
  'Sector Suite'
);

CREATE TYPE bed_type_enum AS ENUM (
  'Single',
  'Double'
);

CREATE TYPE room_status_enum AS ENUM (
  'Available',
  'Occupied',
  'Maintenance',
  'Reserved'
);

CREATE TYPE payment_mode_enum AS ENUM (
  'Cash',
  'Credit Card',
  'Cheque',
  'Room Charge'
);

CREATE TYPE folio_status_enum AS ENUM (
  'Open',
  'Settled',
  'Closed'
);

CREATE TYPE menu_category_enum AS ENUM (
  'Breakfast',
  'Lunch',
  'Dinner',
  'Snack',
  'Beverage'
);

CREATE TYPE report_status_enum AS ENUM (
  'Submitted',
  'Approved',
  'Disputed'
);

-- =============================================================================
-- 2. USER AUTHENTICATION & SECURITY
-- =============================================================================

CREATE TABLE user_accounts (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role user_role NOT NULL DEFAULT 'FRONT_DESK',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 3. FRONT DESK & ACCOMMODATIONS
-- =============================================================================

CREATE TABLE guests (
  guest_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  address TEXT,
  contact_no VARCHAR(20) NOT NULL,
  email VARCHAR(100),
  nationality VARCHAR(50),
  city VARCHAR(50),
  passport_no VARCHAR(50),
  postal_code VARCHAR(20),
  vehicle_registration VARCHAR(20),
  business_name VARCHAR(100),
  business_address TEXT,
  business_position VARCHAR(50),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guest_accounts (
  guest_account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  guest_id UUID NOT NULL REFERENCES guests(guest_id) ON DELETE CASCADE,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rooms (
  room_number VARCHAR(20) PRIMARY KEY,
  room_type room_type_enum NOT NULL DEFAULT 'Standard',
  bed_type bed_type_enum NOT NULL DEFAULT 'Single',
  room_rate NUMERIC(12, 2) NOT NULL CHECK (room_rate >= 0),
  status room_status_enum NOT NULL DEFAULT 'Available',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
  card_no VARCHAR(50) PRIMARY KEY, -- Guest Registration Card No (e.g., CARD 001)
  guest_id UUID NOT NULL REFERENCES guests(guest_id) ON DELETE RESTRICT,
  room_number VARCHAR(20) NOT NULL REFERENCES rooms(room_number) ON DELETE RESTRICT,
  arrival_date DATE NOT NULL,
  departure_date DATE NOT NULL,
  number_of_days INT GENERATED ALWAYS AS (departure_date - arrival_date) STORED,
  mode_of_payment payment_mode_enum NOT NULL DEFAULT 'Cash',
  cheque_no VARCHAR(50),
  advance_paid NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (advance_paid >= 0),
  front_officer_clerk VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_booking_dates CHECK (departure_date >= arrival_date)
);

CREATE TABLE folios (
  folio_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  card_no VARCHAR(50) NOT NULL UNIQUE REFERENCES bookings(card_no) ON DELETE CASCADE,
  guest_id UUID NOT NULL REFERENCES guests(guest_id) ON DELETE RESTRICT,
  room_charges NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (room_charges >= 0),
  incidental_charges NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (incidental_charges >= 0),
  total_balance NUMERIC(12, 2) GENERATED ALWAYS AS (room_charges + incidental_charges) STORED,
  payment_status folio_status_enum NOT NULL DEFAULT 'Open',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 4. BAR DEPARTMENT INVENTORY & AUDIT
-- =============================================================================

CREATE TABLE bar_items (
  bar_item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  item_name VARCHAR(100) NOT NULL UNIQUE,
  unit_price NUMERIC(10, 2) NOT NULL CHECK (unit_price >= 0),
  cost_price NUMERIC(10, 2) NOT NULL CHECK (cost_price >= 0),
  stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bar_ledgers (
  bar_ledger_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  bar_item_id UUID NOT NULL REFERENCES bar_items(bar_item_id) ON DELETE RESTRICT,
  audit_date DATE NOT NULL DEFAULT CURRENT_DATE,
  opening_stock NUMERIC(10, 2) NOT NULL CHECK (opening_stock >= 0),
  add_stock NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (add_stock >= 0),
  total_stock NUMERIC(10, 2) GENERATED ALWAYS AS (opening_stock + add_stock) STORED,
  sold_quantity NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (sold_quantity >= 0),
  closing_stock NUMERIC(10, 2) NOT NULL CHECK (closing_stock >= 0),
  sales_revenue NUMERIC(12, 2) NOT NULL CHECK (sales_revenue >= 0),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_bar_stock_balance CHECK (closing_stock <= (opening_stock + add_stock))
);

CREATE TABLE bar_reports (
  bar_report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES user_accounts(user_id) ON DELETE RESTRICT,
  report_date DATE NOT NULL DEFAULT CURRENT_DATE,
  total_bar_sales NUMERIC(12, 2) NOT NULL CHECK (total_bar_sales >= 0),
  total_cogs NUMERIC(12, 2) NOT NULL CHECK (total_cogs >= 0),
  gross_margin NUMERIC(12, 2) GENERATED ALWAYS AS (total_bar_sales - total_cogs) STORED,
  status report_status_enum NOT NULL DEFAULT 'Submitted',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 5. KITCHEN DEPARTMENT & RECIPE BOM
-- =============================================================================

CREATE TABLE kitchen_stock (
  kitchen_stock_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  ingredient_name VARCHAR(100) NOT NULL,
  audit_date DATE NOT NULL DEFAULT CURRENT_DATE,
  opening_stock NUMERIC(10, 2) NOT NULL CHECK (opening_stock >= 0),
  add_stock NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (add_stock >= 0),
  total_stock NUMERIC(10, 2) GENERATED ALWAYS AS (opening_stock + add_stock) STORED,
  sold_usage NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (sold_usage >= 0),
  closing_balance NUMERIC(10, 2) NOT NULL CHECK (closing_balance >= 0),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_kitchen_stock_balance CHECK (closing_balance <= (opening_stock + add_stock))
);

CREATE TABLE kitchen_reports (
  kitchen_report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES user_accounts(user_id) ON DELETE RESTRICT,
  report_date DATE NOT NULL DEFAULT CURRENT_DATE,
  stock_summary_json JSONB,
  status report_status_enum NOT NULL DEFAULT 'Submitted',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 6. RESTAURANT & DINING MODULE
-- =============================================================================

CREATE TABLE menu_items (
  menu_item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  item_name VARCHAR(100) NOT NULL UNIQUE,
  category menu_category_enum NOT NULL DEFAULT 'Lunch',
  price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
  cost_price NUMERIC(10, 2) NOT NULL DEFAULT 0.00 CHECK (cost_price >= 0),
  is_available BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction Table: Bill of Materials (BOM) linking dishes to kitchen ingredients
CREATE TABLE menu_item_ingredients (
  relation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  menu_item_id UUID NOT NULL REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
  kitchen_stock_id UUID NOT NULL REFERENCES kitchen_stock(kitchen_stock_id) ON DELETE RESTRICT,
  quantity_required NUMERIC(10, 3) NOT NULL CHECK (quantity_required > 0),
  CONSTRAINT unique_menu_ingredient UNIQUE (menu_item_id, kitchen_stock_id)
);

CREATE TABLE restaurant_dockets (
  docket_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  room_number VARCHAR(20) REFERENCES rooms(room_number) ON DELETE SET NULL,
  table_number VARCHAR(20),
  folio_id UUID REFERENCES folios(folio_id) ON DELETE SET NULL,
  sale_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_amount >= 0),
  payment_method payment_mode_enum NOT NULL DEFAULT 'Cash',
  notes TEXT NOT NULL DEFAULT '',
  kitchen_status TEXT NOT NULL DEFAULT 'Pending'
    CHECK (kitchen_status IN ('Pending', 'Cooking', 'Ready', 'Served')),
  settled_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT restaurant_dockets_settlement_requires_served
    CHECK (settled_at IS NULL OR kitchen_status = 'Served')
);

CREATE TABLE docket_items (
  docket_item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  docket_id UUID NOT NULL REFERENCES restaurant_dockets(docket_id) ON DELETE CASCADE,
  menu_item_id UUID NOT NULL REFERENCES menu_items(menu_item_id) ON DELETE RESTRICT,
  quantity INT NOT NULL CHECK (quantity > 0),
  subtotal NUMERIC(10, 2) NOT NULL CHECK (subtotal >= 0)
);

CREATE TABLE restaurant_reports (
  restaurant_report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES user_accounts(user_id) ON DELETE RESTRICT,
  report_date DATE NOT NULL DEFAULT CURRENT_DATE,
  total_restaurant_sales NUMERIC(12, 2) NOT NULL CHECK (total_restaurant_sales >= 0),
  total_breakfast_covers INT NOT NULL DEFAULT 0 CHECK (total_breakfast_covers >= 0),
  total_lunch_dinner_covers INT NOT NULL DEFAULT 0 CHECK (total_lunch_dinner_covers >= 0),
  status report_status_enum NOT NULL DEFAULT 'Submitted',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 7. EXECUTIVE REPORTING & FILE POINTER ARCHIVAL
-- =============================================================================

CREATE TABLE report_documents (
  report_doc_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  department VARCHAR(50) NOT NULL,
  report_title VARCHAR(150) NOT NULL,
  storage_provider VARCHAR(50) NOT NULL DEFAULT 'Google Drive',
  file_pointer_url TEXT NOT NULL,
  generated_by UUID REFERENCES user_accounts(user_id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE executive_summaries (
  summary_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  audit_date DATE NOT NULL UNIQUE DEFAULT CURRENT_DATE,
  total_room_sales_cash NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_room_sales_cash >= 0),
  total_room_sales_cheque NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_room_sales_cheque >= 0),
  total_room_sales_visa NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_room_sales_visa >= 0),
  total_room_sales_credit NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_room_sales_credit >= 0),
  restaurant_sales_total NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (restaurant_sales_total >= 0),
  bar_sales_total NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (bar_sales_total >= 0),
  petty_cash NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (petty_cash >= 0),
  banking_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (banking_amount >= 0),
  grand_total NUMERIC(12, 2) GENERATED ALWAYS AS (
    total_room_sales_cash + total_room_sales_cheque + total_room_sales_visa +
    total_room_sales_credit + restaurant_sales_total + bar_sales_total - petty_cash
  ) STORED,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE consolidated_reports (
  report_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  manager_id UUID NOT NULL REFERENCES user_accounts(user_id) ON DELETE RESTRICT,
  report_date DATE NOT NULL DEFAULT CURRENT_DATE,
  total_revenue NUMERIC(12, 2) NOT NULL CHECK (total_revenue >= 0),
  total_expenses NUMERIC(12, 2) NOT NULL DEFAULT 0.00 CHECK (total_expenses >= 0),
  department_breakdown JSONB,
  report_doc_id UUID REFERENCES report_documents(report_doc_id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 8. INDEXES FOR PERFORMANCE OPTIMIZATION
-- =============================================================================

CREATE INDEX idx_bookings_guest_id ON bookings(guest_id);
CREATE INDEX idx_bookings_room_number ON bookings(room_number);
CREATE INDEX idx_bookings_dates ON bookings(arrival_date, departure_date);
CREATE INDEX idx_folios_card_no ON folios(card_no);
CREATE INDEX idx_restaurant_dockets_room ON restaurant_dockets(room_number);
CREATE INDEX idx_restaurant_dockets_folio ON restaurant_dockets(folio_id);
CREATE INDEX idx_docket_items_docket ON docket_items(docket_id);
CREATE INDEX idx_bar_ledgers_date ON bar_ledgers(audit_date);
CREATE INDEX idx_kitchen_stock_date ON kitchen_stock(audit_date);
CREATE INDEX idx_executive_summaries_date ON executive_summaries(audit_date);

-- =============================================================================
-- 9. AUTOMATED TRIGGERS & BUSINESS LOGIC FUNCTIONS
-- =============================================================================

-- Trigger Function: Auto-update room status to 'Occupied' upon booking check-in
CREATE OR REPLACE FUNCTION update_room_on_booking()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE rooms SET status = 'Occupied', updated_at = CURRENT_TIMESTAMP WHERE room_number = NEW.room_number;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_after_booking_insert ON bookings;
CREATE TRIGGER trg_after_booking_insert
AFTER INSERT ON bookings
FOR EACH ROW
EXECUTE FUNCTION update_room_on_booking();

-- Trigger Function: Auto-post room-charged restaurant bills directly to Guest Folio
CREATE OR REPLACE FUNCTION post_docket_to_folio()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.payment_method = 'Room Charge' AND NEW.folio_id IS NOT NULL THEN
    UPDATE folios
    SET incidental_charges = incidental_charges + NEW.total_amount,
        updated_at = CURRENT_TIMESTAMP
    WHERE folio_id = NEW.folio_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_after_docket_insert ON restaurant_dockets;
CREATE TRIGGER trg_after_docket_insert
AFTER INSERT ON restaurant_dockets
FOR EACH ROW
EXECUTE FUNCTION post_docket_to_folio();

-- Optional Initial Seed Data for Rooms
INSERT INTO rooms (room_number, room_type, bed_type, room_rate, status)
VALUES
  ('101', 'Standard', 'Single', 45000.00, 'Available'),
  ('102', 'Standard', 'Double', 45000.00, 'Occupied'),
  ('103', 'Standard', 'Single', 45000.00, 'Available'),
  ('201', 'Executive', 'Double', 75000.00, 'Available'),
  ('202', 'Executive', 'Double', 75000.00, 'Available'),
  ('203', 'Executive', 'Double', 75000.00, 'Maintenance'),
  ('301', 'Apartment', 'Double', 120000.00, 'Available'),
  ('302', 'Apartment', 'Double', 120000.00, 'Occupied')
ON CONFLICT (room_number) DO NOTHING;
