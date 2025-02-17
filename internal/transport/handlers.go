package transport

import (
	"Backend-trainee-assignment/internal/auth/autherrors"
	transportModel "Backend-trainee-assignment/internal/model/transport"
	"Backend-trainee-assignment/internal/service/serverrors"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type shopService interface {
	GetInfo(ctx context.Context) (*transportModel.InfoResponse, error)
	SendCoins(ctx context.Context, target string, amount int) error
	BuyItem(ctx context.Context, id string) error
}

type authService interface {
	GetOrCreateTokenByCredentials(ctx context.Context, username, providedPassword string) (string, error)
	ParseToken(ctx context.Context, token string) (string, error)
}

type Handler struct {
	shopService shopService
	authService authService
}

func NewHandler(shop shopService, auth authService) *Handler {
	return &Handler{
		shopService: shop,
		authService: auth,
	}
}

func (h *Handler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Use(loggingMiddleware)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(router chi.Router) {

		router.Post("/auth", h.PostAuthEndpoint)

		router.Group(func(routerWithAuth chi.Router) {
			routerWithAuth.Use(h.authMiddleware)

			routerWithAuth.Get("/info", h.GetInfo)
			routerWithAuth.Post("/sendCoin", h.PostSendCoin)
			routerWithAuth.Get("/buy/{item}", h.GetBuyItem)

		})
	})

	return router
}

func (h *Handler) PostAuthEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var authRequest transportModel.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		slog.Error("Error decoding auth request:", err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Failed to decode auth request",
		})
		return
	}

	err := validate.Struct(authRequest)
	if err != nil {
		slog.Error("Error validating auth request:", err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Failed to validate auth request",
		})
		return
	}

	token, err := h.authService.GetOrCreateTokenByCredentials(r.Context(), authRequest.Username, authRequest.Password)
	if err != nil {
		if errors.Is(err, autherrors.ErrInvalidPassword) {

			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
				Errors: "Invalid password",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Failed to create token",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(transportModel.AuthResponse{
		Token: token,
	})

	return
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {

	infoResponse, err := h.shopService.GetInfo(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Internal server error with storage",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(infoResponse)
	return
}

func (h *Handler) PostSendCoin(w http.ResponseWriter, r *http.Request) {

	var coinRequest transportModel.SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&coinRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Failed to decode request body",
		})
		return
	}

	err := validate.Struct(coinRequest)
	if err != nil {
		slog.Error("Error validating coin request:", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Failed to validate coin request",
		})
		return
	}

	err = h.shopService.SendCoins(r.Context(), coinRequest.ToUser, coinRequest.Amount)
	if err != nil {
		if errors.Is(serverrors.ErrBalanceNotEnough, err) {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
				Errors: "User balance not enough",
			})
			return
		}

		if errors.Is(serverrors.ErrInvalidTarget, err) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
				Errors: "Invalid or cant found target",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(transportModel.ErrorResponse{
			Errors: "Internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBuyItem(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
