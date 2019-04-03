package daylevels

import (
	"errors"
	"fmt"
	"github.com/mmarzio67/ml/config"
	"net/http"
	"strconv"
	"time"
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

type Session struct {
	Un           string
	LastActivity time.Time
}

func AllDL() ([]DayLevel, error) {
	rows, err := config.DB.Query("SELECT id, focus, arrabiato, depresso, created_on FROM daylevels")
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

	row := config.DB.QueryRow("SELECT id, focus, arrabiato, depresso, created_on FROM daylevels WHERE id = $1", id)

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

	/*
		// validate form values
		if dl.Isbn == "" || dl.Title == "" || dl.Author == "" || p == "" {
			return dl, errors.New("400. Bad request. All fields must be complete.")
		}

		// convert form values
		f64, err := strconv.ParseFloat(p, 32)
		if err != nil {
			return dl, errors.New("406. Not Acceptable. Price must be a number.")
		}
		dl.Price = float32(f64)
	*/

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
			meditazione, created_on) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = config.DB.Exec(queryDL, dl.Focus, dl.Fischio_orecchie, dl.Power_energy, dl.Dormito, dl.PR, dl.Ansia, dl.Arrabiato, dl.Irritato, dl.Depresso, dl.Cinque_tibetani, dl.Meditazione, dl.CreatedOn)
	if err != nil {
		return dl, errors.New("500. Internal Server Error." + err.Error())
	}
	//return the newly created ID

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

	/*
		if dl.Isbn == "" || dl.Title == "" || dl.Author == "" || p == "" {
			return dl, errors.New("400. Bad Request. Fields can't be empty.")
		}

		// convert form values
		f64, err := strconv.ParseFloat(p, 32)
		if err != nil {
			return dl, errors.New("406. Not Acceptable. Enter number for price.")
		}
		dl.Price = float32(f64)
	*/

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
