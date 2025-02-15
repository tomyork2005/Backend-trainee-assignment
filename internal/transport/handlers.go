package transport

import (
	"context"
	"encoding/json"
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
	GetInfo(ctx context.Context, id string) (string, error)
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

	router.Post("/auth", h.PostAuthEndpoint)

	router.Route("/api", func(routerWithAuth chi.Router) {

		routerWithAuth.Use(h.authMiddleware)

		routerWithAuth.Get("/info", h.GetInfo)
		routerWithAuth.Post("/sendCoin", h.PostSendCoin)
		routerWithAuth.Get("/buy/{item}", h.GetBuyItem)
	})

	return router
}

func (h *Handler) PostAuthEndpoint(w http.ResponseWriter, r *http.Request) {

	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		slog.Error("Error decoding auth request:", err)
		http.Error(w, "Failed to decode auth request", http.StatusBadRequest)
	}

	err := validate.Struct(&authRequest)
	if err != nil {
		slog.Error("Error validating auth request:", err)
		http.Error(w, "Failed to validate auth request", http.StatusBadRequest)
	}

}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	resp := InfoResponse{
		Coins: 1000,
		Inventory: []Item{
			{Type: "heelo my elizabet", Quantity: 2},
			{Type: "cup", Quantity: 1},
		},
		CoinHistory: CoinHistory{
			Received: []CoinReceived{
				{FromUser: "alice", Amount: 50},
			},
			Sent: []CoinSent{
				{ToUser: "bob", Amount: 20},
			},
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) PostSendCoin(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Errors: "invalid request body"})
		return
	}

	if req.Amount <= 0 || req.ToUser == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Errors: "invalid parameters"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetBuyItem - покупка предмета по названию
func (h *Handler) GetBuyItem(w http.ResponseWriter, r *http.Request) {
	item := chi.URLParam(r, "item")
	if item == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Errors: "item parameter is required"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PostAuth - авторизация/регистрация и выдача JWT
func (h *Handler) PostAuth(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Errors: "invalid request body"})
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Errors: "username and password required"})
		return
	}

	token := "MOCKED_JWT_TOKEN"

	resp := AuthResponse{Token: token}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding json response: %v", err)
	}
}
