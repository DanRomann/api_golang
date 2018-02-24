package Route

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"docnota/Usecases"
	"docnota/Utils"
)

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