package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
	"errors"
)

func GetCompanyList(db *sql.DB)([]Models.Company, error){
	companies, err := Models.CompanyList(db)
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func GetCompany(companyId int, token *string, db *sql.DB) (*Models.Company, error){
	var allowed bool
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}
	company := new(Models.Company)
	company.Id = companyId

	permissions, err := company.UserPermissions(userId, db)
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0{
		allowed = company.IsPublic(db)
	}else {
		for _, curPerm := range permissions{
			if curPerm.CompanyId == companyId{
				allowed = true
			}
		}
	}

	if !allowed{
		return nil, errors.New("access denied")
	}

	err = company.Get(userId, db)
	if err != nil {
		return nil, err
	}
	return company, nil
}


func GetCompanyDoc(companyId int, token *string, db *sql.DB) ([]Models.Document, error){
	var hasPermissions bool

	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	company := new(Models.Company)
	company.Id = companyId

	permissions, err := company.UserPermissions(userId, db)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0{
		if !company.IsPublic(db){
			return nil, errors.New("access denied")
		}
	}else {
		for _, curPerm := range permissions{
			if curPerm.CompanyId == companyId{
				hasPermissions = curPerm.Read
			}
		}
	}

	documents, err := company.Docs(hasPermissions, db)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func SearchCompanyByQuery(query *string, token *string, db *sql.DB) ([]Models.Company, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	if len(*query) == 0{
		return nil, errors.New("empty query")
	}

	companies, err := Models.SearchCompany(query, userId, db)
	if err != nil {
		return nil, err
	}

	return companies, nil


}