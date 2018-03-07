package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
	"errors"
)

func getCompanies(w http.ResponseWriter, r *http.Request){
	companies, err := Usecases.GetCompanyList(Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(companies)
	DataResponse(result, w)
}

func getCompany(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["companyId"])
	if err != nil {
		ErrResponse(errors.New("bad company id"), w)
		return
	}

	company, err := Usecases.GetCompany(companyId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(company)
	DataResponse(result, w)
}

func companyDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["companyId"])
	if err != nil {
		ErrResponse(errors.New("bad company id"), w)
		return
	}

	documents, err := Usecases.GetCompanyDoc(companyId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(documents)

	DataResponse(result, w)
}

func createCompany(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()


	var curRequest map[string]interface{}
	err := decoder.Decode(&curRequest)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	err = Usecases.CreateCompany(curRequest, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}
	SuccessResponse("ok", w)
}

func confirmCompany(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	sha, ok := vars["sha"]
	if !ok{
		ErrResponse(errors.New("bad sha"), w)
		return
	}

	err := Usecases.ConfirmCompany(&sha, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	SuccessResponse("ok", w)
}

func getMetaForCreate(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)

	countryId, err := strconv.Atoi(vars["countryId"])
	if err != nil {
		ErrResponse(errors.New("bad country id"), w)
		return
	}

	meta, err := Usecases.GetCompanyMetaByCountry(countryId, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, err := json.Marshal(meta)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
