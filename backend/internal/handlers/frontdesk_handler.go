// backend/internal/handlers/frontdesk_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FrontDeskHandler struct {
	db *sql.DB
}

func NewFrontDeskHandler(db *sql.DB) *FrontDeskHandler {
	return &FrontDeskHandler{db: db}
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
	w.Header().Set("Content-Type", "application/json")
}

// HealthCheck reports connection state to Supabase.
func (h *FrontDeskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	dbStatus := "In-Memory Cache (Disconnected)"
	if h.db != nil {
		if err := h.db.Ping(); err == nil {
			dbStatus = "Connected to Supabase PostgreSQL"
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "healthy",
		"database_status": dbStatus,
		"timestamp":       time.Now().Format(time.RFC3339),
	})
}

// GetRooms returns all rooms from Supabase.
func (h *FrontDeskHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	rows, err := h.db.Query(
		`SELECT room_number, room_type::text, bed_type::text, room_rate, status::text
		   FROM rooms
		  ORDER BY room_number ASC`,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query rooms: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type RoomResp struct {
		ID            string  `json:"id"`
		RoomNumber    string  `json:"room_number"`
		Type          string  `json:"type"`
		BedType       string  `json:"bed_type"`
		PricePerNight float64 `json:"price_per_night"`
		Status        string  `json:"status"`
	}

	var results []RoomResp
	for rows.Next() {
		var rm RoomResp
		if err := rows.Scan(&rm.RoomNumber, &rm.Type, &rm.BedType, &rm.PricePerNight, &rm.Status); err != nil {
			continue
		}
		rm.ID = "r-" + rm.RoomNumber
		results = append(results, rm)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"row iteration failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

// UpdateRoomStatus updates the status of a room in Supabase.
func (h *FrontDeskHandler) UpdateRoomStatus(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		RoomNumber string `json:"room_number"`
		Status     string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Map frontend UI statuses to Postgres room_status_enum ('Available', 'Occupied', 'Maintenance')
	pgStatus := "Available"
	switch strings.ToLower(req.Status) {
	case "occupied":
		pgStatus = "Occupied"
	case "maintenance", "cleaning", "outofservice":
		pgStatus = "Maintenance"
	case "reserved":
		pgStatus = "Reserved"
	default:
		pgStatus = "Available"
	}

	_, err := h.db.Exec(
		`UPDATE rooms SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE room_number = $2`,
		pgStatus, req.RoomNumber,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to update room status: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"room_number": req.RoomNumber,
		"status":      req.Status,
	})
}

// CheckIn handles walk-ins and guest check-ins with atomic transactional safety.
func (h *FrontDeskHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		GuestName       string  `json:"guest_name"`
		GuestPhone      string  `json:"guest_phone"`
		GuestIDDoc      string  `json:"guest_id_doc"`
		RoomNumber      string  `json:"room_number"`
		Nights          int     `json:"nights"`
		DepositPaid     float64 `json:"deposit_paid"`
		PaymentMode     string  `json:"payment_mode"`
		ReservationCode string  `json:"reservation_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Nights <= 0 {
		req.Nights = 1
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"failed to begin transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Handle existing reservation / online booking check-in
	if req.ReservationCode != "" {
		if req.RoomNumber == "" {
			_ = tx.QueryRow(`SELECT room_number FROM bookings WHERE card_no = $1`, req.ReservationCode).Scan(&req.RoomNumber)
		}
		if req.GuestName == "" {
			_ = tx.QueryRow(`SELECT g.first_name || ' ' || g.last_name FROM guests g JOIN bookings b ON g.guest_id = b.guest_id WHERE b.card_no = $1`, req.ReservationCode).Scan(&req.GuestName)
		}

		// Guard: prevent checking in if reservation itself is already checked in
		var thisCheckedIn bool
		_ = tx.QueryRow(`SELECT checked_in FROM bookings WHERE card_no = $1`, req.ReservationCode).Scan(&thisCheckedIn)
		if thisCheckedIn {
			http.Error(w, `{"error":"This reservation has already been checked in."}`, http.StatusConflict)
			return
		}

		// Guard: prevent double check-in if guest already has another active checked-in booking
		var otherCheckedIn int
		_ = tx.QueryRow(`
			SELECT COUNT(*) FROM bookings b
			  LEFT JOIN folios f ON b.card_no = f.card_no
			 WHERE b.guest_id = (SELECT guest_id FROM bookings WHERE card_no = $1)
			   AND b.card_no != $1
			   AND b.checked_in = TRUE
			   AND (f.payment_status IS NULL OR f.payment_status != 'Closed')
		`, req.ReservationCode).Scan(&otherCheckedIn)
		if otherCheckedIn > 0 {
			http.Error(w, `{"error":"This guest already has an active checked-in booking in another room. They must check out first before checking in again."}`, http.StatusConflict)
			return
		}

		_, err = tx.Exec(`UPDATE rooms SET status = 'Occupied', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1`, req.RoomNumber)
		if err != nil {
			http.Error(w, `{"error":"failed to update room status"}`, http.StatusInternalServerError)
			return
		}

		// Mark booking as checked in
		_, _ = tx.Exec(`UPDATE bookings SET checked_in = TRUE WHERE card_no = $1`, req.ReservationCode)

		if req.DepositPaid > 0 {
			_, _ = tx.Exec(`UPDATE bookings SET advance_paid = advance_paid + $1 WHERE card_no = $2`, req.DepositPaid, req.ReservationCode)
		}

		var folioID string
		var roomCharges, incidental, balance float64
		var payStatus string
		err = tx.QueryRow(`
			SELECT folio_id::text, room_charges, incidental_charges, total_balance, payment_status::text
			  FROM folios
			 WHERE card_no = $1
		`, req.ReservationCode).Scan(&folioID, &roomCharges, &incidental, &balance, &payStatus)
		if err != nil {
			var roomRate float64
			_ = tx.QueryRow(`SELECT room_rate FROM rooms WHERE room_number = $1`, req.RoomNumber).Scan(&roomRate)
			totalCharge := roomRate * float64(req.Nights)
			_ = tx.QueryRow(
				`INSERT INTO folios (card_no, guest_id, room_charges, incidental_charges, payment_status)
				 SELECT $1, guest_id, $2, 0.00, 'Open' FROM bookings WHERE card_no = $1
				 RETURNING folio_id::text`,
				req.ReservationCode, totalCharge,
			).Scan(&folioID)
			roomCharges = totalCharge
			balance = totalCharge
			payStatus = "Open"
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Guest successfully checked in to Room %s!", req.RoomNumber),
			"folio": map[string]interface{}{
				"id":            folioID,
				"guest_name":    req.GuestName,
				"guest_phone":   req.GuestPhone,
				"guest_id_doc":  req.GuestIDDoc,
				"room_number":   req.RoomNumber,
				"total_charges": roomCharges,
				"total_paid":    req.DepositPaid,
				"balance":       balance - req.DepositPaid,
			},
		})
		return
	}

	// 1. Fetch Room rate & check availability for new walk-in
	var roomRate float64
	var roomStatus, roomType string
	err = tx.QueryRow(
		`SELECT room_rate, status::text, room_type::text FROM rooms WHERE room_number = $1 FOR UPDATE`,
		req.RoomNumber,
	).Scan(&roomRate, &roomStatus, &roomType)
	if err != nil {
		http.Error(w, `{"error":"room not found"}`, http.StatusNotFound)
		return
	}
	if roomStatus == "Occupied" {
		http.Error(w, `{"error":"room is currently occupied"}`, http.StatusConflict)
		return
	}

	// 2. Parse guest first and last name
	parts := strings.SplitN(strings.TrimSpace(req.GuestName), " ", 2)
	firstName := parts[0]
	lastName := "Guest"
	if len(parts) > 1 {
		lastName = parts[1]
	}

	// Guard: prevent check-in if guest already has an active stay in the hotel
	phoneClean := strings.TrimSpace(req.GuestPhone)
	docClean := strings.TrimSpace(req.GuestIDDoc)
	if phoneClean != "" || docClean != "" {
		var activeStays int
		_ = tx.QueryRow(`
			SELECT COUNT(*) FROM bookings b
			  JOIN guests g ON b.guest_id = g.guest_id
			  LEFT JOIN folios f ON b.card_no = f.card_no
			 WHERE b.checked_in = TRUE
			   AND (f.payment_status IS NULL OR f.payment_status != 'Closed')
			   AND (
			     (NULLIF($1, '') IS NOT NULL AND g.contact_no = $1)
			     OR (NULLIF($2, '') IS NOT NULL AND g.passport_no = $2)
			   )
		`, phoneClean, docClean).Scan(&activeStays)
		if activeStays > 0 {
			http.Error(w, `{"error":"This guest already has an active checked-in stay in the hotel. They must check out first before checking in again."}`, http.StatusConflict)
			return
		}
	}

	var guestID string
	err = tx.QueryRow(
		`INSERT INTO guests (first_name, last_name, contact_no, passport_no)
		 VALUES ($1, $2, $3, $4)
		 RETURNING guest_id::text`,
		firstName, lastName, req.GuestPhone, req.GuestIDDoc,
	).Scan(&guestID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to create guest record: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// 3. Create Booking
	cardNo := fmt.Sprintf("CARD-%s-%d", req.RoomNumber, time.Now().Unix()%10000)
	arrivalDate := time.Now().Format("2006-01-02")
	departureDate := time.Now().AddDate(0, 0, req.Nights).Format("2006-01-02")

	payMode := "Cash"
	switch strings.ToLower(req.PaymentMode) {
	case "credit card", "credit_card", "visa":
		payMode = "Credit Card"
	case "cheque":
		payMode = "Cheque"
	case "room charge", "room_charge":
		payMode = "Room Charge"
	}

	_, err = tx.Exec(
		`INSERT INTO bookings (card_no, guest_id, room_number, arrival_date, departure_date, mode_of_payment, advance_paid, front_officer_clerk, checked_in)
		 VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, TRUE)`,
		cardNo, guestID, req.RoomNumber, arrivalDate, departureDate, payMode, req.DepositPaid, "Front Desk Officer",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to insert booking: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// 4. Create Folio
	totalCharge := roomRate * float64(req.Nights)
	var folioID string
	err = tx.QueryRow(
		`INSERT INTO folios (card_no, guest_id, room_charges, incidental_charges, payment_status)
		 VALUES ($1, $2::uuid, $3, 0.00, 'Open')
		 RETURNING folio_id::text`,
		cardNo, guestID, totalCharge,
	).Scan(&folioID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to initialize folio: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// 5. Explicitly ensure room is Occupied
	_, _ = tx.Exec(`UPDATE rooms SET status = 'Occupied', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1`, req.RoomNumber)

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	balance := totalCharge - req.DepositPaid

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Guest checked in successfully",
		"folio": map[string]interface{}{
			"id":             folioID,
			"card_no":        cardNo,
			"guest_id":       guestID,
			"guest_name":     req.GuestName,
			"room_number":    req.RoomNumber,
			"check_in_date":  arrivalDate,
			"check_out_date": departureDate,
			"total_charges":  totalCharge,
			"total_paid":     req.DepositPaid,
			"balance":        balance,
			"items": []map[string]interface{}{
				{
					"id":          "acc-1",
					"description": fmt.Sprintf("%d Nights Accommodation (%s)", req.Nights, roomType),
					"amount":      totalCharge,
					"category":    "Accommodation",
					"timestamp":   time.Now().Format(time.RFC3339),
				},
			},
		},
	})
}

// GetFolio returns the live itemized folio for a room or card number.
func (h *FrontDeskHandler) GetFolio(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	roomNo := r.URL.Query().Get("room_number")
	cardNo := r.URL.Query().Get("card_no")

	var folioID, card, guestID, gFirst, gLast, gPhone, rmNo, arr, dep string
	var roomCharges, incidentalCharges, totalBalance, advancePaid float64
	var pStatus string

	query := `
		SELECT f.folio_id::text, f.card_no, f.guest_id::text, g.first_name, g.last_name, COALESCE(g.contact_no, ''),
		       b.room_number, b.arrival_date::text, b.departure_date::text,
		       f.room_charges, f.incidental_charges, f.total_balance, b.advance_paid, f.payment_status::text
		  FROM folios f
		  JOIN bookings b ON f.card_no = b.card_no
		  JOIN guests g ON f.guest_id = g.guest_id
		 WHERE (b.room_number = $1 OR f.card_no = $2)
		   AND f.payment_status != 'Closed'
		 ORDER BY f.created_at DESC
		 LIMIT 1
	`
	err := h.db.QueryRow(query, roomNo, cardNo).Scan(
		&folioID, &card, &guestID, &gFirst, &gLast, &gPhone,
		&rmNo, &arr, &dep,
		&roomCharges, &incidentalCharges, &totalBalance, &advancePaid, &pStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"active folio not found for room"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"database error: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Pull restaurant dockets charged to this room/folio
	items := []map[string]interface{}{
		{
			"id":          "item-rm",
			"description": "Room Accommodation Charges",
			"amount":      roomCharges,
			"category":    "Accommodation",
			"timestamp":   arr,
		},
	}
	if advancePaid > 0 {
		items = append(items, map[string]interface{}{
			"id":          "item-pay",
			"description": "Payments & Advance Deposits",
			"amount":      -advancePaid,
			"category":    "Payment",
			"timestamp":   arr,
		})
	}

	docketRows, err := h.db.Query(
		`SELECT docket_id::text, total_amount, sale_date::text
		   FROM restaurant_dockets
		  WHERE (folio_id = $1::uuid OR room_number = $2) AND payment_method = 'Room Charge'`,
		folioID, rmNo,
	)
	if err == nil {
		defer docketRows.Close()
		for docketRows.Next() {
			var dID, dDate string
			var dAmount float64
			if err := docketRows.Scan(&dID, &dAmount, &dDate); err == nil {
				items = append(items, map[string]interface{}{
					"id":          "dkt-" + dID,
					"description": fmt.Sprintf("Restaurant Dining Docket (%s)", dID[:8]),
					"amount":      dAmount,
					"category":    "Restaurant",
					"timestamp":   dDate,
				})
			}
		}
	}

	balance := totalBalance - advancePaid
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                 folioID,
		"card_no":            card,
		"guest_id":           guestID,
		"guest_name":         fmt.Sprintf("%s %s", gFirst, gLast),
		"guest_phone":        gPhone,
		"room_number":        rmNo,
		"check_in_date":      arr,
		"check_out_date":     dep,
		"room_charges":       roomCharges,
		"incidental_charges": incidentalCharges,
		"total_charges":      totalBalance,
		"total_paid":         advancePaid,
		"balance":            balance,
		"status":             pStatus,
		"items":              items,
	})
}

// RecordPayment records a payment toward a folio and updates advance_paid.
func (h *FrontDeskHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		CardNo     string  `json:"card_no"`
		RoomNumber string  `json:"room_number"`
		Amount     float64 `json:"amount"`
		Method     string  `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"failed to begin transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Locate booking by room or card
	var cardNo string
	var totalBalance, advancePaid float64
	err = tx.QueryRow(
		`SELECT b.card_no, f.total_balance, b.advance_paid
		   FROM bookings b
		   JOIN folios f ON b.card_no = f.card_no
		  WHERE b.room_number = $1 OR b.card_no = $2
		  ORDER BY b.created_at DESC LIMIT 1 FOR UPDATE`,
		req.RoomNumber, req.CardNo,
	).Scan(&cardNo, &totalBalance, &advancePaid)
	if err != nil {
		http.Error(w, `{"error":"active folio not found"}`, http.StatusNotFound)
		return
	}

	newAdvance := advancePaid + req.Amount
	_, err = tx.Exec(
		`UPDATE bookings SET advance_paid = $1 WHERE card_no = $2`,
		newAdvance, cardNo,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to update advance paid"}`, http.StatusInternalServerError)
		return
	}

	// If balance is now settled, mark folio status as 'Settled'
	if totalBalance-newAdvance <= 0.01 {
		_, _ = tx.Exec(`UPDATE folios SET payment_status = 'Settled' WHERE card_no = $1`, cardNo)
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"card_no":     cardNo,
		"total_paid":  newAdvance,
		"new_balance": totalBalance - newAdvance,
	})
}

