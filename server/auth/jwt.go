package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
	"x-tract/models"
)

type contextKey string

const claimsContextKey contextKey = "jwt_claims"

// Claims is the JWT payload used by the server.
type Claims struct {
	Subject string      `json:"sub"`
	Email   string      `json:"email"`
	Role    models.Role `json:"role"`
	Issued  int64       `json:"iat"`
	Expires int64       `json:"exp"`
}

// Manager issues and verifies HMAC-signed JWT tokens.
type Manager struct {
	secret     []byte
	ttl        time.Duration
	now        func() time.Time
	issuerName string
}

// Config configures the JWT manager.
type Config struct {
	Secret     string
	TTL        time.Duration
	IssuerName string
}

// NewManager creates a new JWT manager.
func NewManager(cfg Config) *Manager {
	ttl := cfg.TTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &Manager{
		secret:     []byte(cfg.Secret),
		ttl:        ttl,
		now:        time.Now,
		issuerName: cfg.IssuerName,
	}
}

// Issue creates a signed JWT for the provided user.
func (m *Manager) Issue(user *models.User) (string, time.Time, error) {
	now := m.now().UTC()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		Subject: user.ID.String(),
		Email:   user.Email,
		Role:    user.Role,
		Issued:  now.Unix(),
		Expires: expiresAt.Unix(),
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", time.Time{}, err
	}

	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	header := encodeJWTPart(headerJSON)
	payload := encodeJWTPart(payloadJSON)
	signingInput := header + "." + payload
	signature := m.sign(signingInput)

	return signingInput + "." + signature, expiresAt, nil
}

// Verify parses and validates a JWT.
func (m *Manager) Verify(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := m.sign(signingInput)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	if m.now().UTC().Unix() >= claims.Expires {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// ClaimsFromContext returns previously-validated JWT claims from a request context.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

// ContextWithClaims stores claims on a request context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func (m *Manager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(value))
	return encodeJWTPart(mac.Sum(nil))
}

func encodeJWTPart(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}
