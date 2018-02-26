package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
	"errors"
	"log"
)

func GetDocument(requestId int, token *string, db *sql.DB) (*Models.Document, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	doc := new(Models.Document)
	doc.ID = requestId

	if doc.BelongToUserOrPublic(userId, db){
		err := doc.Get(db)
		if err != nil {
			return nil, err
		}
		err = doc.GetOwner(db)
		if err != nil {
			return nil, err
		}

		if doc.UserId !=  userId{
			doc.ReadOnly = true
		}

	}else {
		return nil, errors.New("document don't exists or not public")
	}
	return doc, nil
}

func GetPublicDocuments(isTemplate bool, db *sql.DB) ([]Models.Document, error){
	documents, err := Models.PublicDocuments(isTemplate, db)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func CopyDocument(docId int, token *string, db *sql.DB) error{
	//userId, _ := Utils.ParseToken(token)
	return nil
}

func CreateDocument(doc *Models.Document, token *string, db *sql.DB) (*Models.Document, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	if userId == 0{
		return nil, errors.New("access denied")
	}

	if doc.Name == ""{
		return nil, errors.New("empty doc name")
	}

	doc.UserId = userId
	err = doc.Create(db)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ChangeDoc(newBlock *Models.Block, token *string, db *sql.DB) (*Models.Block, error){
	document := new(Models.Document)
	document.ID = newBlock.DocId
	if document.ID == 0{
		return nil, errors.New("document id is empty")
	}

	user := new(Models.User)
	userId, err := Utils.ParseToken(token)
	user.ID = userId
	if err != nil {
		return nil, err
	}
	if user.ID == 0{
		return nil, errors.New("access denied")
	}

	if newBlock.Order == 0{
		return nil, errors.New("empty order")
	}

	if len(newBlock.Content) == 0 && len(newBlock.Name) == 0{
		return nil, errors.New("empty name and content block")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.Document.ChangeDoc ", err)
		return nil, errors.New("something wrong")
	}


	//block already exists
	if newBlock.Id != 0{
		//Check if the block belong to doc and user
		if !newBlock.BelongToDocumentAndUser(user.ID, document.ID, db){
			tx.Rollback()
			return nil, errors.New("access denied")
		}

		//Check if the block is deleted
		if newBlock.Meta.Deleted{
			err = newBlock.Delete(tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		//Update block
		err = newBlock.Update(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

	}


	//New block
	//Check permission for user to doc
	if !document.BelongToUser(user.ID, db){
		tx.Rollback()
		return nil, errors.New("access denied")
	}

	err = newBlock.Create(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("Usecases.Document.ChangeDoc ", err)
		return nil, errors.New("something wrong")
	}
	return newBlock, nil
}

func CopyDoc(docId int, token *string, db *sql.DB) error{
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	document := new(Models.Document)
	document.ID = docId

	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.Document.CopyDoc ", err)
		return errors.New("something wrong")
	}

	err = document.Copy(userId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
