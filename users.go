package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
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
		email    string
		password string
		hash     string
		uid      int
		sess     *sessions.Session
		err      error
	)

	sess, err = Store.Get(r, "logged")

	if err != nil {
		LogError(w, err)
		return
	}

	tmp := sess.Values["uid"]

	if tmp != nil {
		// Already logged in
		fmt.Fprintf(w, "{\"status\": \"ok\", \"uid\": %d}", tmp)
		return
	}

	email = r.FormValue("email")
	password = r.FormValue("password")

	err = db.Conn.QueryRow("SELECT uid, password FROM users WHERE email = $1 AND enabled = TRUE", email).Scan(&uid, &hash)

	if err == sql.ErrNoRows {
		// Timed derivation of valid email possible
		fmt.Fprintf(w, "{\"status\": \"Error\", \"uid\": -1}")
		return
	} else if err != nil {
		LogError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		fmt.Fprintf(w, "{\"status\": \"Error\", \"uid\": -1}", uid)
		return
	}

	sess.Values["uid"] = uid
	sess.Values["isAdmin"] = true
	sess.Save(r, w)

	response := make(map[string]string)

	response["status"] = "ok"
	response["message"] = "User created"
	response["uid"] = strconv.Itoa(uid)

	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var (
		sess *sessions.Session
		err  error
	)

	sess, err = Store.Get(r, "logged")

	if err != nil {
		LogError(w, err)
		return
	}

	sess.Options.MaxAge = -1
	sess.Save(r, w)

	response := make(map[string]interface{})

	response["status"] = "ok"
	response["message"] = "Logged out"
	response["uid"] = -1

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
		fmt.Fprintf(w, "{\"status\": \"error\", \"message\":\"Missing parameters\"}")
		return
	} else if len(password) < 8 {
		fmt.Fprintf(w, "{\"status\": \"error\", \"message\":\"Password too short\"}")
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
