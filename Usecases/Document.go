package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
	"github.com/pkg/errors"
	//"log"
	//"strconv"
	//"log"
	//"strconv"
)

func GetDocument(requestId int, token *string, db *sql.DB) (*Models.Document, error){
	userId, _ := Utils.ParseToken(token)

	doc := new(Models.Document)
	doc.ID = requestId

	if doc.BelongToUserOrPublic(userId, db){
		err := doc.Get(db)
		if err != nil {
			return nil, err
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

//ToDo change comparsion block content
//func ChangeDoc(newBlock *Models.Block, token *string, db *sql.DB) error{
//	user := new(Models.User)
//	userId := Utils.ParseToken(token)
//	user.ID = userId
//
//	tmpId, err := strconv.Atoi(newBlock.Id)
//	if err != nil {
//		//Block already exists
//
//		//Check if the newBlock belong to doc and user
//		if !newBlock.BelongToDocumentAndUser(db){
//			return errors.New("access denied")
//		}
//
//		//Check if the newBlock is deleted
//		if newBlock.Meta.Deleted{
//			err = newBlock.Delete(db)
//			if err != nil {
//				return err
//			}
//		}
//
//		//Get block from db to comparsion
//		oldBlock := new(Models.Block)
//		oldBlock.Id = newBlock.Id
//		err = oldBlock.Get(db)
//		if err != nil {
//			return err
//		}
//
//		//Check newBlock changed parent
//			//Check Permissions
//				//Change newBlock
//		if newBlock.ParentID != oldBlock.ParentID{
//
//		}
//
//
//		//Check if newBlock is changed content
//		if newBlock.Content != oldBlock.Content || newBlock.Name != oldBlock.Name{
//			err = newBlock.Update(db)
//			if err != nil {
//				return err
//			}
//		}
//	}else {
//		//New newBlock
//			//Copied
//				//Add newBlock to newBlock table and doc_block
//			//Created
//	}
//}
