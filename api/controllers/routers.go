package controllers

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gmarshall142/services/api/middlewares"
	"net/http"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// Bike routes
	//s.Router.Handle("/bikes", middlewares.ValidateToken(s.GetBikes)).Methods("GET")
	s.Router.Handle("/bikes", middlewares.ValidateTokenAndPerm(s.GetBikes, "write:bikes")).Methods("GET")
	s.Router.Handle("/bikes", middlewares.ValidateTokenAndPerm(s.CreateBike, "write:bikes")).Methods("POST")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareJSON(s.GetBike)).Methods("GET")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateBike))).Methods("PUT")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteBike)).Methods("DELETE")

	//s.Router.Handle("/bikerims", middlewares.SetMiddlewareValidToken(s.GetBikeRims)).Methods("GET")
	s.Router.Handle("/bikerims", middlewares.ValidateToken(s.GetBikeRims)).Methods("GET")
	s.Router.HandleFunc("/bikerims/{id}", middlewares.SetMiddlewareJSON(s.GetBikeRim)).Methods("GET")

	// This route is always accessible.
	s.Router.Handle("/api/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello from a public endpoint! You don't need to be authenticated to see this."}`))
	}))

	// This route is only accessible if the user has a valid access_token.
	s.Router.Handle("/api/private", middlewares.EnsureValidTokenSV()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private endpoint! You need to be authenticated to see this."}`))
		}),
	))

	// This route is only accessible if the user has a
	// valid access_token with the read:messages scope.
	s.Router.Handle("/api/private-scoped", middlewares.EnsureValidTokenSV()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			claims := token.CustomClaims.(*middlewares.CustomClaims)
			if !claims.HasScope("read:messages") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private endpoint! You need to be authenticated to see this."}`))
		}),
	))
}
