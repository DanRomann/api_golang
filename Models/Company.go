package Models

import (
	"database/sql"
	"log"
	"errors"
)

type Company struct {
	Id 					int 		`json:"id,omitempty"`
	Name 				string 		`json:"name,omitempty"`
	Description			string  	`json:"description,omitempty"`
	Public				bool		`json:"public,omitempty"`
	Country				string		`json:"country,omitempty"`
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

func CompanyList(db *sql.DB)([]Company, error){
	rows, err := db.Query(`SELECT company.id, company.name, description, country.name FROM company
 								 JOIN country ON country.id = company.country_id
 								 WHERE pub = TRUE`)
	if err != nil {
		log.Println("Models.Company.CompanyList ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	companies := make([]Company, 0)

	for rows.Next(){
		company := new(Company)
		err = rows.Scan(&company.Id, &company.Name, &company.Description, &company.Country)
		if err != nil {
			log.Println("Models.Company.CompanyList ", err)
			return nil, errors.New("something wrong")
		}

		companies = append(companies, *company)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Company.CompanyList ", err)
		return nil, errors.New("something wrong")
	}

	return companies, nil

}