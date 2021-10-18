package routes

import (
	"goultimo/pkg/controllers"
	"goultimo/pkg/middleware"

	"github.com/gorilla/mux"
)

var RegistrarUsuarioRoute = func(router *mux.Router) {
	router.HandleFunc("/usuario/", middleware.ChequeoBD(controllers.CreateUser)).Methods("POST")
	router.HandleFunc("/usuario/", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/usuario/{id}",  middleware.ChequeoBD(middleware.ValidoJWT(controllers.GetUserById)).Methods("GET")
	router.HandleFunc("/usuario/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/usuario/{id}", controllers.DeleteUser).Methods("DELETE")
	//	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/signin/", controllers.SignIn).Methods("POST")
	router.HandleFunc("/expediente/", controllers.CreateExpediente).Methods("POST")
}
