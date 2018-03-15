package Models

import (
	"github.com/satori/go.uuid"
	"database/sql"
	"log"
	"errors"
	"time"
)

type Meta struct {
	Deleted		bool	`json:"deleted,omitempty"`
}

type Block struct{
	Id          	int			`json:"block_id,omitempty"`
	ParentID    	int			`json:"parent_id,omitempty"`
	Name        	string		`json:"name,omitempty"`
	Content     	string		`json:"content,omitempty"`
	Order       	int			`json:"order,omitempty"`
	LastUpdated 	string		`json:"date,omitempty"`
	DocId			int			`json:"doc_id,omitempty"`
	Ltree			string		`json:"ltree,omitempty"`
	Meta        	Meta		`json:"meta,omitempty"`
	RelationsCount	int			`json:"relations_count,omitempty"`
	UUID			string		`json:"uuid"`
}

type Relations struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

type BlockInteraction interface{
	Get(db *sql.DB) error
	GetChain(db *sql.DB) ([]Block, error)
	GetBlockHistory(docId int, db *sql.DB) ([]Block, error)

	Create(db *sql.DB) error

	Update(db *sql.DB) error

	Delete(db *sql.DB) error
	DeleteFromDocument(docId int, db *sql.DB) error
	DeleteFromGroup(groupId int, tx *sql.Tx) error

	SecureGet(userId int, db *sql.DB) bool
}

func (block *Block) Get(db *sql.DB) error{
	var parentId sql.NullInt64
	err := db.QueryRow(`SELECT id, parent_id, name, content, ord FROM block WHERE id = $1`, block.Id).Scan(
		&block.Id, &parentId, &block.Name, &block.Content, &block.Order)
	if err != nil {
		log.Println("Models.Block.Get ", err)
		return errors.New("something wrong")
	}
	if parentId.Valid{
		block.ParentID = int(parentId.Int64)
	}
	return nil
}

func (block *Block) Update(tx *sql.Tx) (int, error){
	var err 	error
	var newId 	int
	if block.ParentID == 0 {
		err = tx.QueryRow(`SELECT update_block($1, $2, $3, $4, $5, $6)`, block.Id, nil, block.Name,
			block.Content, block.Order, block.UUID).Scan(&newId)
	}else {
		err = tx.QueryRow(`SELECT update_block($1, $2, $3, $4, $5, $6)`, block.Id, block.ParentID,
			block.Name, block.Content, block.Order, block.UUID).Scan(&newId)
	}
	if err != nil {
		log.Println("Models.Block.Update ", err)
		return 0, errors.New("something wrong")
	}
	return newId, nil
}

