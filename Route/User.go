package Route

import (
	"net/http"
	"io/ioutil"
	"errors"
	"docnota/Usecases"
	"docnota/Utils"
	"docnota/Models"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"log"
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

	data, err := ioutil.ReadAll(r.Body)
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

	w.Header().Set("Authorization", *token)
	SuccessResponse("ok", w)
}

func getUserDocs(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		ErrResponse(errors.New("bad userId"), w)
		return
	}

	documents, err := Usecases.GetUserDocuments(userId, false, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(documents)
	DataResponse(result, w)
}

func getUserTemplate(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")


	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		ErrResponse(errors.New("bad userId"), w)
		return
	}

	documents, err := Usecases.GetUserDocuments(userId, true, &token, Utils.Connect)
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

func getUser(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)

	requestUser, err := strconv.Atoi(vars["userId"])
	if err != nil {
		ErrResponse(errors.New("bad user id"), w)
		return
	}

	user, err := Usecases.GetUser(requestUser, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(user)

	DataResponse(result, w)
}

func getUsers(w http.ResponseWriter, r *http.Request){
	users, err := Usecases.GetUsers(Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(users)
	DataResponse(result, w)

}

func getUserInbox(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)

	requestUser, err := strconv.Atoi(vars["userId"])
	if err != nil {
		ErrResponse(errors.New("bad user id"), w)
		return
	}

	documents, err := Usecases.GetUserInboxDocument(requestUser, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(documents)
	DataResponse(result, w)
}

func sendDocToUser(w http.ResponseWriter, r *http.Request){
	var curRequest struct{
		UserId	int		`json:"user_id"`
		DocId	int		`json:"doc_id"`
	}

	token := r.Header.Get("Authorization")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("something wrong"), w)
		return
	}

	err = json.Unmarshal(data, &curRequest)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	err = Usecases.SendDocumentToUser(curRequest.DocId, curRequest.UserId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func acceptDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)

	docId, err := strconv.Atoi(vars["docId"])
	if err != nil {
		ErrResponse(errors.New("bad docId"), w)
		return
	}

	err = Usecases.AcceptDoc(docId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)

}

func uploadUserAvatar(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	file, handler, err := r.FormFile("photo")
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}
	log.Println(handler.Filename)
	defer file.Close()

	fileContent := make([]byte, handler.Size)
	_, err = file.Read(fileContent)
	if err != nil {
		log.Println("Route.User.uploadUserAvatar ", err)
		ErrResponse(errors.New("something wrong"), w)
		return
	}

	user, err := Usecases.UploadUserAvatar(fileContent, handler, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(user)
	DataResponse(result, w)
}