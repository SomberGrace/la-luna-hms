// backend/cmd/frontdesk/main.go
package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"backend/internal/config"
	"backend/internal/handlers"

	_ "github.com/lib/pq"
)

// Data Models (Matching Supabase Schema & Mobile Front Desk UX)
type Room struct {
	ID            string  `json:"id"`
	RoomNumber    string  `json:"room_number"`
	Type          string  `json:"type"`            // Standard, Executive, Apartment, Sector Suite
	BedType       string  `json:"bed_type"`        // Single, Double
	PricePerNight float64 `json:"price_per_night"` // In MWK
	Status        string  `json:"status"`          // Available, Occupied, Cleaning, OutOfService, Maintenance
	Floor         int     `json:"floor"`
	Housekeeper   string  `json:"housekeeper,omitempty"`
	Priority      string  `json:"priority,omitempty"`
	GuestName     string  `json:"guest_name,omitempty"`
}

type FolioItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"` // Accommodation, Restaurant, Bar, Incidental, Payment
	Timestamp   string  `json:"timestamp"`
}

type Folio struct {
	ID           string      `json:"id"`
	CardNo       string      `json:"card_no"`
	GuestID      string      `json:"guest_id"`
	GuestName    string      `json:"guest_name"`
	GuestPhone   string      `json:"guest_phone,omitempty"`
	GuestIDDoc   string      `json:"guest_id_doc,omitempty"`
	RoomNumber   string      `json:"room_number"`
	CheckInDate  string      `json:"check_in_date"`
	CheckOutDate string      `json:"check_out_date"`
	TotalCharges float64     `json:"total_charges"`
	TotalPaid    float64     `json:"total_paid"`
	Balance      float64     `json:"balance"`
	Status       string      `json:"status"` // Open, Settled, Closed
	Items        []FolioItem `json:"items"`
}

type Reservation struct {
	ID           string  `json:"id"`
	ResCode      string  `json:"res_code"`
	GuestName    string  `json:"guest_name"`
	GuestPhone   string  `json:"guest_phone"`
	GuestIDDoc   string  `json:"guest_id_doc"`
	RoomType     string  `json:"room_type"`
	RoomNumber   string  `json:"room_number"`
	CheckInDate  string  `json:"check_in_date"`
	CheckOutDate string  `json:"check_out_date"`
	Nights       int     `json:"nights"`
	DepositPaid  float64 `json:"deposit_paid"`
	Status       string  `json:"status"`
}

// In-Memory Fallback Store (for offline resiliency when grid/internet is down)
type Store struct {
	sync.RWMutex
	Rooms        map[string]*Room
	Folios       map[string]*Folio
	Reservations map[string]*Reservation
}

var (
	sqlDB *sql.DB
	memDB = &Store{
		Rooms:        make(map[string]*Room),
		Folios:       make(map[string]*Folio),
		Reservations: make(map[string]*Reservation),
	}
	isPgConnected = false
)

// Load .env file if present
func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

// Init Database: Try Supabase/PostgreSQL, fallback to resilient memory
func initDatabase() {
	loadEnvFile(".env")
	loadEnvFile("../../.env")
	loadEnvFile("../../../.env")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("SUPABASE_DB_URL")
	}

	if dbURL != "" {
		log.Printf("Connecting to PostgreSQL / Supabase...")
		var err error
		sqlDB, err = sql.Open("postgres", dbURL)
		if err == nil {
			sqlDB.SetMaxOpenConns(20)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(15 * time.Minute)

			pingErr := sqlDB.Ping()
			if pingErr == nil {
				isPgConnected = true
				log.Println("✓ Successfully connected to PostgreSQL / Supabase Database!")
				seedPostgreSQLRooms()
				return
			} else {
				log.Printf("PostgreSQL ping failed (%v). Falling back to resilient in-memory mode.", pingErr)
			}
		} else {
			log.Printf("PostgreSQL open failed (%v). Falling back to resilient in-memory mode.", err)
		}
	} else {
		log.Println("No DATABASE_URL found in environment. Running with in-memory resilient storage.")
	}

	initMemorySeed()
}

