package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
)

func GetCompanyList(db *sql.DB)([]Models.Company, error){
	companies, err := Models.CompanyList(db)
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func GetCompany(companyId int, token *string, db *sql.DB) (*Models.Company, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}
	company := new(Models.Company)
	company.Id = companyId

	err = company.Get(userId, db)
	if err != nil {
		return nil, err
	}
	return company, nil
}