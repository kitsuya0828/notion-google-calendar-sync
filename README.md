# notion-googlecalendar-sync

**notion-googlecalendar-sync** is a tool for two-way synchronisation between Notion and Google Calendar.


## Run locally
https://cloud.google.com/docs/authentication/production?hl=ja

```bash
gcloud auth application-default login --impersonate-service-account xxxx-compute@developer.gserviceaccount.com
source .env
go run cmd/main.go
```