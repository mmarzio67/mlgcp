package daylevels

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/mmarzio67/ml/config"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func Bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "007" {
		http.Error(w, "You must be 007 to enter the bar", http.StatusForbidden)
		return
	}
	showSessions() // for demonstration purposes
	config.TPL.ExecuteTemplate(w, "bar.html", u)
}

func Signup(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var us User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("Username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		r := req.FormValue("role")

		bs := []byte(p)
		us = User{un, bs, f, l, r}

		usertaken := SignupAuth(&us)

		if usertaken != nil {
			fmt.Println(usertaken)
			return
		}

		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		// store User in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		us = User{un, bs, f, l, r}
		dbUsers[un] = us

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return

	}

	showSessions() // for demonstration purposes
	config.TPL.ExecuteTemplate(w, "signup.html", us)
}

func Login(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var pc Credentials
	var ur *User
	var errAuth error
	var un string
	// process form submission
	if req.Method == http.MethodPost {
		uf := req.FormValue("Username")
		p := req.FormValue("password")

		pc = Credentials{p, uf}
		ur, errAuth = pc.LoginCred()
		if errAuth != nil {
			http.Error(w, "Something wrong with the user authentication", http.StatusForbidden)
			return
		}
		dbUsers[un] = *ur
		// does the entered password match the stored password?
		err := bcrypt.CompareHashAndPassword(ur.Password, []byte(p))
		if err != nil {
			fmt.Println(ur.Password)
			http.Error(w, "Username and/or password do not match, dude", http.StatusForbidden)
			return
		}

		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		//un := ur.UserName
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return

	}
	showSessions() // for demonstration purposes
	config.TPL.ExecuteTemplate(w, "login.html", ur)
}

func Logout(w http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}