func (block *Block) Delete(tx *sql.Tx) error{
	_, err := tx.Exec(`SELECT delete_block($1, $2, $3)`, block.Id, block.ParentID, block.Order)
	if err != nil {
		log.Println("Models.Block.Delete ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (block *Block) BelongToDocumentAndUser(userId, docId int, db *sql.DB) bool{
	var err  error
	var name string
	if block.ParentID != 0 {
		err = db.QueryRow(`SELECT block.name FROM block
								 JOIN document ON document.id = block.doc_id
								 JOIN client ON client.id = document.client_id
								 WHERE document.id = $1
								 AND client.id = $2
								 AND block.id = $3`,
								docId, userId, block.Id).Scan(&name)
		if err != nil {
			log.Println("Models.Block.BelongToDocumentAndUser ", err)
			return false
		}
		err = db.QueryRow(`SELECT block.name FROM block
								 JOIN document ON document.id = block.doc_id
								 JOIN client ON client.id = document.client_id
								 WHERE document.id = $1
								 AND client.id = $2
								 AND block.id = $3`,
			docId, userId, block.ParentID).Scan(&name)
		if err != nil {
			log.Println("Models.Block.BelongToDocumentAndUser ", err)
			return false
		}
	}
	return true
}

func (block *Block)	Create(tx *sql.Tx) error{
	var err error
	uid := uuid.NewV4()
	curTime := time.Now()
	if block.ParentID == 0 {
		err = tx.QueryRow(`INSERT INTO block(name, content, last_updated, parent_id, ord, doc_id, uid, start_date)
							  VALUES ($1, $2, $3, $4, $5, $6, &7, &8) RETURNING id`, block.Name, block.Content, curTime,
			nil, block.Order, block.DocId, uid, curTime).Scan(&block.Id)
	}else {
		err = tx.QueryRow(`INSERT INTO block(name, content, last_updated, parent_id, ord, doc_id, uid, start_date)
							  VALUES ($1, $2, $3, $4, $5, $6, &7, &8) RETURNING id`, block.Name, block.Content, curTime,
			block.ParentID, block.Order, block.DocId, uid, curTime).Scan(&block.Id)
	}
	if err != nil {
		log.Println("Models.Block.Create ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (block *Block) GetParentChildren(db *sql.DB) ([]Block, error){
	rows, err := db.Query(`SELECT block.id ,block.name, block.content, block.parent_id, block.last_updated, block.ord
 								 FROM block
								 JOIN block b ON block.parent_id = b.id
								 WHERE b.id = $1 ORDER BY ord`, block.ParentID)
	if err != nil {
		log.Println("Models.Block.GetParentChilder ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	blocks := make([]Block, 0)
	for rows.Next(){
		curBlock := new(Block)
		err = rows.Scan(&curBlock.Id, &curBlock.Name, &curBlock.Content, &curBlock.ParentID, &curBlock.LastUpdated, &curBlock.Order)
		if err != nil {
			log.Println("Models.Block.GetParentChilder ", err)
			return nil, errors.New("something wrong")
		}
		blocks = append(blocks, *curBlock)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Block.GetParentChilder ", err)
		return nil, errors.New("something wrong")
	}
	return blocks, nil
}

func SearchBlock(query *string, db *sql.DB) ([]Block, error){
	rows, err := db.Query(`SELECT id, name, content FROM block 
								 WHERE fts @@ plainto_tsquery('ru', $1) LIMIT 200`, query)
	if err != nil{
		log.Println("Search block ", err)
		return nil, err
	}
	defer rows.Close()
	blocks := make([]Block, 0)
	for rows.Next(){
		curBlock := Block{}
		err := rows.Scan(&curBlock.Id, &curBlock.Name, &curBlock.Content)
		if err != nil{
			log.Println("Search block ", err)
			return nil, err
		}
		blocks = append(blocks, curBlock)
	}
	if err = rows.Err(); err != nil{
		log.Println("Search block ", err)
		return nil, err
	}
	return blocks, nil
}

func (block *Block) SecureGet(userId int, db *sql.DB) error{
	var parentBlock	sql.NullInt64
	err := db.QueryRow(`SELECT b.name, b.content, b.last_updated, b.parent_id, b.ord, b.doc_id FROM block b
							  JOIN document doc ON doc.id = b.doc_id
							  WHERE b.id = $1 AND (doc.public = TRUE OR client_id = $2)`, block.Id, userId).Scan(&block.Name,
							  	&block.Content, &block.LastUpdated, &parentBlock, &block.Order, &block.DocId)
	if err != nil {
		return errors.New("access denied")
	}
	if parentBlock.Valid{
		block.ParentID = int(parentBlock.Int64)
	}
	return nil
}

func BlockBelongToUser(userId, blockId int, db *sql.DB) bool{
	var blockName string
	err := db.QueryRow(`SELECT b.name FROM block b
							  JOIN document doc ON doc.id = b.doc_id
							  WHERE b.id = $1 and doc.client_id = $2`, blockId, userId).Scan(&blockName)
	if err != nil {
		log.Println("Models.Block.BlockBelongToUser ", err)
		return false
	}
	return true
}

func BlockBelongOrPublic(userId, blockId int, db *sql.DB) bool{
	var blockName string
	err := db.QueryRow(`SELECT b.name FROM block b
							  JOIN document doc ON doc.id = b.doc_id
							  WHERE b.id = $1 and (doc.client_id = $2 OR doc.public = TRUE)`, blockId, userId).Scan(&blockName)
	if err != nil {
		log.Println("Models.Block.BlockBelongOrPublic ", err)
		return false
	}
	return true
}

func AddRelation(blockId, relationId int, db *sql.DB) error{
	_, err := db.Exec(`INSERT INTO block_relation(block_id, relation_block) VALUES ($1, $2)`, blockId, relationId)
	if err != nil {
		return errors.New("something wrong")
	}
	return nil
}

func SecureRelationDelete(blockId, relationId int, db *sql.DB) error{
	_, err := db.Exec(`DELETE FROM block_relation WHERE block_id = $1 AND relation_block = $2`, blockId, relationId)
	if err != nil {
		log.Println("Models.Block.SecureRelationDelete ", err)
		return errors.New("something wrong")
	}
	return nil
}

func BlockRelations(blockId int, db *sql.DB) ([]Block, error){
	rows, err := db.Query(`SELECT id, name, content, start_date, doc_id FROM block b
								 JOIN block_relation br ON br.relation_block = b.id
								 WHERE br.block_id = $1`, blockId)
	if err != nil {
		log.Println("Models.Block.BlockRelations ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	blocks := make([]Block, 0)
	for rows.Next(){
		block := new(Block)
		err := rows.Scan(&block.Id, &block.Name, &block.Content, &block.LastUpdated, &block.DocId)
		if err != nil {
			log.Println("Models.Block.BlockRelations ", err)
			return nil, errors.New("something wrong")
		}
		blocks = append(blocks, *block)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Block.BlockRelations ", err)
		return nil, errors.New("something wrong")
	}

	return blocks, nil
}