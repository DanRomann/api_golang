package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
	"errors"
	"encoding/json"
	"log"
	"fmt"
)

func CreateCompany(curGroup *Models.Company, token *string, db *sql.DB) error{
	var countInpData, countNecessaryData int
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	if userId == 0{
		return errors.New("access denied")
	}

	var metaInfo map[string]interface{}
	countInpData = len(metaInfo)
	countryId, err := Models.GetCountryIdByName(&curGroup.Country, db)
	if err != nil {
		return err
	}
	companyMetaByCountry, err := Models.GetCountryMeta(countryId, db)
	if err != nil {
		return err
	}

	var curClaims []map[string]interface{}
	err = json.Unmarshal(*companyMetaByCountry, &curClaims)
	if err != nil {
		return err
	}

	for _, claims := range curClaims{
		for name, _ := range metaInfo{
			if name == claims["name"]{
				countNecessaryData++
			}
		}
	}

	if countInpData == countNecessaryData{
		tx, err := db.Begin()
		if err != nil {
			log.Println("Useceases.Company.CreateCompany ", err)
			return nil
		}
		result, _ := json.Marshal(curGroup.Meta)
		sha, err := Models.PrepareRegMessage(userId, countryId, &curGroup.Name, &curGroup.Description,
											curGroup.Public, json.RawMessage(result), tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		body := fmt.Sprintf("<html><body><h1><a href='http://%s%s'> Enter for confirm </a> \n" +
									" Company %s description %s \n meta %s", Utils.MainConfig.ServerConf.Address, *sha,
									curGroup.Name, curGroup.Description, metaInfo)

		err = Utils.SendEmail("bbshk@rsrch.ru", Utils.MainConfig.EmailConf.Addr, body, "Company registration")
		if err != nil {
			log.Println("Useceases.Company.CreateCompany ", err)
			return nil
		}

		err = tx.Commit()
		if err != nil {
			log.Println("Useceases.Company.CreateCompany ", err)
			return nil
		}
	}else {
		return errors.New("not enough necessary data for company registration")
	}

	return nil
}

func ConfirmCompany(sha *string, db *sql.DB)error{
	err := Models.ConfirmCompany(sha, db)
	if err != nil {
		return err
	}
	return nil
}

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

func GetCompanyMetaByCountry(countryId int, token *string, db *sql.DB)(*json.RawMessage, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	if userId == 0{
		return nil, errors.New("access denied")
	}

	meta, err := Models.GetCountryMeta(countryId, db)
	if err != nil {
		return nil, err
	}

	return meta, nil
}