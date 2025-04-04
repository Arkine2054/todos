package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"todos3/todos-api/token"
)

func (h handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	refreshClaims, err := h.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "error verifying refresh token", http.StatusInternalServerError)
		return
	}

	session, err := h.server.GetSession(h.ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		http.Error(w, "session is revoked", http.StatusUnauthorized)
		return
	}
	if session.UserEmail != refreshClaims.Email {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	res := RenewAccessTokenRes{
		accessToken,
		accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	err := h.server.RevokeSession(h.ctx, claims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "error revoking session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
