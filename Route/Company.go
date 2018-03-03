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
