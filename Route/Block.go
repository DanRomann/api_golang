package Route

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"errors"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
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

