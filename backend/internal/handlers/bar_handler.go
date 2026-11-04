// backend/internal/handlers/bar_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BarHandler struct {
	db *sql.DB
}

func NewBarHandler(db *sql.DB) *BarHandler {
	return &BarHandler{db: db}
}

// GetItems returns all bar items with stock levels.
func (h *BarHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	rows, err := h.db.Query(
		`SELECT bar_item_id::text, item_name, unit_price, cost_price, stock_quantity
		   FROM bar_items
		  ORDER BY item_name ASC`,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query bar items: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BarItemDTO struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		UnitPrice float64 `json:"unitPrice"`
		CostPrice float64 `json:"costPrice"`
		Stock     int     `json:"stock"`
	}

	items := make([]BarItemDTO, 0)
	for rows.Next() {
		var it BarItemDTO
		if err := rows.Scan(&it.ID, &it.Name, &it.UnitPrice, &it.CostPrice, &it.Stock); err != nil {
			http.Error(w, `{"error":"failed to read bar items"}`, http.StatusInternalServerError)
			return
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error":"failed to read bar items"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

// RecordSale records a beverage sale and dynamically decrements product stock.
func (h *BarHandler) RecordSale(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Item string `json:"item"`
		Qty  int    `json:"qty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Qty <= 0 || req.Item == "" {
		http.Error(w, `{"error":"invalid sale payload"}`, http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"transaction failed"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Use the authoritative database price and stock for the selected item.
	var itemID, itemName string
	var currentStock int
	var unitPrice float64
	err = tx.QueryRow(
		`SELECT bar_item_id::text, item_name, stock_quantity, unit_price
		   FROM bar_items
		  WHERE item_name = $1
		  FOR UPDATE`,
		req.Item,
	).Scan(&itemID, &itemName, &currentStock, &unitPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"bar item not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"failed to load bar item"}`, http.StatusInternalServerError)
		}
		return
	}
	if currentStock < req.Qty {
		http.Error(w, `{"error":"insufficient stock"}`, http.StatusConflict)
		return
	}
	if _, err := tx.Exec(
		`UPDATE bar_items SET stock_quantity = stock_quantity - $1 WHERE bar_item_id = $2::uuid`,
		req.Qty, itemID,
	); err != nil {
		http.Error(w, `{"error":"failed to update bar stock"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Drink sale recorded and stock updated",
		"item":      itemName,
		"qty":       req.Qty,
		"unitPrice": unitPrice,
		"stock":     currentStock - req.Qty,
	})
}

// GetLedger returns opening, added, sold, and closing balances from bar_ledgers.
func (h *BarHandler) GetLedger(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	query := `
		SELECT l.bar_ledger_id::text, i.bar_item_id::text, i.item_name, i.unit_price, i.cost_price,
		       i.stock_quantity, l.opening_stock, l.add_stock, l.sold_quantity, l.closing_stock, l.sales_revenue
		  FROM bar_items i
		  LEFT JOIN bar_ledgers l ON i.bar_item_id = l.bar_item_id AND l.audit_date = $1
		 ORDER BY i.item_name ASC
	`
	rows, err := h.db.Query(query, date)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"failed to query bar ledger: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type LedgerDTO struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		UnitPrice    float64 `json:"unitPrice"`
		CostPrice    float64 `json:"costPrice"`
		Stock        int     `json:"stock"`
		OpeningStock float64 `json:"openingStock"`
		AddStock     float64 `json:"addStock"`
		SoldQty      float64 `json:"soldQty"`
		ClosingStock float64 `json:"closingStock"`
		SalesRevenue float64 `json:"salesRevenue"`
	}

	list := make([]LedgerDTO, 0)
	for rows.Next() {
		var d LedgerDTO
		var ledgerID sql.NullString
		var op, add, sold, closing, rev sql.NullFloat64

		if err := rows.Scan(&ledgerID, &d.ID, &d.Name, &d.UnitPrice, &d.CostPrice, &d.Stock, &op, &add, &sold, &closing, &rev); err != nil {
			http.Error(w, `{"error":"failed to read bar ledger"}`, http.StatusInternalServerError)
			return
		}
		if op.Valid {
			d.OpeningStock = op.Float64
			d.AddStock = add.Float64
			d.SoldQty = sold.Float64
			d.ClosingStock = closing.Float64
			d.SalesRevenue = rev.Float64
		} else {
			d.OpeningStock = float64(d.Stock)
			d.ClosingStock = float64(d.Stock)
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error":"failed to read bar ledger"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

// SaveLedger saves daily bar audit records ensuring constraint check_bar_stock_balance.
func (h *BarHandler) SaveLedger(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		AuditDate string `json:"audit_date"`
		Items     []struct {
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			UnitPrice    float64 `json:"unitPrice"`
			CostPrice    float64 `json:"costPrice"`
			OpeningStock float64 `json:"openingStock"`
			AddStock     float64 `json:"addStock"`
			SoldQty      float64 `json:"soldQty"`
			ClosingStock float64 `json:"closingStock"`
			SalesRevenue float64 `json:"salesRevenue"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid ledger body"}`, http.StatusBadRequest)
		return
	}

	date := req.AuditDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"error":"transaction failed"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, it := range req.Items {
		var itemID string
		_ = tx.QueryRow("SELECT bar_item_id::text FROM bar_items WHERE item_name = $1 LIMIT 1", it.Name).Scan(&itemID)
		if itemID == "" {
			_ = tx.QueryRow("INSERT INTO bar_items (item_name, unit_price, cost_price, stock_quantity) VALUES ($1, $2, $3, $4) RETURNING bar_item_id::text",
				it.Name, it.UnitPrice, it.CostPrice, int(it.ClosingStock)).Scan(&itemID)
		}

		if itemID != "" {
			// Enforce closing_stock <= (opening_stock + add_stock)
			closing := it.ClosingStock
			if closing > (it.OpeningStock + it.AddStock) {
				closing = it.OpeningStock + it.AddStock
			}

			var exists int
			_ = tx.QueryRow("SELECT COUNT(*) FROM bar_ledgers WHERE bar_item_id = $1::uuid AND audit_date = $2", itemID, date).Scan(&exists)

			if exists > 0 {
				_, err = tx.Exec(
					`UPDATE bar_ledgers
					    SET opening_stock = $1, add_stock = $2, sold_quantity = $3, closing_stock = $4, sales_revenue = $5
					  WHERE bar_item_id = $6::uuid AND audit_date = $7`,
					it.OpeningStock, it.AddStock, it.SoldQty, closing, it.SalesRevenue, itemID, date,
				)
			} else {
				_, err = tx.Exec(
					`INSERT INTO bar_ledgers (bar_item_id, audit_date, opening_stock, add_stock, sold_quantity, closing_stock, sales_revenue)
					 VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)`,
					itemID, date, it.OpeningStock, it.AddStock, it.SoldQty, closing, it.SalesRevenue,
				)
			}
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"failed to save bar ledger: %v"}`, err), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// Reports handles GET/POST for bar financial sales reports.
func (h *BarHandler) Reports(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if h.db == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			ReportDate    string  `json:"report_date"`
			TotalBarSales float64 `json:"total_bar_sales"`
			TotalCogs     float64 `json:"total_cogs"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		date := req.ReportDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		var uID string
		_ = h.db.QueryRow("SELECT user_id::text FROM user_accounts WHERE role = 'BAR_ATTENDANT' LIMIT 1").Scan(&uID)
		if uID == "" {
			_ = h.db.QueryRow("SELECT user_id::text FROM user_accounts LIMIT 1").Scan(&uID)
		}

		// Note: gross_margin is GENERATED ALWAYS AS (total_bar_sales - total_cogs) STORED, so do not insert it!
		_, err := h.db.Exec(
			`INSERT INTO bar_reports (user_id, report_date, total_bar_sales, total_cogs, status)
			 VALUES ($1::uuid, $2, $3, $4, 'Submitted')`,
			uID, date, req.TotalBarSales, req.TotalCogs,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"failed to save bar report: %v"}`, err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		return
	}

	// GET report totals from the day's persisted inventory ledger.
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var sales, cogs, gross float64
	rows, err := h.db.Query(
		`SELECT i.item_name, COALESCE(l.sold_quantity, 0), COALESCE(l.sales_revenue, 0),
		        COALESCE(l.sold_quantity * i.cost_price, 0)
		   FROM bar_items i
		   LEFT JOIN bar_ledgers l ON l.bar_item_id = i.bar_item_id AND l.audit_date = $1
		  ORDER BY i.item_name`,
		date,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to query bar ledger"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type LedgerEntry struct {
		Name    string  `json:"name"`
		Sold    float64 `json:"sold"`
		Revenue float64 `json:"revenue"`
		Cogs    float64 `json:"cogs"`
	}
	ledger := make([]LedgerEntry, 0)
	for rows.Next() {
		var entry LedgerEntry
		if err := rows.Scan(&entry.Name, &entry.Sold, &entry.Revenue, &entry.Cogs); err != nil {
			http.Error(w, `{"error":"failed to read bar ledger"}`, http.StatusInternalServerError)
			return
		}
		sales += entry.Revenue
		cogs += entry.Cogs
		ledger = append(ledger, entry)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error":"failed to read bar ledger"}`, http.StatusInternalServerError)
		return
	}
	gross = sales - cogs

	status := "Not submitted"
	err = h.db.QueryRow(
		`SELECT status::text
		   FROM bar_reports
		  WHERE report_date = $1
		  ORDER BY created_at DESC
		  LIMIT 1`,
		date,
	).Scan(&status)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, `{"error":"failed to query bar report status"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_bar_sales": sales,
		"total_cogs":      cogs,
		"gross_margin":    gross,
		"status":          status,
		"ledger":          ledger,
	})
}
