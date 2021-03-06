package Models

import (
	"database/sql"
	"log"
	"errors"
	"crypto/sha256"
	"time"
	"encoding/base64"
	"docnota/Utils"
	"strconv"
)

type User struct {
	ID			int	   `json:"id,omitempty"`
	Avatar		string `json:"avatar,omitempty"`
	Email 		string `json:"email,omitempty"`
	Password 	string `json:"password,omitempty"`
	FirstName	string `json:"first_name,omitempty"`
	LastName	string `json:"last_name,omitempty"`
	Public		bool   `json:"public,omitempty"`
	Verified	bool   `json:"verified,omitempty"`
	Country		string `json:"country"`
}

type UserInteraction interface{
	Create(tx *sql.DB) error
	CopyDocument(docId int, tx *sql.Tx) error
	AcceptDoc(docId, db sql.DB) error
	CreateUidForEmail(tx *sql.Tx) (string, error)
	Confirm(tx *sql.Tx) error

	SendCompanyRequest(companyId int, db *sql.DB) error

	Get(isOwner bool, db *sql.DB) error
	GetUserCompany(isOwner bool, db *sql.DB) ([]Company, error)
	GetByEmail(db *sql.DB) error
	GetDocuments(isOwner bool, isTemplate bool, db *sql.DB) ([]Document, error)
	InboxDocuments(db *sql.DB) ([]Document, error)

	UpdatePassword(passw *string, db *sql.DB) error
	UpdateAvatar(path string, db *sql.DB) error

	CheckPassword(db *sql.DB) bool
	Exist(db *sql.DB) bool

	Delete(db *sql.DB) error
	RejectDocument(docId int, db *sql.DB) error
}

func (user *User) Create(tx *sql.Tx) error{
	hashedPassword := Utils.HashAndSalt([]byte(user.Password))
	err := tx.QueryRow(`INSERT INTO client(email, pass) VALUES ($1, $2) RETURNING id`, user.Email, hashedPassword).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return errors.New("can't create user")
	}
	return nil
}

func (user *User) Exist(db *sql.DB) bool{
	err := db.QueryRow(`SELECT email FROM client WHERE lower(email) = lower($1)`, user.Email).Scan(&user.Email)
	if err != nil {
		return false
	}
	return true
}

func (user *User) Get(isOwner bool, db *sql.DB) error{
	var err 	error
	var confirm	bool
	var avatar 	sql.NullString

	if isOwner{
		err = db.QueryRow(` SELECT email, first_name, last_name, verified, pub, country.name, avatar, confirmed
								 FROM client
								 JOIN country ON country.id = client.country_id 
								 WHERE client.id = $1`, user.ID).Scan(&user.Email, &user.FirstName, &user.LastName, &user.Verified,
								 &user.Public, &user.Country, &avatar, &confirm)
		if avatar.Valid{
			user.Avatar = avatar.String
		}
		if err != nil {
			log.Println(err)
			return errors.New("user not exists")
		}
		if !confirm{
			return errors.New("account not confirmed")
		}
	}else {
		err = db.QueryRow(" SELECT email, first_name, last_name, verified, pub, country, avatar, confirmed " +
								" FROM client " +
								" JOIN country ON country.id = client.country_id " +
								" WHERE client.id = $1 AND " +
								" client.pub = TRUE AND " +
								" client.confirmed = TRUE", user.ID).Scan(&user.Email, &user.FirstName, &user.LastName, &user.Verified,
								&user.Public, &user.Country, &avatar, &confirm)
		if err != nil {
			return errors.New("user not exists or access denied")
		}

		if avatar.Valid{
			user.Avatar = avatar.String
		}
	}
	return nil
}

func GetUsers(db *sql.DB) ([]User, error){
	users := make([]User, 0)
	rows, err := db.Query("SELECT id, email, first_name, last_name, verified FROM client WHERE pub = TRUE AND confirmed = TRUE")
	defer rows.Close()
	if err != nil{
		log.Println("Model.User.GetUsers ", err)
		return nil, errors.New("something wrong")
	}
	for rows.Next(){
		user := new(User)
		err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Verified)
		if err != nil{
			log.Println("Model.User.GetUsers ", err)
			return nil, errors.New("something wrong")
		}
		users = append(users, *user)
	}
	if err = rows.Err(); err != nil{
		log.Println("Model.User.GetUsers ", err)
		return nil, errors.New("something wrong")
	}
	return users, err
}

