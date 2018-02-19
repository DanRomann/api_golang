package Models

import (
	"time"
	"database/sql"
	"log"
	"github.com/pkg/errors"
)

type Document struct {
	ID          int       `json:"doc_id,omitempty"`
	UserId		int		  `json:"user_id,omitempty q"`
	Name        string    `json:"name,omitempty"`
	Blocks      []Block   `json:"blocks,omitempty"`
	Template    bool      `json:"template,omitempty"`
	Public      bool      `json:"public,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}

type DocumentInteraction interface{
	BelongToUser(userId int, db *sql.DB) bool
	BelongToUserOrPublic(userId int, db *sql.DB) bool
	BelongToGroup(groupId int, db *sql.DB) bool
	IsPublic(db *sql.DB) bool

	Create(userId int, name *string, isTemplate bool, isPublic bool, db *sql.DB) (*Document, error)
	SaveFillTemplate(userId int, tx *sql.Tx) error
	ReferToBlock(blockId, order, parent int, tx *sql.Tx) (*Document, error)

	SendDocumentToUser(userId int, db *sql.DB) error


	Get(db *sql.DB) (*Document, error)
	GetDocumentHistory(db *sql.DB) ([]Block, error)
	Search(query string, db *sql.DB) ([]Document, error)

	Delete(db *sql.DB) error
}

func (doc *Document) BelongToUser(userId int, db *sql.DB) bool{
	err := db.QueryRow("SELECT name FROM document WHERE id = $1 AND client_id = $2", doc.ID, userId).Scan(&doc.Name)
	if err != nil {
		return false
	}
	return true
}

func (doc *Document) BelongToUserOrPublic(userId int, db *sql.DB) bool{
	err := db.QueryRow("SELECT name FROM document WHERE id = $1 AND (client_id = $2 OR public = TRUE)",
		doc.ID, userId).Scan(&doc.Name)
	if err != nil{
		return false
	}
	return true
}

func (doc *Document) IsPublic(db *sql.DB) bool{
	err := db.QueryRow("SELECT name FROM document WHERE id = $1 AND public = TRUE", doc.ID).Scan(&doc.Name)
	if err != nil {
		return false
	}
	return true
}

func (doc *Document) Create(name *string, isTemplate bool, isPublic bool, db *sql.DB) error{
	curTime := time.Now()
	db.QueryRow("INSERT INTO document(name, client_id, template, public, created) VALUES($1, $2, $3, $4, $5)" +
							" RETURNING id", name, doc.UserId, isTemplate, isPublic, curTime).Scan(&doc.ID)
	return nil
}

func PublicDocuments(isTemplate bool, db *sql.DB) ([]Document, error){
	var rows	*sql.Rows
	var err		error

	if isTemplate{
		rows, err = db.Query("SELECT id, name, template, last_updated, created FROM document WHERE public = TRUE" +
									" AND template = TRUE")
	}else {
		rows, err = db.Query("SELECT id, name, template, last_updated, created FROM document WHERE public = TRUE" +
									" AND template = FALSE")
	}
	if err != nil {
		log.Println("Models.Document.PublicDocuments ", err)
		return nil, errors.New("something wrong")
	}

	defer rows.Close()

	docs := make([]Document, 0)

	for rows.Next(){
		doc := new(Document)
		err = rows.Scan(&doc.ID, &doc.Name, &doc.Template, &doc.LastUpdated, &doc.Created)
		if err != nil{
			return nil, err
		}
		docs = append(docs, *doc)
	}
	if len(docs) == 0{
		return nil, errors.New("no docs")
	}
	if 	err = rows.Err(); err != nil{
		log.Println("Models.Document.PublicDocuments ", err)
		return nil, errors.New("something wrong")
	}
	return docs, nil
}

func (doc *Document) Get(db *sql.DB) error{
	rows, err := db.Query("SELECT id, name, content, ord, parent_id FROM get_all_blocks($1)", doc.ID)
	if err != nil {
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}

	blocks := make([]Block, 0)
	for rows.Next() {
		block := new(Block)
		err = rows.Scan(&block.BlockId, &block.Name, &block.Content, &block.Order, &block.ParentID)
		if err != nil {
			log.Println("Models.Document.Get ", err)
			return errors.New("something wrong")
		}
		blocks = append(blocks, *block)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}

	err = db.QueryRow("SELECT name, template, last_updated, created FROM document WHERE id = $1", doc.ID).Scan(&doc.Name,
													&doc.Template, &doc.LastUpdated, &doc.Created)

	doc.Blocks = blocks
	//err = doc.sortBlocks()
	if err != nil{
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}
	defer rows.Close()
	return nil
}



