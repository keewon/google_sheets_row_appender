# Setup
There are several ways to access Google Sheets via API,
but I recommend creating Service Account.
With Service Account, API can access specific documents only.
Following Youtube video is helpful to set it up.

https://www.youtube.com/watch?v=sVURhxyc6jE

Download JSON file from Google cloud, and save it as service_account.json.

# How to Build
## Native build
```
# From Mac
$ go build -o bin/app-macos

# From Windows
> go build -o bin/app-amd64.exe
```

## Cross compile example
```
$ GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe
```

# Usage
```
app-amd64.exe service_account.json test.csv SheetID SheetName
```
 - service_account.json : Download JSON file for the Service Account from Google Cloud
 - test.csv : The CSV file you want to upload
 - SheetID
   - `1VRCpMDBJy539vz2rzdEqEpOqSNJznOXNDFmXF--_DZE` is the SheetID when https://docs.google.com/spreadsheets/d/1VRCpMDBJy539vz2rzdEqEpOqSNJznOXNDFmXF--_DZE/edit#gid=0 is the URL
 - SheetName : The name of the sheet in the document

