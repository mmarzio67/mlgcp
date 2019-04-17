package daylevels

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mmarzio67/ml/config"
	"golang.org/x/crypto/bcrypt"
)

type DayLevel struct {
	Id               int64
	Focus            int64
	Fischio_orecchie int64
	Power_energy     int64
	Dormito          int64
	PR               int64
	Ansia            int64
	Arrabiato        int64
	Irritato         int64
	Depresso         int64
	Cinque_tibetani  bool
	Meditazione      bool
	CreatedOn        time.Time
}

type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

type Session struct {
	Un           string
	LastActivity time.Time
}

func SignupAuth(u *User) error {
	// Parse and decode the request body into a new `Credentials` instance
	Password := u.Password
	fmt.Println(u.UserName)

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), 8)
	fmt.Println(string(hashedPassword))

	// Next, insert the username, along with the hashed password into the database
	if _, err = config.DB.Query("insert into users (user_name, user_pwd, first_name, last_name, idrole) values ($1, $2,$3,$4,$5)", u.UserName, string(hashedPassword), u.First, u.Last, 1); err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		return err
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	return nil
}

func SigninAuth(w http.ResponseWriter, r *http.Request) {

	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the existing entry present in the database for the given username
	result := config.DB.QueryRow("select password from users where username=$1", creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in `storedCreds`
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
}

func AllDL() ([]DayLevel, error) {
	rows, err := config.DB.Query("SELECT id, focus, arrabiato, depresso, createdon FROM daylevels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dls := make([]DayLevel, 0)
	for rows.Next() {
		dl := DayLevel{}
		err := rows.Scan(
			&dl.Id,
			&dl.Focus,
			&dl.Arrabiato,
			&dl.Depresso,
			&dl.CreatedOn)
		if err != nil {
			return nil, err
		}
		dls = append(dls, dl)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dls, nil
}

func OneDL(r *http.Request) (DayLevel, error) {

	var err error
	dl := DayLevel{}
	id := r.FormValue("id")

	/*
		focus,err := strconv.ParseInt(r.FormValue("Focus"), 10, 64)
		if err != nil {
			// handle the error in some way
		}
	*/

	row := config.DB.QueryRow("SELECT id, focus, arrabiato, depresso, createdon FROM daylevels WHERE id = $1", id)

	err = row.Scan(
		&dl.Id,
		&dl.Focus,
		&dl.Arrabiato,
		&dl.Depresso,
		&dl.CreatedOn)
	if err != nil {
		return dl, err
	}

	return dl, nil
}

func PutDL(r *http.Request) (DayLevel, error) {
	var err error

	// get form values
	dl := DayLevel{}

	dl.Focus, err = strconv.ParseInt(r.FormValue("focus"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("focus entry not accepted")
		fmt.Println(err)
	}

	dl.Fischio_orecchie, err = strconv.ParseInt(r.FormValue("fischio_orecchie"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Fischio orecchie entry not accepted")
		fmt.Println(err)
	}

	dl.Power_energy, err = strconv.ParseInt(r.FormValue("power_energy"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Power energy entry not accepted")
		fmt.Println(err)
	}

	dl.Dormito, err = strconv.ParseInt(r.FormValue("dormito"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Dormito entry not accepted")
		fmt.Println(err)
	}

	dl.PR, err = strconv.ParseInt(r.FormValue("pr"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Public relations entry not accepted")
		fmt.Println(err)
	}

	dl.Ansia, err = strconv.ParseInt(r.FormValue("ansia"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Ansia entry not accepted")
		fmt.Println(err)
	}

	dl.Arrabiato, err = strconv.ParseInt(r.FormValue("arrabiato"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Arrabiato entry not accepted")
		fmt.Println(err)
	}

	dl.Irritato, err = strconv.ParseInt(r.FormValue("irritato"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("irritato entry not accepted")
		fmt.Println(err)
	}

	dl.Depresso, err = strconv.ParseInt(r.FormValue("depresso"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Depresso entry not accepted")
		fmt.Println(err)
	}

	dl.Cinque_tibetani, err = strconv.ParseBool(r.FormValue("cinque_tibetani"))
	if err != nil {
		// handle the error in some way
		fmt.Println("Cinque tibetani entry not accepted")
		fmt.Println(err)
	}

	dl.Meditazione, err = strconv.ParseBool(r.FormValue("meditazione"))
	if err != nil {
		// handle the error in some way
		fmt.Println("Meditazione entry not accepted")
		fmt.Println(err)
	}

	dl.CreatedOn = time.Now()
	fmt.Println(dl.CreatedOn)

	// insert values
	queryDL := `INSERT INTO daylevels (
			focus, 
			fischio_orecchie,
			power_energy,
			dormito,
			pr,
			ansia,
			arrabiato,
			irritato,
			depresso,
			cinque_tibetani,
			meditazione) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = config.DB.Exec(queryDL, dl.Focus, dl.Fischio_orecchie, dl.Power_energy, dl.Dormito, dl.PR, dl.Ansia, dl.Arrabiato, dl.Irritato, dl.Depresso, dl.Cinque_tibetani, dl.Meditazione)
	if err != nil {
		fmt.Println(dl)
		return dl, errors.New("500. Internal Server Error." + err.Error())
	}
	//return the newly created ID
	fmt.Println("sono le 9 e tutto va bene")
	lastidrow := config.DB.QueryRow("SELECT max(id) FROM daylevels")
	err = lastidrow.Scan(&dl.Id)
	if err != nil {
		return dl, err
	}
	fmt.Println(dl)
	return dl, nil
}

func UpdateDL(r *http.Request) (DayLevel, error) {
	// get form values
	var err error

	dl := DayLevel{}

	dl.Id, err = strconv.ParseInt(r.FormValue("Id"), 10, 64)

	dl.Focus, err = strconv.ParseInt(r.FormValue("Focus"), 10, 64)
	if err != nil {
		// handle the error in some way
	}

	dl.Arrabiato, err = strconv.ParseInt(r.FormValue("Arrabiato"), 10, 64)
	if err != nil {
		// handle the error in some way
	}

	dl.Depresso, err = strconv.ParseInt(r.FormValue("Depresso"), 10, 64)
	if err != nil {
		// handle the error in some way
	}
	dl.CreatedOn = time.Now()

	// insert values
	_, err = config.DB.Exec("UPDATE daylevels SET focus= $1, arrabiato=$2, depresso=$3 WHERE id=$4;", dl.Focus, dl.Arrabiato, dl.Depresso, dl.Id)
	if err != nil {
		return dl, err
	}
	return dl, nil
}

func DeleteDL(r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return errors.New("400. Bad Request.")
	}

	_, err := config.DB.Exec("DELETE FROM daylevels WHERE id=$1;", id)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
