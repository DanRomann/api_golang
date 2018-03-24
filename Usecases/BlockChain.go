package Usecases

import (
	"net/http"
	"net"
	"time"
	"bytes"
	"log"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strconv"
	"docnota/Models"
)

var Client *http.Client

type ErrorParam struct {
	Error string	`json:"error,omitempty"`
}

type UploadedDocument struct{
	DocID	string			`json:"doc_id,omitempty"`
	Name	string			`json:"name,omitempty"`
	Content	[]Models.Block	`json:"content,omitempty"`

}
type DocRequest struct{
	UserID   string           `json:"user_id,omitempty"`
	Document UploadedDocument `json:"doc,omitempty"`
	ErrorParam
}



func init(){
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:	30 * time.Second,
			KeepAlive: 	30 * time.Second,
			DualStack:  true,
		}).DialContext,
		MaxIdleConns:			100,
		IdleConnTimeout: 		90 * time.Second,
		TLSHandshakeTimeout: 	10 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
	}

	Client = &http.Client{
		Timeout: time.Second * 10,
		Transport: transport,
	}
}

func RegisterUser(u *Models.User) error{
	var curRequestParams struct{
		UserID	string	`json:"user_id"`
	}
	curResponse := new(ErrorParam)
	curRequestParams.UserID = strconv.Itoa(u.ID)

	b, _ := json.Marshal(curRequestParams)

	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "http://159.89.105.110:3523/register", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := Client.Do(req)
	if err != nil {
		log.Println("Services.BlockChain.RegisterUser send request ", err)
		return errors.New("can't register user in blockchain :C")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Services.BlockChain.RegisterUser read response ", err)
		return errors.New("something wrong with user register in blockchain")
	}
	err = json.Unmarshal(respBody, curResponse)
	if err != nil {
		log.Println("RServices.BlockChain.RegisterUser convert response", err)
		return errors.New("something wrong")
	}
	if curResponse.Error != ""{
		log.Println("Services.BlockChain.RegisterUser read result response ", curResponse.Error)
		return errors.New("something wrong")
	}
	return nil
}

func UploadDocument(u *Models.User, doc *Models.Document) error{
	curRequestParams := new(DocRequest)
	curResponse := new(ErrorParam)

	curRequestParams.UserID = strconv.Itoa(u.ID)
	curRequestParams.Document.DocID = strconv.Itoa(doc.ID)
	curRequestParams.Document.Name = doc.Name
	curRequestParams.Document.Content = doc.Blocks

	b, _ := json.Marshal(curRequestParams)

	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "http://159.89.105.110:3523/create", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := Client.Do(req)
	if err != nil {
		log.Println("Services.BlockChain.RegisterUser send request ", err)
		return errors.New("can't UploadDocument user in blockchain :C")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Services.BlockChain.UploadDocument read response ", err)
		return errors.New("something wrong with user register in blockchain")
	}
	err = json.Unmarshal(respBody, curResponse)
	if err != nil {
		log.Println("Services.BlockChain.UploadDocument convert response", err)
		return errors.New("something wrong")
	}
	if curResponse.Error != ""{
		log.Println("Services.BlockChain.UploadDocument read result response ", curResponse.Error)
		return errors.New("something wrong")
	}
	return nil

}

func GetDocumentFromBlockChain(docID int) ([]byte, error){
	var curDocument struct{
		DocID 	string	`json:"doc_id"`
	}
	curDocument.DocID = strconv.Itoa(docID)
	b, _ := json.Marshal(curDocument)

	body := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "http://159.89.105.110:3523/query", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := Client.Do(req)
	if err != nil {
		log.Println("Services.BlockChain.GetDocument send request ", err)
		return nil, errors.New("can't UploadDocument user in blockchain :C")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Services.BlockChain.GetDocument read response ", err)
		return nil, errors.New("something wrong with user register in blockchain")
	}
	return respBody, nil
}
