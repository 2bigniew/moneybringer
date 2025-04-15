package main

import (
	"encoding/json"
	"flag"
	"fmt"
	InvoiceGenerator "moneybringer/invoice-generator"
	InvoiceManager "moneybringer/invoice-manager"
	Invoice "moneybringer/invoice-manager/invoice"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	/* TODO
	   - save invoice as a raw json
	   - save invoice as an pdf
	*/

	fmt.Println("Moneybringer - let's make some money, baby! Prepare new invoice")

	customer := flag.String("customer", "default", "Customer name")

	flag.Parse()

	invoice := InvoiceManager.CreateInvoice(*customer)

	fmt.Println("Invoice data:")
	fmt.Println(invoice)

	SaveInvoiceRaw(invoice)

	invoicePath := getInvoicePdfName(invoice)
	InvoiceGenerator.GenerateInvoicePDF(invoice, invoicePath)
}

func getInvoicePdfName(payload InvoiceManager.InvoiceCreatedData) string {
	dirPath := Invoice.GetInvoiceDirPath()
	invoiceNo := strings.ReplaceAll(payload.InvoiceNo, "/", "_")
	fileName := fmt.Sprintf("%s_%s_%s.pdf", invoiceNo, payload.AuthorFirstName, payload.AuthorLastName)
	filePath := filepath.Join(dirPath, fileName)
	return filePath
}

func SaveInvoiceRaw(payload InvoiceManager.InvoiceCreatedData) bool {
	dirPath := Invoice.GetInvoiceDirPath()
	rawFileDirPath := filepath.Join(dirPath, "raw")
	invoiceNo := strings.ReplaceAll(payload.InvoiceNo, "/", "_")
	fileName := fmt.Sprintf("%s_%s_%s.json", invoiceNo, payload.AuthorFirstName, payload.AuthorLastName)
	filePath := filepath.Join(rawFileDirPath, fileName)

	err := os.MkdirAll(rawFileDirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return false
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return false
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return false
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return false
	}

	fmt.Println("JSON data successfully saved to %s", filePath)

	return true
}

/*
Numeric Types	int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128
Boolean Type	bool
String Type	string
Composite Types	array, slice, struct, map
Reference Types	pointer, channel, function, interface{}
Special Types	nil, byte (alias for uint8), rune (alias for int32)
*/
