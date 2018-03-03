package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
	"io/ioutil"
	"errors"
)

func CountryList(w http.ResponseWriter, r *http.Request){

	countries, err := Usecases.GetCountries(Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(countries)
	DataResponse(result, w)
}

func searchBlock(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	curQuery := new(SearchQuery)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(data, &curQuery)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}
	blocks, err := Usecases.SearchBlockByQuery(&curQuery.Query, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(blocks)
	DataResponse(result, w)
}


func searchDoc(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	curQuery := new(SearchQuery)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(data, &curQuery)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}
	docs, err := Usecases.SearchDocumentsByQuery(&curQuery.Query, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(docs)
	DataResponse(result, w)
}

func searchCompany(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Authorization")

	curQuery := new(SearchQuery)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrResponse(errors.New("bad data"), w)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(data, &curQuery)
	if err != nil {
		ErrResponse(errors.New("bad json"), w)
		return
	}

	companies, err := Usecases.SearchCompanyByQuery(&curQuery.Query, &token, Utils.Connect)
	if err != nil {
		ErrResponse(err, w)
		return
	}

	result, _ := json.Marshal(companies)
	DataResponse(result, w)
}