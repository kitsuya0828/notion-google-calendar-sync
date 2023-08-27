PROJECT_ID="xxxxxx-xxxxxxxx-xxxxxx"
gcloud auth login
gcloud services enable compute.googleapis.com cloudscheduler.googleapis.com logging.googleapis.com cloudfunctions.googleapis.com eventarc.googleapis.com run.googleapis.com calendar-json.googleapis.com firestore.googleapis.com file.googleapis.com --project "${PROJECT_ID}"