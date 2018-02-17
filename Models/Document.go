package Models

import (
	"time"
	"database/sql"
)

type Document struct {
	ID          int       `json:"doc_id,omitempty"`
	UserId		int		  `json:"user_id"`
	Name        string    `json:"name,omitempty"`
	Blocks      []Block   `json:"blocks,omitempty"`
	Template    bool      `json:"template"`
	Public      bool      `json:"public"`
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
	GetPublicDocuments(db *sql.DB) ([]Document, error)
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




