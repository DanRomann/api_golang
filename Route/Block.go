package Route

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
	"html/template"
	"docnota/Models"
)

func getBlock(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	blockId, err := strconv.Atoi(vars["blockId"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	block, err := Usecases.GetBlockById(blockId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(block)
	DataResponse(result, w)
}

func getBlockRelations(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	blockId, err := strconv.Atoi(vars["blockId"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	relations, err := Usecases.GetBlockRelations(blockId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(relations)
	DataResponse(result, w)
}

func addRelation(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	blockId, err := strconv.Atoi(vars["blockId"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	relationBlock, err := strconv.Atoi(vars["relationBlock"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	err = Usecases.AddBlockRelation(blockId, relationBlock, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func deleteRelation(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	blockId, err := strconv.Atoi(vars["blockId"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	relationBlock, err := strconv.Atoi(vars["relationBlock"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	err = Usecases.DeleteRelation(blockId, relationBlock, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func getBlockIFrame(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	tmpl, err := template.ParseFiles("Template/block-iframe.html")
	if err != nil {
		ErrResponse(errors.New("template not found"), w)
		return
	}

	vars := mux.Vars(r)
	blockId, err := strconv.Atoi(vars["blockId"])
	if err != nil {
		ErrResponse(errors.New("bad block id"), w)
		return
	}

	block, err := Usecases.GetBlockById(blockId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	document, err := Usecases.GetDocument(block.DocId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	type Iframe struct {
		Document Models.Document
		Block Models.Block
	}

	iframe := Iframe{Document:*document, Block:*block}

	tmpl.Execute(w, iframe)
}
