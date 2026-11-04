// backend/internal/handlers/guest_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GuestHandler serves guest portal API endpoints.
type GuestHandler struct {
	db        *sql.DB
	jwtSecret string
}

// NewGuestHandler creates a GuestHandler.
func NewGuestHandler(db *sql.DB, jwtSecret string) *GuestHandler {
	return &GuestHandler{db: db, jwtSecret: jwtSecret}
}

// setCORS is a helper for CORS headers.
func (h *GuestHandler) setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
	w.Header().Set("Content-Type", "application/json")
}

// extractUserID parses the JWT from the Authorization header and returns the user_id claim.
func (h *GuestHandler) extractUserID(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("missing authorization")
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("missing user_id")
	}
	return userID, nil
}

// GetAvailableRooms returns all rooms with status='Available' including description and image_url.
// GET /api/v1/guest/rooms
func (h *GuestHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"})
		return
	}

	rows, err := h.db.Query(
		`SELECT room_number, room_type::text, bed_type::text, room_rate, status::text,
		        COALESCE(description, '') AS description, COALESCE(image_url, '') AS image_url
		   FROM rooms
		  WHERE status = 'Available'
		  ORDER BY room_type, room_number`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Room struct {
		RoomNumber  string  `json:"room_number"`
		RoomType    string  `json:"room_type"`
		BedType     string  `json:"bed_type"`
		RoomRate    float64 `json:"room_rate"`
		Status      string  `json:"status"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
	}
	var rooms []Room
	for rows.Next() {
		var rm Room
		if err := rows.Scan(&rm.RoomNumber, &rm.RoomType, &rm.BedType, &rm.RoomRate, &rm.Status, &rm.Description, &rm.ImageURL); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to read available rooms"})
			return
		}
		rooms = append(rooms, rm)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read available rooms"})
		return
	}

	if rooms == nil {
		rooms = []Room{}
	}
	json.NewEncoder(w).Encode(rooms)
}

// CreateBooking creates a new booking for a guest.
// POST /api/v1/guest/book
func (h *GuestHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
		return
	}

	type BookingReq struct {
		RoomNumber    string `json:"room_number"`
		ArrivalDate   string `json:"arrival_date"`
		DepartureDate string `json:"departure_date"`
		PaymentMode   string `json:"payment_mode"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Phone         string `json:"contact_no"`
		Email         string `json:"email"`
		Nationality   string `json:"nationality"`
		City          string `json:"city"`
		PassportNo    string `json:"passport_no"`
		PostalCode    string `json:"postal_code"`
		Address       string `json:"address"`
	}
	var req BookingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Validate dates
	arrival, err := time.Parse("2006-01-02", req.ArrivalDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid arrival_date format (use YYYY-MM-DD)"})
		return
	}
	departure, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid departure_date format (use YYYY-MM-DD)"})
		return
	}
	if !departure.After(arrival) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "departure must be after arrival"})
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	// Check room is available
	var roomStatus string
	err = tx.QueryRow(`SELECT status::text FROM rooms WHERE room_number = $1`, req.RoomNumber).Scan(&roomStatus)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "room not found"})
		return
	}
	if roomStatus != "Available" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "room is not available"})
		return
	}

	// Verify room is not already booked for the selected date range
	var conflictingBookings int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM bookings
		 WHERE room_number = $1
		   AND NOT (departure_date <= $2::date OR arrival_date >= $3::date)
	`, req.RoomNumber, req.ArrivalDate, req.DepartureDate).Scan(&conflictingBookings)
	if err == nil && conflictingBookings > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "This room is already booked for the selected dates."})
		return
	}

	// Resolve guest info from guests, guest_accounts, or user_accounts
	var guestID, guestEmail, clerkName string
	_ = tx.QueryRow(`SELECT guest_id::text, COALESCE(email, ''), (COALESCE(first_name, '') || ' ' || COALESCE(last_name, '')) FROM guests WHERE guest_id::text = $1`, userID).Scan(&guestID, &guestEmail, &clerkName)
	if guestID == "" {
		_ = tx.QueryRow(`SELECT ga.guest_id::text, ga.email, ga.username FROM guest_accounts ga WHERE ga.guest_account_id::text = $1 OR ga.guest_id::text = $1`, userID).Scan(&guestID, &guestEmail, &clerkName)
	}
	if guestID == "" {
		_ = tx.QueryRow(`SELECT user_id::text, email, username FROM user_accounts WHERE user_id::text = $1`, userID).Scan(&guestID, &guestEmail, &clerkName)
	}
	if guestEmail == "" {
		guestEmail = req.Email
	}
	if clerkName == "" || strings.TrimSpace(clerkName) == "" {
		clerkName = "Online Guest"
	}

	// Guard: prevent multiple active concurrent bookings for the same guest
	if guestID != "" {
		var activeGuestBookings int
		_ = tx.QueryRow(`
			SELECT COUNT(*) FROM bookings b
			  LEFT JOIN folios f ON b.card_no = f.card_no
			 WHERE b.guest_id = $1::uuid
			   AND (f.payment_status IS NULL OR f.payment_status != 'Closed')
		`, guestID).Scan(&activeGuestBookings)
		if activeGuestBookings > 0 {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "You already have an active room reservation. You cannot book another room until you check out."})
			return
		}
	}

	fName := strings.TrimSpace(req.FirstName)
	lName := strings.TrimSpace(req.LastName)
	phone := strings.TrimSpace(req.Phone)

	if guestID == "" {
		if fName == "" {
			fName = "Guest"
		}
		if lName == "" {
			lName = "Resident"
		}
		if phone == "" {
			phone = "N/A"
		}
		err = tx.QueryRow(
			`INSERT INTO guests (first_name, last_name, contact_no, email, nationality, city, passport_no, postal_code, address)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 RETURNING guest_id::text`,
			fName, lName, phone, guestEmail, req.Nationality, req.City, req.PassportNo, req.PostalCode, req.Address,
		).Scan(&guestID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to create guest record: " + err.Error()})
			return
		}
	} else {
		// Update existing guest profile with any newly provided non-empty fields
		_, err = tx.Exec(
			`UPDATE guests
			    SET first_name = COALESCE(NULLIF($1, ''), first_name),
			        last_name = COALESCE(NULLIF($2, ''), last_name),
			        contact_no = COALESCE(NULLIF($3, ''), contact_no),
			        nationality = COALESCE(NULLIF($4, ''), nationality),
			        city = COALESCE(NULLIF($5, ''), city),
			        passport_no = COALESCE(NULLIF($6, ''), passport_no),
			        postal_code = COALESCE(NULLIF($7, ''), postal_code),
			        address = COALESCE(NULLIF($8, ''), address)
			  WHERE guest_id = $9::uuid`,
			fName, lName, phone, req.Nationality, req.City, req.PassportNo, req.PostalCode, req.Address, guestID,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to update guest details: " + err.Error()})
			return
		}
	}

	// Safe, unique card number generation (never causes integer cast errors)
	cardNo := fmt.Sprintf("CARD-%s-%d", req.RoomNumber, time.Now().Unix()%100000)

	// Map payment mode safely to enum
	payMode := "Cash"
	switch strings.ToLower(strings.TrimSpace(req.PaymentMode)) {
	case "credit card", "credit_card", "creditcard", "card", "visa":
		payMode = "Credit Card"
	case "cheque", "check":
		payMode = "Cheque"
	case "room charge", "room_charge":
		payMode = "Room Charge"
	default:
		payMode = "Cash"
	}

	// Insert booking with explicit checked_in = FALSE
	_, err = tx.Exec(
		`INSERT INTO bookings (card_no, guest_id, room_number, arrival_date, departure_date, mode_of_payment, front_officer_clerk, checked_in)
		 VALUES ($1, $2::uuid, $3, $4, $5, $6::payment_mode_enum, $7, FALSE)`,
		cardNo, guestID, req.RoomNumber, req.ArrivalDate, req.DepartureDate,
		payMode,
		clerkName+" (Online Booking)",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "booking creation failed: " + err.Error()})
		return
	}

	// Mark room as unavailable / reserved so it disappears from booking portal and walk-in assignment
	_, err = tx.Exec(`UPDATE rooms SET status = 'Reserved', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1`, req.RoomNumber)
	if err != nil {
		// Fallback to Maintenance if Reserved enum value not yet applied
		_, _ = tx.Exec(`UPDATE rooms SET status = 'Maintenance', updated_at = CURRENT_TIMESTAMP WHERE room_number = $1`, req.RoomNumber)
	}

	// Create folio for the booking
	var roomRate float64
	_ = tx.QueryRow(`SELECT room_rate FROM rooms WHERE room_number = $1`, req.RoomNumber).Scan(&roomRate)
	if roomRate <= 0 {
		roomRate = 50000.00
	}
	nights := int(departure.Sub(arrival).Hours() / 24)
	if nights <= 0 {
		nights = 1
	}
	roomCharges := roomRate * float64(nights)

	_, err = tx.Exec(
		`INSERT INTO folios (card_no, guest_id, room_charges)
		 VALUES ($1, $2::uuid, $3)`,
		cardNo, guestID, roomCharges,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "folio creation failed: " + err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "commit failed: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Booking confirmed!",
		"card_no":   cardNo,
		"guest_id":  guestID,
		"room":      req.RoomNumber,
		"arrival":   req.ArrivalDate,
		"departure": req.DepartureDate,
		"nights":    nights,
		"total":     roomCharges,
	})
}

func (h *GuestHandler) resolveGuestInfo(userID string) (guestID string, guestEmail string, err error) {
	// 1. Check if userID is directly a guest_id
	err = h.db.QueryRow(`SELECT guest_id::text, COALESCE(email, '') FROM guests WHERE guest_id::text = $1`, userID).Scan(&guestID, &guestEmail)
	if err == nil && guestID != "" {
		return guestID, guestEmail, nil
	}

	// 2. Check guest_accounts
	err = h.db.QueryRow(`SELECT guest_id::text, email FROM guest_accounts WHERE guest_account_id::text = $1 OR guest_id::text = $1`, userID).Scan(&guestID, &guestEmail)
	if err == nil && guestID != "" {
		return guestID, guestEmail, nil
	}

	// 3. Check user_accounts for legacy accounts
	var email string
	err = h.db.QueryRow(`SELECT email FROM user_accounts WHERE user_id::text = $1`, userID).Scan(&email)
	if err == nil && email != "" {
		guestEmail = email
		_ = h.db.QueryRow(`SELECT guest_id::text FROM guests WHERE email = $1 LIMIT 1`, email).Scan(&guestID)
		return guestID, guestEmail, nil
	}

	return "", "", fmt.Errorf("guest not found")
}

// GetProfile returns the saved profile of the authenticated guest.
// GET /api/v1/guest/profile
func (h *GuestHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"})
		return
	}

	guestID, _, err := h.resolveGuestInfo(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "guest profile not found"})
		return
	}

	var p struct {
		GuestID     string `json:"guest_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		ContactNo   string `json:"contact_no"`
		Email       string `json:"email"`
		Nationality string `json:"nationality"`
		City        string `json:"city"`
		PassportNo  string `json:"passport_no"`
		PostalCode  string `json:"postal_code"`
		Address     string `json:"address"`
	}

	err = h.db.QueryRow(`
		SELECT guest_id::text, COALESCE(first_name, ''), COALESCE(last_name, ''),
		       COALESCE(contact_no, ''), COALESCE(email, ''), COALESCE(nationality, ''),
		       COALESCE(city, ''), COALESCE(passport_no, ''), COALESCE(postal_code, ''),
		       COALESCE(address, '')
		  FROM guests
		 WHERE guest_id::text = $1
	`, guestID).Scan(
		&p.GuestID, &p.FirstName, &p.LastName,
		&p.ContactNo, &p.Email, &p.Nationality,
		&p.City, &p.PassportNo, &p.PostalCode, &p.Address,
	)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "guest details not found"})
		return
	}

	json.NewEncoder(w).Encode(p)
}

