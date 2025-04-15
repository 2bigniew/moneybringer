package InvoiceManager

import (
	"fmt"
	CompanyData "moneybringer/invoice-manager/company"
	CustomerData "moneybringer/invoice-manager/customer"
	Invoice "moneybringer/invoice-manager/invoice"
	TimeUtils "moneybringer/utils/time"
	"os"
	"strings"
)

type InvoicePayment struct {
	deadline string
	method   string
}

type InvoiceFrom struct {
	FullName  string
	Address   string
	TaxNumber string
	Email     string
}

type CustomerAddress struct {
	StreetAddress string
	State         string
	Number        string
	ZipCode       string
	City          string
}

type InvoiceTo struct {
	FullName string
	Address  CustomerAddress
}

type InvoiceSummary struct {
	TotalAmount     float32
	TotalTaxAmount  float32
	TotalGrossValue float32
}

type InvoiceCreatedData struct {
	InvoiceNo        string
	DateOfIssue      string
	PlaceOfIssue     string
	ServiceStartDate string
	ServiceEndDate   string
	Payment          InvoicePayment
	InvoiceFrom      InvoiceFrom
	InvoiceTo        InvoiceTo
	IBAN             string
	SWIFT            string
	InvoicePositions []Invoice.InvoicePosition
	InvoiceSummary   InvoiceSummary
	Notes            string
	IssuedAnInvoice  string
	AuthorFirstName  string
	AuthorLastName   string
}

func CreateInvoice(customerName string) InvoiceCreatedData {
	var invoicePayment InvoicePayment

	customer := getCustomerData(customerName)
	companyData := getCompanyData()
	dateOfIssue := getDateOfIssue()
	serviceStartDate := getServiceStartDate(companyData.InvoiceDetails.DefaultServiceStartDay)
	serviceEndDate := getServiceEndDate(companyData.InvoiceDetails.DefaultServiceEndDay)
	invoiceNumber := getInvoiceNumber()
	invoicePayment.deadline = getPaymentDeadline(companyData.Payment.PeriodInDays, dateOfIssue)
	invoicePayment.method = fmt.Sprintf("%s (%d days)", companyData.Payment.Method, companyData.Payment.PeriodInDays)
	invoicePositions := Invoice.GetInvoicePositions(companyData.InvoicePosition)
	invoiceFrom := getInvoiceFrom(companyData)
	invoiceTo := getInvoiceTo(customer)
	invoiceSummary := getInvoiceSummary(invoicePositions)

	return InvoiceCreatedData{
		InvoiceNo:        invoiceNumber,
		DateOfIssue:      dateOfIssue,
		PlaceOfIssue:     companyData.InvoiceDetails.DefaultPlaceOfIssue,
		ServiceStartDate: serviceStartDate,
		ServiceEndDate:   serviceEndDate,
		Payment:          invoicePayment,
		InvoiceFrom:      invoiceFrom,
		InvoiceTo:        invoiceTo,
		IBAN:             companyData.CompanyDetails.IBAN,
		SWIFT:            companyData.CompanyDetails.SWIFT,
		InvoicePositions: invoicePositions,
		InvoiceSummary:   invoiceSummary,
		Notes:            strings.Join(companyData.InvoiceDetails.DefaultNotes, ", "),
		IssuedAnInvoice:  fmt.Sprintf("%s %s", companyData.PersonalDetails.FirstName, companyData.PersonalDetails.LastName),
		AuthorFirstName:  companyData.PersonalDetails.FirstName,
		AuthorLastName:   companyData.PersonalDetails.LastName,
	}
}

func getInvoiceSummary(positions []Invoice.InvoicePosition) InvoiceSummary {
	var totalAmount float32 = 0
	var totalTaxAmount float32 = 0
	var totalGrossValue float32 = 0

	for _, position := range positions {
		totalAmount += position.NetValue
		totalTaxAmount += position.TaxAmount
		totalGrossValue += position.GrossValue
	}

	return InvoiceSummary{
		TotalAmount:     totalAmount,
		TotalTaxAmount:  totalTaxAmount,
		TotalGrossValue: totalGrossValue,
	}
}

