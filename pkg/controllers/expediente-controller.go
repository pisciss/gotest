package controllers

import (
	"encoding/json"
	"goultimo/pkg/models"
	"goultimo/pkg/utils"
	"net/http"
	//"github.com/dgrijalva/jwt-go"
)

func CreateExpediente(w http.ResponseWriter, r *http.Request) {
	exp := &models.Expediente{}
	utils.ParseBody(r, exp)

	err := exp.Validate()
	if err != nil {

		http.Error(w, "Error en la validacion"+err.Error(), http.StatusBadRequest)
		return
	}

	//	existe, err := models.CheckEmail(user.Email)

	e := exp.CreateUser()
	res, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