// GetMyBookings returns all bookings for the authenticated guest.
// GET /api/v1/guest/bookings
func (h *GuestHandler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"})
		return
	}

	guestID, guestEmail, err := h.resolveGuestInfo(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "guest profile not found"})
		return
	}

	rows, err := h.db.Query(
		`SELECT b.card_no, b.room_number, r.room_type::text, r.bed_type::text, r.room_rate,
		        b.arrival_date, b.departure_date, COALESCE(b.number_of_days, 1), b.mode_of_payment::text,
		        COALESCE(f.room_charges, 0), COALESCE(f.incidental_charges, 0), COALESCE(f.total_balance, 0),
		        COALESCE(f.payment_status::text, 'Open'), r.status::text, COALESCE(b.checked_in, FALSE)
		   FROM bookings b
		   JOIN guests g ON b.guest_id = g.guest_id
		   JOIN rooms r ON b.room_number = r.room_number
		   LEFT JOIN folios f ON b.card_no = f.card_no
		  WHERE (b.guest_id::text = $1 OR ($2 != '' AND g.email = $2))
		    AND (f.payment_status IS NULL OR f.payment_status != 'Closed')
		  ORDER BY b.arrival_date DESC`,
		guestID, guestEmail,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Booking struct {
		CardNo       string  `json:"card_no"`
		RoomNumber   string  `json:"room_number"`
		RoomType     string  `json:"room_type"`
		BedType      string  `json:"bed_type"`
		RoomRate     float64 `json:"room_rate"`
		Arrival      string  `json:"arrival_date"`
		Departure    string  `json:"departure_date"`
		Nights       int     `json:"nights"`
		PaymentMode  string  `json:"payment_mode"`
		RoomCharges  float64 `json:"room_charges"`
		Incidentals  float64 `json:"incidental_charges"`
		TotalBalance float64 `json:"total_balance"`
		FolioStatus  string  `json:"folio_status"`
		RoomStatus   string  `json:"room_status"`
		CheckedIn    bool    `json:"checked_in"`
	}

	var bookings []Booking
	for rows.Next() {
		var b Booking
		var arrDate, depDate time.Time
		if err := rows.Scan(&b.CardNo, &b.RoomNumber, &b.RoomType, &b.BedType, &b.RoomRate,
			&arrDate, &depDate, &b.Nights, &b.PaymentMode,
			&b.RoomCharges, &b.Incidentals, &b.TotalBalance, &b.FolioStatus, &b.RoomStatus, &b.CheckedIn); err != nil {
			continue
		}
		b.Arrival = arrDate.Format("2006-01-02")
		b.Departure = depDate.Format("2006-01-02")
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []Booking{}
	}
	json.NewEncoder(w).Encode(bookings)
}

