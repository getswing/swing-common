package usecase

import (
	"net/http"

	"getswing.app/player-service/internal/app/models"
	"getswing.app/player-service/internal/app/repository"
	"getswing.app/player-service/internal/infrastructure/config"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	cfg        config.Config
	playerRepo repository.PlayerRepository
}

func NewAuthHandler(cfg config.Config, playerRepo repository.PlayerRepository) *AuthHandler {
	return &AuthHandler{cfg: cfg, playerRepo: playerRepo}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email and password required"})
	}
	// hash password
	// hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to hash password"})
	// }
	u := &models.Player{
		Email: req.Email,
		// PasswordHash: string(hash),
	}
	if err := h.playerRepo.Create(c.Request().Context(), u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"id": u.ID, "email": u.Email})
}

type loginRequest struct {
	PhoneNumber      string `json:"phone_number"`
	PhoneCountryCode string `json:"phone_country_code"`
	Password         string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil || req.PhoneNumber == "" || req.PhoneCountryCode == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "phone number is required"})
	}
	u, err := h.playerRepo.FindByPhone(c.Request().Context(), req.PhoneNumber, req.PhoneCountryCode)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}
	// if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
	// 	return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	// }

	// now := time.Now()
	// claims := jwt.RegisteredClaims{
	// 	Subject:   string(rune(u.ID)),
	// 	Issuer:    h.cfg.JWTIssuer,
	// 	ExpiresAt: jwt.NewNumericDate(now.Add(h.cfg.JWTAccessTTL)),
	// 	IssuedAt:  jwt.NewNumericDate(now),
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ss, err := token.SignedString([]byte(h.cfg.JWTSecret))
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to sign token"})
	// }

	return c.JSON(http.StatusOK, echo.Map{
		// "access_token": ss,
		"token_type": "Bearer",
		"expires_in": int(h.cfg.JWTAccessTTL.Seconds()),
		"user":       echo.Map{"id": u.ID, "email": u.Email},
	})
}
