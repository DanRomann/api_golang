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

//ToDo
func GetCompanyDoc(companyId int, token *string, db *sql.DB) ([]Models.Document, error){
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
		if !company.IsPublic(db){
			return nil, errors.New("access denied")
		}
	}else {
		for _, curPerm := range permissions{
			if curPerm.CompanyId == companyId{
				allowed = curPerm.Read
				if allowed{}
			}
		}
	}


	return nil, nil
}