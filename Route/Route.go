package Route

import (
	"github.com/gorilla/mux"
	"docnota/Utils"
)


func Route() *mux.Router {
	Utils.ConnectToDB()
	r := mux.NewRouter()


	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/user/confirm", confirmUser).Methods("POST")
	r.HandleFunc("/user/auth", authUser).Methods("POST")
	r.HandleFunc("/user/{userId}/docs", getUserDocs).Methods("GET")

	r.HandleFunc("/country", CountryList).Methods("GET")

	return r
}