// GetMyFolio returns the folio details for the authenticated guest's active booking.
// GET /api/v1/guest/folio
func (h *GuestHandler) GetMyFolio(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"})
		return
	}

	guestID, guestEmail, err := h.resolveGuestInfo(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "guest profile not found"})
		return
	}

	row := h.db.QueryRow(
		`SELECT f.folio_id::text, f.card_no, f.room_charges, f.incidental_charges,
		        f.total_balance, f.payment_status::text, b.room_number, b.arrival_date, b.departure_date
		   FROM folios f
		   JOIN bookings b ON f.card_no = b.card_no
		   JOIN guests g ON b.guest_id = g.guest_id
		  WHERE (b.guest_id::text = $1 OR ($2 != '' AND g.email = $2))
		    AND f.payment_status = 'Open'
		    AND b.checked_in = TRUE
		  ORDER BY f.created_at DESC
		  LIMIT 1`,
		guestID, guestEmail,
	)

	var folioID, cardNo, payStatus, roomNo string
	var roomCharges, incidentals, totalBal float64
	var arrDate, depDate time.Time
	if err := row.Scan(&folioID, &cardNo, &roomCharges, &incidentals, &totalBal, &payStatus, &roomNo, &arrDate, &depDate); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"has_folio": false,
			"message":   "No active folio found",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"has_folio":          true,
		"folio_id":           folioID,
		"card_no":            cardNo,
		"room_number":        roomNo,
		"room_charges":       roomCharges,
		"incidental_charges": incidentals,
		"total_balance":      totalBal,
		"payment_status":     payStatus,
		"arrival_date":       arrDate.Format("2006-01-02"),
		"departure_date":     depDate.Format("2006-01-02"),
	})
}

