package auth

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

// function ini untuk melakukan init terhadap semua
// hal yang dibutuhkan oleh Auth Services
func Register2(router *chi.Mux, db *sql.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// seperti grouping endpoint
	// jadi yang ada di dalamnya sudah memiliki endpoint /api/auth
	// sebagai dasarnya
	router.Route("/ecommerce/auth", func(r chi.Router) {
		r.Post("/signup", handler.Register)
		r.Post("/signin", handler.Login)
	})
}
