package mail

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"tender-scraper/config"
	"tender-scraper/database"
	"tender-scraper/types"
	"text/template"
	"strings"
	"github.com/mailgun/mailgun-go/v4"
)

func SendEmailsToUsers(cfg *config.Config, users []types.UserInfo, tenders []types.Tender, db *sql.DB) {
	mg := mailgun.NewMailgun(cfg.Domain, cfg.MailgunAPIKey)
	tmpl := CompileEmailTemplate()
	
	for _, user := range users {
		emailBody, err := RenderEmailTemplate(tmpl, user, tenders)
		if err != nil {
			log.Printf("Error rendering email template for %s: %v", user.Email, err)
			continue
		}

		message := mg.NewMessage("TenderSpot@bina.cloud", "G7 Contract Opportunities", "", user.Email)
		message.SetHtml(emailBody)

		_, _, err = mg.Send(context.Background(), message)
		if err != nil {
			log.Printf("Error sending email to %s: %v", user.Email, err)
		} else {
			log.Printf("Email sent successfully to %s", user.Email)
		}
	}

	ids, err := getSentTenderIDs(tenders)

	if err != nil {
		log.Printf("Error getting sent tender IDs: %v", err)
		return
	}

	database.MarkTendersAsNotified(db, ids)
	fmt.Println("Tenders marked as notified.")
}

func SendEmailToUserRegister(cfg *config.Config, user types.UserInfo) {
	mg := mailgun.NewMailgun(cfg.Domain, cfg.MailgunAPIKey)
	tmpl := CompileEmailTemplateUserRegister()

	emailBody, err := RenderEmailTemplateRegister(tmpl, user)

	if err != nil {
		log.Printf("Error rendering email template for %s: %v", user.Email, err)
		
	}
	message := mg.NewMessage("TenderSpot@bina.cloud", "Welcome to TenderSpot – Your Free Tender Information Tool!", "", user.Email)
	message.SetHtml(emailBody)

	_, _, err = mg.Send(context.Background(), message)
	if err != nil {
		log.Printf("Error sending email to %s: %v", user.Email, err)
	} else {
		log.Printf("Email sent successfully to %s", user.Email)
	}
		
	fmt.Println("Done send email after user register.")
}

func CompileEmailTemplate() *template.Template {
	const htmlTemplate = `
	<div style="margin: 10px; padding: 20px; border: 1px solid #ccc; border-radius: 5px; background-color: #ffffff; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); color: #000; font-family: Arial, sans-serif; ">
    <img alt="Tender Opportunities" src="https://bina-prod-storage.obs.my-kualalumpur-1.alphaedge.tmone.com.my/TenderSpot%20logo%20%281%29.png"/>
    <p>
        Hey <b>{{.Username}}!</b>
    </p>
    <p>
        <b style="color: #2a64f9">TenderSpot</b> found new opportunities that might be of interest to you. Here they are:
    </p>
    <h2>G7 Contract Opportunities</h2>
    <h3>JKR E-TENDER</h3>
    <div style="margin-top: 20px; margin-bottom: 20px;">
        {{ range .Tenders }}
            <div style="margin-top: 20px; margin-bottom: 20px; ">
                {{ .Index }}. <a href="{{ .Link }}" style="text-decoration: none; color: #2a64f9;">{{ .Name }}</a><br>
                - Kod Bidang: {{ .KodBidang }}<br>
                - Kebenaran Khas: {{ .KebenaranKhas }}<br>
                - Taraf: {{ .Taraf }}<br>
            </div>
        {{ end }}
    </div>
</div>

`

	tmpl, err := template.New("emailTemplate").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error compiling HTML template: %v", err)
	}
	return tmpl
}

func CompileEmailTemplateUserRegister() *template.Template {
	const htmlTemplate = `
	<div style="margin: 10px; padding: 20px; border: 1px solid #ccc; border-radius: 5px; background-color: #ffffff; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); color: #000; font-family: Arial, sans-serif; ">
    <img alt="Tender Opportunities" src="https://bina-prod-storage.obs.my-kualalumpur-1.alphaedge.tmone.com.my/TenderSpot%20logo%20%281%29.png"/>	
    <p>
	Dear <b>{{.Username}}!</b>
    </p>
	<p>Thank you for subscribing to TenderSpot, Bina Cloud's free tool for getting the latest tender information. We are excited to have you with us!<br><br>
	TenderSpot will email you whenever new tenders are available.<br><br>
	Thank you for being part of our community.<br><br>
	<br>
	Best regards,<br>
	Bina Cloud Team<br>
	</p>
	</div>
`

	tmpl, err := template.New("emailTemplate").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error compiling HTML template: %v", err)
	}
	return tmpl
}

func RenderEmailTemplate(tmpl *template.Template, user types.UserInfo, tenders []types.Tender) (string, error) {
	var emailBody bytes.Buffer
	data := struct {
		Username string
		Tenders  []types.Tender
	}{
		Username: user.Name,
		Tenders:  tenders,
	}
	err := tmpl.Execute(&emailBody, data)
	if err != nil {
		return "", err
	}
	return emailBody.String(), nil
}

func RenderEmailTemplateRegister(tmpl *template.Template, user types.UserInfo) (string, error) {
	var emailBody bytes.Buffer
	// Split the full name to get the first name
	firstName := strings.Split(user.Name, " ")[0]

	data := struct {
		Username string		
	}{
		Username: firstName,		
	}
	err := tmpl.Execute(&emailBody, data)
	if err != nil {
		return "", err
	}
	return emailBody.String(), nil
}

func getSentTenderIDs(tenders []types.Tender) ([]int, error) {
	var ids []int
	for _, tender := range tenders {
		ids = append(ids, tender.ID)
	}
	return ids, nil
}
