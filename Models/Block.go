package Models

import "database/sql"

type Block struct{
	BlockId		string		`json:"block_id,omitempty"`
	ParentID	string		`json:"parent_id,omitempty"`
	Name		string		`json:"name,omitempty"`
	Content		string		`json:"content,omitempty"`
	Order		int			`json:"order,omitempty"`
	Date		string		`json:"date,omitempty"`
}

type BlockInteraction interface{
	Get(db *sql.DB) (*Block, error)
	GetChain(db *sql.DB) ([]Block, error)
	GetBlockHistory(docId int, db *sql.DB) ([]Block, error)
	Search(query string, db *sql.DB) ([]Block, error)

	Update(name, content string, db *sql.DB) error

	Delete(db *sql.DB) error
	DeleteFromDocument(docId int, db *sql.DB) error
	DeleteFromGroup(groupId int, tx *sql.Tx) error

	BelongToDocument(docId int, db *sql.DB) bool
}