package Usecases

import (
	"database/sql"
	"docnota/Models"
	"docnota/Utils"
	"github.com/pkg/errors"
)

func GetDocument(requestId int, token *string, db *sql.DB) (*Models.Document, error){
	userId := Utils.ParseToken(token)

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
