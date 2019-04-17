package meditation

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mmarzio67/ml/config"
)

type Meditation struct {
	Id         int64
	Meditation string
	Author     string
	PrefDay    int64
	PrefMonth  string
	Timeused   int64
	CreatedOn  time.Time
}

type ActionMed struct {
	Id        int64
	Action    string
	Idmed     int64
	Idusr     int64
	CreatedOn time.Time
}


func CreateMeditation(r *http.Request) (Meditation, error) {
	var err error

	// get form values
	mt := Meditation{}

	mt.Meditation = r.FormValue("pensiero")
	if err != nil {
		// handle the error in some way
		fmt.Println("meditation entry not accepted")
		fmt.Println(err)
	}

	mt.Author = r.FormValue("autore")
	if err != nil {
		// handle the error in some way
		fmt.Println("author entry not accepted")
		fmt.Println(err)
	}

	mt.PrefDay, err = strconv.ParseInt(r.FormValue("giorno"), 10, 64)
	if err != nil {
		// handle the error in some way
		fmt.Println("preferred day entry not accepted")
		fmt.Println(err)
	}

	mt.PrefMonth = r.FormValue("mese")
	if err != nil {
		// handle the error in some way
		fmt.Println("preferred month entry not accepted")
		fmt.Println(err)
	}

	mt.CreatedOn = time.Now()
	mt.Timeused = 0

	// insert values
	queryMt := `INSERT INTO meditations (
			meditation,
			timesused,
			pref_month,
			pref_day,
			author
			) 
			VALUES ($1, $2, $3, $4, $5)`
	_, err = config.DB.Exec(queryMt, mt.Meditation, mt.Timeused, mt.PrefMonth, mt.PrefDay, mt.Author)
	if err != nil {
		return mt, errors.New("500. Internal Server Error." + err.Error())
	}
	//return the newly created ID
	lastidrow := config.DB.QueryRow("SELECT max(id) FROM meditations")
	err = lastidrow.Scan(&mt.Id)
	if err != nil {
		return mt, err
	}
	fmt.Println(mt)
	return mt, nil
}

func AllMt() ([]Meditation, error) {
	rows, err := config.DB.Query("SELECT id, meditation, timesused, createdon FROM meditations")
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mts := make([]Meditation, 0)
	for rows.Next() {
		mt := Meditation{}
		err := rows.Scan(
			&mt.Id,
			&mt.Meditation,
			&mt.Timeused,
			&mt.CreatedOn)
		if err != nil {
			return nil, err
		}
		mts = append(mts, mt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mts, nil
}
