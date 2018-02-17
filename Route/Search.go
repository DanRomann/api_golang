package Route

import (
	"net/http"
	"docnota/Usecases"
	"docnota/Utils"
	"encoding/json"
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
