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
	Id              int64
	Focus           int64
	FischioOrecchie int64
	PowerEnergy     int64
	Dormito         int64
	PR              int64
	Ansia           int64
	Arrabiato       int64
	Irritato        int64
	Depresso        int64
	CinqueTib       bool
	Meditazione     bool
	CreatedOn       time.Time
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
	Password string
	Username string
}

type session struct {
	un           string
	lastActivity time.Time
}

func SignupAuth(u *User) error {
	// Parse and decode the request body into a new `Credentials` instance
	Password := u.Password

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), 8)

	// Next, insert the username, along with the hashed password into the database
	if _, err = config.DB.Query("insert into users (user_name, user_pwd, first_name, last_name, idrole) values ($1, $2,$3,$4,$5)", u.UserName, string(hashedPassword), u.First, u.Last, 1); err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		return err
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	return err
}

func (creds *Credentials) LoginCred() (u *User, e error) {

	var err error
	result := config.DB.QueryRow("select first_name, last_name, user_name, user_pwd, idrole from users where user_name=$1", creds.Username)
	if err != nil {
		fmt.Println("something wrong with the query to the credentials persistance")
		return nil, err
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	su := &User{}
	// Store the obtained password in `storedCreds`
	err = result.Scan(&su.First, &su.Last, &su.UserName, &su.Password, &su.Role)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			fmt.Println("username does not exit")
			return nil, err
		}
		// If the error is of any other type, send a 500 status
		fmt.Println("something wrong with the query to the credentials persistance")
		return nil, err
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(su.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		fmt.Println("password seem do not match")
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
	return su, nil
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

func AllDL() (*[]DayLevel, error) {

	queryAllDL := `SELECT id, 
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
					meditazione,
					createdon
					FROM daylevels
					WHERE uid=$1`

	rows, err := config.DB.Query(queryAllDL, 6)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dl := DayLevel{}
	dls := make([]DayLevel, 0)
	for rows.Next() {
		err := rows.Scan(
			&dl.Id,
			&dl.Focus,
			&dl.FischioOrecchie,
			&dl.PowerEnergy,
			&dl.Dormito,
			&dl.PR,
			&dl.Ansia,
			&dl.Arrabiato,
			&dl.Irritato,
			&dl.Depresso,
			&dl.CinqueTib,
			&dl.Meditazione,
			&dl.CreatedOn)

		if err != nil {
			return nil, err
		}
		dls = append(dls, dl)
		fmt.Println("all'interno di allDL")

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &dls, nil
}

func OneDL(id int64) (*DayLevel, error) {

	var err error
	dl := DayLevel{}
	fmt.Println(id)

	oneQueryDL := `SELECT id,
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
	meditazione
	FROM daylevels
	WHERE id=$1`

	row := config.DB.QueryRow(oneQueryDL, id)

	err = row.Scan(
		&dl.Id,
		&dl.Focus,
		&dl.FischioOrecchie,
		&dl.PowerEnergy,
		&dl.Dormito,
		&dl.PR,
		&dl.Ansia,
		&dl.Arrabiato,
		&dl.Irritato,
		&dl.Depresso,
		&dl.CinqueTib,
		&dl.Meditazione)

	if err != nil {
		fmt.Println(err)
		return &dl, err
	}
	return &dl, nil
}

func PutDL(r *http.Request) (*DayLevel, error) {
	var err error

	// get form values
	dl := DayLevel{}

	dl.Focus, err = strconv.ParseInt(r.FormValue("focus"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("focus entry not accepted")
		fmt.Println(err)
	}

	dl.FischioOrecchie, err = strconv.ParseInt(r.FormValue("fischio_orecchie"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Fischio orecchie entry not accepted")
		fmt.Println(err)
	}

	dl.PowerEnergy, err = strconv.ParseInt(r.FormValue("power_energy"), 10, 64)
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

	dl.CinqueTib, err = strconv.ParseBool(r.FormValue("cinque_tibetani"))
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
			meditazione,
			uid) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,$12)`
	var uid int
	uid = 6
	_, err = config.DB.Exec(queryDL,
		&dl.Focus,
		&dl.FischioOrecchie,
		&dl.PowerEnergy,
		&dl.Dormito,
		&dl.PR,
		&dl.Ansia,
		&dl.Arrabiato,
		&dl.Irritato,
		&dl.Depresso,
		&dl.CinqueTib,
		&dl.Meditazione,
		uid)
	if err != nil {
		return &dl, errors.New("500. Internal Server Error." + err.Error())
	}
	//return the newly created ID
	fmt.Println("sono le 9 e tutto va bene")
	lastidrow := config.DB.QueryRow("SELECT max(id) FROM daylevels")
	err = lastidrow.Scan(&dl.Id)
	if err != nil {
		return &dl, err
	}
	fmt.Println(&dl)
	return &dl, nil
}

func UpdateDL(r *http.Request) (*DayLevel, error) {
	// get form values
	var err error

	dl := DayLevel{}

	dl.Id, err = strconv.ParseInt(r.FormValue("Id"), 10, 64)

	dl.Focus, err = strconv.ParseInt(r.FormValue("focus"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("focus entry not accepted")
		fmt.Println(err)
	}

	dl.FischioOrecchie, err = strconv.ParseInt(r.FormValue("fischio_orecchie"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("Fischio orecchie entry not accepted")
		fmt.Println(err)
	}

	dl.PowerEnergy, err = strconv.ParseInt(r.FormValue("power_energy"), 10, 64)
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

	dl.CinqueTib, err = strconv.ParseBool(r.FormValue("cinque_tibetani"))
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

	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	dl.CreatedOn = t
	fmt.Println("all'interno di UpdateDL")

	updateQuery := `UPDATE daylevels 
			  SET focus= $1,
			  fischio_orecchie=$2,
			  power_energy=$3,
			  dormito=$4,
			  pr=$5,
			  ansia=$6,
			  arrabiato=$7,
			  irritato=$8,
			  depresso=$9,
			  cinque_tibetani=$10,
			  meditazione=$11
			  WHERE id=$12`

	// insert values
	_, err = config.DB.Exec(updateQuery,
		&dl.Focus,
		&dl.FischioOrecchie,
		&dl.PowerEnergy,
		&dl.Dormito,
		&dl.PR,
		&dl.Ansia,
		&dl.Arrabiato,
		&dl.Irritato,
		&dl.Depresso,
		&dl.CinqueTib,
		&dl.Meditazione,
		&dl.Id)

	if err != nil {
		return &dl, err
	}
	return &dl, nil
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
