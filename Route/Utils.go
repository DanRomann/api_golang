package Route

import (
	"net/http"
	"encoding/json"
)

type TokenResponse struct{
	Token	string	`json:"token"`
}

type SearchQuery struct{
	Query	string	`json:"query"`
}


func ErrResponse(err error, w http.ResponseWriter){
	var curResponse struct{
		Error string `json:"error"`
	}
	curResponse.Error = err.Error()
	result, _ := json.Marshal(curResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func SuccessResponse(result string, w http.ResponseWriter){
	var curResponse struct{
		Result string `json:"result"`
	}
	curResponse.Result = result
	response, _ := json.Marshal(curResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func DataResponse(result []byte, w http.ResponseWriter){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
