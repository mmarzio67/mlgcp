package daylevels

import (
	"database/sql"
	"fmt"
	"log"
	"mmarzio/mlGCP/config"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	dls, err := AllDL()
	if err != nil {
		http.Redirect(w, r, "/dls/create", http.StatusSeeOther)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "daylevels.html", dls)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	config.TPL.ExecuteTemplate(w, "create.html", nil)
}

func CreateProcess(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	dl, err := PutDL(r)
	fmt.Println(dl)

	if err != nil {
		println("error in processing PutDL")
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	config.TPL.ExecuteTemplate(w, "created.html", dl)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'Id' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.

	id, err := strconv.ParseInt(keys[0], 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("id parameter reading accepted")
		fmt.Println(err)
	}

	dl, err := OneDL(id)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "update.html", dl)
}

func UpdateProcess(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	dl, err := UpdateDL(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	config.TPL.ExecuteTemplate(w, "updated.html", dl)
}

func DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	err := DeleteDL(r)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/daylevels", http.StatusSeeOther)
}
