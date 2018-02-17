package Utils

import (
	sp "github.com/SparkPost/gosparkpost"
	"log"
)

var client sp.Client

func init(){
	apiKey := "3339359bcf1583ec17c0b2e78091561a1bb9ee4e"
	cfg := &sp.Config{
		BaseUrl:	"https://api.sparkpost.com",
		ApiKey:		apiKey,
		ApiVersion: 1,
	}
	err := client.Init(cfg)
	if err != nil {
		log.Println(err)
	}
}

func SendEmail(recipient string, from string, body string, subject string) error{

	tx := &sp.Transmission{
		Recipients: []string{recipient},
		Content: sp.Content{
			HTML: 		body,
			From: 		from,
			Subject:	subject,
		},
	}
	id, _, err := client.Send(tx)
	if err != nil {
		log.Printf("%+v\n", err, id)
		return err
	}
	return nil
}