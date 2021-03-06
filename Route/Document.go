package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"io/ioutil"
	"docnota/Models"
	"html/template"
	"strings"
)

func getPublicDocs(w http.ResponseWriter, r *http.Request){
	documents, err := Usecases.GetPublicDocuments(false, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}
	result, _ := json.Marshal(documents)
	DataResponse(result, w)
}

func getPublicTemplates(w http.ResponseWriter, r *http.Request){
	documents, err := Usecases.GetPublicDocuments(true, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}
	result, _ := json.Marshal(documents)
	DataResponse(result, w)
}

func getDoc(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)
	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad doc id"), w)
		return
	}

	document, err := Usecases.GetDocument(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(document)
	DataResponse(result, w)
}

func copyDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)
	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad doc id"), w)
		return
	}

	err = Usecases.CopyDocument(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}
	SuccessResponse("ok", w)
}

func createDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	doc := new(Models.Document)
	err = json.Unmarshal(data, &doc)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	doc, err = Usecases.CreateDocument(doc, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(doc)
	DataResponse(result, w)
}

func commitDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	defer r.Body.Close()

	block := new(Models.Block)
	err = json.Unmarshal(data, &block)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}


	block, err = Usecases.ChangeDoc(block, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(block)
	DataResponse(result, w)
}

func fillTemplate(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	defer r.Body.Close()

	doc := new(Models.Document)
	err = json.Unmarshal(data, &doc)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	err = Usecases.FillTemplate(doc, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func metaEdit(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	defer r.Body.Close()

	doc := new(Models.Document)
	err = json.Unmarshal(data, &doc)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	err = Usecases.DocMetaEdit(doc, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func blockChainUploadDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad doc id"), w)
		return
	}

	err = Usecases.BlockChainUpload(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
	}
	SuccessResponse("ok", w)
}

func getDocIFrame(w http.ResponseWriter, r *http.Request){
	token := ""
	vars := mux.Vars(r)
	tmpl, err := template.ParseFiles("Templates/fb-iframe.html")
	if err != nil {
		ErrResponse(errors.New("template not found"), w)
		return
	}
	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad doc id"), w)
		return
	}

	document, err := Usecases.GetDocument(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	type Iframe struct {
		Document Models.Document
		User Models.User
		Company Models.Company
	}

	usr, err := Usecases.GetUser(document.UserId, &token, Utils.Connect)
	if err != nil {
		println("User not find")
		//ErrResponse(err, w)
		//return
		usr = new(Models.User)
	}

	company, err := Usecases.GetCompanyByDocument(document.ID, &token, Utils.Connect)
	if err != nil {
		println("company not find")
		//ErrResponse(err, w)
		//return
		company = new(Models.Company)
	}

	iframe := Iframe{Document:*document, User:*usr, Company:*company}

	//result, _ := json.Marshal(document)
	//DataResponse(result, w)
	tmpl.Execute(w, iframe)
}

func getFullDocIFrame(w http.ResponseWriter, r *http.Request){
	token := ""
	vars := mux.Vars(r)
	tmpl, err := template.ParseFiles("Templates/fb-iframe--full.html")
	if err != nil {
		ErrResponse(errors.New("template not found"), w)
		return
	}
	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad doc id"), w)
		return
	}

	document, err := Usecases.GetDocument(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	type Block struct{
		Name string
		Content template.HTML
		Ltree string
	}

	type Iframe struct {
		Document Models.Document
		User Models.User
		Company Models.Company
		Blocks []Block
	}

	usr, err := Usecases.GetUser(document.UserId, &token, Utils.Connect)
	if err != nil {
		println("User not find")
		//ErrResponse(err, w)
		//return
		usr = new(Models.User)
	}

	company, err := Usecases.GetCompanyByDocument(document.ID, &token, Utils.Connect)
	if err != nil {
		println("company not find")
		//ErrResponse(err, w)
		//return
		company = new(Models.Company)
	}

	blocks := make([]Block, 0)

	for _, b := range document.Blocks{
		ltree := strings.Split(b.Ltree[1:len(b.Ltree)-1], ",")
		block := Block{Name:b.Name, Content:template.HTML(b.Content), Ltree:strings.Join(ltree, ".")}
		blocks = append(blocks, block)
	}

	iframe := Iframe{Document:*document, User:*usr, Company:*company, Blocks:blocks}

	//result, _ := json.Marshal(document)
	//DataResponse(result, w)
	tmpl.Execute(w, iframe)
}