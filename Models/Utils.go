package Models

import (
	"database/sql"
	"log"
	"errors"
)

type Country struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
}

func GetCountryList(db *sql.DB) ([]Country, error){
	rows, err := db.Query("SELECT id, name FROM country")
	if err != nil {
		log.Println("Models.Utils.GetCountryList", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	countries := make([]Country, 0)
	for rows.Next(){
		curCountry := new(Country)
		err = rows.Scan(&curCountry.ID, &curCountry.Name)
		if err != nil {
			log.Println("Models.Utils.GetCountryList", err)
			return nil, errors.New("something wrong")
		}
		countries = append(countries, *curCountry)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Utils.GetCountryList", err)
		return nil, errors.New("something wrong")
	}
	return countries, nil
}