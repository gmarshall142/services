package controllers

import "github.com/gmarshall142/services/api/middlewares"

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
	s.Router.HandleFunc("/bikes", middlewares.SetMiddlewareJSON(s.GetBikes)).Methods("GET")
	s.Router.HandleFunc("/bikes", middlewares.SetMiddlewareJSON(s.CreateBike)).Methods("POST")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareJSON(s.GetBike)).Methods("GET")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateBike))).Methods("PUT")
	s.Router.HandleFunc("/bikes/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteBike)).Methods("DELETE")

	s.Router.HandleFunc("/bikerims", middlewares.SetMiddlewareJSON(s.GetBikeRims)).Methods("GET")
	s.Router.HandleFunc("/bikerims/{id}", middlewares.SetMiddlewareJSON(s.GetBikeRim)).Methods("GET")

	//Posts routes
	//s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	//s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	//s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	//s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	//s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
}
