resource "google_storage_bucket" "gcs" {
  name          = "test"
  project       = "yokobot-dev"
  location      = "asia"
  force_destroy = true
}