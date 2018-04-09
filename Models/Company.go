package Models

import (
	"database/sql"
	"log"
	"errors"
	"encoding/json"
	"crypto/sha256"
	"time"
	"encoding/base64"
)

type Company struct {
	Id 					int 				`json:"id,omitempty"`
	Name 				string 				`json:"name,omitempty"`
	Description			string  			`json:"description,omitempty"`
	Public				bool				`json:"public,omitempty"`
	Country				string				`json:"country,omitempty"`
	Meta				*json.RawMessage	`json:"meta"`
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
	var rows *sql.Rows
	var err	 error

	if permissions{
		rows, err = db.Query(`SELECT id, name, description FROM document
									JOIN company_doc cd ON document.id = cd.doc_id
									WHERE cd.company_id = $1`, company.Id)
	}else {
		rows, err = db.Query(`SELECT id, name, description FROM document
									JOIN company_doc cd ON document.id = cd.doc_id
									WHERE cd.company_id = $1 AND public = TRUE`, company.Id)
	}
	if err != nil {
		return nil, errors.New("no docs")
	}

	defer rows.Close()

	documents := make([]Document, 0)
	for rows.Next(){
		doc := new(Document)
		err = rows.Scan(&doc.ID, &doc.Name, &doc.Description)
		if err != nil {
			log.Println("Models.Company.Docs ", err)
			return nil, err
		}

		documents = append(documents, *doc)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Company.Docs ", err)
		return nil, err
	}
	return documents, nil
}

func SearchCompany(query *string, userId int, db *sql.DB) ([]Company, error){
	rows, err := db.Query(`SELECT id, name, description FROM company
							 	 	WHERE pub = TRUE AND lower(name) %> lower($1)
								 UNION DISTINCT
								 SELECT id, name, description FROM company
  									JOIN client_company cc ON cc.company_id = company.id
									WHERE cc.client_id = $2 AND lower(name) %> lower($1)`, query, userId)
	if err != nil {
		log.Println("Models.Company.SearchCompany", err)
		return nil, err
	}

	defer rows.Close()

	companies := make([]Company, 0)
	for rows.Next(){
		company := new(Company)
		err := rows.Scan(&company.Id, &company.Name, &company.Description)
		if err != nil {
			log.Println("Models.Company.SearchCompany", err)
			return nil, err
		}
		companies = append(companies, *company)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Company.SearchCompany", err)
		return nil, err
	}

	return companies, nil
}

func GetCountryMeta(countryId int, db *sql.DB)(*json.RawMessage, error){
	var result *json.RawMessage
	err := db.QueryRow(`SELECT fields FROM company_meta WHERE country_id = $1`, countryId).Scan(&result)
	if err != nil {
		return nil, errors.New("something wrong")
	}
	return result, nil
}

func PrepareRegMessage(clientId, countryId int, name, description *string, isPub bool, metaInfo json.RawMessage, tx *sql.Tx) (*string, error) {
	var companyId int
	err := tx.QueryRow(`INSERT INTO company(name, description, country_id, pub, meta) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
							name, description, countryId, isPub, metaInfo).Scan(&companyId)
	if err != nil {
		log.Println("Models.Company.PrepareRegMessage insert company", err)
		return nil, errors.New("something wrong")
	}
	_, err = tx.Exec(`INSERT INTO client_company(client_id, company_id, gr_admin, gr_invite, gr_kick,
							gr_read, gr_write, gr_update, gr_delete, responsible, client_confirm, company_confirm) VALUES (
							$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
							clientId, companyId, true, true, true, true, true, true, true, true, true, true)
	if err != nil {
		log.Println("Models.Company.PrepareRegMessage insert client_company ", err)
		return nil, errors.New("something wrong")
	}

	hash := sha256.New()
	curTime := time.Now()
	hash.Write([]byte(*name + *description + curTime.String()))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return &sha, nil
}

func ConfirmCompany(sha *string, db *sql.DB) error{
	_, err := db.Exec("UPDATE company SET confirm = $1, created = $2", true, time.Now())
	if err != nil {
		log.Println("Models.Company.ConfirmCompany ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (company *Company) GetByDocument(docId int, db *sql.DB) error{
	err := db.QueryRow(`SELECT cm.id, cm.name, cm.description FROM company cm
  							  JOIN company_doc cd on cm.id = cd.company_id
  							  WHERE cd.doc_id = $1 AND cm.pub = TRUE AND cm.confirm = TRUE;`,
		docId).Scan(&company.Id, &company.Name, &company.Description)
	if err != nil {
		//log.Println("Models.Company.Get ", err)
		return errors.New("something wrong")
	}
	return nil
}



