package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrigangha/GoAuthApi/internals/api"
	"github.com/mrigangha/GoAuthApi/internals/services"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(services.CorsMiddleware)
	r.With(services.RateLimiterMiddleware).Get("/show", api.GetUsers)
	r.Post("/signup", api.CreateUser)
	http.ListenAndServe(":8080", r)
	os.Exit(0)
}
