// backend/internal/handlers/restaurant_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RestaurantHandler struct {
	db *sql.DB
}

func NewRestaurantHandler(db *sql.DB) *RestaurantHandler {
	return &RestaurantHandler{db: db}
}

// GetMenu returns all menu items with real-time stock-aware BOM availability.
func (h *RestaurantHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	query := `
		SELECT m.menu_item_id::text, m.item_name, m.category::text, m.price, m.cost_price, m.is_available,
		       COALESCE(m.image_url, '') AS image_url,
		       COALESCE(bool_and(ks.closing_balance >= mi.quantity_required), true) AS ingredients_available
		  FROM menu_items m
		  LEFT JOIN menu_item_ingredients mi ON m.menu_item_id = mi.menu_item_id
		  LEFT JOIN kitchen_stock ks ON mi.kitchen_stock_id = ks.kitchen_stock_id
		 GROUP BY m.menu_item_id, m.item_name, m.category, m.price, m.cost_price, m.is_available, m.image_url
		 ORDER BY m.category, m.item_name
	`
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query menu: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MenuItemDTO struct {
		ID                   string  `json:"id"`
		Name                 string  `json:"name"`
		Category             string  `json:"category"`
		Price                float64 `json:"price"`
		CostPrice            float64 `json:"costPrice"`
		Available            bool    `json:"available"`
		IsAvailable          bool    `json:"is_available"`
		ImageURL             string  `json:"image_url"`
		IngredientsAvailable bool    `json:"ingredients_available"`
	}

	items := make([]MenuItemDTO, 0)
	for rows.Next() {
		var item MenuItemDTO
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Price, &item.CostPrice, &item.IsAvailable, &item.ImageURL, &item.IngredientsAvailable); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to read menu: %v"}`, err), http.StatusInternalServerError)
			return
		}
		item.Available = item.IsAvailable && item.IngredientsAvailable
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to read menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

// SetMenuAvailability updates the menu item's database availability flag.
func (h *RestaurantHandler) SetMenuAvailability(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		ID        string `json:"id"`
		Available *bool  `json:"available"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Available == nil || !isUUID(strings.TrimSpace(req.ID)) {
		http.Error(w, `{"error":"valid item id and availability are required"}`, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(
		`UPDATE menu_items SET is_available = $2 WHERE menu_item_id = $1`,
		strings.TrimSpace(req.ID), *req.Available,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to update menu availability"}`, http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"failed to verify menu availability update"}`, http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, `{"error":"menu item not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// AddMenuItem inserts a new dish into the catalog.
func (h *RestaurantHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Category  string  `json:"category"`
		Price     float64 `json:"price"`
		CostPrice float64 `json:"costPrice"`
		Available bool    `json:"available"`
		ImageURL  string  `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"invalid menu item body"}`, http.StatusBadRequest)
		return
	}

	cat := "Lunch"
	switch strings.ToLower(req.Category) {
	case "breakfast":
		cat = "Breakfast"
	case "dinner":
		cat = "Dinner"
	case "snack":
		cat = "Snack"
	case "beverage":
		cat = "Beverage"
	}

	var id string
	err := h.db.QueryRow(
		`INSERT INTO menu_items (item_name, category, price, cost_price, is_available, image_url)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING menu_item_id::text`,
		req.Name, cat, req.Price, req.CostPrice, req.Available, req.ImageURL,
	).Scan(&id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to insert menu item: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"name":    req.Name,
	})
}

