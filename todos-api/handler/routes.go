package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var r *chi.Mux

func RegisterRouter(handler *handler) *chi.Mux {
	r = chi.NewRouter()
	tokenMaker := handler.TokenMaker

	//toDos

	r.Route("/todos", func(r chi.Router) {
		r.With(GetAdminMiddlewareFunc(tokenMaker)).Get("/listall", handler.listTodos)

		r.Route("/", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(GetAuthMiddlewareFunc(tokenMaker))
				r.Post("/", handler.createTodos)
				r.Get("/list", handler.listTodos)
				r.Get("/listuser", handler.listUserTodos)
				r.Get("/get", handler.getTodos)
				r.Patch("/", handler.updateTodos)
				r.Delete("/", handler.deleteTodos)
			})

		})
	})

	//Users

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.createUser)
		r.Post("/login", handler.loginUser)

		r.Group(func(r chi.Router) {
			r.Use(GetAdminMiddlewareFunc(tokenMaker))
			r.Get("/", handler.listUsers)
			r.Route("/{id}", func(r chi.Router) {
				r.Patch("/", handler.updateUser)
				r.Delete("/", handler.deleteUser)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(GetAuthMiddlewareFunc(tokenMaker))
			r.Patch("/", handler.updateUser)
			r.Post("/logout", handler.logoutUser)
		})

	})

	//Token

	r.Group(func(r chi.Router) {
		r.Use(GetAuthMiddlewareFunc(tokenMaker))
		r.Route("/tokens", func(r chi.Router) {
			r.Post("/renew", handler.renewAccessToken)
			r.Post("/revoke", handler.revokeSession)
		})
	})

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, r)
}
