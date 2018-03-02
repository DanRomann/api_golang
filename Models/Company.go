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
	CompanyId	int	 `json:"company_id"`
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
	Get(companyId int, db *sql.DB)

	Create(tx *sql.Tx) error

	Delete(db *sql.DB) error

	GetDocuments(hasPermissions bool, db *sql.DB)

	SendInvite(userId int, db *sql.DB) error

	PublicOrUserIsMember(userId, db *sql.DB) error

	IsPublic(db *sql.DB) bool

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

func (company *Company)	Get(userId int, db *sql.DB) error{
	err := db.QueryRow(`SELECT cm.name, cm.description, c.name  FROM company cm
								JOIN country c ON cm.country_id = c.id
								JOIN client_company cc ON cc.client_id = $1
								WHERE cm.id = $2 AND ((cc.client_confirm = TRUE AND cc.company_confirm = TRUE) OR cm.pub = TRUE);`,
								userId,	company.Id).Scan(&company.Name, &company.Description, &company.Country)
	if err != nil {
		log.Println("Models.Company.Get ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (company *Company)	SendInvite(userId int, db *sql.DB) error{
	return nil
}

func (company *Company) UserPermissions(userId int, db *sql.DB) ([]Permissions, error){
	rows, err := db.Query(`SELECT company_id, gr_admin, gr_invite, gr_kick, gr_read, gr_write, gr_update, gr_delete FROM client_company cc
								 JOIN company ON company.id = cc.company_id
								 JOIN client ON client.id = cc.client_id
								 WHERE client_id = $1`, userId)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	permissions := make([]Permissions, 0)
	for rows.Next(){
		curPermission := new(Permissions)
		err = rows.Scan(&curPermission.CompanyId, &curPermission.Admin, &curPermission.Invite, &curPermission.Kick,
						&curPermission.Read, &curPermission.Write, &curPermission.Update, &curPermission.Delete)
		if err != nil {
			log.Println("Models.Company.UserPermissions ", err)
			return nil, errors.New("something wrong")
		}

		permissions = append(permissions, *curPermission)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Company.UserPermissions ", err)
		return nil, errors.New("something wrong")
	}

	return permissions, nil
}

func (company *Company)	IsPublic(db *sql.DB) bool{
	var pub	bool
	err := db.QueryRow(`SELECT pub FROM company WHERE id = $1`, company.Id).Scan(&pub)
	if err != nil {
		return false
	}
	return pub
}

func (company *Company) Docs(permissions bool, db *sql.DB)([]Document, error){
	//var rows *sql.Row
	//var err	 error
	//if permissions{
	//	rows, err = db.Query(`SELECT `)
	//}
	return nil, nil
}





