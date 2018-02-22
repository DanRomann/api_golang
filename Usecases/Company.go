package Usecases

import (
	"database/sql"
	"docnota/Models"
)

func GetCompanyList(db *sql.DB)([]Models.Company, error){
	companies, err := Models.CompanyList(db)
	if err != nil {
		return nil, err
	}

	return companies, nil
}