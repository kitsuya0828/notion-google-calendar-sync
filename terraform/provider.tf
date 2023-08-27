provider "google" {
  project = local.project_id
  region  = local.region
}

data "google_project" "default" {
}

data "google_compute_default_service_account" "default" {
}