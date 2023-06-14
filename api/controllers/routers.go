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
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users", middlewares.ValidateTokenAndPerm(s.CreateUser, "admin")).Methods("POST")
	s.Router.HandleFunc("/users/{id}", middlewares.ValidateToken(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.ValidateTokenAndPerm(s.UpdateUser, "admin")).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.ValidateTokenAndPerm(s.DeleteUser, "admin")).Methods("DELETE")

	// Bike routes
	s.Router.Handle("/bikes", middlewares.ValidateToken(s.GetBikes)).Methods("GET")
	s.Router.Handle("/bikes", middlewares.ValidateTokenAndPerm(s.CreateBike, "write:bikes")).Methods("POST")
	s.Router.HandleFunc("/bikes/{id}", middlewares.ValidateToken(s.GetBike)).Methods("GET")
	s.Router.HandleFunc("/bikes/{id}", middlewares.ValidateTokenAndPerm(s.UpdateBike, "write:bikes")).Methods("PUT")
	s.Router.HandleFunc("/bikes/{id}", middlewares.ValidateTokenAndPerm(s.DeleteBike, "write:bikes")).Methods("DELETE")

	s.Router.Handle("/bikerims", middlewares.ValidateToken(s.GetBikeRims)).Methods("GET")
	s.Router.HandleFunc("/bikerims/{id}", middlewares.ValidateToken(s.GetBikeRim)).Methods("GET")

	// Audio routes
	s.Router.Handle("/audio", middlewares.ValidateTokenAndPerm(s.CreateAudio, "write:bikes")).Methods("POST")
	s.Router.HandleFunc("/audio/{id}", middlewares.ValidateTokenAndPerm(s.UpdateAudio, "write:bikes")).Methods("PUT")
	s.Router.HandleFunc("/audio/{id}", middlewares.ValidateTokenAndPerm(s.DeleteAudio, "write:bikes")).Methods("DELETE")
	s.Router.Handle("/audio/formats", middlewares.ValidateToken(s.GetAudioFormats)).Methods("GET")
	s.Router.HandleFunc("/audio/discogs", middlewares.ValidateToken(s.GetAudioData)).Methods("GET")
	s.Router.HandleFunc("/audio/title/{title}", middlewares.ValidateToken(s.GetAudiosByTitle)).Methods("GET")
	s.Router.HandleFunc("/audio/tracks/{id}", middlewares.ValidateToken(s.GetAudioTrackssById)).Methods("GET")

	// Video routes
	s.Router.Handle("/video/formats", middlewares.ValidateToken(s.GetVideoFormats)).Methods("GET")
	s.Router.Handle("/video", middlewares.ValidateToken(s.GetVideos)).Methods("GET")
	s.Router.HandleFunc("/video/moviesdb/{id}", middlewares.ValidateToken(s.GetVideoData)).Methods("GET")
	s.Router.Handle("/video", middlewares.ValidateTokenAndPerm(s.CreateVideo, "write:bikes")).Methods("POST")
	s.Router.HandleFunc("/video/{id}", middlewares.ValidateTokenAndPerm(s.UpdateVideo, "write:bikes")).Methods("PUT")
	s.Router.HandleFunc("/video/{id}", middlewares.ValidateTokenAndPerm(s.DeleteVideo, "write:bikes")).Methods("DELETE")
	s.Router.HandleFunc("/video/title/{title}", middlewares.ValidateToken(s.GetVideosByTitle)).Methods("GET")

	// Test Routes ==============================================================
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