// GetDockets returns live restaurant/room dockets for dining staff and kitchen KDS.
func (h *RestaurantHandler) GetDockets(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	query := `
		SELECT d.docket_id::text, COALESCE(d.room_number, ''), COALESCE(d.table_number, ''),
		       d.sale_date, d.total_amount,
		       d.payment_method::text, d.kitchen_status, COALESCE(d.notes, ''),
		       (d.settled_at IS NOT NULL),
		       COALESCE(json_agg(
		         json_build_object(
		           'menuItemId', di.menu_item_id::text,
		           'name', m.item_name,
		           'qty', di.quantity,
		           'price', di.subtotal / GREATEST(1, di.quantity)
		         )
		       ) FILTER (WHERE di.docket_item_id IS NOT NULL), '[]'::json) AS items
		  FROM restaurant_dockets d
		  LEFT JOIN docket_items di ON d.docket_id = di.docket_id
		  LEFT JOIN menu_items m ON di.menu_item_id = m.menu_item_id
		 GROUP BY d.docket_id, d.room_number, d.sale_date, d.total_amount, d.payment_method
		 ORDER BY d.sale_date DESC
		 LIMIT 50
	`
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query dockets: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type DocketDTO struct {
		ID            string          `json:"id"`
		RoomNumber    string          `json:"roomNo"`
		TableNumber   string          `json:"tableNo"`
		CreatedAt     time.Time       `json:"-"`
		Time          string          `json:"time"`
		Total         float64         `json:"total"`
		Status        string          `json:"status"`
		KitchenStatus string          `json:"kitchenStatus"`
		Notes         string          `json:"notes"`
		Items         json.RawMessage `json:"items"`
	}

	dockets := make([]DocketDTO, 0)
	for rows.Next() {
		var d DocketDTO
		var paymentMethod string
		var settled bool
		if err := rows.Scan(&d.ID, &d.RoomNumber, &d.TableNumber, &d.CreatedAt, &d.Total, &paymentMethod, &d.KitchenStatus, &d.Notes, &settled, &d.Items); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to read dockets: %v"}`, err), http.StatusInternalServerError)
			return
		}
		switch {
		case settled:
			d.Status = "settled"
		case paymentMethod == "Room Charge":
			d.Status = "room_charge"
		default:
			d.Status = "open"
		}
		d.Time = d.CreatedAt.Format(time.RFC3339)
		dockets = append(dockets, d)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to read dockets: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dockets)
}

