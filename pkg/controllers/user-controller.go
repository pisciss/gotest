package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"goultimo/pkg/models"
	"goultimo/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type error interface {
	Error() string
}

var User models.User

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := models.GetUsers()
	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	ID, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	users, error := models.GetUserById(ID)
	if error.Error != nil {
		http.Error(w, "No encontrado", 400)
		return
	}

	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	utils.ParseBody(r, user)

	if len(user.Email) == 0 {
		http.Error(w, "El email es requerido", 400)
		return

	}

	existe, err := models.CheckEmail(user.Email)

	if err != nil {
		user.Password, err = models.GeneratehashPassword(user.Password)
		if err != nil {
			log.Fatalln("error in password hash")
		}
		u := user.CreateUser()
		res, _ := json.Marshal(u)
		w.Header().Set("Content-Type", "aplication/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
	if existe {
		http.Error(w, "Ya existe ese email registrado", 400)
		return
	}

}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	ID, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	users := models.DeleteUser(ID)
	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user = &models.User{}
	utils.ParseBody(r, user)
	vars := mux.Vars(r)

	userId := vars["id"]
	ID, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	users, db := models.GetUserById(ID)
	if users.Username != "" {
		users.Username = user.Username
	}
	if users.Email != "" {
		users.Email = user.Email
	}
	if users.Password != "" {
		users.Password = user.Password
	}
	db.Save(&users)
	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
func SignIn(w http.ResponseWriter, r *http.Request) {
	/*	connection := GetDatabase()
		defer Closedatabase(connection)
	*/

	var authdetails models.Authentication

	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		log.Printf("entra al primer")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)

		return
	}

	//var authuser models.User
	usr, _ := models.CheckUser(authdetails.Email)
	if usr.Email == "" {

		log.Printf("usuario no encontrado")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	log.Printf(usr.Password)

	check := models.CheckPasswordHash(authdetails.Password, usr.Password)

	if !check {
		log.Printf("contraseña incorrecta")
		//		err = errors.New("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(usr.Email, usr.Role)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token models.Token
	token.Email = usr.Email
	token.Role = usr.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte("captainjacksparrowsayshi6576657")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ProcesoToken(tk string) (*models.Token, bool, error) {
	miClave := []byte("captainjacksparrowsayshi6576657")
	claims := &models.Token{}

	splitToken := strings.Split(tk, "Bearer")
	if len(splitToken) != 2 {
		return claims, false, errors.New("formato de token invalido")

	}
	tk = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
		return miClave, nil
	})
	if err == nil {

		encontrado, _ := models.CheckEmail(claims.Email)
		if encontrado == true {
			Email = claims.Email
		}
		return claims, encontrado, nil
	}
	if !tkn.Valid {
		return claims, false, errors.New("token invalido")
	}
	return claims, false, err
}
