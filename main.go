package main

import (
  "fmt"
  "log"
  "os"
  "encoding/csv"

  "golang.org/x/net/context"
  "google.golang.org/api/sheets/v4"
  "google.golang.org/api/option"
)

func main() {
  if len(os.Args) != 5 {
    fmt.Println("Usage: go run main.go <credential.json> <csv_file_path> <spreadsheet_id> <sheet_name>")
    return
  }

  credentialsFile := os.Args[1]
  csvFilePath := os.Args[2]
  spreadsheetId := os.Args[3]
  sheetName := os.Args[4]

  csvData := ReadCsvFile(csvFilePath)

  srv, err := sheets.NewService(context.Background(),
    option.WithCredentialsFile(credentialsFile),
    option.WithScopes(sheets.SpreadsheetsScope))

  if err != nil {
    log.Fatalf("Unable to open service with give credentialsFile: %v", err)
  }

  AppendCsvDataToGoogleSheets(srv, spreadsheetId, sheetName, csvData)

  fmt.Printf("Done!")
}

func ReadCsvFile(csvFilePath string) [][]string {
  csvFile, err := os.Open(csvFilePath)
  if err != nil {
    log.Fatalf("Unable to open csv file: %v", err)
  }

  csvReader := csv.NewReader(csvFile)
  data, err := csvReader.ReadAll()
  if err != nil {
    log.Fatalf("Unabled to read csv file: %v", err)
  }

  csvFile.Close()
  return data
}

func AppendCsvDataToGoogleSheets(srv *sheets.Service, spreadsheetId string, sheetName string, csvData [][]string) {
  var interfaceRecord [][]interface{}

  for i, line := range csvData {
    if i > 0 {// Omit header line
      fmt.Printf("%+v\n", line)

      var interfaceLine []interface{}

      for _, field := range line {
        interfaceLine = append(interfaceLine, field)
      }

      interfaceRecord = append(interfaceRecord, interfaceLine)
    }
  }

  valueInputOption := "USER_ENTERED"
  insertDataOption := "INSERT_ROWS"

  rb := &sheets.ValueRange{
    Values: interfaceRecord,
  }
  response2, err := srv.Spreadsheets.Values.Append(spreadsheetId, sheetName, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(context.Background()).Do()
  if err != nil || response2.HTTPStatusCode != 200 {
    log.Fatalf("Append error: %v", err)
  }
}
// vim:ts=2 sts=2 sw=2 et
