// backend/internal/services/auth_service.go
package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// AuthService handles user authentication with optional development-only demo accounts.
type AuthService struct {
	jwtSecret      string
	db             *sql.DB
	allowDemoUsers bool
}

// NewAuthService constructs an AuthService with the given JWT signing secret.
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{jwtSecret: jwtSecret}
}

// SetDB injects a live database connection (called after PostgreSQL connects).
func (s *AuthService) SetDB(db *sql.DB) {
	s.db = db
}

func (s *AuthService) SetDemoUsersEnabled(enabled bool) {
	s.allowDemoUsers = enabled
}

// LoginRequest is the JSON body sent by the frontend.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserInfo is the user data embedded in the login response.
type UserInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// jwtClaims are the custom claims stored in the JWT.
type jwtClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func normalizeRole(role string) string {
	role = strings.ToUpper(strings.TrimSpace(role))
	switch role {
	case "FRONTDESK", "FRONT_DESK":
		return "FRONT_DESK"
	case "CHEF", "KITCHEN":
		return "CHEF"
	case "BAR", "BAR_ATTENDANT":
		return "BAR_ATTENDANT"
	case "WAITER", "RESTAURANT":
		return "WAITER"
	case "MANAGER":
		return "MANAGER"
	case "GUEST":
		return "GUEST"
	case "ADMIN":
		return "ADMIN"
	default:
		return role
	}
}

// Login authenticates a user against PostgreSQL. Demo fallbacks are disabled.
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	normRole := normalizeRole(req.Role)

	// If logging in as GUEST, check the dedicated guest_accounts table
	if normRole == "GUEST" {
		if s.db == nil {
			return nil, fmt.Errorf("database connection unavailable")
		}
		var guestID, username, email, passwordHash, fullName string
		row := s.db.QueryRow(
			`SELECT ga.guest_id::text, ga.username, ga.email, ga.password_hash,
			        (g.first_name || ' ' || g.last_name) AS full_name
			   FROM guest_accounts ga
			   JOIN guests g ON ga.guest_id = g.guest_id
			  WHERE (LOWER(ga.username) = LOWER($1) OR LOWER(ga.email) = LOWER($1))
			    AND ga.is_active = true`,
			strings.TrimSpace(req.Username),
		)
		err := row.Scan(&guestID, &username, &email, &passwordHash, &fullName)
		if err != nil {
			// Fallback check for legacy guest in user_accounts
			var uID, uName, uMail, uHash, uRole string
			uErr := s.db.QueryRow(
				`SELECT user_id::text, username, email, password_hash, role::text
				   FROM user_accounts
				  WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1))
				    AND role = 'GUEST' AND is_active = true`,
				strings.TrimSpace(req.Username),
			).Scan(&uID, &uName, &uMail, &uHash, &uRole)
			if uErr == nil {
				guestID = uID
				username = uName
				email = uMail
				passwordHash = uHash
				fullName = uName
				err = nil
			}
		}

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("invalid credentials")
			}
			return nil, fmt.Errorf("authentication query failed: %w", err)
		}

		if bcryptErr := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); bcryptErr != nil {
			return nil, fmt.Errorf("invalid credentials")
		}

		return s.buildResponse(guestID, username, fullName, "GUEST")
	}

	// Staff and Admin login path (strictly from user_accounts table)
	if s.db != nil {
		var userID, username, email, passwordHash, dbRole string
		row := s.db.QueryRow(
			`SELECT user_id::text, username, email, password_hash, role::text
			   FROM user_accounts
			  WHERE (LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($1))
			    AND role != 'GUEST' AND is_active = true`,
			strings.TrimSpace(req.Username),
		)
		err := row.Scan(&userID, &username, &email, &passwordHash, &dbRole)
		if err == nil {
			if bcryptErr := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); bcryptErr != nil {
				return nil, fmt.Errorf("invalid credentials")
			}
			normalizedRole := normalizeRole(dbRole)
			if req.Role != "" && normalizeRole(req.Role) != normalizedRole {
				return nil, fmt.Errorf("role mismatch")
			}
			return s.buildResponse(userID, username, username, normalizedRole)
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("authentication database query failed: %w", err)
		}
	}

	return nil, fmt.Errorf("invalid credentials")
}

// RegisterRequest is the JSON body for guest self-registration.
type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Nationality string `json:"nationality"`
	City        string `json:"city"`
	PassportNo  string `json:"passport_no"`
	PostalCode  string `json:"postal_code"`
	Address     string `json:"address"`
}

// Register creates a new guest account in guests and guest_accounts tables.
func (s *AuthService) Register(req RegisterRequest) (*LoginResponse, error) {
	if s.db == nil {
		return nil, fmt.Errorf("registration requires database connection")
	}

	// Validate required fields
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		return nil, fmt.Errorf("username, email, password, first name and last name are required")
	}

	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		phone = "N/A"
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert into guests table with all comprehensive details
	var guestID string
	err = tx.QueryRow(
		`INSERT INTO guests (first_name, last_name, contact_no, email, nationality, city, passport_no, postal_code, address)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING guest_id::text`,
		strings.TrimSpace(req.FirstName),
		strings.TrimSpace(req.LastName),
		phone,
		strings.TrimSpace(req.Email),
		strings.TrimSpace(req.Nationality),
		strings.TrimSpace(req.City),
		strings.TrimSpace(req.PassportNo),
		strings.TrimSpace(req.PostalCode),
		strings.TrimSpace(req.Address),
	).Scan(&guestID)
	if err != nil {
		return nil, fmt.Errorf("failed to create guest profile: %w", err)
	}

	// 2. Insert into guest_accounts table (separate from staff user_accounts)
	var guestAccountID string
	err = tx.QueryRow(
		`INSERT INTO guest_accounts (guest_id, username, email, password_hash)
		 VALUES ($1::uuid, $2, $3, $4)
		 RETURNING guest_account_id::text`,
		guestID,
		strings.TrimSpace(req.Username),
		strings.TrimSpace(req.Email),
		string(hash),
	).Scan(&guestAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("username or email is already registered")
		}
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	fullName := strings.TrimSpace(req.FirstName) + " " + strings.TrimSpace(req.LastName)
	return s.buildResponse(guestID, strings.TrimSpace(req.Username), fullName, "GUEST")
}

// buildResponse signs a JWT and assembles a LoginResponse.
func (s *AuthService) buildResponse(userID, username, fullName, role string) (*LoginResponse, error) {
	claims := jwtClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "la-luna-hms",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("token signing failed: %w", err)
	}
	return &LoginResponse{
		Token: signed,
		User: UserInfo{
			UserID:   userID,
			Username: username,
			FullName: fullName,
			Role:     role,
		},
	}, nil
}
