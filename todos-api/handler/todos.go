package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
	"todos3/todos-api/storer"
	"todos3/todos-api/token"
)

func (h *handler) createTodos(w http.ResponseWriter, r *http.Request) {
	var t ToDosReq
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)
	st := toStorerTodos(t)
	st.UserID = claims.ID

	todo, err := h.server.CreateTodos(h.ctx, st)
	if err != nil {
		http.Error(w, "error creating todos", http.StatusInternalServerError)
		return
	}
	res := toTodosRes(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getTodos(w http.ResponseWriter, r *http.Request) {
	var t ToDosReq
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	todo, err := h.server.GetTodos(h.ctx, claims.ID, t.ID)
	if err != nil {
		http.Error(w, "error getting todos", http.StatusInternalServerError)
		return
	}
	res := toTodosRes(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listUserTodos(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	list, err := pageSortFilter(r)
	if err != nil {
		http.Error(w, "error sort todos", http.StatusBadRequest)
		return
	}

	todos, err := h.server.ListUserTodos(h.ctx, claims.ID, list)
	if err != nil {
		http.Error(w, "error listing todos", http.StatusInternalServerError)
		return
	}

	var res ListToDosRes
	res.UserID = claims.ID
	for _, t := range todos {
		res.ToDos = append(res.ToDos, toTodosRes(&t))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) listTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := h.server.ListTodos(h.ctx)
	if err != nil {
		http.Error(w, "error listing todos", http.StatusInternalServerError)
		return
	}
	var res []ToDosRes
	for _, t := range todos {
		res = append(res, toTodosRes(&t))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateTodos(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	var t ToDosReq
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	todo, err := h.server.GetTodos(h.ctx, claims.ID, t.ID)
	if err != nil {
		http.Error(w, "error getting todos", http.StatusInternalServerError)
		return
	}

	//patch our todos request
	patchTodosReq(todo, t)

	updated, err := h.server.UpdateTodos(h.ctx, todo)
	if err != nil {
		http.Error(w, "error updating todos", http.StatusInternalServerError)
		return
	}

	res := toTodosRes(updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteTodos(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid ID, have no access for this content", http.StatusForbidden)
		return
	}

	if err := h.server.DeleteTodos(h.ctx, id); err != nil {
		http.Error(w, "error deleting todos", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pageSortFilter(r *http.Request) (list storer.List, err error) {

	list.Sort = r.URL.Query().Get("sort")
	list.Order = r.URL.Query().Get("order")
	list.Title = r.URL.Query().Get("title")
	list.Title = "%" + list.Title + "%"

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	limitString, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	if err != nil || limitString < 1 {
		limitString = 10 // Default to 1 items per page
	}

	list.Limit = int(limitString)

	// Calculate the OFFSET
	list.Offset = (int(page) - 1) * list.Limit

	return list, err
}

func getID(r *http.Request) (id int, err error) {
	i, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return int(i), fmt.Errorf("error parse ID: %v", err)

	}

	claims := r.Context().Value(authKey{}).(*token.UserClaims)

	if claims.ID != int(i) {
		return int(i), fmt.Errorf("invalid ID")

	}
	return int(i), err
}

func toStorerTodos(t ToDosReq) *storer.ToDos {
	return &storer.ToDos{
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

func toTodosRes(t *storer.ToDos) ToDosRes {
	return ToDosRes{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func patchTodosReq(todo *storer.ToDos, t ToDosReq) {

	if t.Title != "" {
		todo.Title = t.Title
	}
	if t.Description != "" {
		todo.Description = t.Description
	}

	if t.Completed != false {
		todo.Completed = t.Completed
	}

	todo.UpdatedAt = toTimePtr(time.Now())
}