func getInvoiceFrom(companyData CompanyData.Company) InvoiceFrom {
	street := companyData.CompanyDetails.Address.Street
	homeNumber := companyData.CompanyDetails.Address.Number
	zipcode := companyData.CompanyDetails.Address.ZipCode
	city := companyData.CompanyDetails.Address.City
	addres := fmt.Sprintf("%s %s, %s %s", street, homeNumber, zipcode, city)

	return InvoiceFrom{
		FullName:  companyData.CompanyDetails.FullName,
		Address:   addres,
		TaxNumber: companyData.CompanyDetails.TaxNumber,
		Email:     companyData.CompanyDetails.Email,
	}
}

func getInvoiceTo(customer CustomerData.Customer) InvoiceTo {

	return InvoiceTo{
		FullName: customer.FullName,
		Address:  CustomerAddress(customer.Address),
	}
}

func getInvoiceNumber() string {
	/* number of invoice in month/current month/current year */
	currentTime := TimeUtils.GetCurrentTime()
	currentYear := currentTime.Year()
	currentMonthNumber := int(currentTime.Month())
	monthDirPath := Invoice.GetInvoiceDirPath()

	files, err := os.ReadDir(monthDirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	fileCount := 0
	for _, file := range files {
		if !file.IsDir() { // Only count files, skip directories
			fileCount++
		}
	}

	nextInvoiceNumber := fileCount + 1

	return fmt.Sprintf("%d/%d/%d", nextInvoiceNumber, currentMonthNumber, currentYear)
}

func getPaymentDeadline(defaultPaymentPeriodInDays int, dateOfIssue string) string {
	issueProposedTime, isOriginalTime := TimeUtils.GetDataFromDdMmYyyyFormat(dateOfIssue)

	if isOriginalTime == false {
		fmt.Printf("Using current timestamp to calculate deadline day")
	}

	deadlineProposedDay := issueProposedTime.AddDate(0, 0, defaultPaymentPeriodInDays)
	formated := TimeUtils.FormatToDdMmYyyy(deadlineProposedDay)

	fmt.Printf("Enter date of payment deadline (or press Enter to use the default: %s):", formated)

	var input string
	_, err := fmt.Scanln(&input)

	if err != nil || input == "" {
		input = formated
	}

	return input
}

func getDateOfIssue() string {
	currentTime := TimeUtils.GetCurrentTime()
	formated := TimeUtils.FormatToDdMmYyyy(currentTime)

	fmt.Printf("Enter date of issue (or press Enter to use the default: %s):", formated)

	var input string
	_, err := fmt.Scanln(&input)

	if err != nil || input == "" {
		input = formated
	}

	return input
}

func getServiceStartDate(defaultServiceStartDay int) string {
	previousMonthTime := TimeUtils.GetPreviousMonthTime()
	adjustedTime := TimeUtils.SetDayOfMonth(previousMonthTime, defaultServiceStartDay)
	formated := TimeUtils.FormatToDdMmYyyy(adjustedTime)

	fmt.Printf("Enter service start date (or press Enter to use the default: %s):", formated)

	var input string
	_, err := fmt.Scanln(&input)

	if err != nil || input == "" {
		input = formated
	}

	return input
}

func getServiceEndDate(defaultServiceEndDay int) string {
	currentTime := TimeUtils.GetCurrentTime()
	adjustedTime := TimeUtils.SetDayOfMonth(currentTime, defaultServiceEndDay)
	formated := TimeUtils.FormatToDdMmYyyy(adjustedTime)

	fmt.Printf("Enter service end date (or press Enter to use the default: %s):", formated)

	var input string
	_, err := fmt.Scanln(&input)

	if err != nil || input == "" {
		input = formated
	}

	return input
}

func getCompanyData() CompanyData.Company {
	company := CompanyData.GetCompanyData()

	return company
}

func getCustomerData(customerName string) CustomerData.Customer {
	customer := CustomerData.GetCustomerData(customerName)

	return customer
}
