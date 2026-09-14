package http

import (
	"net/http"
	"strconv"

	"forum/internal/service"
)

func InteractionHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		userID, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		value, err := strconv.Atoi(r.FormValue("value"))
		if err != nil {
			http.Error(w, "Bad Request: invalid reaction value", http.StatusBadRequest)
			return
		}
		if err := c.InteractionService.SetReaction(r.Context(), userID, r.FormValue("target_id"), r.FormValue("target_type"), value); err != nil {
			if err == service.ErrEmptyInteractionTarget || err == service.ErrInvalidInteractionTarget || err == service.ErrInvalidInteractionValue {
				http.Error(w, "Bad Request: invalid interaction", http.StatusBadRequest)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
