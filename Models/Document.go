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
	Description	string	  `json:"description"`
	Blocks      []Block   `json:"blocks,omitempty"`
	Template    bool      `json:"template"`
	Public      bool      `json:"public"`
	ReadOnly	bool	  `json:"read_only,omitempty"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}

type DocumentInteraction interface{
	BelongToUser(userId int, db *sql.DB) bool
	BelongToUserOrPublic(userId int, db *sql.DB) bool
	BelongToGroup(groupId int, db *sql.DB) bool
	IsPublic(db *sql.DB) bool

	Create(userId int, name *string, isTemplate bool, isPublic bool, db *sql.DB) (*Document, error)
	Copy()
	SaveFillTemplate(userId int, tx *sql.Tx) error

	SendDocumentToUser(userId int, db *sql.DB) error


	Get(db *sql.DB) error
	GetOwner(db *sql.DB) error
	GetDocumentHistory(db *sql.DB) ([]Block, error)
	Search(query string, db *sql.DB) ([]Document, error)

	Delete(db *sql.DB) error

	Sort() error
}

func (doc *Document) BelongToUser(userId int, db *sql.DB) bool{
	var name string
	err := db.QueryRow(`SELECT name FROM document WHERE id = $1 AND client_id = $2`, doc.ID, userId).Scan(&name)
	if err != nil {
		return false
	}
	return true
}

func (doc *Document) BelongToUserOrPublic(userId int, db *sql.DB) bool{
	err := db.QueryRow(`SELECT name FROM document WHERE id = $1 AND (client_id = $2 OR public = TRUE)`,
		doc.ID, userId).Scan(&doc.Name)
	if err != nil{
		return false
	}
	return true
}

func (doc *Document) IsPublic(db *sql.DB) bool{
	err := db.QueryRow(`SELECT name FROM document WHERE id = $1 AND public = TRUE`, doc.ID).Scan(&doc.Name)
	if err != nil {
		return false
	}
	return true
}

func (doc *Document) Create(db *sql.DB) error{
	curTime := time.Now()
	err := db.QueryRow(`INSERT INTO document(name, description, client_id, template, public, last_updated,
 							  created) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, doc.Name, doc.Description,
 							  			doc.UserId, doc.Template, doc.Public, curTime, curTime).Scan(&doc.ID)
	if err != nil {
		log.Println("Models.Document.Create ", err)
		return errors.New("something wrong")
	}
	return nil
}

