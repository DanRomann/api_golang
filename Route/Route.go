package Route

import (
	"github.com/gorilla/mux"
	"docnota/Utils"
)


func Route() *mux.Router {
	Utils.ConnectToDB()
	r := mux.NewRouter()

	//User
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
	r.HandleFunc("/user/uploadImage", uploadUserAvatar).Methods("POST")


	//Document
	r.HandleFunc("/document/templates", getPublicTemplates).Methods("GET")
	r.HandleFunc("/document/public", getPublicDocs).Methods("GET")
	r.HandleFunc("/document/{docId}", getDoc).Methods("GET")
	r.HandleFunc("/document/{docId}/copy", copyDoc).Methods("POST")
	r.HandleFunc("/document", createDoc).Methods("POST")
	r.HandleFunc("/document/fillTemplate", fillTemplate).Methods("POST")
	r.HandleFunc("/document/edit", commitDoc).Methods("POST")
	r.HandleFunc("/document/meta/edit", metaEdit).Methods("POST")

	//Company
	r.HandleFunc("/companies", getCompanies).Methods("GET")
	r.HandleFunc("/company/{companyId}", getCompany).Methods("GET")
	r.HandleFunc("/company/{companyId}/doc", companyDoc).Methods("GET")
	r.HandleFunc("/company", createCompany).Methods("POST")
	r.HandleFunc("/company/confirm/{sha}", confirmCompany).Methods("GET")
	r.HandleFunc("/company/meta_for_create/{countryId}", getMetaForCreate).Methods("GET")

	//Block
	r.HandleFunc("/block/{blockId}/addRelation/{relationBlock}", addRelation).Methods("POST")
	r.HandleFunc("/block/{blockId}/deleteRelation/{relationBlock}", deleteRelation).Methods("POST")
	r.HandleFunc("/block/{blockId}/relations", getBlockRelations).Methods("GET")
	r.HandleFunc("/block/{blockId}", getBlock).Methods("GET")


	//Search
	r.HandleFunc("/document/search", searchDoc).Methods("POST")
	r.HandleFunc("/block/search", searchBlock).Methods("POST")
	r.HandleFunc("/company/search", searchCompany).Methods("POST")
	r.HandleFunc("/user/search", searchUser).Methods("POST")

	//Utils
	r.HandleFunc("/country", CountryList).Methods("GET")

	//BlockChain
	r.HandleFunc("/document/{docID}/blockChain/upload", blockChainUploadDoc).Methods("POST")

	//iFrame
	r.HandleFunc("/document/{docId}/iframe", getDocIFrame).Methods("GET")
	r.HandleFunc("/block/{blockId}/iframe", getBlockIFrame).Methods("GET")

	return r
}