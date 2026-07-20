package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
	"github.com/naputt/cinema-ticket-backend/internal/middleware"
	"github.com/naputt/cinema-ticket-backend/pkg/config"
)

type AuthUsecase struct {
	userRepo domain.UserRepository
	cfg      *config.Config
	oauthCfg *oauth2.Config
}

func NewAuthUsecase(userRepo domain.UserRepository, cfg *config.Config) *AuthUsecase {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	return &AuthUsecase{
		userRepo: userRepo,
		cfg:      cfg,
		oauthCfg: oauthCfg,
	}
}

func (uc *AuthUsecase) GetAuthURL() string {
	return uc.oauthCfg.AuthCodeURL(uuid.New().String(), oauth2.AccessTypeOffline)
}

func (uc *AuthUsecase) HandleCallback(ctx context.Context, code string) (string, string, error) {
	exchangeCtx := context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: 10 * time.Second})
	token, err := uc.oauthCfg.Exchange(exchangeCtx, code)
	if err != nil {
		slog.ErrorContext(ctx, "google oauth exchange failed", "error", err)
		return "", "", fmt.Errorf("token exchange failed: %w", err)
	}

	idTokenRaw, ok := token.Extra("id_token").(string)
	if !ok {
		slog.ErrorContext(ctx, "id_token missing from token response")
		return "", "", errors.New("id_token not found in token response")
	}

	payload, err := idtoken.Validate(ctx, idTokenRaw, uc.cfg.Google.ClientID)
	if err != nil {
		slog.ErrorContext(ctx, "id token validation failed", "error", err)
		return "", "", fmt.Errorf("id token validation failed: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		slog.ErrorContext(ctx, "email claim missing or invalid", "claims", payload.Claims)
		return "", "", errors.New("email not found in id token")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		slog.ErrorContext(ctx, "name claim missing or invalid", "email", email)
		return "", "", errors.New("name not found in id token")
	}

	avatar := ""
	if av, ok := payload.Claims["picture"]; ok {
		avatar, _ = av.(string)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		slog.ErrorContext(ctx, "find user by email failed", "email", email, "error", err)
		return "", "", fmt.Errorf("find user failed: %w", err)
	}

	if user == nil {
		isAdmin := false
		for _, ae := range uc.cfg.App.AdminEmails {
			if ae == email {
				isAdmin = true
				break
			}
		}

		role := domain.RoleUser
		if isAdmin {
			role = domain.RoleAdmin
		}

		user = &domain.User{
			ID:        uuid.New().String(),
			Email:     email,
			Name:      name,
			AvatarURL: avatar,
			Role:      role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			slog.ErrorContext(ctx, "create user failed", "email", email, "error", err)
			return "", "", fmt.Errorf("create user failed: %w", err)
		}

		slog.InfoContext(ctx, "new user created via google oauth", "email", email, "role", role)
	} else {
		user.Name = name
		user.AvatarURL = avatar
		user.UpdatedAt = time.Now()
		if err := uc.userRepo.Update(ctx, user); err != nil {
			slog.ErrorContext(ctx, "update user failed", "user_id", user.ID, "error", err)
			return "", "", fmt.Errorf("update user failed: %w", err)
		}

		slog.InfoContext(ctx, "existing user logged in via google oauth", "email", email)
	}

	accessToken, refreshToken, err := uc.generateTokens(user)
	if err != nil {
		slog.ErrorContext(ctx, "token generation failed", "user_id", user.ID, "error", err)
		return "", "", fmt.Errorf("token generation failed: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (string, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return "", err
	}

	accessToken, _, err := uc.generateTokens(user)
	return accessToken, err
}

func (uc *AuthUsecase) generateTokens(user *domain.User) (string, string, error) {
	now := time.Now()

	accessClaims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.cfg.JWT.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(uc.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := middleware.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.cfg.JWT.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(uc.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

func (uc *AuthUsecase) GetAccessExpiry() time.Duration {
	return uc.cfg.JWT.AccessExpiry
}

func (uc *AuthUsecase) IsProduction() bool {
	return uc.cfg.App.Production
}
