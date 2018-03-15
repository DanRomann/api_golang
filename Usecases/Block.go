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

func GetBlockById(blockId int, token *string, db *sql.DB) (*Models.Block, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	block := new(Models.Block)
	block.Id = blockId

	err = block.SecureGet(userId, db)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func AddBlockRelation(blockId, relationBlock int, token *string, db *sql.DB) error{
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	if userId == 0{
		return errors.New("access denied")
	}


	blockBelong := Models.BlockBelongToUser(userId, blockId, db)
	relationAccess := Models.BlockBelongOrPublic(userId, blockId, db)

	if !(blockBelong && relationAccess){
		return errors.New("access denied")
	}

	err = Models.AddRelation(blockId, relationBlock, db)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRelation(blockId, relationBlock int, token *string, db *sql.DB) error{
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return err
	}

	if userId == 0{
		return errors.New("access denied")
	}

	blockBelong := Models.BlockBelongToUser(userId, blockId, db)

	if !blockBelong{
		return errors.New("access denied")
	}

	err = Models.SecureRelationDelete(blockId, relationBlock, db)
	if err != nil {
		return err
	}
	return nil
}

func GetBlockRelations(blockId int, token *string, db *sql.DB) ([]Models.Block, error){
	userId, err := Utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	permissions := Models.BlockBelongOrPublic(userId, blockId, db)

	if !permissions{
		return nil, errors.New("access denied")
	}

	relations, err := Models.BlockRelations(blockId, db)
	if err != nil {
		return nil, err
	}

	return relations, nil
}