// Ensure initial rooms exist in Supabase/PostgreSQL
func seedPostgreSQLRooms() {
	if !isPgConnected || sqlDB == nil {
		return
	}

	var count int
	err := sqlDB.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&count)
	if err != nil {
		log.Printf("Warning checking rooms count in PostgreSQL: %v", err)
		return
	}

	if count == 0 {
		log.Println("Seeding initial room inventory into Supabase...")
		initialRooms := []struct {
			number string
			rType  string
			bed    string
			rate   float64
			status string
		}{
			{"101", "Standard", "Single", 45000, "Available"},
			{"102", "Standard", "Double", 45000, "Occupied"},
			{"103", "Standard", "Single", 45000, "Available"},
			{"201", "Executive", "Double", 75000, "Available"},
			{"202", "Executive", "Double", 75000, "Available"},
			{"203", "Executive", "Double", 75000, "Maintenance"},
			{"301", "Apartment", "Double", 120000, "Available"},
			{"302", "Apartment", "Double", 120000, "Occupied"},
		}

		for _, r := range initialRooms {
			_, err := sqlDB.Exec(
				"INSERT INTO rooms (room_number, room_type, bed_type, room_rate, status) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (room_number) DO NOTHING",
				r.number, r.rType, r.bed, r.rate, r.status,
			)
			if err != nil {
				log.Printf("Error inserting room %s: %v", r.number, err)
			}
		}
		log.Println("Initial room inventory seeded into Supabase.")
	}
}

