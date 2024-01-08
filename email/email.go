package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

type Application struct {
	Location           string             `json:"location"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	StreetAddress      string             `json:"streetAddress"`
	Apt                string             `json:"apt"`
	City               string             `json:"city"`
	State              string             `json:"state"`
	Zip                string             `json:"zip"`
	PhoneNumber        string             `json:"phoneNumber"`
	Email              string             `json:"email"`
	GrossMonthlyIncome string             `json:"grossMonthlyIncome"`
	PreviousLandlords  []PreviousLandlord `json:"previousLandlord"`
	Employer           Employer           `json:"employer"`
	Notes              string             `json:"notes"`
}

type PreviousLandlord struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	ReasonForLeaving string `json:"reasonForLeaving"`
	RentPaid         string `json:"rentPaid"`
	Dates            string `json:"dates"`
	OkToContact      string `json:"okToContact"`
}

type Employer struct {
	Name              string `json:"name"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	DatesOfEmployment string `json:"datesOfEmployment"`
}

var tmpl = `
Location: {{.Location}}
First Name: {{.FirstName}}
Last Name: {{.LastName}}
Street Address: {{.StreetAddress}}
Apt: {{.Apt}}
City: {{.City}}
State: {{.State}}
Zip: {{.Zip}}
Phone Number: {{.PhoneNumber}}
Email: {{.Email}}
Gross Monthly Income: {{.GrossMonthlyIncome}}
Previous Landlords:
{{range .PreviousLandlords}}
	Name: {{.Name}}
	Phone: {{.Phone}}
	Email: {{.Email}}
	Address: {{.Address}}
	Reason For Leaving: {{.ReasonForLeaving}}
	Rent Paid: {{.RentPaid}}
	Dates: {{.Dates}}
	Ok To Contact: {{.OkToContact}}
{{end}}
Employer:
	Name: {{.Employer.Name}}
	Phone: {{.Employer.Phone}}
	Email: {{.Employer.Email}}
	Dates Of Employment: {{.Employer.DatesOfEmployment}}
Notes: {{.Notes}}
`

func SendEmail(app Application) error {
	auth := smtp.PlainAuth("", os.Getenv("GMAIL_EMAIL"), os.Getenv("GMAIL_PASSWORD"), "smtp.gmail.com")
	to := strings.Split(os.Getenv("GMAIL_DESTINATION"), ",")
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return err
	}
	b := &bytes.Buffer{}
	_, err = b.WriteString(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n", os.Getenv("GMAIL_DESTINATION"), "Rental Application"))
	err = t.ExecuteTemplate(b, "email", app)
	if err != nil {
		return err
	}
	err = smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL_EMAIL"), to, b.Bytes())
	if err != nil {
		return err
	}
	return nil
}
