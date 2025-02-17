package transport

import (
	"Backend-trainee-assignment/internal/auth/autherrors"
	transportModel "Backend-trainee-assignment/internal/model/transport"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log"
	"log/slog"
	"net/http"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type shopService interface {
	GetInfo(ctx context.Context) (*transportModel.InfoResponse, error)
	SendCoins(ctx context.Context, target string) error
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

		resp := transportModel.ErrorResponse{
			Errors: "Failed to decode auth request",
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	err := validate.Struct(authRequest)
	if err != nil {
		slog.Error("Error validating auth request:", err)

		resp := transportModel.ErrorResponse{
			Errors: "Failed to validate auth request",
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	token, err := h.authService.GetOrCreateTokenByCredentials(r.Context(), authRequest.Username, authRequest.Password)
	if err != nil {
		if errors.Is(err, autherrors.ErrInvalidPassword) {

			resp := transportModel.ErrorResponse{
				Errors: "Invalid password",
			}
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		resp := transportModel.ErrorResponse{
			Errors: "Failed to create token",
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := transportModel.AuthResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
	return
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {

	infoResponse, err := h.shopService.GetInfo(r.Context())
	if err != nil {
		resp := transportModel.ErrorResponse{
			Errors: "Internal server error with storage",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(infoResponse)
	return
}

func (h *Handler) PostSendCoin(w http.ResponseWriter, r *http.Request) {
	var req transportModel.SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, transportModel.ErrorResponse{Errors: "invalid request body"})
		return
	}

	if req.Amount <= 0 || req.ToUser == "" {
		writeJSON(w, http.StatusBadRequest, transportModel.ErrorResponse{Errors: "invalid parameters"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBuyItem(w http.ResponseWriter, r *http.Request) {
	item := chi.URLParam(r, "item")
	if item == "" {
		writeJSON(w, http.StatusBadRequest, transportModel.ErrorResponse{Errors: "item parameter is required"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostAuth(w http.ResponseWriter, r *http.Request) {
	var req transportModel.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, transportModel.ErrorResponse{Errors: "invalid request body"})
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, transportModel.ErrorResponse{Errors: "username and password required"})
		return
	}

	token := "MOCKED_JWT_TOKEN"

	resp := transportModel.AuthResponse{Token: token}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding json response: %v", err)
	}
}
