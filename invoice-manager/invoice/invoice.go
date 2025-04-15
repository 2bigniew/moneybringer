package Invoice

import (
	"bufio"
	"fmt"
	CompanyData "moneybringer/invoice-manager/company"
	TimeUtils "moneybringer/utils/time"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type InvoicePosition struct {
	ItemNo                                 int
	ProductOrServiceName                   string
	PolishClassificationOfGoodsAndServices string
	Unit                                   string
	Quantity                               int
	NetPrice                               float32
	NetValue                               float32
	TaxRate                                int
	TaxAmount                              float32
	GrossValue                             float32
	Currency                               string
}

const INVOICES_DIR_PATH = "./invoices"

func GetInvoiceDirPath() string {
	currentTime := TimeUtils.GetCurrentTime()
	currentYear := strconv.Itoa(currentTime.Year())
	yearDirPath := filepath.Join(INVOICES_DIR_PATH, currentYear)

	if _, yearErr := os.Stat(yearDirPath); os.IsNotExist(yearErr) {
		os.Mkdir(yearDirPath, 0755)
		fmt.Printf("Directory created: %s", yearDirPath)
	}

	monthName := currentTime.Month()

	monthDirPath := filepath.Join(yearDirPath, monthName.String())

	if _, monthErr := os.Stat(monthDirPath); os.IsNotExist(monthErr) {
		os.Mkdir(monthDirPath, 0755)
		fmt.Printf("Directory created: %s", monthDirPath)
	}

	return monthDirPath
}

func GetInvoicePositions(defaultPosition CompanyData.InvoicePosition) []InvoicePosition {
	var positionsCounter int = 0
	var shouldAddNewPosition bool = true
	var invoicePositionsSlice []InvoicePosition

	for {
		positionsCounter++
		position := createInvoicePosition(positionsCounter, defaultPosition)
		invoicePositionsSlice = append(invoicePositionsSlice, position)

		fmt.Printf("Shoul add another position? Y/n (yes, no)")

		var input string
		_, err := fmt.Scanln(&input)

		if err != nil {
			shouldAddNewPosition = false
		}

		if input != "Y" {
			shouldAddNewPosition = false
		}

		if shouldAddNewPosition == false {
			break
		}
	}

	return invoicePositionsSlice
}

func createInvoicePosition(itemNo int, defaultPosition CompanyData.InvoicePosition) InvoicePosition {
	fmt.Printf("Enter product (or press Enter to use the default: %s):", defaultPosition.DefaultProduct)
	productOrServiceName := createStringPosition(defaultPosition.DefaultProduct)

	fmt.Printf("Enter unit (or press Enter to use the default: %s):", defaultPosition.DefaultUnit)
	unit := createStringPosition(defaultPosition.DefaultUnit)

	fmt.Printf("Enter net price (or press Enter to use the default: %f):", defaultPosition.DefaultNetPrice)
	netPrice := createFloatPosition(defaultPosition.DefaultNetPrice)

	fmt.Printf("Enter tax rate (or press Enter to use the default: %f):", defaultPosition.DefaultTaxRate)
	taxRate := createFloatPosition(defaultPosition.DefaultTaxRate)

	fmt.Printf("Enter polish classification of goods and services (or press Enter to use the default: %s):", defaultPosition.PolishClassificationOfGoodsAndServices)
	polishClassificationOfGoodsAndServices := createStringPosition(defaultPosition.PolishClassificationOfGoodsAndServices)

	fmt.Printf("Enter quantity (or press Enter to use the default: %d):", 160)
	quantity := createIntPosition(160)

	netValue := netPrice * float64(quantity)
	taxAmount := calculateTaxAmount(taxRate, netValue)
	grossValue := netValue + taxAmount

	fmt.Printf("Enter currency (or press Enter to use the default: %s):", defaultPosition.DefaultCurrency)
	Currency := createStringPosition(defaultPosition.DefaultCurrency)

	return InvoicePosition{
		ItemNo:                                 itemNo,
		ProductOrServiceName:                   productOrServiceName,
		PolishClassificationOfGoodsAndServices: polishClassificationOfGoodsAndServices,
		Unit:                                   unit,
		Quantity:                               quantity,
		NetPrice:                               float32(netPrice),
		NetValue:                               float32(netValue),
		TaxRate:                                int(taxRate),
		TaxAmount:                              float32(taxAmount),
		GrossValue:                             float32(grossValue),
		Currency:                               Currency,
	}
}

func createStringPosition(defaultValue string) string {
	var input string
	var value string
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		value = defaultValue
	}

	// Trim spaces and newline characters
	value = strings.TrimSpace(input)

	if err != nil || value == "" {
		value = defaultValue
	}

	return value
}

func createIntPosition(defaultValue int) int {

	var input string
	_, err := fmt.Scanln(&input)

	// Handle case when Scanln does not read anything
	if err != nil && err.Error() == "expected newline" {
		input = "" // Set to empty if newline is detected
	}

	// Convert input to int, or use default value if empty
	number, convErr := strconv.Atoi(strings.TrimSpace(input))
	if convErr != nil || input == "" {
		number = defaultValue
	}

	return number
}

func createFloatPosition(defaultValue float64) float64 {

	var input string
	_, err := fmt.Scanln(&input)

	// Handle case when Scanln does not read anything
	if err != nil && err.Error() == "expected newline" {
		input = "" // Set to empty if newline is detected
	}

	// Convert input to float64, or use default value if empty
	number, convErr := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if convErr != nil || input == "" {
		number = defaultValue
	}

	return number
}

func calculateTaxAmount(taxRate float64, netValue float64) float64 {
	if taxRate > 0 {
		taxDecimal := taxRate / 100
		return netValue * taxDecimal
	}

	return 0
}
