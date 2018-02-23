package Usecases

import (
	"database/sql"
	"docnota/Models"
	"errors"
	"docnota/Utils"
)

func SearchBlockByQuery(query *string, token *string, db *sql.DB) ([]Models.Block, error){
	if len(*query) == 0{
		return nil, errors.New("empty query")
	}
	_, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	blocks, err := Models.SearchBlock(query, db)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}