func PublicDocuments(isTemplate bool, db *sql.DB) ([]Document, error){
	var rows		*sql.Rows
	var err			error

	if isTemplate{
		rows, err = db.Query(`SELECT id, name, template, last_updated, created FROM document WHERE public = TRUE
									AND template = TRUE`)
	}else {
		rows, err = db.Query(`SELECT id, name, template, last_updated, created FROM document WHERE public = TRUE
									AND template = FALSE`)
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
	var description sql.NullString
	rows, err := db.Query(`
								WITH RECURSIVE tree AS(
							  		SELECT
    								ARRAY[]::INTEGER[] || b.ord as path,
									b.id,
									b.parent_id,
									b.name,
									b.content,
									b.ord,
									b.last_updated,
                					(SELECT count(relation_block) FROM block_relation br WHERE br.block_id = b.id) as br
									FROM block b
										WHERE b.doc_id = $1 AND b.parent_id IS NULL AND b.last_date IS NULL

							  		UNION ALL

							  		SELECT
    								path || b.ord,
									b.id,
									b.parent_id,
									b.name,
									b.content,
									b.ord,
									b.last_updated,
                					(SELECT count(relation_block) FROM block_relation br WHERE br.block_id = b.id AND b.last_date IS NULL) as br
									FROM block b
										JOIN tree ON tree.id = b.parent_id
									) SELECT id, parent_id, name, content, ord, last_updated, br, path FROM tree ORDER BY path;



`, doc.ID)
	if err != nil {
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}
	defer rows.Close()
	blocks := make([]Block, 0)
	for rows.Next() {
		var parentId 		sql.NullInt64
		var blockContent 	sql.NullString

		block := new(Block)
		err = rows.Scan(&block.Id, &parentId, &block.Name, &blockContent, &block.Order, &block.LastUpdated, &block.RelationsCount, &block.Ltree)
		block.DocId = doc.ID
		if err != nil {
			log.Println("Models.Document.Get ", err)
			return errors.New("something wrong")
		}


		if parentId.Valid{
			block.ParentID = int(parentId.Int64)
		}
		if blockContent.Valid{
			block.Content = blockContent.String
		}
		blocks = append(blocks, *block)
	}

	if err = rows.Err(); err != nil{
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}

	err = db.QueryRow(`SELECT name, description, template, last_updated, created, public FROM document WHERE id = $1`,
							doc.ID).Scan(&doc.Name, &description, &doc.Template, &doc.LastUpdated, &doc.Created, &doc.Public)
	if description.Valid{
		doc.Description = description.String
	}

	doc.Blocks = blocks
	if err != nil{
		log.Println("Models.Document.Get ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (doc *Document) Commit(db *sql.DB) error  {
	return nil
}

func (doc *Document) GetOwner(db *sql.DB) error{
	err := db.QueryRow(`SELECT client.id FROM client 
							  	JOIN document doc ON doc.client_id = client.id
							  		WHERE doc.id = $1`, doc.ID).Scan(&doc.UserId)
	if err != nil {
		return errors.New("something wrong")
	}
	return nil
}

func (doc *Document) Copy(userId int, tx *sql.Tx) error{
	_, err := tx.Exec(`SELECT copy_doc($1, $2)`, doc.ID, userId)
	if err != nil {
		log.Println("Models.Document.Copy ", err)
		return errors.New("something wrong")
	}
	return nil
}

func (doc *Document) SendDocumentToUser(userId int, db *sql.DB) error{
	_, err := db.Exec("INSERT INTO receive_document (client_id, document_id) VALUES ($1, $2)", userId, doc.ID)
	if err != nil{
		log.Println("Models.Document.SendDocumentToUser ", err)
		return errors.New("something wrong")
	}
	return nil
}

func SearchDocs(query *string, db *sql.DB) ([]Document, error){
	rows, err := db.Query(`SELECT id, name, template, last_updated, created 
								  FROM document WHERE lower(name) %> lower($1) AND public = TRUE LIMIT 200 `, query)
	if err != nil {
		log.Println("Models.Document.SearchDocs ", err)
		return nil, errors.New("something wrong")
	}
	defer rows.Close()

	docs := make([]Document, 0)

	for rows.Next(){
		curDoc := Document{}
		err := rows.Scan(&curDoc.ID, &curDoc.Name, &curDoc.Template, &curDoc.LastUpdated, &curDoc.Created)
		if err != nil{
			log.Println("Models.Document.Search ", err)
			return nil, errors.New("something wrong")
		}
		docs = append(docs, curDoc)
	}
	if err = rows.Err(); err != nil{
		log.Println("Search block ", err)
		return nil, errors.New("something wrong")
	}
	return docs, nil
}


func (doc *Document) SaveFillTemplate(userId int, tx *sql.Tx) error{
	dict := make(map[int]int)
	err := tx.QueryRow(`INSERT INTO document(name, client_id, last_updated, created) VALUES($1, $2, $3, $4)
 								RETURNING id`, doc.Name, userId, time.Now(), time.Now()).Scan(&doc.ID)
	if err != nil {
		tx.Rollback()
		log.Println("Modesl.Document.SaveFillTemplate can't insert doc ", err)
		return errors.New("something wrong")
	}
	for _, block := range doc.Blocks{
		var newId int
		err := tx.QueryRow(`INSERT INTO block(name, content, date, client_id) VALUES($1, $2, $3, $4) RETURNING id`,
									block.Name,	block.Content, time.Now(), userId).Scan(&newId)
		dict[block.Id] = newId
		if err != nil {
			tx.Rollback()
			log.Println("Models.Document.SaveFillTemplate can't insert block ", err)
			return errors.New("something wrong")
		}
	}
	for _, block := range doc.Blocks{
		var curParent, curId	int
		var err 				error

		curParent = dict[block.ParentID]
		if curParent != 0 {
			_, err = tx.Exec(`UPDATE block SET parent_id = $1 WHERE id = $2`, curParent, curId)
		}

		if err != nil {
			log.Printf("%+v\n", block)
			tx.Rollback()
			log.Println("Model.Document.SaveFillTemplate can't insert doc_block ", err)
			return errors.New("something wrong")
		}
	}
	return nil
}

func (doc *Document) Edit(db *sql.DB) error{
	_, err := db.Exec(`UPDATE document SET name = $1, description = $2, public = $3, last_updated = $4, template = $5
							 WHERE id = $6`, doc.Name, doc.Description, doc.Public, time.Now(), doc.Template, doc.ID)
	if err != nil {
		log.Println("Models.Document.Edit ", err)
		return errors.New("something wrong")
	}
	return nil
}