// CheckOut enforces the Zero-Balance Lock before closing a folio and releasing the room.
func (h *FrontDeskHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		CardNo     string `json:"card_no"`
		RoomNumber string `json:"room_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"transaction start failed"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var cardNo, roomNumber string
	var totalBalance, advancePaid float64
	err = tx.QueryRow(
		`SELECT b.card_no, b.room_number, f.total_balance, b.advance_paid
		   FROM bookings b
		   JOIN folios f ON b.card_no = f.card_no
		  WHERE b.room_number = $1 OR b.card_no = $2
		  ORDER BY b.created_at DESC LIMIT 1 FOR UPDATE`,
		req.RoomNumber, req.CardNo,
	).Scan(&cardNo, &roomNumber, &totalBalance, &advancePaid)
	if err != nil {
		http.Error(w, `{"error":"booking not found"}`, http.StatusNotFound)
		return
	}

	// Zero-Balance Lock verification!
	outstanding := totalBalance - advancePaid
	if outstanding > 0.01 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Zero-Balance Lock Enforced",
			"message": fmt.Sprintf("Guest folio has an outstanding balance of MWK %.2f. All charges must be settled prior to checkout.", outstanding),
			"balance": outstanding,
		})
		return
	}

	// Close Folio, release room to 'Available', and mark booking checked_in = FALSE
	_, _ = tx.Exec(`UPDATE folios SET payment_status = 'Closed', updated_at = CURRENT_TIMESTAMP WHERE card_no = $1`, cardNo)
	_, _ = tx.Exec(`UPDATE rooms SET status = 'Available', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1`, roomNumber)
	_, _ = tx.Exec(`UPDATE bookings SET checked_in = FALSE, departure_date = CURRENT_DATE WHERE card_no = $1`, cardNo)

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Checkout completed successfully. Room released for housekeeping inspection.",
		"room_number": roomNumber,
		"card_no":     cardNo,
	})
}

// ExtendStay adds extra nights to a guest stay and recalculates folio room charges.
func (h *FrontDeskHandler) ExtendStay(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		RoomNumber string `json:"room_number"`
		CardNo     string `json:"card_no"`
		ExtraDays  int    `json:"extra_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ExtraDays <= 0 {
		http.Error(w, `{"error":"invalid extra days"}`, http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"transaction failed"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var cardNo, depDate string
	var roomRate, roomCharges float64
	err = tx.QueryRow(
		`SELECT b.card_no, b.departure_date::text, r.room_rate, f.room_charges
		   FROM bookings b
		   JOIN rooms r ON b.room_number = r.room_number
		   JOIN folios f ON b.card_no = f.card_no
		  WHERE b.room_number = $1 OR b.card_no = $2
		  ORDER BY b.created_at DESC LIMIT 1 FOR UPDATE`,
		req.RoomNumber, req.CardNo,
	).Scan(&cardNo, &depDate, &roomRate, &roomCharges)
	if err != nil {
		http.Error(w, `{"error":"booking not found"}`, http.StatusNotFound)
		return
	}

	currentDep, _ := time.Parse("2006-01-02", depDate)
	newDep := currentDep.AddDate(0, 0, req.ExtraDays).Format("2006-01-02")
	additionalCharge := roomRate * float64(req.ExtraDays)
	newRoomCharges := roomCharges + additionalCharge

	_, err = tx.Exec(`UPDATE bookings SET departure_date = $1 WHERE card_no = $2`, newDep, cardNo)
	if err != nil {
		http.Error(w, `{"error":"failed to update departure date"}`, http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`UPDATE folios SET room_charges = $1, payment_status = 'Open', updated_at = CURRENT_TIMESTAMP WHERE card_no = $2`, newRoomCharges, cardNo)
	if err != nil {
		http.Error(w, `{"error":"failed to update folio charges"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":           true,
		"new_checkout_date": newDep,
		"added_charges":     additionalCharge,
		"new_room_charges":  newRoomCharges,
	})
}

// GetReservations lists all bookings.
func (h *FrontDeskHandler) GetReservations(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	rows, err := h.db.Query(
		`SELECT b.card_no, g.first_name, g.last_name, COALESCE(g.contact_no, ''),
		        b.room_number, r.room_type::text, b.arrival_date::text, b.departure_date::text,
		        COALESCE(b.number_of_days, 1), COALESCE(b.advance_paid, 0), COALESCE(g.passport_no, ''),
		        COALESCE(g.email, ''), COALESCE(g.address, ''), COALESCE(g.city, ''),
		        COALESCE(g.nationality, ''), b.mode_of_payment::text, r.status::text,
		        COALESCE(f.payment_status::text, 'Open'), b.front_officer_clerk, COALESCE(b.checked_in, FALSE)
		   FROM bookings b
		   JOIN guests g ON b.guest_id = g.guest_id
		   JOIN rooms r ON b.room_number = r.room_number
		   LEFT JOIN folios f ON b.card_no = f.card_no
		  ORDER BY b.created_at DESC`,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query bookings: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var card, fName, lName, phone, rmNo, rType, arr, dep string
		var passport, email, addr, city, nat, payMode, rmStatus, folStatus, clerk string
		var days int
		var adv float64
		var checkedIn bool
		if err := rows.Scan(&card, &fName, &lName, &phone, &rmNo, &rType, &arr, &dep, &days, &adv,
			&passport, &email, &addr, &city, &nat, &payMode, &rmStatus, &folStatus, &clerk, &checkedIn); err == nil {
			list = append(list, map[string]interface{}{
				"id":           card,
				"res_code":     card,
				"card_no":      card,
				"guest_name":   fmt.Sprintf("%s %s", fName, lName),
				"first_name":   fName,
				"last_name":    lName,
				"guest_phone":  phone,
				"guest_id_doc": passport,
				"email":        email,
				"address":      addr,
				"city":         city,
				"nationality":  nat,
				"room_number":  rmNo,
				"room_type":    rType,
				"check_in":     arr,
				"check_out":    dep,
				"nights":       days,
				"deposit_paid": adv,
				"payment_mode": payMode,
				"room_status":  rmStatus,
				"folio_status": folStatus,
				"clerk":        clerk,
				"is_online":    strings.Contains(clerk, "Online"),
				"checked_in":   checkedIn,
			})
		}
	}

	json.NewEncoder(w).Encode(list)
}

// DailyReport generates check-in, occupancy and room sales analytics.
func (h *FrontDeskHandler) DailyReport(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var totalRooms, occupiedRooms int
	_ = h.db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&totalRooms)
	_ = h.db.QueryRow("SELECT COUNT(*) FROM rooms WHERE status = 'Occupied'").Scan(&occupiedRooms)

	var totalRoomRevenue, totalBalance float64
	_ = h.db.QueryRow("SELECT COALESCE(SUM(room_charges), 0), COALESCE(SUM(total_balance), 0) FROM folios WHERE payment_status != 'Closed'").Scan(&totalRoomRevenue, &totalBalance)

	occRate := 0.0
	if totalRooms > 0 {
		occRate = (float64(occupiedRooms) / float64(totalRooms)) * 100.0
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_rooms":        totalRooms,
		"occupied_rooms":     occupiedRooms,
		"available_rooms":    totalRooms - occupiedRooms,
		"occupancy_rate":     fmt.Sprintf("%.1f%%", occRate),
		"total_room_revenue": totalRoomRevenue,
		"total_balance_due":  totalBalance,
		"timestamp":          time.Now().Format("2006-01-02 15:04:05"),
	})
}
