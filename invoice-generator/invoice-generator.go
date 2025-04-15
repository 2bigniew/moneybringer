package InvoiceGenerator

import (
	"fmt"
	"log"
	InvoiceManager "moneybringer/invoice-manager"

	"github.com/phpdave11/gofpdf"
)

func GenerateInvoicePDF(invoice InvoiceManager.InvoiceCreatedData, outputPath string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	basicDocumentSetup(pdf)
	createHeaderSection(pdf, invoice)
	createCompanyAndCustomerSection(pdf, invoice)
	createPositionsSection(pdf, invoice)
	createSummarySection(pdf, invoice)

	// Save PDF
	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		log.Fatalf("Error saving PDF: %v", err)
	}

	log.Println("Invoice PDF generated successfully at", outputPath)
}

func basicDocumentSetup(pdf *gofpdf.Fpdf) {
	pdf.AddUTF8Font("Inter", "", "assets/fonts/Inter-VariableFont_opsz,wght.ttf")
	pdf.AddUTF8Font("InterItalic", "", "assets/fonts/Inter-Italic-VariableFont_opsz,wght.ttf")
	pdf.AddUTF8Font("Inter", "B", "assets/fonts/static/Inter_18pt-Bold.ttf")

	pdf.SetMargins(2, 10, 2)
	pdf.AddPage()
	pdf.SetFont("Inter", "B", 16)
}

func createHeaderSection(pdf *gofpdf.Fpdf, invoice InvoiceManager.InvoiceCreatedData) {
	// Title
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(10)

	// Invoice Information
	pdf.SetFont("Inter", "", 12)
	pdf.Cell(0, 10, "Invoice Number: "+invoice.InvoiceNo)
	pdf.Ln(6)
	pdf.Cell(0, 10, "Date of Issue: "+invoice.DateOfIssue)
	pdf.Ln(6)
	pdf.Cell(0, 10, "Place of Issue: "+invoice.PlaceOfIssue)
	pdf.Ln(6)
	pdf.Cell(0, 10, "Service Start date: "+invoice.ServiceStartDate)
	pdf.Ln(6)
	pdf.Cell(0, 10, "Service End date: "+invoice.ServiceEndDate)
	pdf.Ln(16)
}

func createCompanyAndCustomerSection(pdf *gofpdf.Fpdf, invoice InvoiceManager.InvoiceCreatedData) {
	// From and To Addresses
	pdf.SetFont("Inter", "B", 16)
	pdf.Cell(0, 10, "Details")
	pdf.Ln(16)

	pdf.SetFont("Inter", "B", 12)

	pdf.CellFormat(95, 6, "From:", "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "To:", "", 1, "R", false, 0, "")

	pdf.SetFont("Inter", "", 12)

	pdf.CellFormat(95, 6, invoice.InvoiceFrom.FullName, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, invoice.InvoiceTo.FullName, "", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, invoice.InvoiceFrom.Address, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, invoice.InvoiceTo.Address.StreetAddress+", "+invoice.InvoiceTo.Address.Number, "", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "Tax Number: "+invoice.InvoiceFrom.TaxNumber, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, invoice.InvoiceTo.Address.City+", "+invoice.InvoiceTo.Address.State+" - "+invoice.InvoiceTo.Address.ZipCode, "", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "Email: "+invoice.InvoiceFrom.Email, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "", "", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "IBAN: "+invoice.IBAN, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "", "", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "SWIFT: "+invoice.SWIFT, "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "", "", 1, "R", false, 0, "")

	pdf.Ln(10)
}

func createPositionsSection(pdf *gofpdf.Fpdf, invoice InvoiceManager.InvoiceCreatedData) {
	// Table Header
	pdf.SetFont("Inter", "B", 16)
	pdf.Cell(0, 10, "Positions")
	pdf.Ln(16)

	pdf.SetFont("Inter", "B", 8)
	pdf.SetFillColor(200, 200, 200)
	headers := []string{"No.", "Product / Service name", "Symbol PKWiU", "Unit", "Qt", "Net price", "Net value", "Tax rate", "Tax amount", "Gross value", "Currency"}
	colWidths := []float64{10, 50, 24, 8, 8, 19, 19, 15, 19, 19, 15} // Column widths

	// Draw headers
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}

	pdf.Ln(-1) // Move to next line

	// Table Content
	pdf.SetFont("Inter", "", 8)
	for i, pos := range invoice.InvoicePositions {
		pdf.CellFormat(10, 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")                     // Item No.
		pdf.CellFormat(50, 10, pos.ProductOrServiceName, "1", 0, "C", false, 0, "")                   // Product / Service name
		pdf.CellFormat(24, 10, pos.PolishClassificationOfGoodsAndServices, "1", 0, "C", false, 0, "") // Symbol PKWiU
		pdf.CellFormat(8, 10, pos.Unit, "1", 0, "C", false, 0, "")                                    // Unit
		pdf.CellFormat(8, 10, fmt.Sprintf("%d", pos.Quantity), "1", 0, "C", false, 0, "")             // Quantity
		pdf.CellFormat(19, 10, fmt.Sprintf("%.2f", pos.NetPrice), "1", 0, "C", false, 0, "")          // Net price
		pdf.CellFormat(19, 10, fmt.Sprintf("%.2f", pos.NetValue), "1", 0, "C", false, 0, "")          // Net value
		pdf.CellFormat(15, 10, fmt.Sprintf("%d", pos.TaxRate), "1", 0, "C", false, 0, "")             // Tax rate
		pdf.CellFormat(19, 10, fmt.Sprintf("%.2f", pos.TaxAmount), "1", 0, "C", false, 0, "")         // Tax amount
		pdf.CellFormat(19, 10, fmt.Sprintf("%.2f", pos.GrossValue), "1", 0, "C", false, 0, "")        // Gross value
		pdf.CellFormat(15, 10, pos.Currency, "1", 1, "C", false, 0, "")                               // Currency
	}
}

func createSummarySection(pdf *gofpdf.Fpdf, invoice InvoiceManager.InvoiceCreatedData) {
	// Check if array has elements and avoid out-of-range error
	var currency string = "PLN"
	if len(invoice.InvoicePositions) > 0 && invoice.InvoicePositions[0].Currency != "" {
		currency = invoice.InvoicePositions[0].Currency
	}

	// Summary
	pdf.Ln(16)
	pdf.SetFont("Inter", "B", 16)
	pdf.Cell(0, 10, "Summary")
	pdf.Ln(8)
	pdf.SetFont("Inter", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Total Amount: %.2f %s", invoice.InvoiceSummary.TotalAmount, currency))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Tax Amount: %.2f %s", invoice.InvoiceSummary.TotalTaxAmount, currency))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total Gross Value: %.2f %s", invoice.InvoiceSummary.TotalGrossValue, currency))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Issued An Invoice: %s %s", invoice.AuthorFirstName, invoice.AuthorLastName))
	pdf.Ln(20)

	pdf.CellFormat(150, 10, "Notes", "1", 0, "C", false, 0, "")
	pdf.Ln(-1) // Move to next line

	pdf.SetFont("Inter", "", 10)
	pdf.MultiCell(150, 20, invoice.Notes, "1", "C", false)

	pdf.Ln(8)
}
