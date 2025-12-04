package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GoogleLoginURL(state string) string
	GoogleCallback(ctx context.Context, code string) (string, *User, error) // returns JWT
	ParseToken(tokenString string) (*jwt.RegisteredClaims, string /* userID */, error)
	Register(ctx context.Context, req RegisterRequest) (string, *User, error)
	Login(ctx context.Context, req LoginRequest) (string, *User, error)
}

type service struct {
	repo   Repository
	google *GoogleOAuth
	jwtKey []byte
}

func NewService(repo Repository, google *GoogleOAuth, jwtSecret string) Service {
	return &service{
		repo:   repo,
		google: google,
		jwtKey: []byte(jwtSecret),
	}
}

// user-visible URL
func (s *service) GoogleLoginURL(state string) string {
	return s.google.AuthCodeURL(state)
}

func (s *service) GoogleCallback(ctx context.Context, code string) (string, *User, error) {
	info, err := s.google.ExchangeAndFetchUser(ctx, code)
	if err != nil {
		return "", nil, err
	}
	// find or create user
	u, err := s.repo.FindByProvider(ctx, "google", info.Sub)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		u = &User{
			Email:      info.Email,
			Name:       info.Name,
			AvatarURL:  info.Picture,
			Provider:   "google",
			ProviderID: info.Sub,
		}
		if err := s.repo.Create(ctx, u); err != nil {
			return "", nil, err
		}
	}

	// create JWT
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   u.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", nil, err
	}
	return signed, u, nil
}

func (s *service) ParseToken(tokenString string) (*jwt.RegisteredClaims, string, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, "", err
	}
	claims, ok := tok.Claims.(*jwt.RegisteredClaims)
	if !ok || !tok.Valid {
		return nil, "", err
	}
	return claims, claims.Subject, nil
}

// --- Implementations ---

func (s *service) Register(ctx context.Context, req RegisterRequest) (string, *User, error) {
	// 1. Check if user exists
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, fmt.Errorf("email already in use")
	}

	// 2. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	// 3. Create User
	u := &User{
		Email:        req.Email,
		PasswordHash: string(hashed), // Store Hash
		Name:         req.Name,
		Provider:     "local",
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return "", nil, err
	}

	// 4. Generate Token
	return s.generateJWT(u)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, *User, error) {
	// 1. Find User
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// 2. Check Password
	if u.Provider == "google" {
		return "", nil, fmt.Errorf("please login with google")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// 3. Generate Token
	return s.generateJWT(u)
}

// Helper to avoid code duplication
func (s *service) generateJWT(u *User) (string, *User, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   u.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtKey)
	return signed, u, err
}
