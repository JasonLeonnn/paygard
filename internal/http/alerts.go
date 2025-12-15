package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JasonLeonnn/paygard/internal/services"
)

func GetAlertsHandler(service *services.AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		limit, _ := strconv.Atoi(limitStr)

		alerts, err := service.GetAlerts(r.Context(), limit)
		if err != nil {
			http.Error(w, "Failed to retrieve alerts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alerts)
	}
}
