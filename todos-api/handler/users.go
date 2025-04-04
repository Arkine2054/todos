package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todos3/todos-api/storer"
	"todos3/todos-api/token"
	"todos3/todos-api/util"
)

var sqlErrCreateUser = errors.New("error getting user: sql: no rows in result set")

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u UsersReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// validation email and password
	if err := util.ValidateEmail(u.Email); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	if err := util.ValidatePassword(u.PasswordHash); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	storerUser, err := h.server.GetUser(h.ctx, u.Email)
	if err != nil {
		if errors.Is(err, sqlErrCreateUser) {
			http.Error(w, "error get data for creating user", http.StatusInternalServerError)
			return
		}
	}

	// email exist
	if u.Email == storerUser.Email {
		h.logger.WithField("email", storerUser.Email).Warn("user is already exist")
		http.Error(w, "email is already exist", http.StatusConflict)
		return
	}

	// hash password
	hashed, err := util.PasswordHash(u.PasswordHash)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	u.PasswordHash = hashed

	created, err := h.server.CreateUser(h.ctx, toStorerUser(u))
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listUsers(w http.ResponseWriter, r *http.Request) {
	list, err := pageSortFilter(r)
	if err != nil {
		http.Error(w, "error sort users", http.StatusInternalServerError)
	}

	users, err := h.server.ListUsers(h.ctx, list)
	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}

	var res ListUserRes
	for _, u := range users {
		res.Users = append(res.Users, toUserRes(&u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {

	var u UsersReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	user, err := h.server.GetUser(h.ctx, claims.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	patchUserReq(user, u)

	if user.Email == "" {
		user.Email = claims.Email
	}

	updated, err := h.server.UpdateUser(h.ctx, user)
	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid ID, have no access for this content", http.StatusForbidden)
		return
	}

	err = h.server.DeleteUser(h.ctx, id)
	if err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var u LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	gu, err := h.server.GetUser(h.ctx, u.Email)

	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	err = util.CheckPassword(u.PasswordHash, gu.PasswordHash)

	if err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}
	//	create token and return it as response

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 15*time.Hour)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 24*time.Hour)
	if err != nil {
		http.Error(w, "error creating refresh token", http.StatusInternalServerError)
		return
	}
	session, err := h.server.CreateSession(h.ctx, &storer.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    gu.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	res := LoginUserRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  toUserRes(gu),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

func (h handler) logoutUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	err := h.server.DeleteSession(h.ctx, claims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "error deleting session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toStorerUser(u UsersReq) *storer.Users {
	return &storer.Users{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsAdmin:      u.IsAdmin,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
func toUserRes(u *storer.Users) UsersRes {
	return UsersRes{
		ID:        u.ID,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func patchUserReq(user *storer.Users, u UsersReq) {

	if u.PasswordHash != "" {
		hashed, err := util.PasswordHash(u.PasswordHash)
		if err != nil {
			panic(err)
		}
		user.PasswordHash = hashed
	}

	user.UpdatedAt = toTimePtr(time.Now())
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}
