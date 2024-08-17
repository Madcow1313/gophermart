package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Madcow1313/gophermart/internal/auth"
	"github.com/Madcow1313/gophermart/internal/dbconnector"
	"github.com/Madcow1313/gophermart/internal/handler"
	"github.com/go-chi/chi"
)

type Server struct {
	Host        string
	DatabaseDSN string
	SecretKey   string
	HH          *handler.Handler
}

func New(host string) *Server {
	return &Server{Host: host}
}

func (s *Server) Serve() {
	r := chi.NewRouter()
	s.HH = handler.New()
	ba := auth.NewBasicAuth(s.SecretKey)

	s.HH.DBConnector = dbconnector.NewConnector(s.DatabaseDSN)
	err := s.HH.DBConnector.ConnectToGophermartDB(func(db *sql.DB, args ...interface{}) error {
		return s.HH.DBConnector.CreateTable(db, dbconnector.CreateGophermart)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = s.HH.DBConnector.ConnectToRegisterDB(func(db *sql.DB, args ...interface{}) error {
		return s.HH.DBConnector.CreateTable(db, dbconnector.CreateUsersQuery)
	})
	if err != nil {
		log.Fatal(err)
	}
	r.Post("/api/user/register", ba.AddUserCookie(s.HH.RegisterUser()))
	// r.Post("/api/user/login", nil)
	// r.Post("/api/user/orders", nil)
	// r.Get("/api/user/orders", nil)
	// r.Get("/api/user/balance", nil)
	// r.Post("/api/user/balance/withdraw", nil)
	// r.Get("/api/user/balance/withdrawals", nil)
	err = http.ListenAndServe(s.Host, r)
	if err != nil {
		log.Fatal(err)
	}
}
