package main

import (
	"mmarzio/medaliving/webapp/daylevels"
	"net/http"
)

func main() {
	http.HandleFunc("/", daylevels.Index)
	http.HandleFunc("/dls", daylevels.Index)
	http.HandleFunc("/dls/show", daylevels.Show)
	http.HandleFunc("/dls/create", daylevels.Create)
	http.HandleFunc("/dls/create/process", daylevels.CreateProcess)
	http.HandleFunc("/dls/update", daylevels.Update)
	http.HandleFunc("/dls/update/process", daylevels.UpdateProcess)
	http.HandleFunc("/dls/delete/process", daylevels.DeleteProcess)
	http.HandleFunc("/signup", daylevels.Signup)
	http.HandleFunc("/login", daylevels.Login)
	http.HandleFunc("/logout", daylevels.Logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dls", http.StatusSeeOther)
}
