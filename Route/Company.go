package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
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
