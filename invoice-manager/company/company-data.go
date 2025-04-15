package CompanyData

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	ZipCode string `json:"zipCode"`
	City    string `json:"city"`
}

type Payment struct {
	Method       string `json:"method"`
	PeriodInDays int    `json:"periodInDays"`
}

type PersonalDetails struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phome     string `json:"phome"`
}

type CompanyDetails struct {
	FullName  string  `json:"fullName"`
	Address   Address `json:"address"`
	TaxNumber string  `json:"taxNumber"`
	Email     string  `json:"email"`
	Phome     string  `json:"phome"`
	IBAN      string  `json:"IBAN"`
	SWIFT     string  `json:"SWIFT"`
}

type InvoicePosition struct {
	DefaultProduct                         string  `json:"defaultProduct"`
	DefaultUnit                            string  `json:"defaultUnit"`
	DefaultNetPrice                        float64 `json:"defaultNetPrice"`
	DefaultTaxRate                         float64 `json:"defaultTaxRate"`
	PolishClassificationOfGoodsAndServices string  `json:"polishClassificationOfGoodsAndServices"`
	DefaultCurrency                        string  `json:"defaultCurrency"`
}

type InvoiceDetails struct {
	DefaultNotes           []string `json:"defaultNotes"`
	DefaultServiceStartDay int      `json:"defaultServiceStartDay"`
	DefaultServiceEndDay   int      `json:"defaultServiceEndDay"`
	DefaultPlaceOfIssue    string   `json:"defaultPlaceOfIssue"`
}

/* TODO - add fields geters */
type Company struct {
	Payment         Payment         `json:"payment"`
	PersonalDetails PersonalDetails `json:"personalDetails"`
	CompanyDetails  CompanyDetails  `json:"companyDetails"`
	InvoicePosition InvoicePosition `json:"invoicePosition"`
	InvoiceDetails  InvoiceDetails  `json:"invoiceDetails"`
}

const COMPANY_JSON_PATH = "./config/company.json"

func GetCompanyData() Company {
	// Open the JSON file
	file, osErr := os.Open(COMPANY_JSON_PATH)
	if osErr != nil {
		log.Fatal(osErr)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file content into a byte slice
	jsonData, readErr := io.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
		os.Exit(1)
	}

	// Unmarshal the JSON into the Company struct
	var company Company
	jsonErr := json.Unmarshal(jsonData, &company)
	if jsonErr != nil {
		fmt.Println("Error unmarshalling JSON:", jsonErr)
		os.Exit(1)
	}

	return company
}
