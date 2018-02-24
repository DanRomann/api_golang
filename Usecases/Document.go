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
	var contentChanged, parentChanged, nameChanged, orderChanged bool

	tx, err := db.Begin()
	if err != nil {
		log.Println("Usecases.Document.ChangeDoc ", err)
		return nil, errors.New("something wrong")
	}

	document := new(Models.Document)
	document.ID = newBlock.DocId
	user := new(Models.User)
	userId, err := Utils.ParseToken(token)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if userId == 0{
		tx.Rollback()
		return nil, errors.New("access denied")
	}
	user.ID = userId


	//block already exists
	if newBlock.Id != 0{
		//Check if the block belong to doc and user
		if !newBlock.BelongToDocumentAndUser(user.ID, document.ID, db){
			tx.Rollback()
			return nil, errors.New("access denied")
		}

		//Check if the block deleted
		if newBlock.Meta.Deleted{
			children, err := newBlock.GetParentChildren(db)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			for id, child := range children{
				if id > newBlock.Id{
					child.Id--
					err = child.Update(tx)
					if err != nil {
						tx.Rollback()
						return nil, err
					}
				}
			}

			err = newBlock.Delete(tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		//Get old version block for compare
		oldBlock := new(Models.Block)
		oldBlock.Id = newBlock.Id
		err = oldBlock.Get(db)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		//Check if the block change parent
		if oldBlock.ParentID != newBlock.ParentID {
			//Check user permission for new parent block
			if !user.HasPermissionToBlock(newBlock.ParentID, db){
				tx.Rollback()
				return nil, errors.New("access denied")
			}
			parentChanged = true

		}

		if oldBlock.Content != newBlock.Content{
			contentChanged = true
		}
		if oldBlock.Name != newBlock.Name{
			nameChanged = true
		}
		if oldBlock.Order != newBlock.Order{
			orderChanged = true
		}
		log.Println(nameChanged, contentChanged, parentChanged, orderChanged)
		log.Println("====================")

		if nameChanged || contentChanged || parentChanged || orderChanged{
			if orderChanged{
				//pizdec
				children, err := newBlock.GetParentChildren(db)
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				for id, child := range children {
					if id >= newBlock.Order{
						child.Order++
						err = child.Update(tx)
						if err != nil {
							tx.Rollback()
							return nil, err
						}
					}
				}

			}
			err = newBlock.Update(tx)
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
		}else {
			return nil, nil
		}

	}


	//Block are created
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
