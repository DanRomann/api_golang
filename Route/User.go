package Route

import (
	"net/http"
	"io/ioutil"
	"errors"
	"docnota/Usecases"
	"docnota/Utils"
	"docnota/Models"
	"encoding/json"
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

func ConfirmUser(w http.ResponseWriter, r *http.Request){
	var curRequest struct{
		UID			string		`json:"uid"`
		User		Models.User	`json:"user"`
	}

	var curResponse struct{
		Token	string	`json:"token"`
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

	curResponse.Token = *token

	result, _ := json.Marshal(curResponse)
	DataResponse(result, w)

}
