package Route

import (
	"net/http"
	"io/ioutil"
	"errors"
	"docnota/Usecases"
	"docnota/Utils"
	"docnota/Models"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"strconv"
)

func createUser(w http.ResponseWriter, r *http.Request){
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	defer r.Body.Close()

	user := new(Models.User)
	err = json.Unmarshal(data, &user)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}


	err = Usecases.CreateUser(user, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func confirmUser(w http.ResponseWriter, r *http.Request){
	var curRequest struct{
		UID			string		`json:"uid"`
		User		Models.User	`json:"user"`
	}

	curResponse := TokenResponse{}

	data, err := ioutil.ReadAll(r.Body)
	log.Println(string(data))
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	err = json.Unmarshal(data, &curRequest)

	token, err := Usecases.ConfirmUser(&curRequest.User, &curRequest.UID, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	curResponse.Token = *token

	result, _ := json.Marshal(curResponse)
	DataResponse(result, w)

}

func getUserDocs(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		ErrResponse(errors.New("bad userId"), w)
	}

	documents, err := Usecases.GetUserDocuments(userId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(documents)
	DataResponse(result, w)
}

func authUser(w http.ResponseWriter, r *http.Request){
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	defer r.Body.Close()

	user := new(Models.User)
	err = json.Unmarshal(data, &user)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	token, err := Usecases.Auth(user, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	w.Header().Set("Authorization", *token)
	SuccessResponse("ok", w)
}