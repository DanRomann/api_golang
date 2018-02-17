package Usecases

import (
	"docnota/Models"
	"database/sql"
	"encoding/json"
	"log"
	"errors"
	"docnota/Utils"
	"unicode/utf8"
)

func CreateUser(user *Models.User, db *sql.DB) error{
	if user.Exist(db){
		return errors.New("user already exists")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.User.CreateUser ", err)
		return errors.New("something wrong")
	}

	err = user.Create(tx)
	if err != nil {
		return err
	}

	uid, err := user.CreateUidForEmail(tx)
	if err != nil {
		return err
	}

	body := "<html><body><h1><a href='" + Utils.MainConfig.FrontConf.Addr + *uid + "'> Enter for confirm </a>"
	subject := "Confirmation Docnota account"
	err = Utils.SendEmail(user.Email, Utils.MainConfig.EmailConf.Addr, body, subject)
	if err != nil {
		tx.Rollback()
		log.Println("Usecases.User.CreateUser ", err)
		return errors.New("something wrong")
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Usecases.User.CreateUser ", err)
		return errors.New("something wrong")
	}
	return nil
}

func ConfirmUser(user *Models.User, uid *string, db *sql.DB) (*token, error){
	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.User.ConfirmUser ", err)
		return nil, errors.New("something wrong")
	}

	err = user.Confirm(uid, tx)
	if err != nil {
		return nil, err
	}

	token, err := Utils.CreateToken(user.ID, 7)
	if err != nil {
		log.Println("Usecases.User.ConfirmUser ", err)
		return nil, errors.New("something wrong")
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Usecases.User.ConfirmUser ", err)
		return nil, errors.New("something wrong")
	}
	return token, nil
}

func SearchUser(data []byte, db *sql.DB) ([]byte, error){
	var curRequest struct{
		Query	string	`json:"query"`
	}

	err := json.Unmarshal(data, &curRequest)
	if err != nil {
		log.Println("Usecases.User.SearchUsers ", err)
		return nil, errors.New("bad json")
	}

	if curRequest.Query == ""{
		return nil, errors.New("empty query parameter")
	}
	users, err := Models.SearchUsers(curRequest.Query, db)
	if err != nil {
		return nil, err
	}
	data, _ = json.Marshal(users)
	return data, nil
}

func GetUser(requestUser int, token *string, db *sql.DB) ([]byte, error){
	var isOwner bool
	userId := Utils.ParseToken(token)
	user := Models.User{
		ID: userId,
	}
	if userId == requestUser{
		isOwner = true
	}else {
		isOwner = false
	}

	err := user.Get(isOwner, db)
	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(user)
	return result, nil
}

func GetUsers(db *sql.DB) ([]byte, error){
	users, err := Models.GetUsers(db)
	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(users)
	return result, nil
}

func Auth(data []byte, db *sql.DB) (*string, error){
	user := new(Models.User)
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, errors.New("bad json")
	}

	if !Utils.EmailValid.MatchString(user.Email){
		return nil, errors.New("invalid email")
	}

	if utf8.RuneCountInString(user.Password) < 8 || utf8.RuneCountInString(user.Password) > 80{
		return nil, errors.New("password must be longer than 8 and less than 80 characters")
	}

	if user.CheckPassword(db){
		token, err := Utils.CreateToken(user.ID, 7)
		if err != nil {
			log.Println("Usecases.User.CheckPassword ", err)
		}
		return token, nil
	}
	return nil, errors.New("password or email invalid")
}

func GetUserCompany(requestUser int, token *string, db *sql.DB) ([]byte, error){
	var isOwner bool

	userId := Utils.ParseToken(token)
	user := new(Models.User)
	user.ID = userId


	if user.ID == requestUser{
		isOwner = true
	}else {
		isOwner = false
	}

	companies, err := user.GetUserCompany(isOwner, db)
	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(companies)
	return result, nil
}

func GetUserDocuments(){}
