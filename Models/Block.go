package Models

import (
	"database/sql"
	"log"
	"errors"
	"time"
)

type Meta struct {
	Deleted		bool	`json:"deleted,omitempty"`
}

type Block struct{
	Id          int		`json:"block_id,omitempty"`
	ParentID    int		`json:"parent_id,omitempty"`
	Name        string	`json:"name,omitempty"`
	Content     string	`json:"content,omitempty"`
	Order       int		`json:"order,omitempty"`
	LastUpdated string	`json:"date,omitempty"`
	DocId		int		`json:"doc_id,omitempty"`
	Ltree		string	`json:"ltree,omitempty"`
	Meta        Meta	`json:"meta,omitempty"`
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

	BelongToDocumentAndUserOrPublic(db *sql.DB) bool
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

func (block *Block) Update(tx *sql.Tx) error{
	var parentId sql.NullInt64
	if block.ParentID == 0{
		parentId.Valid = false
	}else {
		parentId.Int64 = int64(block.ParentID)
	}
	_, err := tx.Exec(`SELECT update_block($1, $2, $3, $4, $5)`, block.Id, block.ParentID, block.Name, block.Content, block.Order)
	if err != nil {
		log.Println("Models.Block.Update ", err)
		return errors.New("something wrong")
	}
	return nil
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
	var parentId sql.NullInt64
	if block.ParentID == 0{
		parentId.Valid = false
	}else {
		parentId.Int64 = int64(block.ParentID)
	}
	err := tx.QueryRow(`INSERT INTO block(name, content, last_updated, parent_id, ord, doc_id)
							  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, block.Name, block.Content, time.Now(),
							  parentId.Int64, block.Order, block.DocId).Scan(&block.Id)
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
