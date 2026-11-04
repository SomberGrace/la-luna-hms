// backend/cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/services"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Load configuration settings (reads .env, DATABASE_URL, SUPABASE_DB_URL, JWT_SECRET, PORT)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be set to a secret of at least 32 bytes")
	}

	// 2. Initialize auth service
	authService := services.NewAuthService(cfg.JWTSecret)
	authService.SetDemoUsersEnabled(cfg.AllowDemoUsers)

	// 3. Connect to Supabase PostgreSQL with automated baseline seeder
	var sqlDB *sql.DB
	if cfg.DatabaseURL != "" {
		var err error
		sqlDB, err = db.InitDB(cfg)
		if err == nil {
			authService.SetDB(sqlDB)
		} else {
			log.Printf("⚠️ Supabase connection failed (%v).", err)
		}
	} else {
		log.Println("⚠️ No DATABASE_URL found. Running without a database.")
	}

	// 4. Initialize Handlers for each Operational Module
	authHandler := handlers.NewAuthHandler(authService)
	frontDeskHandler := handlers.NewFrontDeskHandler(sqlDB)
	restaurantHandler := handlers.NewRestaurantHandler(sqlDB)
	kitchenHandler := handlers.NewKitchenHandler(sqlDB)
	barHandler := handlers.NewBarHandler(sqlDB)
	managerHandler := handlers.NewManagerHandler(sqlDB)
	guestHandler := handlers.NewGuestHandler(sqlDB, cfg.JWTSecret)

	// =========================================================================
	// ROUTE REGISTRATIONS
	// =========================================================================

	// General Health
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		dbStatus := "In-Memory / Offline"
		if sqlDB != nil && sqlDB.Ping() == nil {
			dbStatus = "Connected to Supabase PostgreSQL"
		}
		fmt.Fprintf(w, `{"status":"ok","service":"La Luna HMS Unified API","port":"%s","database":"%s"}`, cfg.Port, dbStatus)
	})

	// 1. Authentication & User Management
	http.HandleFunc("/api/v1/auth/login", authHandler.HandleLogin)
	http.HandleFunc("/api/v1/auth/register", authHandler.HandleRegister)

	// 2. Front Desk & Room Management (PMS) & Customer Portal
	http.Handle("/api/v1/frontdesk/health", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.HealthCheck))
	http.Handle("/api/v1/frontdesk/rooms", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.GetRooms))
	http.Handle("/api/v1/frontdesk/rooms/status", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.UpdateRoomStatus))
	http.Handle("/api/v1/frontdesk/checkin", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.CheckIn))
	http.Handle("/api/v1/frontdesk/walkin", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.CheckIn))
	http.Handle("/api/v1/frontdesk/folio", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.GetFolio))
	http.Handle("/api/v1/frontdesk/payment", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.RecordPayment))
	http.Handle("/api/v1/frontdesk/payments", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.RecordPayment))
	http.Handle("/api/v1/frontdesk/checkout", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.CheckOut))
	http.Handle("/api/v1/frontdesk/extension", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.ExtendStay))
	http.Handle("/api/v1/frontdesk/reservations", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.GetReservations))
	http.Handle("/api/v1/frontdesk/reports", handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")(frontDeskHandler.DailyReport))

	// 3. Restaurant POS & Table Management
	http.Handle("/api/v1/restaurant/menu", handlers.RequireRole(cfg.JWTSecret, "WAITER")(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			restaurantHandler.AddMenuItem(w, r)
		} else {
			restaurantHandler.GetMenu(w, r)
		}
	}))
	http.Handle("/api/v1/restaurant/menu/availability", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.SetMenuAvailability))
	http.Handle("/api/v1/restaurant/dockets", handlers.RequireRole(cfg.JWTSecret, "WAITER")(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			restaurantHandler.CreateDocket(w, r)
		} else {
			restaurantHandler.GetDockets(w, r)
		}
	}))
	http.Handle("/api/v1/restaurant/dockets/settle", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.SettleDocket))
	http.Handle("/api/v1/restaurant/occupied-rooms", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.GetOccupiedRooms))
	http.Handle("/api/v1/restaurant/summary", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.Summary))
	http.Handle("/api/v1/restaurant/reports", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.Reports))
	http.Handle("/api/v1/restaurant/breakfast", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.GetBreakfastGuests))
	http.Handle("/api/v1/restaurant/breakfast/reconcile", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.ReconcileBreakfast))
	http.Handle("/api/v1/restaurant/dinner/reconcile", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.ReconcileDinner))
	http.Handle("/api/v1/restaurant/reconcile", handlers.RequireRole(cfg.JWTSecret, "WAITER")(restaurantHandler.ReconcileMeal))

	// 4. Kitchen Display System (KDS) & Stock Management
	http.Handle("/api/v1/kitchen/stock", handlers.RequireRole(cfg.JWTSecret, "CHEF")(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			kitchenHandler.SaveStock(w, r)
		} else {
			kitchenHandler.GetStock(w, r)
		}
	}))
	http.Handle("/api/v1/kitchen/recipes", handlers.RequireRole(cfg.JWTSecret, "CHEF")(kitchenHandler.GetRecipes))
	http.Handle("/api/v1/kitchen/reports", handlers.RequireRole(cfg.JWTSecret, "CHEF")(kitchenHandler.SubmitReport))
	http.Handle("/api/v1/kitchen/dockets", handlers.RequireRole(cfg.JWTSecret, "CHEF")(kitchenHandler.GetDockets))
	http.Handle("/api/v1/kitchen/dockets/status", handlers.RequireRole(cfg.JWTSecret, "CHEF")(kitchenHandler.UpdateDocketStatus))

	// 5. Bar POS & Stock Ledger
	http.Handle("/api/v1/bar/items", handlers.RequireRole(cfg.JWTSecret, "BAR_ATTENDANT")(barHandler.GetItems))
	http.Handle("/api/v1/bar/sales", handlers.RequireRole(cfg.JWTSecret, "BAR_ATTENDANT")(barHandler.RecordSale))
	http.Handle("/api/v1/bar/ledger", handlers.RequireRole(cfg.JWTSecret, "BAR_ATTENDANT")(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			barHandler.SaveLedger(w, r)
		} else {
			barHandler.GetLedger(w, r)
		}
	}))
	http.Handle("/api/v1/bar/reports", handlers.RequireRole(cfg.JWTSecret, "BAR_ATTENDANT")(barHandler.Reports))

	// 6. Centralized Management Dashboard & Analytics
	http.Handle("/api/v1/manager/summary", handlers.RequireRole(cfg.JWTSecret, "MANAGER", "ADMIN")(managerHandler.Summary))
	http.Handle("/api/v1/manager/audits", handlers.RequireRole(cfg.JWTSecret, "MANAGER", "ADMIN")(managerHandler.GetAudits))
	http.Handle("/api/v1/manager/audits/action", handlers.RequireRole(cfg.JWTSecret, "MANAGER", "ADMIN")(managerHandler.AuditAction))
	http.Handle("/api/v1/manager/consolidated", handlers.RequireRole(cfg.JWTSecret, "MANAGER", "ADMIN")(managerHandler.ConsolidatedReport))

	// 7. Guest Booking Portal (Public & Authenticated)
	http.HandleFunc("/api/v1/guest/rooms", guestHandler.GetAvailableRooms)
	http.Handle("/api/v1/guest/profile", handlers.RequireRole(cfg.JWTSecret, "GUEST")(guestHandler.GetProfile))
	http.Handle("/api/v1/guest/book", handlers.RequireRole(cfg.JWTSecret, "GUEST")(guestHandler.CreateBooking))
	http.Handle("/api/v1/guest/bookings", handlers.RequireRole(cfg.JWTSecret, "GUEST")(guestHandler.GetMyBookings))
	http.Handle("/api/v1/guest/folio", handlers.RequireRole(cfg.JWTSecret, "GUEST")(guestHandler.GetMyFolio))
	http.Handle("/api/v1/guest/checkin", handlers.RequireRole(cfg.JWTSecret, "GUEST")(guestHandler.CheckIn))
	http.HandleFunc("/api/v1/guest/resident-login", guestHandler.ResidentLogin)

	fmt.Printf("✓ La Luna Hotel Management System Unified API active on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handlers.CORSMiddleware(http.DefaultServeMux)))
}
