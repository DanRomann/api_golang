package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
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

//func commitDoc(w http.ResponseWriter, r *http.Request){
//	data, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		ErrResponse(errors.New("bad data"), w)
//		return
//	}
//	defer r.Body.Close()
//
//	document := new(Models.Document)
//	err = json.Unmarshal(data, &document)
//	if err != nil {
//		ErrResponse(errors.New("bad json"), w)
//		return
//	}
//
//	token := r.Header.Get("Authorization")
//
//
//}