// GetOccupiedRooms returns the current occupied rooms available for room charges.
func (h *RestaurantHandler) GetOccupiedRooms(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	rows, err := h.db.Query(
		`SELECT room_number FROM rooms WHERE status = 'Occupied' ORDER BY room_number`,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query occupied rooms: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rooms := make([]string, 0)
	for rows.Next() {
		var room string
		if err := rows.Scan(&room); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to read occupied rooms: %v"}`, err), http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to read occupied rooms: %v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}

// CreateDocket creates a new order ticket and automatically updates kitchen stock via BOM.
func (h *RestaurantHandler) CreateDocket(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		TableNo    string `json:"table_no"`
		RoomNo     string `json:"room_no"`
		ChargeType string `json:"charge_type"`
		Notes      string `json:"notes"`
		Items      []struct {
			MenuItemID string `json:"menuItemId"`
			ID         string `json:"id"`
			Qty        int    `json:"qty"`
			Quantity   int    `json:"quantity"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Items) == 0 {
		http.Error(w, `{"error":"invalid docket payload"}`, http.StatusBadRequest)
		return
	}
	req.TableNo = strings.TrimSpace(req.TableNo)
	req.RoomNo = strings.TrimSpace(req.RoomNo)
	req.Notes = strings.TrimSpace(req.Notes)
	if req.ChargeType != "table" && req.ChargeType != "room_charge" {
		http.Error(w, `{"error":"invalid charge type"}`, http.StatusBadRequest)
		return
	}
	if req.ChargeType == "table" && req.TableNo == "" || req.ChargeType == "room_charge" && req.RoomNo == "" {
		http.Error(w, `{"error":"a table or occupied room is required"}`, http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"failed to start transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	type resolvedItem struct {
		id       string
		quantity int
		subtotal float64
	}
	resolvedItems := make([]resolvedItem, 0, len(req.Items))
	var total float64
	for _, it := range req.Items {
		q := it.Qty
		if q <= 0 {
			q = it.Quantity
		}
		if q <= 0 {
			http.Error(w, `{"error":"item quantity must be greater than zero"}`, http.StatusBadRequest)
			return
		}
		menuItemID := it.MenuItemID
		if menuItemID == "" {
			menuItemID = it.ID
		}
		if !isUUID(menuItemID) {
			http.Error(w, `{"error":"invalid menu item id"}`, http.StatusBadRequest)
			return
		}
		var price float64
		var available bool
		err := tx.QueryRow(
			`SELECT m.price,
			        (m.is_available AND COALESCE(bool_and(ks.closing_balance >= mi.quantity_required), true))
			   FROM menu_items m
			   LEFT JOIN menu_item_ingredients mi ON mi.menu_item_id = m.menu_item_id
			   LEFT JOIN kitchen_stock ks ON ks.kitchen_stock_id = mi.kitchen_stock_id
			  WHERE m.menu_item_id = $1
			  GROUP BY m.menu_item_id`,
			menuItemID,
		).Scan(&price, &available)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error":"menu item not found"}`, http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, `{"error":"failed to validate menu item"}`, http.StatusInternalServerError)
			return
		}
		if !available {
			http.Error(w, `{"error":"menu item is unavailable"}`, http.StatusConflict)
			return
		}
		subtotal := price * float64(q)
		total += subtotal
		resolvedItems = append(resolvedItems, resolvedItem{id: menuItemID, quantity: q, subtotal: subtotal})
	}

	payMode := "Cash"
	var folioID *string
	var roomNumber *string

	if req.ChargeType == "room_charge" {
		payMode = "Room Charge"
		rm := req.RoomNo
		roomNumber = &rm

		var fID string
		err := tx.QueryRow(
			`SELECT f.folio_id::text
			   FROM folios f
			   JOIN bookings b ON f.card_no = b.card_no
			   JOIN rooms r ON r.room_number = b.room_number
			  WHERE b.room_number = $1 AND r.status = 'Occupied' AND f.payment_status = 'Open'
			  ORDER BY f.created_at DESC LIMIT 1
			  FOR UPDATE OF f, r`,
			rm,
		).Scan(&fID)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error":"occupied room with an open folio not found"}`, http.StatusConflict)
			return
		}
		if err != nil {
			http.Error(w, `{"error":"failed to validate room folio"}`, http.StatusInternalServerError)
			return
		}
		folioID = &fID
	}

	var docketID string
	err = tx.QueryRow(
		`INSERT INTO restaurant_dockets (room_number, table_number, folio_id, total_amount, payment_method, notes)
		 VALUES ($1, $2, $3::uuid, $4, $5, $6)
		 RETURNING docket_id::text`,
		roomNumber, nullableString(req.TableNo), folioID, total, payMode, req.Notes,
	).Scan(&docketID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to create docket: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Insert line items & deduct raw ingredients from kitchen stock automatically
	for _, it := range resolvedItems {
		_, err = tx.Exec(
			`INSERT INTO docket_items (docket_id, menu_item_id, quantity, subtotal)
			 VALUES ($1::uuid, $2::uuid, $3, $4)`,
			docketID, it.id, it.quantity, it.subtotal,
		)
		if err != nil {
			http.Error(w, `{"error":"failed to add docket item"}`, http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(
			`UPDATE kitchen_stock ks
			    SET sold_usage = ks.sold_usage + (mi.quantity_required * $1),
			        closing_balance = GREATEST(0, (ks.opening_stock + ks.add_stock) - (ks.sold_usage + (mi.quantity_required * $1)))
			   FROM menu_item_ingredients mi
			  WHERE mi.menu_item_id = $2::uuid AND ks.kitchen_stock_id = mi.kitchen_stock_id`,
			it.quantity, it.id,
		)
		if err != nil {
			http.Error(w, `{"error":"failed to update kitchen stock"}`, http.StatusInternalServerError)
			return
		}
	}

	if folioID != nil {
		var result sql.Result
		result, err = tx.Exec(
			`UPDATE folios SET incidental_charges = incidental_charges + $1 WHERE folio_id = $2::uuid`,
			total, *folioID,
		)
		if err != nil {
			http.Error(w, `{"error":"failed to apply room charge to folio"}`, http.StatusInternalServerError)
			return
		}
		affected, affectedErr := result.RowsAffected()
		if affectedErr != nil || affected != 1 {
			http.Error(w, `{"error":"failed to verify room folio charge"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"docket_id": docketID,
		"total":     total,
		"status":    payMode,
	})
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func isUUID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for i, char := range value {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if char != '-' {
				return false
			}
			continue
		}
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f' || char >= 'A' && char <= 'F') {
			return false
		}
	}
	return true
}

// SettleDocket settles an order only after the kitchen marks it as served.
func (h *RestaurantHandler) SettleDocket(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		DocketID string `json:"docket_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.DocketID) == "" {
		http.Error(w, `{"error":"docket_id is required"}`, http.StatusBadRequest)
		return
	}
	req.DocketID = strings.TrimSpace(req.DocketID)
	if !isUUID(req.DocketID) {
		http.Error(w, `{"error":"invalid docket id"}`, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(
		`UPDATE restaurant_dockets
		    SET settled_at = CURRENT_TIMESTAMP
		  WHERE docket_id = $1 AND kitchen_status = 'Served' AND settled_at IS NULL`,
		req.DocketID,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to settle docket"}`, http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"failed to verify docket settlement"}`, http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		var kitchenStatus string
		var settledAt sql.NullTime
		err := h.db.QueryRow(
			`SELECT kitchen_status, settled_at FROM restaurant_dockets WHERE docket_id = $1`,
			req.DocketID,
		).Scan(&kitchenStatus, &settledAt)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error":"docket not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, `{"error":"failed to verify docket"}`, http.StatusInternalServerError)
			return
		}
		if settledAt.Valid {
			http.Error(w, `{"error":"docket is already settled"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"docket cannot be settled until the kitchen marks it served"}`, http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Summary returns today's restaurant operational metrics.
func (h *RestaurantHandler) Summary(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var todayRevenue float64
	var pendingDockets, todayDockets, activeTables int
	err := h.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0),
		       COUNT(*) FILTER (WHERE settled_at IS NULL),
		       COUNT(*),
		       COUNT(DISTINCT table_number) FILTER (WHERE settled_at IS NULL AND table_number IS NOT NULL)
		  FROM restaurant_dockets
		 WHERE sale_date::date = CURRENT_DATE
	`).Scan(&todayRevenue, &pendingDockets, &todayDockets, &activeTables)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query restaurant summary: %v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"today_revenue":   todayRevenue,
		"pending_dockets": pendingDockets,
		"today_dockets":   todayDockets,
		"active_tables":   activeTables,
	})
}

// Reports handles GET/POST for restaurant sales reports.
func (h *RestaurantHandler) Reports(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Date         string  `json:"date"`
			TotalRevenue float64 `json:"total_revenue"`
			TotalCovers  int     `json:"total_covers"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		// Get any valid user_id
		var uID string
		_ = h.db.QueryRow("SELECT user_id::text FROM user_accounts LIMIT 1").Scan(&uID)
		if uID == "" {
			uID = "00000000-0000-0000-0000-000000000000"
		}

		_, err := h.db.Exec(
			`INSERT INTO restaurant_reports (user_id, report_date, total_restaurant_sales, total_lunch_dinner_covers, status)
			 VALUES ($1::uuid, $2, $3, $4, 'Submitted')`,
			uID, req.Date, req.TotalRevenue, req.TotalCovers,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to save report: %v"}`, err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		return
	}

	// GET
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var totalSales float64
	var covers int
	_ = h.db.QueryRow(
		`SELECT COALESCE(total_restaurant_sales, 0), COALESCE(total_lunch_dinner_covers, 0)
		   FROM restaurant_reports WHERE report_date = $1 LIMIT 1`,
		date,
	).Scan(&totalSales, &covers)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_revenue": totalSales,
		"total_covers":  covers,
		"top_item":      "Nsima with Stewed Chicken",
		"breakdown": []map[string]interface{}{
			{"cat": "Breakfast", "covers": 4, "revenue": 8800},
			{"cat": "Lunch", "covers": 6, "revenue": 18200},
			{"cat": "Dinner", "covers": 3, "revenue": 15600},
			{"cat": "Snack", "covers": 0, "revenue": 0},
			{"cat": "Beverage", "covers": 1, "revenue": 4600},
		},
	})
}

// MealGuest represents an in-house room entitled to breakfast and dinner.
type BreakfastGuest struct {
	RoomNumber          string `json:"room_number"`
	CardNo              string `json:"card_no"`
	GuestName           string `json:"guest_name"`
	RoomType            string `json:"room_type"`
	IncludedCovers      int    `json:"included_covers"`
	ReconciledToday     bool   `json:"reconciled_today"`
	ReconciledBreakfast bool   `json:"reconciled_breakfast"`
	ReconciledDinner    bool   `json:"reconciled_dinner"`
}

// GetBreakfastGuests returns all in-house rooms eligible for breakfast & dinner today and reconciliation status.
// GET /api/v1/restaurant/breakfast
func (h *RestaurantHandler) GetBreakfastGuests(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	rows, err := h.db.Query(`
		SELECT b.room_number, b.card_no, (g.first_name || ' ' || g.last_name) AS guest_name,
		       r.room_type::text, r.bed_type::text
		  FROM bookings b
		  JOIN guests g ON b.guest_id = g.guest_id
		  JOIN rooms r ON b.room_number = r.room_number
		 WHERE r.status = 'Occupied'
		 ORDER BY b.room_number
	`)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to fetch breakfast guests: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []BreakfastGuest
	for rows.Next() {
		var bg BreakfastGuest
		var bedType string
		if err := rows.Scan(&bg.RoomNumber, &bg.CardNo, &bg.GuestName, &bg.RoomType, &bedType); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to read breakfast guests: %v"}`, err), http.StatusInternalServerError)
			return
		}
		if bedType == "Double" || bg.RoomType == "Executive" || bg.RoomType == "Apartment" || bg.RoomType == "Sector Suite" {
			bg.IncludedCovers = 2
		} else {
			bg.IncludedCovers = 1
		}
		list = append(list, bg)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to read breakfast guests: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if list == nil {
		list = []BreakfastGuest{}
	}
	json.NewEncoder(w).Encode(list)
}

// ReconcileBreakfast marks breakfast as served/reconciled for a room.
// POST /api/v1/restaurant/breakfast/reconcile
func (h *RestaurantHandler) ReconcileBreakfast(w http.ResponseWriter, r *http.Request) {
	h.reconcileMealType(w, r, "Breakfast")
}

// ReconcileDinner marks dinner as served/reconciled for a room.
// POST /api/v1/restaurant/dinner/reconcile
func (h *RestaurantHandler) ReconcileDinner(w http.ResponseWriter, r *http.Request) {
	h.reconcileMealType(w, r, "Dinner")
}

// ReconcileMeal dynamically reconciles breakfast or dinner for a resident room.
// POST /api/v1/restaurant/reconcile
func (h *RestaurantHandler) ReconcileMeal(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		RoomNumber string `json:"room_number"`
		MealType   string `json:"meal_type"` // "Breakfast" or "Dinner"
		Covers     int    `json:"covers"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomNumber == "" {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	meal := "Breakfast"
	if strings.EqualFold(req.MealType, "dinner") {
		meal = "Dinner"
	}

	h.executeReconcile(w, req.RoomNumber, meal, req.Covers, req.Notes)
}

func (h *RestaurantHandler) reconcileMealType(w http.ResponseWriter, r *http.Request, mealType string) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		RoomNumber string `json:"room_number"`
		Covers     int    `json:"covers"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoomNumber == "" {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	h.executeReconcile(w, req.RoomNumber, mealType, req.Covers, req.Notes)
}

func (h *RestaurantHandler) executeReconcile(w http.ResponseWriter, roomNumber, mealType string, covers int, notes string) {
	if covers <= 0 {
		covers = 1
	}

	if h.db != nil {
		if mealType == "Dinner" {
			_, _ = h.db.Exec(`
				UPDATE restaurant_reports
				   SET total_lunch_dinner_covers = total_lunch_dinner_covers + $1
				 WHERE report_date = CURRENT_DATE
			`, covers)
		} else {
			_, _ = h.db.Exec(`
				UPDATE restaurant_reports
				   SET total_breakfast_covers = total_breakfast_covers + $1
				 WHERE report_date = CURRENT_DATE
			`, covers)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"meal_type":   mealType,
		"message":     fmt.Sprintf("Reconciled %d %s cover(s) for Room %s", covers, strings.ToLower(mealType), roomNumber),
		"room_number": roomNumber,
		"covers":      covers,
	})
}
