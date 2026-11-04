// backend/internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	"backend/internal/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// InitDB initializes PostgreSQL connection pool and verifies connectivity.
func InitDB(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil || cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("empty database URL")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("✓ Connected to Supabase PostgreSQL Database successfully")
	if cfg.RunStartupMigrations {
		runMigrations(db)
	}
	if cfg.SeedDatabase {
		SeedDatabase(db)
	}
	return db, nil
}

func runMigrations(db *sql.DB) {
	_, _ = db.Exec(`ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'GUEST'`)
	_, _ = db.Exec(`ALTER TYPE room_status_enum ADD VALUE IF NOT EXISTS 'Reserved'`)
	_, _ = db.Exec(`ALTER TABLE rooms ADD COLUMN IF NOT EXISTS description TEXT DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE rooms ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE menu_items ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT ''`)

	// Create dedicated guest_accounts table separate from staff/admin user_accounts
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS guest_accounts (
			guest_account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			guest_id UUID NOT NULL REFERENCES guests(guest_id) ON DELETE CASCADE,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_guest_accounts_guest_id ON guest_accounts(guest_id);
		CREATE INDEX IF NOT EXISTS idx_guest_accounts_username ON guest_accounts(username);
		CREATE INDEX IF NOT EXISTS idx_guest_accounts_email ON guest_accounts(email);
	`)
	if err != nil {
		log.Printf("guest_accounts migration notice: %v", err)
	}

	// Add checked_in flag to bookings for explicit front desk check-in tracking
	_, err = db.Exec(`ALTER TABLE bookings ADD COLUMN IF NOT EXISTS checked_in BOOLEAN NOT NULL DEFAULT FALSE`)
	if err != nil {
		log.Printf("checked_in column migration notice: %v", err)
	}

	// Ensure orphaned Occupied rooms (from demo seeds without active booking) are reset to Available
	_, _ = db.Exec(`
		UPDATE rooms 
		   SET status = 'Available', updated_at = CURRENT_TIMESTAMP
		 WHERE status = 'Occupied'
		   AND room_number NOT IN (
		       SELECT room_number FROM bookings 
		        WHERE departure_date >= CURRENT_DATE
		   )
	`)

	// Ensure rooms with pending (un-checked-in) online bookings are marked Reserved, not Occupied
	_, _ = db.Exec(`
		UPDATE rooms 
		   SET status = 'Reserved', updated_at = CURRENT_TIMESTAMP
		 WHERE room_number IN (
		       SELECT b.room_number FROM bookings b
		        LEFT JOIN folios f ON b.card_no = f.card_no
		       WHERE b.checked_in = FALSE
		         AND (f.payment_status IS NULL OR f.payment_status != 'Closed')
		 )
		 AND status = 'Occupied'
	`)

	// Populate images for existing menu items if empty
	defaultImages := map[string]string{
		"Nsima with Stewed Chicken":      "https://images.unsplash.com/photo-1547496502-affa22d38842?auto=format&fit=crop&w=600&q=80",
		"Grilled Lake Tilapia (Chambo)":  "https://images.unsplash.com/photo-1519708227418-c8fd9a32b7a2?auto=format&fit=crop&w=600&q=80",
		"Prime Beef Tenderloin Steak":    "https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=600&q=80",
		"Full English Resort Breakfast":  "https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?auto=format&fit=crop&w=600&q=80",
		"Warm Mandasi & Local Honey":     "https://images.unsplash.com/photo-1589367920969-ab8e050bbb04?auto=format&fit=crop&w=600&q=80",
		"Crispy Chips & Peri-Peri Sauce": "https://images.unsplash.com/photo-1573080496219-bb080dd4f877?auto=format&fit=crop&w=600&q=80",
		"Fresh Mango & Passion Juice":    "https://images.unsplash.com/photo-1613478223719-2ab802602423?auto=format&fit=crop&w=600&q=80",
		"Carlsberg Green Lager (500ml)":  "https://images.unsplash.com/photo-1608270119358-1c4b92b6a50f?auto=format&fit=crop&w=600&q=80",
		"Malawi Gin & Tonic":             "https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?auto=format&fit=crop&w=600&q=80",
		"Coca-Cola 500ml":                "https://images.unsplash.com/photo-1622483767028-3f66f32aef97?auto=format&fit=crop&w=600&q=80",
	}

	for dish, img := range defaultImages {
		_, _ = db.Exec("UPDATE menu_items SET image_url = $1 WHERE item_name = $2 AND (image_url IS NULL OR image_url = '')", img, dish)
	}
}

// SeedDatabase ensures core baseline records exist in the database for each module.
func SeedDatabase(db *sql.DB) {
	seedUserAccounts(db)
	seedRooms(db)
	seedBarItems(db)
	seedKitchenStock(db)
	seedMenuItems(db)
	seedMenuBOM(db)
}

func seedUserAccounts(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user_accounts").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	// Default password hash for "LaLuna2025!"
	hash, err := bcrypt.GenerateFromPassword([]byte("LaLuna2025!"), bcrypt.DefaultCost)
	if err != nil {
		hash = []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lkii")
	}

	users := []struct {
		username string
		email    string
		role     string
	}{
		{"admin", "admin@lalunahms.mw", "ADMIN"},
		{"reception1", "reception@lalunahms.mw", "FRONT_DESK"},
		{"chef1", "kitchen@lalunahms.mw", "CHEF"},
		{"barman1", "bar@lalunahms.mw", "BAR_ATTENDANT"},
		{"waiter1", "restaurant@lalunahms.mw", "WAITER"},
		{"manager1", "manager@lalunahms.mw", "MANAGER"},
	}

	for _, u := range users {
		_, err := db.Exec(
			`INSERT INTO user_accounts (username, email, password_hash, role)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (username) DO NOTHING`,
			u.username, u.email, string(hash), u.role,
		)
		if err != nil {
			log.Printf("Seed user warning: %v", err)
		}
	}
	log.Println("✓ Baseline user accounts seeded into Supabase")
}

func seedRooms(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	rooms := []struct {
		number string
		rType  string
		bed    string
		rate   float64
		status string
	}{
		{"101", "Standard", "Single", 45000, "Available"},
		{"102", "Standard", "Double", 45000, "Available"},
		{"103", "Standard", "Single", 45000, "Available"},
		{"201", "Executive", "Double", 75000, "Available"},
		{"202", "Executive", "Double", 75000, "Available"},
		{"203", "Executive", "Double", 75000, "Available"},
		{"301", "Apartment", "Double", 120000, "Available"},
		{"302", "Apartment", "Double", 120000, "Available"},
	}

	for _, r := range rooms {
		_, _ = db.Exec(
			`INSERT INTO rooms (room_number, room_type, bed_type, room_rate, status)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (room_number) DO NOTHING`,
			r.number, r.rType, r.bed, r.rate, r.status,
		)
	}
	log.Println("✓ Baseline room inventory seeded into Supabase")
}

func seedBarItems(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM bar_items").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	items := []struct {
		name  string
		unit  float64
		cost  float64
		stock int
	}{
		{"Castle Lager (500ml)", 1500, 900, 48},
		{"Carlsberg Green (500ml)", 1200, 700, 36},
		{"Coca-Cola 500ml", 600, 350, 60},
		{"Sprite 500ml", 600, 350, 48},
		{"Mosi Lager (330ml)", 1000, 600, 36},
		{"Malawi Gin 250ml", 3500, 2200, 18},
	}

	for _, i := range items {
		_, _ = db.Exec(
			`INSERT INTO bar_items (item_name, unit_price, cost_price, stock_quantity)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (item_name) DO NOTHING`,
			i.name, i.unit, i.cost, i.stock,
		)
	}
	log.Println("✓ Baseline bar items seeded into Supabase")
}

func seedKitchenStock(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM kitchen_stock").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	stocks := []struct {
		name    string
		opening float64
		add     float64
		sold    float64
		closing float64
	}{
		{"Maize Flour (Gramil)", 50, 25, 18, 57},
		{"Fresh Tilapia (Lake Malawi)", 22, 10, 14, 18},
		{"Whole Chickens", 15, 0, 11, 4},
		{"Cooking Oil (Sunseed)", 20, 10, 6, 24},
		{"Potatoes (Irish)", 35, 0, 22, 13},
		{"Tomatoes & Onions", 18, 5, 12, 11},
	}

	for _, s := range stocks {
		_, _ = db.Exec(
			`INSERT INTO kitchen_stock (ingredient_name, audit_date, opening_stock, add_stock, sold_usage, closing_balance)
			 VALUES ($1, CURRENT_DATE, $2, $3, $4, $5)`,
			s.name, s.opening, s.add, s.sold, s.closing,
		)
	}
	log.Println("✓ Baseline kitchen stock seeded into Supabase")
}

func seedMenuItems(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM menu_items").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	dishes := []struct {
		name     string
		category string
		price    float64
		cost     float64
	}{
		{"Nsima with Stewed Chicken", "Lunch", 3500, 1200},
		{"Grilled Lake Tilapia (Chambo)", "Dinner", 4500, 1800},
		{"Warm Mandasi & Local Honey", "Breakfast", 800, 250},
		{"Coca-Cola 500ml", "Beverage", 600, 350},
		{"Crispy Chips & Peri-Peri Sauce", "Snack", 1200, 450},
		{"Prime Beef Tenderloin Steak", "Dinner", 8500, 3500},
		{"Fresh Mango & Passion Juice", "Beverage", 1500, 600},
	}

	for _, d := range dishes {
		_, _ = db.Exec(
			`INSERT INTO menu_items (item_name, category, price, cost_price, is_available)
			 VALUES ($1, $2, $3, $4, true)
			 ON CONFLICT (item_name) DO NOTHING`,
			d.name, d.category, d.price, d.cost,
		)
	}
	log.Println("✓ Baseline menu items seeded into Supabase")
}

func seedMenuBOM(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM menu_item_ingredients").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	// Link Nsima with Stewed Chicken to Maize Flour, Chicken, Oil
	var nsimaID, tilapiaID string
	var flourID, chickenID, oilID, fishID string

	_ = db.QueryRow("SELECT menu_item_id FROM menu_items WHERE item_name = 'Nsima with Stewed Chicken'").Scan(&nsimaID)
	_ = db.QueryRow("SELECT menu_item_id FROM menu_items WHERE item_name = 'Grilled Lake Tilapia (Chambo)'").Scan(&tilapiaID)

	_ = db.QueryRow("SELECT kitchen_stock_id FROM kitchen_stock WHERE ingredient_name = 'Maize Flour (Gramil)' LIMIT 1").Scan(&flourID)
	_ = db.QueryRow("SELECT kitchen_stock_id FROM kitchen_stock WHERE ingredient_name = 'Whole Chickens' LIMIT 1").Scan(&chickenID)
	_ = db.QueryRow("SELECT kitchen_stock_id FROM kitchen_stock WHERE ingredient_name = 'Cooking Oil (Sunseed)' LIMIT 1").Scan(&oilID)
	_ = db.QueryRow("SELECT kitchen_stock_id FROM kitchen_stock WHERE ingredient_name = 'Fresh Tilapia (Lake Malawi)' LIMIT 1").Scan(&fishID)

	if nsimaID != "" && flourID != "" {
		_, _ = db.Exec("INSERT INTO menu_item_ingredients (menu_item_id, kitchen_stock_id, quantity_required) VALUES ($1, $2, 0.35) ON CONFLICT DO NOTHING", nsimaID, flourID)
	}
	if nsimaID != "" && chickenID != "" {
		_, _ = db.Exec("INSERT INTO menu_item_ingredients (menu_item_id, kitchen_stock_id, quantity_required) VALUES ($1, $2, 0.25) ON CONFLICT DO NOTHING", nsimaID, chickenID)
	}
	if nsimaID != "" && oilID != "" {
		_, _ = db.Exec("INSERT INTO menu_item_ingredients (menu_item_id, kitchen_stock_id, quantity_required) VALUES ($1, $2, 0.05) ON CONFLICT DO NOTHING", nsimaID, oilID)
	}
	if tilapiaID != "" && fishID != "" {
		_, _ = db.Exec("INSERT INTO menu_item_ingredients (menu_item_id, kitchen_stock_id, quantity_required) VALUES ($1, $2, 1.0) ON CONFLICT DO NOTHING", tilapiaID, fishID)
	}
	log.Println("✓ Baseline recipe BOM links seeded into Supabase")
}
