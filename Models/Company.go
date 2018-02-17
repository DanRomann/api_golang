package Models

import "database/sql"

type Company struct {
	Id 					int 		`json:"id,omitempty"`
	Name 				string 		`json:"name,omitempty"`
	Description			string  	`json:"description,omitempty"`
	Public				bool		`json:"public,omitempty"`
	Country				string		`json:"country,omitempty"`
	City				string		`json:"city,omitempty"`
	Street				string		`json:"street,omitempty"`
}

type Permissions struct {
	Read 		bool `json:"read"`
	Write 		bool `json:"write"`
	Update 		bool `json:"update"`
	Delete 		bool `json:"delete"`
	Invite 		bool `json:"invite"`
	Kick		bool `json:"kick"`
	Admin		bool `json:"admin"`
	Responsible bool `json:"responsible"`
}

type CompanyInteraction interface{
	Get(name *string, db *sql.DB) bool

	Create(tx *sql.Tx) error

	Delete(db *sql.DB) error

	GetDocuments(hasPermissions bool, db *sql.DB)

	SendInvite(db *sql.DB) error

}