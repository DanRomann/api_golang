package Route

import (
	"github.com/gorilla/mux"
	"docnota/Utils"
)


func Route() *mux.Router {
	Utils.ConnectToDB()
	r := mux.NewRouter()


	r.HandleFunc("/user/list", getUsers).Methods("GET")
	r.HandleFunc("/user/{userId}", getUser).Methods("GET")
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/user/confirm", confirmUser).Methods("POST")
	r.HandleFunc("/user/auth", authUser).Methods("POST")
	r.HandleFunc("/user/{userId}/docs", getUserDocs).Methods("GET")
	r.HandleFunc("/user/{userId}/templates", getUserTemplate).Methods("GET")
	r.HandleFunc("/user/{userId}/inbox", getUserInbox).Methods("GET")
	r.HandleFunc("/user/send_doc", sendDocToUser).Methods("POST")
	r.HandleFunc("/user/accept_doc/{docId}", acceptDoc).Methods("POST")



	r.HandleFunc("/document/templates", getPublicTemplates).Methods("GET")
	r.HandleFunc("/document/public", getPublicDocs).Methods("GET")
	r.HandleFunc("/document/{docId}", getDoc).Methods("GET")
	r.HandleFunc("/document/{docId}/copy", copyDoc).Methods("POST")
	r.HandleFunc("/document", createDoc).Methods("POST")
	r.HandleFunc("/document/fillTemplate", fillTemplate).Methods("POST")
	r.HandleFunc("/document/edit", commitDoc).Methods("POST")
	r.HandleFunc("/document/meta/edit", metaEdit).Methods("POST")

	r.HandleFunc("/companies", getCompanies).Methods("GET")
	r.HandleFunc("/company/{companyId}", getCompany).Methods("GET")
	r.HandleFunc("/company/{companyId}/doc", companyDoc).Methods("GET")
	r.HandleFunc("/company/meta_for_create/{countryId}", getMetaForCreate).Methods("GET")


	r.HandleFunc("/document/search", searchDoc).Methods("POST")
	r.HandleFunc("/block/search", searchBlock).Methods("POST")
	r.HandleFunc("/company/search", searchCompany).Methods("POST")
	r.HandleFunc("/user/search", searchUser).Methods("POST")

	r.HandleFunc("/country", CountryList).Methods("GET")

	return r
}