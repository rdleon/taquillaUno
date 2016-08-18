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
	UID      int    `json:"uid"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Users []User

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		hash string
		uid  int
		err  error
	)

	if _, ok := CheckAuth(r); ok {
		fmt.Fprintf(w, "{\"loggedin\": true}")
		return
	}

	credentials := struct {
		email    string
		password string
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&credentials)

	err = db.Conn.QueryRow("SELECT uid, password FROM users WHERE email = $1 AND enabled = TRUE", credentials.email).Scan(&uid, &hash)

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

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(credentials.password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": \"Wrong username/password\"}")
		return
	}

	response := make(map[string]string)

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		Issuer:    "taquilla.uno",
		Subject:   credentials.email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	// TODO: Use a configurable and secret key
	response["token"], err = token.SignedString([]byte("verysecretKey"))

	if err != nil {
		LogError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	response["status"] = "ok"

	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		uname    string
		fname    string
		email    string
		password string
		hash     []byte
		err      error
	)

	uname = r.FormValue("uname")
	fname = r.FormValue("fname")
	email = r.FormValue("email")
	password = r.FormValue("password")

	if len(uname) < 2 && len(fname) < 1 && len(email) < 3 {
		// Set bad request header
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"Missing parameters\"}")
		return
	} else if len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"Password too short\"}")
		return
	}

	hash, err = bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		LogError(w, err)
		return
	}

	stmt, err := db.Conn.Prepare("INSERT INTO users(user_name, full_name, email, password) values($1, $2, $3, $4)")

	if err != nil {
		LogError(w, err)
		return
	}

	res, err := stmt.Exec(uname, fname, email, string(hash))

	if err != nil {
		LogError(w, err)
		return
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		LogError(w, err)
		return
	}

	fmt.Fprintf(w, "{\"status\": \"ok\", \"message\": \"User Created\", \"uid\": %d}", lastId)
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

	rows, err := db.Conn.Query("SELECT uid, user_name, fullname, email FROM users WHERE enabled = 't' ORDER BY uid DESC")
	if err != nil {
		LogError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&uid, &uname, &fname, &email)
		user = User{}
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