// Seed In-Memory Store for fallback / offline simulation
func initMemorySeed() {
	rooms := []*Room{
		{ID: "r101", RoomNumber: "101", Type: "Standard", BedType: "Single", PricePerNight: 45000, Status: "Available", Floor: 1},
		{ID: "r102", RoomNumber: "102", Type: "Standard", BedType: "Double", PricePerNight: 45000, Status: "Occupied", Floor: 1, GuestName: "Alice Banda"},
		{ID: "r103", RoomNumber: "103", Type: "Standard", BedType: "Single", PricePerNight: 45000, Status: "Available", Floor: 1},
		{ID: "r201", RoomNumber: "201", Type: "Executive", BedType: "Double", PricePerNight: 75000, Status: "Available", Floor: 2},
		{ID: "r202", RoomNumber: "202", Type: "Executive", BedType: "Double", PricePerNight: 75000, Status: "Cleaning", Floor: 2, Housekeeper: "Grace Phiri", Priority: "High"},
		{ID: "r203", RoomNumber: "203", Type: "Executive", BedType: "Double", PricePerNight: 75000, Status: "OutOfService", Floor: 2, Housekeeper: "Maintenance", Priority: "Normal"},
		{ID: "r301", RoomNumber: "301", Type: "Apartment", BedType: "Double", PricePerNight: 120000, Status: "Available", Floor: 3},
		{ID: "r302", RoomNumber: "302", Type: "Apartment", BedType: "Double", PricePerNight: 120000, Status: "Occupied", Floor: 3, GuestName: "Chikondi Tembo"},
	}
	for _, r := range rooms {
		memDB.Rooms[r.RoomNumber] = r
	}

	memDB.Folios["102"] = &Folio{
		ID:           "fol-102",
		CardNo:       "CARD-102-001",
		GuestID:      "g-9901",
		GuestName:    "Alice Banda",
		GuestPhone:   "+265 999 123 456",
		GuestIDDoc:   "MW-ID-449102",
		RoomNumber:   "102",
		CheckInDate:  time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
		CheckOutDate: time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		TotalCharges: 90000,
		TotalPaid:    45000,
		Balance:      45000,
		Status:       "Open",
		Items: []FolioItem{
			{ID: "item-1", Description: "2 Nights Accommodation (Standard)", Amount: 90000, Category: "Accommodation", Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			{ID: "item-2", Description: "Initial Cash Deposit", Amount: -45000, Category: "Payment", Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
		},
	}

	memDB.Folios["302"] = &Folio{
		ID:           "fol-302",
		CardNo:       "CARD-302-002",
		GuestID:      "g-9902",
		GuestName:    "Chikondi Tembo",
		GuestPhone:   "+265 888 765 432",
		GuestIDDoc:   "MW-PP-882103",
		RoomNumber:   "302",
		CheckInDate:  time.Now().Add(-48 * time.Hour).Format("2006-01-02"),
		CheckOutDate: time.Now().Format("2006-01-02"),
		TotalCharges: 240000,
		TotalPaid:    240000,
		Balance:      0,
		Status:       "Settled",
		Items: []FolioItem{
			{ID: "item-302-1", Description: "2 Nights Accommodation (Apartment)", Amount: 240000, Category: "Accommodation", Timestamp: time.Now().Add(-48 * time.Hour).Format(time.RFC3339)},
			{ID: "item-302-2", Description: "Visa Payment Settlement", Amount: -240000, Category: "Payment", Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
		},
	}

	memDB.Reservations["RES-2026-001"] = &Reservation{
		ID:           "res-1",
		ResCode:      "RES-2026-001",
		GuestName:    "David Mwale",
		GuestPhone:   "+265 991 555 777",
		GuestIDDoc:   "MW-ID-881230",
		RoomType:     "Standard",
		RoomNumber:   "101",
		CheckInDate:  time.Now().Format("2006-01-02"),
		CheckOutDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		Nights:       2,
		DepositPaid:  45000,
		Status:       "Confirmed",
	}

	memDB.Reservations["RES-2026-002"] = &Reservation{
		ID:           "res-2",
		ResCode:      "RES-2026-002",
		GuestName:    "Fatima Phiri",
		GuestPhone:   "+265 882 111 222",
		GuestIDDoc:   "MW-PP-990112",
		RoomType:     "Executive",
		RoomNumber:   "201",
		CheckInDate:  time.Now().Format("2006-01-02"),
		CheckOutDate: time.Now().AddDate(0, 0, 3).Format("2006-01-02"),
		Nights:       3,
		DepositPaid:  75000,
		Status:       "Confirmed",
	}
}

// CORS Middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// Helper: Map frontend payment strings to Supabase payment_mode_enum
func toPaymentModeEnum(m string) string {
	lower := strings.ToLower(m)
	if strings.Contains(lower, "card") || strings.Contains(lower, "visa") || strings.Contains(lower, "pos") {
		return "Credit Card"
	}
	if strings.Contains(lower, "cheque") {
		return "Cheque"
	}
	if strings.Contains(lower, "room") {
		return "Room Charge"
	}
	return "Cash" // Default cash (or mobile money cash equivalent)
}

// Helper: Map frontend room status to Supabase room_status_enum
func toRoomStatusEnum(s string) string {
	switch s {
	case "Occupied":
		return "Occupied"
	case "Maintenance", "OutOfService":
		return "Maintenance"
	default:
		return "Available"
	}
}

// Handlers

func getRooms(w http.ResponseWriter, r *http.Request) {
	if isPgConnected && sqlDB != nil {
		rows, err := sqlDB.Query(`
			SELECT 
				r.room_number, r.room_type, r.bed_type, r.room_rate, r.status,
				COALESCE(g.first_name || ' ' || g.last_name, '') as guest_name
			FROM rooms r
			LEFT JOIN bookings b ON r.room_number = b.room_number AND b.departure_date >= CURRENT_DATE
			LEFT JOIN guests g ON b.guest_id = g.guest_id
			ORDER BY r.room_number
		`)
		if err == nil {
			defer rows.Close()
			var result []*Room
			for rows.Next() {
				var rNo, rType, bed, status, guestName string
				var rate float64
				if err := rows.Scan(&rNo, &rType, &bed, &rate, &status, &guestName); err == nil {
					floor := 1
					if len(rNo) > 0 {
						if rNo[0] == '2' {
							floor = 2
						} else if rNo[0] == '3' {
							floor = 3
						}
					}
					result = append(result, &Room{
						ID:            "r" + rNo,
						RoomNumber:    rNo,
						Type:          rType,
						BedType:       bed,
						PricePerNight: rate,
						Status:        status,
						Floor:         floor,
						GuestName:     guestName,
					})
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			return
		}
		log.Printf("PostgreSQL getRooms error, using memory fallback: %v", err)
	}

	// Memory fallback
	memDB.RLock()
	defer memDB.RUnlock()

	var result []*Room
	for _, room := range memDB.Rooms {
		result = append(result, room)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getFolio(w http.ResponseWriter, r *http.Request) {
	roomNo := r.URL.Query().Get("room_number")
	if roomNo == "" {
		http.Error(w, "Missing room_number parameter", http.StatusBadRequest)
		return
	}

	if isPgConnected && sqlDB != nil {
		var f Folio
		var arrival, departure time.Time
		var roomCharges, incidentalCharges, totalBalance, advancePaid float64

		err := sqlDB.QueryRow(`
			SELECT 
				f.folio_id, f.card_no, f.guest_id,
				COALESCE(g.first_name || ' ' || g.last_name, 'Guest') as guest_name,
				COALESCE(g.contact_no, ''), COALESCE(g.passport_no, ''),
				b.room_number, b.arrival_date, b.departure_date,
				f.room_charges, f.incidental_charges, f.total_balance, b.advance_paid, f.payment_status
			FROM folios f
			JOIN bookings b ON f.card_no = b.card_no
			JOIN guests g ON f.guest_id = g.guest_id
			WHERE b.room_number = $1 AND f.payment_status != 'Closed'
			ORDER BY f.created_at DESC LIMIT 1
		`, roomNo).Scan(
			&f.ID, &f.CardNo, &f.GuestID, &f.GuestName,
			&f.GuestPhone, &f.GuestIDDoc,
			&f.RoomNumber, &arrival, &departure,
			&roomCharges, &incidentalCharges, &totalBalance, &advancePaid, &f.Status,
		)

		if err == nil {
			f.CheckInDate = arrival.Format("2006-01-02")
			f.CheckOutDate = departure.Format("2006-01-02")
			f.TotalCharges = totalBalance
			f.TotalPaid = advancePaid
			f.Balance = totalBalance - advancePaid

			f.Items = []FolioItem{
				{
					ID:          "item-room",
					Description: fmt.Sprintf("Accommodation Room %s", roomNo),
					Amount:      roomCharges,
					Category:    "Accommodation",
					Timestamp:   arrival.Format(time.RFC3339),
				},
			}

			if advancePaid > 0 {
				f.Items = append(f.Items, FolioItem{
					ID:          "item-advance",
					Description: "Advance Deposit Paid",
					Amount:      -advancePaid,
					Category:    "Payment",
					Timestamp:   arrival.Format(time.RFC3339),
				})
			}

			// Add restaurant incidental dockets if any
			docketRows, dErr := sqlDB.Query(`
				SELECT docket_id, total_amount, payment_method, sale_date
				FROM restaurant_dockets
				WHERE folio_id = $1
			`, f.ID)
			if dErr == nil {
				defer docketRows.Close()
				for docketRows.Next() {
					var docID, method string
					var amount float64
					var saleDate time.Time
					if err := docketRows.Scan(&docID, &amount, &method, &saleDate); err == nil {
						f.Items = append(f.Items, FolioItem{
							ID:          docID,
							Description: fmt.Sprintf("Restaurant Docket (%s)", method),
							Amount:      amount,
							Category:    "Restaurant",
							Timestamp:   saleDate.Format(time.RFC3339),
						})
					}
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(f)
			return
		}
	}

	// Memory Fallback
	memDB.RLock()
	defer memDB.RUnlock()

	folio, exists := memDB.Folios[roomNo]
	if !exists {
		http.Error(w, "Active folio not found for this room", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folio)
}

func checkIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GuestName       string  `json:"guest_name"`
		GuestPhone      string  `json:"guest_phone"`
		GuestIDDoc      string  `json:"guest_id_doc"`
		RoomNumber      string  `json:"room_number"`
		Nights          int     `json:"nights"`
		DepositPaid     float64 `json:"deposit_paid"`
		PaymentMethod   string  `json:"payment_method,omitempty"`
		ReservationCode string  `json:"reservation_code,omitempty"`
		ClerkName       string  `json:"clerk_name,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Nights <= 0 {
		req.Nights = 1
	}
	if req.ClerkName == "" {
		req.ClerkName = "Faith Mwanza"
	}

	// Split name
	nameParts := strings.Fields(req.GuestName)
	firstName := req.GuestName
	lastName := "Guest"
	if len(nameParts) > 1 {
		firstName = nameParts[0]
		lastName = strings.Join(nameParts[1:], " ")
	}

	cardNo := fmt.Sprintf("CARD-%s-%d", req.RoomNumber, time.Now().Unix()%10000)

	if isPgConnected && sqlDB != nil {
		tx, err := sqlDB.Begin()
		if err == nil {
			var guestID string
			err = tx.QueryRow(`
				INSERT INTO guests (first_name, last_name, contact_no, passport_no)
				VALUES ($1, $2, $3, $4)
				RETURNING guest_id
			`, firstName, lastName, req.GuestPhone, req.GuestIDDoc).Scan(&guestID)

			if err == nil {
				var roomRate float64
				err = tx.QueryRow("SELECT room_rate FROM rooms WHERE room_number = $1", req.RoomNumber).Scan(&roomRate)
				if err == nil {
					totalRoomCharges := roomRate * float64(req.Nights)
					arrival := time.Now()
					departure := arrival.AddDate(0, 0, req.Nights)

					paymentMode := toPaymentModeEnum(req.PaymentMethod)

					_, err = tx.Exec(`
						INSERT INTO bookings (
							card_no, guest_id, room_number, arrival_date, departure_date,
							mode_of_payment, advance_paid, front_officer_clerk
						) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
					`, cardNo, guestID, req.RoomNumber, arrival, departure, paymentMode, req.DepositPaid, req.ClerkName)

					if err == nil {
						folioStatus := "Open"
						if req.DepositPaid >= totalRoomCharges {
							folioStatus = "Settled"
						}

						var folioID string
						err = tx.QueryRow(`
							INSERT INTO folios (card_no, guest_id, room_charges, incidental_charges, payment_status)
							VALUES ($1, $2, $3, 0.00, $4)
							RETURNING folio_id
						`, cardNo, guestID, totalRoomCharges, folioStatus).Scan(&folioID)

						if err == nil {
							tx.Commit()

							newFolio := &Folio{
								ID:           folioID,
								CardNo:       cardNo,
								GuestID:      guestID,
								GuestName:    req.GuestName,
								GuestPhone:   req.GuestPhone,
								GuestIDDoc:   req.GuestIDDoc,
								RoomNumber:   req.RoomNumber,
								CheckInDate:  arrival.Format("2006-01-02"),
								CheckOutDate: departure.Format("2006-01-02"),
								TotalCharges: totalRoomCharges,
								TotalPaid:    req.DepositPaid,
								Balance:      totalRoomCharges - req.DepositPaid,
								Status:       folioStatus,
								Items: []FolioItem{
									{
										ID:          "item-accom",
										Description: fmt.Sprintf("%d Night(s) Accommodation - Room %s", req.Nights, req.RoomNumber),
										Amount:      totalRoomCharges,
										Category:    "Accommodation",
										Timestamp:   arrival.Format(time.RFC3339),
									},
								},
							}

							if req.DepositPaid > 0 {
								newFolio.Items = append(newFolio.Items, FolioItem{
									ID:          "item-deposit",
									Description: fmt.Sprintf("Initial Deposit (%s)", paymentMode),
									Amount:      -req.DepositPaid,
									Category:    "Payment",
									Timestamp:   arrival.Format(time.RFC3339),
								})
							}

							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(map[string]interface{}{
								"message": "Check-in committed to Supabase PostgreSQL database successfully",
								"folio":   newFolio,
							})
							return
						}
					}
				}
			}
			tx.Rollback()
			log.Printf("PostgreSQL checkIn failed, falling back to memory: %v", err)
		}
	}

	// Memory Fallback
	memDB.Lock()
	defer memDB.Unlock()

	room, exists := memDB.Rooms[req.RoomNumber]
	if !exists || room.Status != "Available" {
		http.Error(w, "Room is unavailable or does not exist", http.StatusBadRequest)
		return
	}

	room.Status = "Occupied"
	room.GuestName = req.GuestName
	totalCharge := room.PricePerNight * float64(req.Nights)
	balance := totalCharge - req.DepositPaid

	newFolio := &Folio{
		ID:           fmt.Sprintf("fol-%s", req.RoomNumber),
		CardNo:       cardNo,
		GuestID:      fmt.Sprintf("g-%d", time.Now().Unix()),
		GuestName:    req.GuestName,
		GuestPhone:   req.GuestPhone,
		GuestIDDoc:   req.GuestIDDoc,
		RoomNumber:   req.RoomNumber,
		CheckInDate:  time.Now().Format("2006-01-02"),
		CheckOutDate: time.Now().AddDate(0, 0, req.Nights).Format("2006-01-02"),
		TotalCharges: totalCharge,
		TotalPaid:    req.DepositPaid,
		Balance:      balance,
		Status:       "Open",
		Items: []FolioItem{
			{
				ID:          fmt.Sprintf("item-%d", time.Now().Unix()),
				Description: fmt.Sprintf("%d Night(s) Accommodation - %s", req.Nights, room.Type),
				Amount:      totalCharge,
				Category:    "Accommodation",
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	}

	if req.DepositPaid > 0 {
		newFolio.Items = append(newFolio.Items, FolioItem{
			ID:          fmt.Sprintf("pay-%d", time.Now().Unix()),
			Description: "Initial Deposit Received",
			Amount:      -req.DepositPaid,
			Category:    "Payment",
			Timestamp:   time.Now().Format(time.RFC3339),
		})
	}

	memDB.Folios[req.RoomNumber] = newFolio
	if req.ReservationCode != "" {
		if res, ok := memDB.Reservations[req.ReservationCode]; ok {
			res.Status = "CheckedIn"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Check-in successful (cached locally)",
		"folio":   newFolio,
	})
}

func recordPayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomNumber string  `json:"room_number"`
		Amount     float64 `json:"amount"`
		Method     string  `json:"method"` // Cash, Visa, Airtel Money, TNM Mpamba, Cheque
		Reference  string  `json:"reference,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Payment amount must be greater than zero", http.StatusBadRequest)
		return
	}

	if isPgConnected && sqlDB != nil {
		// Update advance_paid or folio status in Postgres
		var cardNo string
		var advancePaid, totalBalance float64
		err := sqlDB.QueryRow(`
			SELECT b.card_no, b.advance_paid, f.total_balance
			FROM folios f
			JOIN bookings b ON f.card_no = b.card_no
			WHERE b.room_number = $1 AND f.payment_status != 'Closed'
			LIMIT 1
		`, req.RoomNumber).Scan(&cardNo, &advancePaid, &totalBalance)

		if err == nil {
			newAdvance := advancePaid + req.Amount
			newStatus := "Open"
			if newAdvance >= totalBalance {
				newStatus = "Settled"
			}

			_, _ = sqlDB.Exec("UPDATE bookings SET advance_paid = $1 WHERE card_no = $2", newAdvance, cardNo)
			_, _ = sqlDB.Exec("UPDATE folios SET payment_status = $1, updated_at = CURRENT_TIMESTAMP WHERE card_no = $2", newStatus, cardNo)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Payment recorded in Supabase PostgreSQL successfully",
				"folio": map[string]interface{}{
					"room_number":   req.RoomNumber,
					"total_paid":    newAdvance,
					"balance":       totalBalance - newAdvance,
					"total_charges": totalBalance,
				},
			})
			return
		}
	}

	// Memory Fallback
	memDB.Lock()
	defer memDB.Unlock()

	folio, exists := memDB.Folios[req.RoomNumber]
	if !exists {
		http.Error(w, "Folio not found for this room", http.StatusNotFound)
		return
	}

	folio.TotalPaid += req.Amount
	folio.Balance -= req.Amount
	if folio.Balance <= 0 {
		folio.Status = "Settled"
	}

	desc := fmt.Sprintf("Payment Received (%s)", req.Method)
	if req.Reference != "" {
		desc += fmt.Sprintf(" - Ref: %s", req.Reference)
	}

	paymentItem := FolioItem{
		ID:          fmt.Sprintf("pay-%d", time.Now().Unix()),
		Description: desc,
		Amount:      -req.Amount,
		Category:    "Payment",
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	folio.Items = append(folio.Items, paymentItem)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment recorded successfully",
		"folio":   folio,
	})
}

func checkOut(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomNumber string `json:"room_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isPgConnected && sqlDB != nil {
		var cardNo string
		var totalBalance, advancePaid float64
		err := sqlDB.QueryRow(`
			SELECT b.card_no, f.total_balance, b.advance_paid
			FROM folios f
			JOIN bookings b ON f.card_no = b.card_no
			WHERE b.room_number = $1 AND f.payment_status != 'Closed'
			LIMIT 1
		`, req.RoomNumber).Scan(&cardNo, &totalBalance, &advancePaid)

		if err == nil {
			balance := totalBalance - advancePaid
			if balance > 0 {
				http.Error(w, fmt.Sprintf("Zero-Balance Lock Active: Outstanding balance of MWK %.2f remaining", balance), http.StatusForbidden)
				return
			}

			// Mark folio Closed and room Available / Cleaning
			_, _ = sqlDB.Exec("UPDATE folios SET payment_status = 'Closed', updated_at = CURRENT_TIMESTAMP WHERE card_no = $1", cardNo)
			_, _ = sqlDB.Exec("UPDATE rooms SET status = 'Available', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1", req.RoomNumber)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Check-out completed in Supabase. Room status released and handed over to Housekeeping.",
			})
			return
		}
	}

	// Memory Fallback
	memDB.Lock()
	defer memDB.Unlock()

	folio, exists := memDB.Folios[req.RoomNumber]
	if !exists {
		http.Error(w, "No active stay found for this room", http.StatusNotFound)
		return
	}

	if folio.Balance > 0 {
		http.Error(w, fmt.Sprintf("Zero-Balance Lock Active: Outstanding balance of MWK %.2f remaining", folio.Balance), http.StatusForbidden)
		return
	}

	if room, ok := memDB.Rooms[req.RoomNumber]; ok {
		room.Status = "Cleaning"
		room.Priority = "High"
		room.GuestName = ""
	}
	delete(memDB.Folios, req.RoomNumber)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Check-out completed successfully. Room handed over to Housekeeping.",
	})
}

func extendStay(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomNumber       string `json:"room_number"`
		AdditionalNights int    `json:"additional_nights"`
		PayNow           bool   `json:"pay_now"`
		PaymentMethod    string `json:"payment_method,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.AdditionalNights <= 0 {
		http.Error(w, "Additional nights must be at least 1", http.StatusBadRequest)
		return
	}

	if isPgConnected && sqlDB != nil {
		var cardNo string
		var departure time.Time
		var roomRate, roomCharges float64
		err := sqlDB.QueryRow(`
			SELECT b.card_no, b.departure_date, r.room_rate, f.room_charges
			FROM folios f
			JOIN bookings b ON f.card_no = b.card_no
			JOIN rooms r ON b.room_number = r.room_number
			WHERE b.room_number = $1 AND f.payment_status != 'Closed'
			LIMIT 1
		`, req.RoomNumber).Scan(&cardNo, &departure, &roomRate, &roomCharges)

		if err == nil {
			newDeparture := departure.AddDate(0, 0, req.AdditionalNights)
			additionalCharge := roomRate * float64(req.AdditionalNights)
			newRoomCharges := roomCharges + additionalCharge

			_, _ = sqlDB.Exec("UPDATE bookings SET departure_date = $1 WHERE card_no = $2", newDeparture, cardNo)
			_, _ = sqlDB.Exec("UPDATE folios SET room_charges = $1, updated_at = CURRENT_TIMESTAMP WHERE card_no = $2", newRoomCharges, cardNo)

			if req.PayNow {
				_, _ = sqlDB.Exec("UPDATE bookings SET advance_paid = advance_paid + $1 WHERE card_no = $2", additionalCharge, cardNo)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":           "Stay extended in Supabase PostgreSQL",
				"new_checkout_date": newDeparture.Format("2006-01-02"),
			})
			return
		}
	}

	// Memory Fallback
	memDB.Lock()
	defer memDB.Unlock()

	room, roomExists := memDB.Rooms[req.RoomNumber]
	folio, folioExists := memDB.Folios[req.RoomNumber]
	if !roomExists || !folioExists {
		http.Error(w, "Active stay not found for room", http.StatusNotFound)
		return
	}

	additionalCharge := room.PricePerNight * float64(req.AdditionalNights)
	currentCO, err := time.Parse("2006-01-02", folio.CheckOutDate)
	if err != nil {
		currentCO = time.Now()
	}
	newCO := currentCO.AddDate(0, 0, req.AdditionalNights)
	folio.CheckOutDate = newCO.Format("2006-01-02")
	folio.TotalCharges += additionalCharge
	folio.Balance += additionalCharge

	folio.Items = append(folio.Items, FolioItem{
		ID:          fmt.Sprintf("ext-%d", time.Now().Unix()),
		Description: fmt.Sprintf("Stay Extension: +%d Night(s) (%s)", req.AdditionalNights, room.Type),
		Amount:      additionalCharge,
		Category:    "Accommodation",
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	if req.PayNow {
		folio.TotalPaid += additionalCharge
		folio.Balance -= additionalCharge
		method := req.PaymentMethod
		if method == "" {
			method = "Cash"
		}
		folio.Items = append(folio.Items, FolioItem{
			ID:          fmt.Sprintf("pay-ext-%d", time.Now().Unix()),
			Description: fmt.Sprintf("Extension Payment Received (%s)", method),
			Amount:      -additionalCharge,
			Category:    "Payment",
			Timestamp:   time.Now().Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Stay extended successfully",
		"folio":   folio,
	})
}

func updateRoomStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomNumber  string `json:"room_number"`
		Status      string `json:"status"` // Available, Occupied, Cleaning, OutOfService, Maintenance
		Housekeeper string `json:"housekeeper,omitempty"`
		Priority    string `json:"priority,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgStatus := toRoomStatusEnum(req.Status)
	if isPgConnected && sqlDB != nil {
		_, err := sqlDB.Exec("UPDATE rooms SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE room_number = $2", pgStatus, req.RoomNumber)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": fmt.Sprintf("Room %s status updated in Supabase to %s", req.RoomNumber, pgStatus),
				"room": map[string]string{
					"room_number": req.RoomNumber,
					"status":      req.Status,
				},
			})
			return
		}
	}

	// Memory Fallback
	memDB.Lock()
	defer memDB.Unlock()

	room, exists := memDB.Rooms[req.RoomNumber]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	room.Status = req.Status
	if req.Housekeeper != "" {
		room.Housekeeper = req.Housekeeper
	}
	if req.Priority != "" {
		room.Priority = req.Priority
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Room status updated",
		"room":    room,
	})
}

func getReservations(w http.ResponseWriter, r *http.Request) {
	memDB.RLock()
	defer memDB.RUnlock()

	var result []*Reservation
	for _, res := range memDB.Reservations {
		result = append(result, res)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getReports(w http.ResponseWriter, r *http.Request) {
	if isPgConnected && sqlDB != nil {
		var totalRooms, occupied, available, maintenance int
		_ = sqlDB.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&totalRooms)
		_ = sqlDB.QueryRow("SELECT COUNT(*) FROM rooms WHERE status = 'Occupied'").Scan(&occupied)
		_ = sqlDB.QueryRow("SELECT COUNT(*) FROM rooms WHERE status = 'Available'").Scan(&available)
		_ = sqlDB.QueryRow("SELECT COUNT(*) FROM rooms WHERE status = 'Maintenance'").Scan(&maintenance)

		occupancyRate := 0.0
		if totalRooms > 0 {
			occupancyRate = (float64(occupied) / float64(totalRooms)) * 100
		}

		var totalRev, totalBal float64
		_ = sqlDB.QueryRow("SELECT COALESCE(SUM(advance_paid), 0) FROM bookings").Scan(&totalRev)
		_ = sqlDB.QueryRow("SELECT COALESCE(SUM(total_balance - advance_paid), 0) FROM folios f JOIN bookings b ON f.card_no = b.card_no WHERE f.payment_status != 'Closed'").Scan(&totalBal)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_rooms":        totalRooms,
			"occupied_rooms":     occupied,
			"available_rooms":    available,
			"cleaning_rooms":     maintenance,
			"out_of_service":     maintenance,
			"occupancy_rate":     occupancyRate,
			"total_revenue_mwk":  totalRev,
			"total_balance_mwk":  totalBal,
			"active_stays_count": occupied,
			"database_engine":    "PostgreSQL / Supabase",
			"timestamp":          time.Now().Format(time.RFC3339),
		})
		return
	}

	// Memory fallback
	memDB.RLock()
	defer memDB.RUnlock()

	totalRooms := len(memDB.Rooms)
	occupied := 0
	available := 0
	cleaning := 0
	outOfService := 0

	for _, room := range memDB.Rooms {
		switch room.Status {
		case "Occupied":
			occupied++
		case "Available":
			available++
		case "Cleaning":
			cleaning++
		case "OutOfService":
			outOfService++
		}
	}

	occupancyRate := 0.0
	if totalRooms > 0 {
		occupancyRate = (float64(occupied) / float64(totalRooms)) * 100
	}

	totalRevenue := 0.0
	totalOutstanding := 0.0
	for _, folio := range memDB.Folios {
		totalRevenue += folio.TotalPaid
		totalOutstanding += folio.Balance
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_rooms":        totalRooms,
		"occupied_rooms":     occupied,
		"available_rooms":    available,
		"cleaning_rooms":     cleaning,
		"out_of_service":     outOfService,
		"occupancy_rate":     occupancyRate,
		"total_revenue_mwk":  totalRevenue,
		"total_balance_mwk":  totalOutstanding,
		"active_stays_count": len(memDB.Folios),
		"database_engine":    "In-Memory Cache (Offline Resilient)",
		"timestamp":          time.Now().Format(time.RFC3339),
	})
}

// Health check endpoint to inspect database connection status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	status := "In-Memory Cache"
	if isPgConnected {
		status = "Connected to Supabase PostgreSQL"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "ok",
		"database_status": status,
		"is_postgresql":   isPgConnected,
		"timestamp":       time.Now().Format(time.RFC3339),
	})
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be set to a secret of at least 32 bytes")
	}
	initDatabase()
	jwtSecret := cfg.JWTSecret

	frontDeskAuth := func(handler http.HandlerFunc) http.HandlerFunc {
		return enableCORS(handlers.RequireRole(jwtSecret, "FRONT_DESK")(handler))
	}
	http.HandleFunc("/api/v1/frontdesk/health", frontDeskAuth(healthCheck))
	http.HandleFunc("/api/v1/frontdesk/rooms", frontDeskAuth(getRooms))
	http.HandleFunc("/api/v1/frontdesk/folio", frontDeskAuth(getFolio))
	http.HandleFunc("/api/v1/frontdesk/checkin", frontDeskAuth(checkIn))
	http.HandleFunc("/api/v1/frontdesk/payment", frontDeskAuth(recordPayment))
	http.HandleFunc("/api/v1/frontdesk/checkout", frontDeskAuth(checkOut))
	http.HandleFunc("/api/v1/frontdesk/extension", frontDeskAuth(extendStay))
	http.HandleFunc("/api/v1/frontdesk/rooms/status", frontDeskAuth(updateRoomStatus))
	http.HandleFunc("/api/v1/frontdesk/reservations", frontDeskAuth(getReservations))
	http.HandleFunc("/api/v1/frontdesk/reports", frontDeskAuth(getReports))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("La Luna Front Desk Backend Server listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server launch failure: %v", err)
	}
}
