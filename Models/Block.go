package Models

import "database/sql"

type Meta struct {
	Deleted		bool	`json:"deleted"`
}

type Block struct{
	Id       string `json:"block_id,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Content  string `json:"content,omitempty"`
	Order    int    `json:"order,omitempty"`
	Date     string `json:"date,omitempty"`
	Meta     Meta   `json:"meta,omitempty"`
}

type BlockInteraction interface{
	Get(db *sql.DB) error
	GetChain(db *sql.DB) ([]Block, error)
	GetBlockHistory(docId int, db *sql.DB) ([]Block, error)
	Search(query string, db *sql.DB) ([]Block, error)

	Update(db *sql.DB) error

	Delete(db *sql.DB) error
	DeleteFromDocument(docId int, db *sql.DB) error
	DeleteFromGroup(groupId int, tx *sql.Tx) error

	BelongToDocumentAndUserOrPublic(db *sql.DB) bool
}

func (block *Block) Get(db *sql.DB) error{
	return nil
}

func (block *Block) Update(db *sql.DB) error{
	//Update name and content
	return nil
}

func (block *Block) Delete(db *sql.DB) error{
	return nil
}

func (block *Block) BelongToDocumentAndUser(db *sql.DB) bool{
	return false
}