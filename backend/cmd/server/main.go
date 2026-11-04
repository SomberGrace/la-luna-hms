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
	// Load configuration and validate the signing secret.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be set to a secret of at least 32 bytes")
	}

	// Set up authentication.
	authService := services.NewAuthService(cfg.JWTSecret)
	authService.SetDemoUsersEnabled(cfg.AllowDemoUsers)

	// Connect to PostgreSQL when a database URL is available.
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

	// Set up the API handlers.
	authHandler := handlers.NewAuthHandler(authService)
	frontDeskHandler := handlers.NewFrontDeskHandler(sqlDB)
	restaurantHandler := handlers.NewRestaurantHandler(sqlDB)
	kitchenHandler := handlers.NewKitchenHandler(sqlDB)
	barHandler := handlers.NewBarHandler(sqlDB)
	managerHandler := handlers.NewManagerHandler(sqlDB)
	guestHandler := handlers.NewGuestHandler(sqlDB, cfg.JWTSecret)

	frontDeskOnly := handlers.RequireRole(cfg.JWTSecret, "FRONT_DESK")
	restaurantOnly := handlers.RequireRole(cfg.JWTSecret, "WAITER")
	kitchenOnly := handlers.RequireRole(cfg.JWTSecret, "CHEF")
	barOnly := handlers.RequireRole(cfg.JWTSecret, "BAR_ATTENDANT")
	managerOnly := handlers.RequireRole(cfg.JWTSecret, "MANAGER", "ADMIN")
	guestOnly := handlers.RequireRole(cfg.JWTSecret, "GUEST")

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

	// Authentication
	http.HandleFunc("/api/v1/auth/login", authHandler.HandleLogin)
	http.HandleFunc("/api/v1/auth/register", authHandler.HandleRegister)

	// Front desk and room management
	http.Handle("/api/v1/frontdesk/health", frontDeskOnly(frontDeskHandler.HealthCheck))
	http.Handle("/api/v1/frontdesk/rooms", frontDeskOnly(frontDeskHandler.GetRooms))
	http.Handle("/api/v1/frontdesk/rooms/status", frontDeskOnly(frontDeskHandler.UpdateRoomStatus))
	http.Handle("/api/v1/frontdesk/checkin", frontDeskOnly(frontDeskHandler.CheckIn))
	http.Handle("/api/v1/frontdesk/walkin", frontDeskOnly(frontDeskHandler.CheckIn))
	http.Handle("/api/v1/frontdesk/folio", frontDeskOnly(frontDeskHandler.GetFolio))
	http.Handle("/api/v1/frontdesk/payment", frontDeskOnly(frontDeskHandler.RecordPayment))
	http.Handle("/api/v1/frontdesk/payments", frontDeskOnly(frontDeskHandler.RecordPayment))
	http.Handle("/api/v1/frontdesk/checkout", frontDeskOnly(frontDeskHandler.CheckOut))
	http.Handle("/api/v1/frontdesk/extension", frontDeskOnly(frontDeskHandler.ExtendStay))
	http.Handle("/api/v1/frontdesk/reservations", frontDeskOnly(frontDeskHandler.GetReservations))
	http.Handle("/api/v1/frontdesk/reports", frontDeskOnly(frontDeskHandler.DailyReport))

	// Restaurant and table management
	http.Handle("/api/v1/restaurant/menu", restaurantOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			restaurantHandler.AddMenuItem(w, r)
		} else {
			restaurantHandler.GetMenu(w, r)
		}
	}))
	http.Handle("/api/v1/restaurant/menu/availability", restaurantOnly(restaurantHandler.SetMenuAvailability))
	http.Handle("/api/v1/restaurant/dockets", restaurantOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			restaurantHandler.CreateDocket(w, r)
		} else {
			restaurantHandler.GetDockets(w, r)
		}
	}))
	http.Handle("/api/v1/restaurant/dockets/settle", restaurantOnly(restaurantHandler.SettleDocket))
	http.Handle("/api/v1/restaurant/occupied-rooms", restaurantOnly(restaurantHandler.GetOccupiedRooms))
	http.Handle("/api/v1/restaurant/summary", restaurantOnly(restaurantHandler.Summary))
	http.Handle("/api/v1/restaurant/reports", restaurantOnly(restaurantHandler.Reports))
	http.Handle("/api/v1/restaurant/breakfast", restaurantOnly(restaurantHandler.GetBreakfastGuests))
	http.Handle("/api/v1/restaurant/breakfast/reconcile", restaurantOnly(restaurantHandler.ReconcileBreakfast))
	http.Handle("/api/v1/restaurant/dinner/reconcile", restaurantOnly(restaurantHandler.ReconcileDinner))
	http.Handle("/api/v1/restaurant/reconcile", restaurantOnly(restaurantHandler.ReconcileMeal))

	// Kitchen and stock management
	http.Handle("/api/v1/kitchen/stock", kitchenOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			kitchenHandler.SaveStock(w, r)
		} else {
			kitchenHandler.GetStock(w, r)
		}
	}))
	http.Handle("/api/v1/kitchen/recipes", kitchenOnly(kitchenHandler.GetRecipes))
	http.Handle("/api/v1/kitchen/reports", kitchenOnly(kitchenHandler.SubmitReport))
	http.Handle("/api/v1/kitchen/dockets", kitchenOnly(kitchenHandler.GetDockets))
	http.Handle("/api/v1/kitchen/dockets/status", kitchenOnly(kitchenHandler.UpdateDocketStatus))

	// Bar and stock ledger
	http.Handle("/api/v1/bar/items", barOnly(barHandler.GetItems))
	http.Handle("/api/v1/bar/sales", barOnly(barHandler.RecordSale))
	http.Handle("/api/v1/bar/ledger", barOnly(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			barHandler.SaveLedger(w, r)
		} else {
			barHandler.GetLedger(w, r)
		}
	}))
	http.Handle("/api/v1/bar/reports", barOnly(barHandler.Reports))

	// Management dashboard and analytics
	http.Handle("/api/v1/manager/summary", managerOnly(managerHandler.Summary))
	http.Handle("/api/v1/manager/audits", managerOnly(managerHandler.GetAudits))
	http.Handle("/api/v1/manager/audits/action", managerOnly(managerHandler.AuditAction))
	http.Handle("/api/v1/manager/consolidated", managerOnly(managerHandler.ConsolidatedReport))

	// Guest booking portal
	http.HandleFunc("/api/v1/guest/rooms", guestHandler.GetAvailableRooms)
	http.Handle("/api/v1/guest/profile", guestOnly(guestHandler.GetProfile))
	http.Handle("/api/v1/guest/book", guestOnly(guestHandler.CreateBooking))
	http.Handle("/api/v1/guest/bookings", guestOnly(guestHandler.GetMyBookings))
	http.Handle("/api/v1/guest/folio", guestOnly(guestHandler.GetMyFolio))
	http.Handle("/api/v1/guest/checkin", guestOnly(guestHandler.CheckIn))
	http.HandleFunc("/api/v1/guest/resident-login", guestHandler.ResidentLogin)

	fmt.Printf("✓ La Luna Hotel Management System Unified API active on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handlers.CORSMiddleware(http.DefaultServeMux)))
}
