package handler

import (
	"context"
	logging "todos3/todos-api/pkg"
	"todos3/todos-api/server"
	"todos3/todos-api/token"
)

type handler struct {
	ctx        context.Context
	server     *server.Server
	TokenMaker *token.JwtMaker
	logger     logging.Logger
}

func NewHandler(server *server.Server, logger logging.Logger, secretKey string) *handler {
	return &handler{
		ctx:        context.Background(),
		server:     server,
		TokenMaker: token.NewJwtMaker(secretKey),
		logger:     logger,
	}
}

//func (h handler) PaginationTodos(w http.ResponseWriter, r *http.Request) {
//	// Extract 'page' and 'limit' query parameters
//
//	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
//	fmt.Println("page:", page)
//	if err != nil || page < 1 {
//		page = 1 // Default to page 1
//	}
//	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
//	if err != nil || limit < 1 {
//		limit = 1 // Default to 10 items per page
//	}
//	fmt.Println("limit:", limit)
//
//	// Calculate the OFFSET
//	offset := (int(page) - 1) * int(limit)
//	fmt.Println("offset:", offset)
//
//	todos, err := h.server.PaginationTodos(h.ctx, int(limit), offset)
//	fmt.Println("todos:", todos)
//
//	var res []ToDosRes
//	for _, t := range todos {
//		res = append(res, toTodosRes(&t))
//	}
//	fmt.Println("res:", res)
//	if err != nil {
//		http.Error(w, "error pagination todos", http.StatusInternalServerError)
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(res)
//
//}
