package server

import (
	"context"
	"todos3/todos-api/storer"
)

type Server struct {
	storer *storer.PostgresStorer
}

func NewServer(storer *storer.PostgresStorer) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) CreateTodos(ctx context.Context, t *storer.ToDos) (*storer.ToDos, error) {
	return s.storer.CreateTodos(ctx, t)
}

func (s *Server) GetTodos(ctx context.Context, userID, todosID int) (*storer.ToDos, error) {
	return s.storer.GetTodos(ctx, userID, todosID)
}

func (s *Server) ListUserTodos(ctx context.Context, userID int, list storer.List) ([]storer.ToDos, error) {
	return s.storer.ListUserTodos(ctx, userID, list)
}

func (s *Server) ListTodos(ctx context.Context) ([]storer.ToDos, error) {
	return s.storer.ListTodos(ctx)
}

func (s *Server) UpdateTodos(ctx context.Context, t *storer.ToDos) (*storer.ToDos, error) {
	return s.storer.UpdateTodos(ctx, t)
}

func (s *Server) DeleteTodos(ctx context.Context, id int) error {
	return s.storer.DeleteTodos(ctx, id)
}

func (s *Server) CreateUser(ctx context.Context, u *storer.Users) (*storer.Users, error) {
	return s.storer.CreateUser(ctx, u)
}

func (s *Server) GetUser(ctx context.Context, email string) (*storer.Users, error) {
	return s.storer.GetUser(ctx, email)

}

func (s *Server) ListUsers(ctx context.Context, list storer.List) ([]storer.Users, error) {
	return s.storer.ListUsers(ctx, list)

}

func (s *Server) UpdateUser(ctx context.Context, u *storer.Users) (*storer.Users, error) {
	return s.storer.UpdateUser(ctx, u)

}

func (s *Server) DeleteUser(ctx context.Context, id int) error {
	return s.storer.DeleteUser(ctx, id)

}

func (s *Server) CreateSession(ctx context.Context, se *storer.Session) (*storer.Session, error) {
	return s.storer.CreateSession(ctx, se)
}

func (s *Server) GetSession(ctx context.Context, id string) (*storer.Session, error) {
	return s.storer.GetSession(ctx, id)
}

func (s *Server) RevokeSession(ctx context.Context, id string) error {
	return s.storer.RevokeSession(ctx, id)
	//err := s.storer.RevokeSession(ctx, sr.GetId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &pb.SessionRes{}, nil
}

func (s *Server) DeleteSession(ctx context.Context, id string) error {
	return s.storer.DeleteSession(ctx, id)
	//err := s.storer.DeleteSession(ctx, sr.GetId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &pb.SessionRes{}, nil
}
