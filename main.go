package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/luno/luno-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type BalanceData struct {
	Timestamp string
	Asset     string
	AccountID string
	Balance   string
	Reserved  string
}

func main() {
	// Initialize the Luno client
	lunoClient := luno.NewClient()
	lunoClient.SetAuth(os.Getenv("LUNO_API_KEY"), os.Getenv("LUNO_API_SECRET"))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get balances from Luno
	balances, err := getLunoBalances(ctx, lunoClient)
	if err != nil {
		log.Fatal("Error getting Luno balances:", err)
	}

	// Save to Google Sheets
	if err := saveToGoogleSheets(ctx, balances); err != nil {
		log.Fatal("Error saving to Google Sheets:", err)
	}

	log.Println("Successfully saved balances to Google Sheets")
}

func getLunoBalances(ctx context.Context, client *luno.Client) ([]BalanceData, error) {
	request := &luno.GetBalancesRequest{}
	response, err := client.GetBalances(ctx, request)
	if err != nil {
		return nil, err
	}

	var balances []BalanceData
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	for _, balance := range response.Balance {
		balances = append(balances, BalanceData{
			Timestamp: timestamp,
			Asset:     balance.Asset,
			AccountID: balance.AccountId,
			Balance:   balance.Balance.String(),
			Reserved:  balance.Reserved.String(),
		})
		log.Printf("Asset: %s, Account ID: %s, Balance: %s, Reserved: %s\n",
			balance.Asset,
			balance.AccountId,
			balance.Balance.String(),
			balance.Reserved.String(),
		)
	}

	return balances, nil
}

func saveToGoogleSheets(ctx context.Context, balances []BalanceData) error {
	// Read credentials file
	credBytes, err := os.ReadFile(os.Getenv("JSON_CREDENTIALS"))
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %v", err)
	}

	// Create JWT config
	config, err := google.JWTConfigFromJSON(credBytes, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("failed to create JWT config: %v", err)
	}

	// Create client
	client := config.Client(ctx)

	// Create sheets service
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create sheets service: %v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	range_ := "<spreadsheetName!A:E>"

	// First, get existing values to check if headers exist
	existing, err := srv.Spreadsheets.Values.Get(spreadsheetId, range_).Do()
	if err != nil {
		return fmt.Errorf("failed to get existing values: %v", err)
	}

	// Prepare data
	var values [][]interface{}

	// Only add headers if the sheet is empty
	if len(existing.Values) == 0 {
		values = append(values, []interface{}{
			"Timestamp",
			"Asset",
			"Account ID",
			"Balance",
			"Reserved",
		})
	}

	// Add balance data
	for _, balance := range balances {
		values = append(values, []interface{}{
			balance.Timestamp,
			balance.Asset,
			balance.AccountID,
			balance.Balance,
			balance.Reserved,
		})
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	// Append data to sheet
	_, err = srv.Spreadsheets.Values.Append(
		spreadsheetId,
		range_,
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return fmt.Errorf("failed to append data: %v", err)
	}

	return nil
}
