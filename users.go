package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rdleon/taquillaUno/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      int64  `json:"uid"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Users []User

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		hash  string
		uid   int
		err   error
		creds User
	)

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if _, ok := CheckAuth(r); ok {
		fmt.Fprintf(w, "{\"loggedin\": true}")
		return
	}

	creds = User{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&creds)

	err = db.Conn.QueryRow("SELECT uid, password FROM users WHERE email = $1 AND enabled = TRUE", creds.Email).Scan(&uid, &hash)

	if err == sql.ErrNoRows {
		// Timed derivation of valid email possible
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": \"Error\", \"uid\": -1}")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		LogError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": \"Wrong username/password\"}")
		return
	}

	response := make(map[string]string)

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		Issuer:    "taquilla.uno",
		Subject:   creds.Email,
	}

	// TODO: Add an ID to the jwt and save it on the database or
	// in a shared cache
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	response["token"], err = token.SignedString([]byte(Conf["secret"]))

	if err != nil {
		LogError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// TODO: invalidate the JWT

	if _, ok := CheckAuth(r); ok {
		response := map[string]string{
			"logout": "ok",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Fprintf(w, `{"logout": "ok"}`)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		hash  []byte
		err   error
		creds User
	)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&creds)

	if len(creds.Name) < 2 && len(creds.FullName) < 1 && len(creds.Email) < 3 {
		// Set bad request header
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"Missing parameters\"}")
		return
	} else if len(creds.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"Password too short\"}")
		return
	}

	hash, err = bcrypt.GenerateFromPassword([]byte(creds.Password), 10)

	if err != nil {
		LogError(w, err)
		return
	}

	stmt, err := db.Conn.Prepare("INSERT INTO users(user_name, full_name, email, password) values($1, $2, $3, $4)")

	// TODO: email already register error
	if err != nil {
		LogError(w, err)
		return
	}

	res, err := stmt.Exec(creds.Name, creds.FullName, creds.Email, string(hash))

	if err != nil {
		LogError(w, err)
		return
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		LogError(w, err)
		return
	}

	creds.UID = lastId

	resp := map[string]User{
		"created": creds,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var (
		uid   int
		uname string
		fname string
		email string
		user  User
		users Users
	)

	if _, ok := CheckAuth(r); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"error\": \"Unauthorized\"}")
		return
	}

	// TODO: Add pagination
	rows, err := db.Conn.Query("SELECT uid, user_name, fullname, email FROM users WHERE enabled = 't' ORDER BY uid DESC")
	if err != nil {
		LogError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&uid, &uname, &fname, &email)
		user = User{
			UID:      int64(uid),
			Name:     uname,
			FullName: fname,
			Email:    email,
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		LogError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		LogError(w, err)
		return
	}
}
