package workers

import (
	"ametory-cooperative/objects"
	"ametory-cooperative/services"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AMETORY/ametory-erp-modules/context"
)

func SendMail(erpContext *context.ERPContext) {
	appService, ok := erpContext.AppService.(*services.AppService)
	if ok {
		dataSub := appService.Redis.Subscribe(*erpContext.Ctx, "SEND:MAIL")
		for {
			msg, err := dataSub.ReceiveMessage(*erpContext.Ctx)
			if err != nil {
				log.Println(err)
			}
			var emailData objects.EmailData
			err = json.Unmarshal([]byte(msg.Payload), &emailData)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("FullName", emailData.FullName)
			fmt.Println("Email", emailData.Email)
			// fmt.Println("sender", sender)
			subject := "Welcome to " + appService.Config.Server.AppName
			if emailData.Subject != "" {
				subject = emailData.Subject
			}
			erpContext.EmailSender.SetAddress(emailData.FullName, emailData.Email)

			if err := erpContext.EmailSender.SendEmail(subject, emailData, []string{}); err != nil {
				log.Println(err)
				fmt.Println(err)
				continue
			}

		}
	}
}
