# notion-google-calendar-sync

**notion-google-calendar-sync** is a tool for two-way synchronisation between Notion and Google Calendar.

## Features
![notion-google-calendar-sync](https://github.com/Kitsuya0828/notion-google-calendar-sync/assets/60843722/9ec42b2e-9151-4ce4-9647-6c52fe3bc0ac)

* Periodically monitor and synchronize Notion and Google Calendar events
* The tool is deployed to Google Cloud, all using free tier products (Cloud Functions, Cloud FireStore, etc.)
* Terraform code is available

## Prerequisites
* [gcloud CLI](https://cloud.google.com/sdk/docs/install)
* [Terraform](https://developer.hashicorp.com/terraform/downloads)
* [Google Cloud Project](https://cloud.google.com/free)
* [Notion API Integration](https://www.notion.so/help/create-integrations-with-the-notion-api)
* [Google Calendar](https://calendar.google.com/)

## Deploy
Copy the template to `locals.tf` and edit it to match your Google Cloud Project configuration. Be especially careful that `bucket_name` must be globally unique.
```bash
cd terraform
cp locals.tf.tmp locals.tf
```

Enbale the Google Cloud APIs to be used.
You can enable the APIs automatically using Terraform, but it may take some time to be activated, so use the `gcloud` command.
```bash
# infra/init.sh
PROJECT_ID="xxxxxx-xxxxxxxx-xxxxxx" # Change Required
gcloud auth login
gcloud services enable compute.googleapis.com cloudscheduler.googleapis.com logging.googleapis.com cloudfunctions.googleapis.com eventarc.googleapis.com run.googleapis.com calendar-json.googleapis.com firestore.googleapis.com --project "${PROJECT_ID}"
```

Now we can finally deploy the tool to Google Cloud.
You can change it later on the Google Cloud console, but if it bothers you, you can change the runtime environment variables from [`infra/main.tf`](https://github.com/Kitsuya0828/notion-google-calendar-sync/blob/734cd9b1151176eeec4f13b72a536ba942aa2ea9/infra/main.tf#L79) before executing the following Terraform commands.

```bash
gcloud auth application-default login
terraform init
terraform plan
terraform apply
```
Once `terraform apply` is complete, you will see your service account email as follows:
```
Outputs:

service_account_email = "notion-google-calendar-sync@xxxxxx-xxxxxxxx-xxxxxx.iam.gserviceaccount.com"
```
Then, in your Google Calendar, remember to grant the appropriate permissions to the service account you have created.

![google-calendar-grant-permission](https://github.com/Kitsuya0828/notion-google-calendar-sync/assets/60843722/28b921b2-7a25-43a2-a429-04b013f65ab6)

## License
"notion-google-calendar-sync" is under [MIT License](https://opensource.org/license/mit/).
