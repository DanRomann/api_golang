package Usecases

import (
	"database/sql"
	"docnota/Models"
)

func GetCountries(db *sql.DB) ([]Models.Country, error){
	countries, err := Models.GetCountryList(db)
	if err != nil {
		return nil, err
	}
	return countries, nil
}