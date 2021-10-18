package middleware

import (
	"goultimo/pkg/config"
	"goultimo/pkg/controllers"
	"net/http"
)

func ChequeoBD(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.CheckConnection() == 0 {
			http.Error(w, "Conexion perdida con la base", 500)
			return
		}
		next.ServeHTTP(w, r)
	}

}
/*
func ValidoJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, err := controllers.ProcesoToken(r.Header.Get("Authorization"))
		if err != nil {

			http.Error(w, "Error en el token"+err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}

}
*/