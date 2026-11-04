// backend/internal/handlers/kitchen_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type KitchenHandler struct {
	db *sql.DB
}

func NewKitchenHandler(db *sql.DB) *KitchenHandler {
	return &KitchenHandler{db: db}
}

// GetDockets returns current kitchen tickets from restaurant dockets.
func (h *KitchenHandler) GetDockets(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	rows, err := h.db.Query(`
		SELECT d.docket_id::text, COALESCE(d.room_number, ''), COALESCE(d.table_number, ''),
		       d.sale_date, d.kitchen_status,
		       COALESCE(d.notes, ''),
		       COALESCE(json_agg(json_build_object('name', m.item_name, 'qty', di.quantity))
		         FILTER (WHERE di.docket_item_id IS NOT NULL), '[]'::json)
		  FROM restaurant_dockets d
		  LEFT JOIN docket_items di ON di.docket_id = d.docket_id
		  LEFT JOIN menu_items m ON m.menu_item_id = di.menu_item_id
		 WHERE d.kitchen_status <> 'Served'
		 GROUP BY d.docket_id, d.room_number, d.table_number, d.sale_date, d.kitchen_status, d.notes
		 ORDER BY d.sale_date ASC
		 LIMIT 100
	`)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query kitchen dockets: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ticket struct {
		ID            string          `json:"id"`
		RoomNumber    string          `json:"roomNo"`
		TableNumber   string          `json:"tableNo"`
		CreatedAtTime time.Time       `json:"-"`
		CreatedAt     string          `json:"createdAt"`
		KitchenStatus string          `json:"kitchenStatus"`
		Notes         string          `json:"notes"`
		Items         json.RawMessage `json:"items"`
	}
	tickets := make([]ticket, 0)
	for rows.Next() {
		var item ticket
		if err := rows.Scan(&item.ID, &item.RoomNumber, &item.TableNumber, &item.CreatedAtTime, &item.KitchenStatus, &item.Notes, &item.Items); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to read kitchen dockets: %v"}`, err), http.StatusInternalServerError)
			return
		}
		item.CreatedAt = item.CreatedAtTime.Format(time.RFC3339)
		tickets = append(tickets, item)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to read kitchen dockets: %v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tickets)
}

