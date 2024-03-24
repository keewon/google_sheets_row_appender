package main

import (
  "fmt"
  "flag"
  "log"
  "os"
  "encoding/csv"

  "golang.org/x/net/context"
  "google.golang.org/api/sheets/v4"
  "google.golang.org/api/option"

  // "github.com/fsnotify/fsnotify"
)

func main() {
  // var watchFlag bool
  var startRowIndex1 int
  var endRowIndex1 int
  var credentialsFile string
  // var csvFilePath string
  // var sheetName string

  // flag.StringVar(&watchFlag, "watchDir", false, "Watch the csv file for changes")
  flag.IntVar(&startRowIndex1, "start", 1, "Start row to append data")
  flag.IntVar(&endRowIndex1, "end", -1, "End row to append data")
  flag.StringVar(&credentialsFile, "credentialsFile", "service_account.json", "Google Sheets credentials file")
  // flag.StringVar(&csvFilePath, "csv_file", "", "Path to the csv file")
  // flag.StringVar(&sheetName, "sheet_name", "", "Name of the sheet")

  flag.Parse()


  if len(flag.Args()) < 3 {
    fmt.Println("Usage: main.exe [OPTIONS] <csv_file> <spreadsheet_id> <sheet_name>")
    fmt.Println("Options:")
    flag.PrintDefaults()
    return
  }
  csvFilePath := flag.Arg(0)
  spreadsheetId := flag.Arg(1)
  sheetName := flag.Arg(2)

  if csvFilePath != "" || sheetName != "" {
    csvData := ReadCsvFile(csvFilePath)
    srv, err := sheets.NewService(context.Background(),
      option.WithCredentialsFile(credentialsFile),
      option.WithScopes(sheets.SpreadsheetsScope))

    if err != nil {
      log.Fatalf("Unable to open service with give credentialsFile: %v", err)
    }

    startRow := startRowIndex1 - 1
    endRow := endRowIndex1 - 1
    if endRowIndex1 == -1 {
      endRow = len(csvData) - 1
    }

    AppendCsvDataToGoogleSheets(srv, spreadsheetId, sheetName, csvData, startRow, endRow)

    fmt.Printf("Appended!\n")
  }  
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

func AppendCsvDataToGoogleSheets(srv *sheets.Service, spreadsheetId string, sheetName string, csvData [][]string, startRow int, endRow int) {
  var interfaceRecord [][]interface{}

  for i, line := range csvData {
    if startRow <= i && i <= endRow {
      fmt.Printf("%v: %+v\n", i + 1, line)

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