func (user *User) GetByEmail(db *sql.DB) error{
	err := db.QueryRow("SELECT id FROM client WHERE email = $1 AND confirmed = TRUE", user.Email).Scan(&user.ID)
	if err != nil{
		return errors.New("user not exists or check email for confirm your account")
	}
	return nil
}

func (user *User) CreateUidForEmail(tx *sql.Tx) (*string, error){
	hash := sha256.New()
	curTime := time.Now()
	hash.Write([]byte(user.Email + user.Password + curTime.String()))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	dateExp := curTime.Add(time.Hour * 24 * 365)

	_, err := tx.Exec("INSERT INTO client_confirm(client_id, uid, date_exp) VALUES ($1, $2, $3)", user.ID, sha, dateExp)
	if err != nil{
		log.Println("Model.User.CreateEmailContent, ", err)
		tx.Rollback()
		return nil, errors.New("something wrong")
	}
	return &sha, nil
}


func (user *User) Confirm(uid *string, tx *sql.Tx) error{
	var dateExp 	time.Time
	var CountryID	int

	curTime := time.Now()

	err := tx.QueryRow("SELECT date_exp, client_id FROM client_confirm WHERE uid = $1 AND date_confirm is NULL ",
								uid).Scan(&dateExp, &user.ID)

	if err != nil{
		log.Println("model.User.Confirm select client_id by uid", err)
		return errors.New("something wrong")
	}

	if curTime.After(dateExp){
		log.Println("Rotten link")

		return errors.New("rotten link")
	}

	CountryID, _ = strconv.Atoi(user.Country)


	_, err = tx.Exec(" UPDATE client SET confirmed = TRUE, first_name = $1, last_name = $2, pub = $3," +
							" country_id = $4 WHERE id = $5", user.FirstName, user.LastName, user.Public, CountryID, user.ID)
	if err != nil{
		tx.Rollback()
		log.Println("Model.User.Confirm update client, ", err)
		return errors.New("something wrong")
	}

	_, err = tx.Exec("UPDATE client_confirm SET date_confirm = $1 WHERE client_id = $2", curTime, user.ID)
	if err != nil{
		tx.Rollback()
		log.Println("Model.User.Confirm update client_confirm, ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (user *User) Delete(db *sql.DB) error{
	_, err := db.Exec("DELETE FROM client WHERE id = $1", user.ID)
	if err != nil {
		return errors.New("can't delete user")
	}
	return nil
}

func GetPublicUsers(db *sql.DB) ([]User, error){
	users := make([]User, 0)
	rows, err := db.Query("SELECT id, email, first_name, last_name, verified, avatar FROM client WHERE pub = TRUE AND confirmed = TRUE")
	defer rows.Close()
	if err != nil{
		log.Println("Model.User.GetPublicUsers, ", err)
		return nil, errors.New("something wrong")
	}
	for rows.Next(){
		user := new(User)
		err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Verified, &user.Avatar)
		if err != nil{
			log.Println("Model.User.GetPublicUsers, ", err)
			return nil, err
		}
		users = append(users, *user)
	}
	if err = rows.Err(); err != nil{
		log.Println("Model.User.GetPublicUsers error after rows.next(), ", err)
		return nil, err
	}
	return users, err
}

func (user *User) GetUserCompany(isOwner bool, db *sql.DB) ([]Company, error){
	var	rows	*sql.Rows
	var	err		error
	if isOwner{
		rows, err = db.Query("  SELECT cm.id, cm.name, cm.description FROM company cm " +
									" JOIN client_company cc ON cm.id = cc.company_id " +
									" WHERE client_id = $1", user.ID)
		if err != nil {
			log.Println("Model.User.GetUserCompany select for owner ", err)
			return nil, errors.New("something wrong")
		}
	}else {
		rows, err = db.Query(" SELECT cm.id, cm.name, cm.description FROM company cm " +
									" JOIN client_company cc ON cm.id = cc.company_id " +
									" JOIN client cl ON cl.id = cc.client_id " +
									" WHERE cl.id = $1 AND cl.pub = TRUE ", user.ID)
		if err != nil {
			log.Println("Model.User.GetUserCompany select for non onwner ", err)
			return nil, errors.New("something wrong")
		}
	}
	defer rows.Close()

	companies := make([]Company, 0)
	for rows.Next(){
		company := new(Company)
		err = rows.Scan(&company.Id, &company.Name, &company.Description)
		if err != nil{
			log.Println("Model.User.GetUserCompany ", err)
			return nil, errors.New("something wrong")
		}
		companies = append(companies, *company)
	}
	if len(companies) == 0{
		return nil, nil
	}
	if err = rows.Err(); err != nil{
		log.Println("Model.User.GetUserCompany ", err)
		return nil, errors.New("something wrong")
	}
	return companies, nil
}

func (user *User) CopyDocument(docId int, tx *sql.Tx) error{
	var pub bool
	err := tx.QueryRow("SELECT public FROM document WHERE id = $1", docId).Scan(&pub)
	if err != nil{
		log.Println("Model.User.CopyDocument ", err)
		return errors.New("something wrong")
	}
	if !pub{
		err = tx.QueryRow(" SELECT c4.gr_read FROM company_doc JOIN company c2 ON company_doc.company_id = c2.id " +
								" JOIN client_company c4 ON c2.id = c4.company_id " +
								" WHERE doc_id = $1 AND c4.client_id = $2", docId, user.ID).Scan(&pub)
		if err != nil{
			log.Println("Model.User.CopyDocument ", err)
			return errors.New("document is not accept")
		}
		if pub == false{
			return errors.New("document is not accept")
		}
	}
	_, err = tx.Exec("SELECT * FROM copy_doc($1, $2)", docId, user.ID)
	if err != nil{
		tx.Rollback()
		log.Println("Model.User.CopyDocument ", err)
		return errors.New("something wrong")
	}
	return nil
}


func (user *User) UpdatePassword(passw string, db *sql.DB) error{
	hashedPassword := Utils.HashAndSalt([]byte(passw))
	_, err := db.Exec("UPDATE client SET pass = $1 WHERE id = $2", hashedPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func SearchUsers(query string, db *sql.DB) ([]User, error){
	rows, err := db.Query(` SELECT id, first_name, last_name, email
  								  FROM client WHERE lower(client.first_name || client.last_name || client.email) %> 
								  lower($1) AND pub = TRUE`, query)
	if err != nil {
		log.Println("Model.User.Search ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next(){
		var tmpUser User
		err := rows.Scan(&tmpUser.ID, &tmpUser.FirstName, &tmpUser.LastName, &tmpUser.Email)
		if err != nil {
			log.Println("Model.User.Search ", err)
			return nil, errors.New("something wrong")		}
		users = append(users, tmpUser)
	}

	if err = rows.Err(); err != nil{
		log.Println("Model.User.Search ", err)
		return nil, errors.New("something wrong")
	}
	return users, nil
}

func (user *User) UpdateAvatar(path string, db *sql.DB) error{
	_, err := db.Exec("UPDATE client SET avatar = $1 WHERE id = $2", path, user.ID)
	if err != nil{
		log.Println("Model.User.UpdateAvatar ", err)
		return errors.New("can't update avatar")
	}
	return nil
}

func (user *User) SendCompanyRequest(companyId int, db *sql.DB) error{
	_, err := db.Exec(" INSERT INTO client_company(client_id, company_id, confirm, company_conf, gr_read) " +
							" VALUES($1, $2, FALSE, FALSE, TRUE)", user.ID, companyId)
	if err != nil {
		log.Println("Model.User.SendCompanyRequest ", err)
		return errors.New("can't send request")
	}
	return nil
}

func (user *User) RejectDocument(docId int, db *sql.DB) error{
	_, err := db.Exec("DELETE FROM recieve_document WHERE client_id = $2 AND document = $1", docId, user.ID)
	if err != nil{
		log.Println("Model.User.RejectDocument ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (user *User) CheckPassword(db *sql.DB) bool{
	var tmpPwd *string
	err := db.QueryRow("SELECT id, pass FROM client WHERE email = $1", user.Email).Scan(&user.ID, &tmpPwd)
	if err != nil{
		return false
	}
	if Utils.ComparePassword(*tmpPwd, []byte(user.Password)){
		return true
	}
	return false
}

func (user *User) AcceptDoc(docId int, tx *sql.Tx) error{
	_, err := tx.Exec("SELECT * FROM copy_doc($1, $2)", docId, user.ID)
	if err != nil{
		log.Println("Model.User.AcceptDoc ", err)
		return errors.New("something wrong")
	}
	_, err = tx.Exec("DELETE FROM receive_document WHERE client_id = $2 AND document_id = $1", docId, user.ID)
	if err != nil{
		log.Println("Model.User.AcceptDoc ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (user *User) InboxDocuments(db *sql.DB) ([]Document, error){
	var description sql.NullString
	rows, err := db.Query(`SELECT id, name, rd.client_id, public, template, last_updated, created, description
								 FROM document JOIN receive_document rd ON document.id = rd.document_id
								 WHERE rd.client_id = $1`, user.ID)
	if err != nil {
		log.Println("Model.User.GetInboxDocument ", err)
		return nil, errors.New("something wrong")
	}
	docs := make([]Document, 0)
	for rows.Next(){
		doc := new(Document)
		err = rows.Scan(&doc.ID, &doc.Name, &doc.UserId, &doc.Public, &doc.Template, &doc.LastUpdated, &doc.Created, &description)
		if description.Valid{
			doc.Description = description.String
		}
		if err != nil {
			log.Println("Model.User.GetInboxDocument ", err)
			return nil, errors.New("something wrong")
		}
		docs = append(docs, *doc)
	}
	if err = rows.Err(); err != nil{
		log.Println("Model.User.GetInboxDocument ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()
	return docs, nil
}

func (user *User) GetDocuments(isOwner bool, isTemplate bool, db *sql.DB) ([]Document, error){
	var rows	*sql.Rows
	var err		error

	if isOwner{
		if isTemplate{
			rows, err = db.Query("SELECT id, name, last_updated, created FROM document WHERE client_id = $1" +
				" AND template = TRUE", user.ID)
		}else {
			rows, err = db.Query("SELECT id, name, last_updated, created FROM document WHERE client_id = $1"+
				" AND template = FALSE", user.ID)
		}
		if err != nil {
			log.Println("Models.User.GetDocuments ", err)
			return nil, errors.New("something wrong")
		}
	}else {
		if isTemplate{
			rows, err = db.Query("SELECT id, name, last_updated, created FROM document " +
				" WHERE client_id = $1 AND public = TRUE AND template = TRUE", user.ID)
		}else {
			rows, err = db.Query("SELECT id, name, last_updated, created FROM document "+
				" WHERE client_id = $1 AND public = TRUE AND template = FALSE", user.ID)

		}
	}
	defer rows.Close()

	docs := make([]Document, 0)

	for rows.Next(){
		doc := new(Document)
		err = rows.Scan(&doc.ID, &doc.Name, &doc.LastUpdated, &doc.Created)
		if err != nil{
			return nil, err
		}
		docs = append(docs, *doc)
	}
	if len(docs) == 0{
		return nil, errors.New("no docs")
	}
	if 	err = rows.Err(); err != nil{
		log.Println("Models.User.GetDocuments ", err)
		return nil, errors.New("something wrong")
	}
	return docs, nil
}

func (user *User) HasPermissionToBlock(blockId int, db *sql.DB) bool{
	err := db.QueryRow(`SELECT id FROM block
 							  JOIN document ON document.id = block.doc_id
 							  JOIN client ON client.id = document.client_id
 							  WHERE parent_id = $1 AND client.id = $2`, blockId, user.ID).Scan(&blockId)
	if err != nil {
		log.Println("Midels.User.HasPermissionToBlock ", err)
		return false
	}
	return true
}

func (user *User) UploadAvatar(fileName *string, db *sql.DB) error{
	_, err := db.Exec(`UPDATE client SET avatar = $1 WHERE id = $2`, fileName, user.ID)
	if err != nil {
		log.Println("Models.User.UploadAvatar ", err)
		return errors.New("something wrong")
	}
	user.Avatar = *fileName
	return nil
}