# Luno Balance to Google Sheets Exporter

This Go application automatically fetches balance information from your Luno account and exports it to Google Sheets.

## Prerequisites

Before running this application, you need to have:

1. Go installed on your system
2. A Luno account with API credentials 
3. A Google Cloud Project with the Google Sheets API enabled
4. A service account with appropriate permissions for Google Sheets

## Dependencies Installation

1. Initialize your Go module (if not already done):
```bash
   go mod init your-project-name
```

## Setup

1. Install dependencies
```bash
go get github.com/luno/luno-go
go get golang.org/x/oauth2/google
go get google.golang.org/api/option
go get google.golang.org/api/sheets/v4
```

2. Prepare environment and files:
   - Place your Google Sheets credentials JSON file in the project root (do not upload to GitHub)
   - Set up your Google Sheet and note the spreadsheet ID
   - Have your Luno API credentials ready

## Configuration Steps

1. Replace placeholders in the code:
   - `LUNO_API_KEY` and `LUNO_API_SECRET` with your Luno API credentials
   - `JSON_CREDENTIALS` with your JSON Google service account credentials
   - `SPREADSHEET_ID` with your Google Sheet ID
   - `spreadsheetName!A:E` with your desired sheet name and range

2. Ensure your Google Sheet:
   - Is shared with the service account email
   - Has the correct sheet name (default is "spreadsheetName")

## Running the Application

```bash
go run main.go
```

## Data Format

The application writes the following columns to Google Sheets:
- Timestamp (YYYY-MM-DD HH:MM:SS)
- Asset (e.g., XBT, ETH)
- Account ID
- Balance
- Reserved

## Error Handling

The application handles:
- API connection issues
- Authentication failures
- Google Sheets access problems
- Data formatting issues

## Security Notes

- Never commit API credentials to version control
- Store credentials securely
- Use environment variables for sensitive data
- Restrict Google Sheet access appropriately

## Troubleshooting

1. **Google Sheets Access Error**
   - Verify service account permissions
   - Check if the sheet is shared with service account
   - Confirm spreadsheet ID is correct

2. **Luno API Connection Failed**
   - Verify API credentials
   - Check internet connectivity
   - Confirm API endpoints are accessible

## Disclaimer

This tool is not officially associated with Luno or Google. Use at your own risk and verify all data independently.