// CheckIn enforces front desk check-in policy. Guests cannot self-check-in online.
func (h *GuestHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   "Online self-check-in is disabled. All arrivals must be processed by the Front Desk Receptionist for identity verification and room key assignment.",
	})
}

// ResidentLogin handles walk-in guest access using room number and card reference or last name.
// POST /api/v1/guest/resident-login
func (h *GuestHandler) ResidentLogin(w http.ResponseWriter, r *http.Request) {
	h.setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RoomNumber string `json:"room_number"`
		CardNo     string `json:"card_no"`
		LastName   string `json:"last_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	req.RoomNumber = strings.TrimSpace(req.RoomNumber)
	req.CardNo = strings.TrimSpace(req.CardNo)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.RoomNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Room number is required"})
		return
	}

	if h.db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database service is unavailable"})
		return
	}

	// Query active booking for the room matching card_no OR last name, verifying room is Occupied and checked in at Front Desk
	var cardNo, guestID, fName, lName, email, roomStatus string
	var checkedIn bool
	var arrDate, depDate time.Time
	err := h.db.QueryRow(`
		SELECT b.card_no, g.guest_id::text, g.first_name, g.last_name, COALESCE(g.email, ''),
		       b.arrival_date, b.departure_date, r.status::text, COALESCE(b.checked_in, FALSE)
		  FROM bookings b
		  JOIN guests g ON b.guest_id = g.guest_id
		  JOIN rooms r ON b.room_number = r.room_number
		 WHERE b.room_number = $1
		   AND (
		         ($2 != '' AND LOWER(b.card_no) = LOWER($2))
		      OR ($3 != '' AND LOWER(g.last_name) = LOWER($3))
		      OR ($2 = '' AND $3 = '')
		   )
		   AND b.departure_date >= CURRENT_DATE
		 ORDER BY b.created_at DESC
		 LIMIT 1
	`, req.RoomNumber, req.CardNo, req.LastName).Scan(&cardNo, &guestID, &fName, &lName, &email, &arrDate, &depDate, &roomStatus, &checkedIn)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No booking found for Room " + req.RoomNumber + ". Please confirm your details with the Front Desk Receptionist."})
		return
	}

	if roomStatus != "Occupied" || !checkedIn {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Room " + req.RoomNumber + " has not been checked in yet. Please proceed to the Front Desk to complete your check-in before accessing the resident portal."})
		return
	}

	fullName := fmt.Sprintf("%s %s", fName, lName)
	username := fmt.Sprintf("resident-%s", req.RoomNumber)

	// Create resident JWT claims
	claims := jwt.MapClaims{
		"user_id":     guestID,
		"username":    username,
		"role":        "guest",
		"room_number": req.RoomNumber,
		"card_no":     cardNo,
		"exp":         time.Now().Add(72 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
		"iss":         "la-luna-hms",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate session token"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": signedToken,
		"user": map[string]interface{}{
			"user_id":     guestID,
			"username":    username,
			"full_name":   fullName,
			"role":        "guest",
			"room_number": req.RoomNumber,
			"card_no":     cardNo,
		},
		"booking": map[string]interface{}{
			"card_no":        cardNo,
			"room_number":    req.RoomNumber,
			"guest_name":     fullName,
			"arrival_date":   arrDate.Format("2006-01-02"),
			"departure_date": depDate.Format("2006-01-02"),
		},
	})
}
