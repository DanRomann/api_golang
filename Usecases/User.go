package Usecases

import (
	"docnota/Models"
	"database/sql"
	"encoding/json"
	"log"
	"errors"
	"docnota/Utils"
	"unicode/utf8"
	"strconv"
)

func CreateUser(user *Models.User, db *sql.DB) error{
	if !Utils.EmailValid.MatchString(user.Email){
		return errors.New("bad email")
	}

	if utf8.RuneCountInString(user.Password) > 100 || utf8.RuneCountInString(user.Password) < 8{
		return errors.New("password must be longer than 8 and less than 100 symbols")
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

func ConfirmUser(user *Models.User, uid *string, db *sql.DB) (*string, error){
	if !Utils.IsLetter.MatchString(user.FirstName) || !Utils.IsLetter.MatchString(user.LastName){
		return nil, errors.New("first and last name must contains only letters")
	}

	if utf8.RuneCountInString(user.FirstName) > 100 || utf8.RuneCountInString(user.LastName) > 100{
		return nil, errors.New("first and last name must be less than 100 symbols")
	}

	if len(user.FirstName) == 0 || len(user.LastName) == 0{
		return nil, errors.New("first and last name should not be empty")
	}

	_, err := strconv.Atoi(user.Country)
	if err != nil {
		return nil, errors.New("bad country id")
	}

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

func GetUser(requestUser int, token *string, db *sql.DB) (*Models.User, error){
	var isOwner bool

	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}


	user := new(Models.User)
	user.ID = requestUser

	if userId == requestUser{
		isOwner = true
	}else {
		isOwner = false
	}

	err = user.Get(isOwner, db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsers(db *sql.DB) ([]byte, error){
	users, err := Models.GetUsers(db)
	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(users)
	return result, nil
}

func Auth(user *Models.User, db *sql.DB) (*string, error){
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

	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}
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

func GetUserDocuments(requestId int, isTemplate bool, token *string, db *sql.DB)([]Models.Document, error){
	var isOwner bool

	userId, _ := Utils.ParseToken(token)


	user := new(Models.User)
	user.ID = requestId
	if userId == requestId{
		isOwner = true
	}else {
		isOwner = false
	}

	documents, err := user.GetDocuments(isOwner, isTemplate, Utils.Connect)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func GetUserInboxDocument(requetId int, token *string, db *sql.DB)([]Models.Document, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	if userId != requetId{
		return nil, errors.New("access denied")
	}

	user := new(Models.User)
	user.ID = requetId

	documents, err := user.InboxDocuments(db)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func SendDocumentToUser(docId, userId int, token *string, db *sql.DB) error{
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	if userId == 0{
		return errors.New("access denied")
	}

	document := new(Models.Document)
	document.ID = docId

	if !document.BelongToUserOrPublic(userId, db){
		return errors.New("access denied")
	}

	err = document.SendDocumentToUser(userId, db)
	if err != nil {
		return err
	}

	return nil
}

func AcceptDoc(docId int, token *string, db *sql.DB) error{
	var allowed bool
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	if userId == 0{
		return errors.New("access denied")
	}

	user := new(Models.User)
	user.ID = userId

	inboxDoc, err := user.InboxDocuments(db)
	if err != nil {
		return err
	}

	for _, doc := range inboxDoc{
		if doc.ID == docId{
			allowed = true
		}
	}

	if !allowed{
		return errors.New("access denied")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.User.AcceptDoc ", err)
		return errors.New("something wrong")
	}

	err = user.AcceptDoc(docId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Usecases.User.AcceptDoc ", err)
		return errors.New("something wrong")
	}
	return nil
}