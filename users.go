package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rdleon/taquillaUno/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      int    `json:"uid"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Created  int    `json:"created"`
	Enabled  bool   `json:"enabled"`
}

type Users []User

func LoginForm(w http.ResponseWriter, r *http.Request) {
	var (
		tpl  *template.Template
		sess *sessions.Session
		err  error
	)

	sess, err = Store.Get(r, "logged")

	if err != nil {
		LogError(w, err)
		return
	}

	tmp := sess.Values["uid"]

	if tmp != nil {
		// Already logged in
		http.Redirect(w, r, "/", 301)
		return
	}

	tpl, err = template.New("layout.html").ParseFiles("templates/layout.html", "templates/login.html")

	if err != nil {
		LogError(w, err)
		return
	}

	err = tpl.Execute(w, nil)

	if err != nil {
		LogError(w, err)
	}
}

func CreateUserForm(w http.ResponseWriter, r *http.Request) {
	var (
		tpl  *template.Template
		sess *sessions.Session
		err  error
	)

	sess, err = Store.Get(r, "logged")

	if err != nil {
		LogError(w, err)
		return
	}

	tmp := sess.Values["uid"]

	if tmp == nil {
		// Not logged in
		http.Redirect(w, r, "/", 301)
		return
	}

	tpl, err = template.New("create_user.html").ParseFiles("templates/layout.html", "templates/create_user.html")

	if err != nil {
		LogError(w, err)
		return
	}

	err = tpl.Execute(w, nil)

	if err != nil {
		LogError(w, err)
	}
}

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
		http.Redirect(w, r, "/", 301)
		return
	}

	email = r.FormValue("email")
	password = r.FormValue("password")

	err = db.Conn.QueryRow("SELECT uid, password FROM users WHERE email = $1 AND enabled = TRUE", email).Scan(&uid, &hash)

	if err == sql.ErrNoRows {
		// Timed derivation of valid email possible
		http.Redirect(w, r, "/login", 301)
		return
	} else if err != nil {
		LogError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	sess.Values["uid"] = uid
	sess.Save(r, w)

	http.Redirect(w, r, "/", 301)
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

	http.Redirect(w, r, "/login", 301)
}
