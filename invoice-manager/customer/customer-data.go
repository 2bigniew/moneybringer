package CustomerData

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Address struct {
	StreetAddress string `json:"streetAddress"`
	State         string `json:"state"`
	Number        string `json:"country"`
	ZipCode       string `json:"zipCode"`
	City          string `json:"city"`
}

type Customer struct {
	FullName string  `json:"fullName"`
	Address  Address `json:"address"`
}

type CustomersData struct {
	Customers map[string]Customer `json:"customers"`
}

const CUSTOMERS_JSON_PATH = "./config/customers.json"

func GetCustomerData(customerName string) Customer {
	file, osErr := os.Open(CUSTOMERS_JSON_PATH)

	if osErr != nil {
		log.Fatal(osErr)
		os.Exit(1)
	}
	defer file.Close()

	jsonData, readErr := io.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
		os.Exit(1)
	}

	var customersData CustomersData
	jsonErr := json.Unmarshal(jsonData, &customersData)
	if jsonErr != nil {
		fmt.Println("Error unmarshalling JSON:", jsonErr)
		os.Exit(1)
	}

	customer, exists := customersData.Customers[customerName]
	if !exists {
		fmt.Printf("Customer: %s does not exist, check your config customers.json file", customerName)
		os.Exit(1)

	}

	return customer
}
