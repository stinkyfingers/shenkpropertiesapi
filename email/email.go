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
	LandlordName     string `json:"landlordName"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	RentalAddress    string `json:"rentalAddress"`
	ReasonForLeaving string `json:"reasonForLeaving"`
	MonthlyRent      string `json:"monthlyRent"`
	OccupancyStart   string `json:"occupancyStart"`
	OccupancyEnd     string `json:"occupancyEnd"`
	OkToContact      string `json:"okToContact"`
}

type Employer struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phoneNumber"`
	Email          string `json:"email"`
	Position       string `json:"position"`
	SupervisorName string `json:"supervisorName"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
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
	LandlordName: {{.LandlordName}}
	Phone: {{.PhoneNumber}}
	Email: {{.Email}}
	Rental Address: {{.RentalAddress}}
	Reason For Leaving: {{.ReasonForLeaving}}
	Rent Paid: {{.RentPaid}}
	Occupancy Start: {{.OccupancyStart}}
	Occupancy End: {{.OccupancyEnd}}
	Ok To Contact: {{.OkToContact}}
{{end}}
Employer:
	Name: {{.Employer.Name}}
	Address: {{.Employer.Address}}
	Phone Number: {{.Employer.PhoneNumber}}
	Email: {{.Employer.Email}}
	Position: {{.Employer.Position}}	
	Supervisor Name: {{.Employer.SupervisorName}}
	Start Date: {{.Employer.StartDate}}
	End Date: {{.Employer.EndDate}}
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
