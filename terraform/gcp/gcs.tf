resource "google_storage_bucket" "gcs" {
  name          = var.storage_name
  project       = "yokobot-dev"
  location      = "asia"
  force_destroy = true
}