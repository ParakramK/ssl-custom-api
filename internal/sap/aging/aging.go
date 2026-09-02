package aging

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetCustomerAging(w http.ResponseWriter, r *http.Request) {
	cardCode := r.URL.Query().Get("CardCode")
	companyDB := r.URL.Query().Get("CompanyDB")

	if cardCode == "" {
		http.Error(w, "CardCode is required", http.StatusBadRequest)
		return
	}

	if companyDB == "" {
		http.Error(w, "CompanyDB is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	response, err := h.service.GetCustomerAging(
		ctx,
		cardCode,
		companyDB,
	)
	if err != nil {
		http.Error(
			w,
			"failed to fetch customer aging",
			http.StatusInternalServerError,
		)

		log.Printf(
			"customer aging error: cardCode=%s companyDB=%s error=%v",
			cardCode,
			companyDB,
			err,
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}
