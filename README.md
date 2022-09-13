# Run

## Onvista

heroku run onvista -gconfig=/app/cmd/onvista/

## Scheduler

heroku addons:open scheduler

scheudler job calls eg:  
$ onvista -gconfig=/app/cmd/onvista/  
$ etfmatic -gconfig=/app/cmd/etfmatic/  
$ visualvest -gconfig=/app/cmd/visualvest/  

# 2do

refactor google api credentials and token into env

# Google API

https://developers.google.com/sheets/api/guides/values

https://godoc.org/google.golang.org/api/sheets/v4

https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/append

https://stackoverflow.com/questions/39691100/golang-google-sheets-api-v4-write-update-example

# Sheet

https://docs.google.com/spreadsheets/d/12h-p9jDKY6MzoQP-TGG-yKw4WCcJ6Xy2g399RicHnnU/edit#gid=0 