// UpdateDocketStatus advances a kitchen ticket through its service workflow.
func (h *KitchenHandler) UpdateDocketStatus(w http.ResponseWriter, r *http.Request) {
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
		Status   string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	req.DocketID = strings.TrimSpace(req.DocketID)
	if !isUUID(req.DocketID) || req.Status != "Cooking" && req.Status != "Ready" && req.Status != "Served" {
		http.Error(w, `{"error":"invalid docket status update"}`, http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(`
		UPDATE restaurant_dockets
		   SET kitchen_status = $2
		 WHERE docket_id = $1
		   AND ((kitchen_status = 'Pending' AND $2 = 'Cooking')
		     OR (kitchen_status = 'Cooking' AND $2 = 'Ready')
		     OR (kitchen_status = 'Ready' AND $2 = 'Served'))
	`, req.DocketID, req.Status)
	if err != nil {
		http.Error(w, `{"error":"failed to update kitchen docket"}`, http.StatusInternalServerError)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"failed to verify kitchen docket update"}`, http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		var currentStatus string
		err := h.db.QueryRow(
			`SELECT kitchen_status FROM restaurant_dockets WHERE docket_id = $1`,
			req.DocketID,
		).Scan(&currentStatus)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"docket not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, `{"error":"failed to verify kitchen docket"}`, http.StatusInternalServerError)
			return
		}
		http.Error(w, `{"error":"docket status must advance one step at a time"}`, http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GetStock returns raw ingredients from kitchen_stock.
func (h *KitchenHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	rows, err := h.db.Query(
		`SELECT kitchen_stock_id::text, ingredient_name, opening_stock, add_stock, total_stock, sold_usage, closing_balance
		   FROM kitchen_stock
		  WHERE audit_date = $1
		  ORDER BY ingredient_name ASC`,
		date,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query kitchen stock: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type StockItemDTO struct {
		ID             string  `json:"id"`
		Name           string  `json:"name"`
		Unit           string  `json:"unit"`
		OpeningStock   float64 `json:"openingStock"`
		AddStock       float64 `json:"addStock"`
		TotalStock     float64 `json:"totalStock"`
		SoldUsage      float64 `json:"soldUsage"`
		ClosingBalance float64 `json:"closingBalance"`
		MinAlert       float64 `json:"minAlert"`
	}

	var items []StockItemDTO
	for rows.Next() {
		var it StockItemDTO
		if err := rows.Scan(&it.ID, &it.Name, &it.OpeningStock, &it.AddStock, &it.TotalStock, &it.SoldUsage, &it.ClosingBalance); err == nil {
			it.Unit = "kg"
			it.MinAlert = 10
			items = append(items, it)
		}
	}

	json.NewEncoder(w).Encode(items)
}

// SaveStock updates or creates kitchen_stock ledger records.
func (h *KitchenHandler) SaveStock(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		AuditDate string `json:"audit_date"`
		Items     []struct {
			ID             string  `json:"id"`
			Name           string  `json:"name"`
			Unit           string  `json:"unit"`
			OpeningStock   float64 `json:"openingStock"`
			AddStock       float64 `json:"addStock"`
			SoldUsage      float64 `json:"soldUsage"`
			ClosingBalance float64 `json:"closingBalance"`
			MinAlert       float64 `json:"minAlert"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	auditDate := req.AuditDate
	if auditDate == "" {
		auditDate = time.Now().Format("2006-01-02")
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"transaction failed"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range req.Items {
		var exists int
		_ = tx.QueryRow(
			`SELECT COUNT(*) FROM kitchen_stock WHERE ingredient_name = $1 AND audit_date = $2`,
			it.Name, auditDate,
		).Scan(&exists)

		if exists > 0 {
			_, err = tx.Exec(
				`UPDATE kitchen_stock
				    SET opening_stock = $1, add_stock = $2, sold_usage = $3, closing_balance = $4
				  WHERE ingredient_name = $5 AND audit_date = $6`,
				it.OpeningStock, it.AddStock, it.SoldUsage, it.ClosingBalance, it.Name, auditDate,
			)
		} else {
			_, err = tx.Exec(
				`INSERT INTO kitchen_stock (ingredient_name, audit_date, opening_stock, add_stock, sold_usage, closing_balance)
				 VALUES ($1, $2, $3, $4, $5, $6)`,
				it.Name, auditDate, it.OpeningStock, it.AddStock, it.SoldUsage, it.ClosingBalance,
			)
		}
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to save %s: %v"}`, it.Name, err), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// GetRecipes returns dishes and their Bill of Materials (BOM) ingredients breakdown.
func (h *KitchenHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	query := `
		SELECT m.menu_item_id::text, m.item_name, m.category::text, m.price, m.cost_price,
		       COALESCE(json_agg(
		         json_build_object(
		           'name', ks.ingredient_name,
		           'qty', mi.quantity_required,
		           'unit', 'kg',
		           'approxCostMWK', ROUND(mi.quantity_required * 3000)
		         )
		       ) FILTER (WHERE ks.kitchen_stock_id IS NOT NULL), '[]'::json) AS ingredients
		  FROM menu_items m
		  LEFT JOIN menu_item_ingredients mi ON m.menu_item_id = mi.menu_item_id
		  LEFT JOIN kitchen_stock ks ON mi.kitchen_stock_id = ks.kitchen_stock_id
		 GROUP BY m.menu_item_id, m.item_name, m.category, m.price, m.cost_price
		 ORDER BY m.item_name
	`
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query recipes: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type RecipeDTO struct {
		DishID       string          `json:"dishId"`
		DishName     string          `json:"dishName"`
		Category     string          `json:"category"`
		SellingPrice float64         `json:"sellingPrice"`
		CostPrice    float64         `json:"costPrice"`
		Ingredients  json.RawMessage `json:"ingredients"`
	}

	var list []RecipeDTO
	for rows.Next() {
		var rec RecipeDTO
		if err := rows.Scan(&rec.DishID, &rec.DishName, &rec.Category, &rec.SellingPrice, &rec.CostPrice, &rec.Ingredients); err == nil {
			list = append(list, rec)
		}
	}

	json.NewEncoder(w).Encode(list)
}

// SubmitReport handles kitchen shift reports with JSONB stock summary.
func (h *KitchenHandler) SubmitReport(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			ReportDate   string          `json:"report_date"`
			StockSummary json.RawMessage `json:"stock_summary"`
			ChefNotes    string          `json:"chef_notes"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		date := req.ReportDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		var uID string
		_ = h.db.QueryRow("SELECT user_id::text FROM user_accounts WHERE role = 'CHEF' LIMIT 1").Scan(&uID)
		if uID == "" {
			_ = h.db.QueryRow("SELECT user_id::text FROM user_accounts LIMIT 1").Scan(&uID)
		}

		summaryJSON := string(req.StockSummary)
		if summaryJSON == "" {
			summaryJSON = "{}"
		}

		_, err := h.db.Exec(
			`INSERT INTO kitchen_reports (user_id, report_date, stock_summary_json, status)
			 VALUES ($1::uuid, $2, $3::jsonb, 'Submitted')`,
			uID, date, summaryJSON,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to submit kitchen report: %v"}`, err), http.StatusInternalServerError)
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

	var status, summary string
	err := h.db.QueryRow(
		`SELECT status::text, COALESCE(stock_summary_json::text, '{}')
		   FROM kitchen_reports WHERE report_date = $1 LIMIT 1`,
		date,
	).Scan(&status, &summary)
	if err != nil {
		status = "Submitted"
		summary = "{}"
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        status,
		"stock_summary": json.RawMessage(summary),
	})
}
