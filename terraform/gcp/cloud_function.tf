data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = "src"
  output_path = "output/functions.zip"
}

resource "google_storage_bucket" "src_bucket" {
  name          = "${var.storage_name}-src"
  location      = "asia"
  storage_class = "STANDARD"
}

resource "google_storage_bucket_object" "packages" {
  name   = "packages/functions.${data.archive_file.function_archive.output_md5}.zip"
  bucket = google_storage_bucket.src_bucket.name
  source = data.archive_file.function_archive.output_path
}

resource "google_cloudfunctions_function" "finalize_function" {
  name                  = "${var.storage_name}-finalize-function"
  description           = "${var.storage_name}-finalize-function"
  runtime               = "python39"
  source_archive_bucket = google_storage_bucket.src_bucket.name
  source_archive_object = google_storage_bucket_object.packages.name
  available_memory_mb   = 128
  timeout               = 120

  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = google_storage_bucket.gcs.name
  }
}

resource "google_cloudfunctions_function" "delete_function" {
  name                  = "${var.storage_name}-delete-function"
  description           = "${var.storage_name}-delete-function"
  runtime               = "python39"
  source_archive_bucket = google_storage_bucket.src_bucket.name
  source_archive_object = google_storage_bucket_object.packages.name
  available_memory_mb   = 128
  timeout               = 120

  event_trigger {
    event_type = "google.storage.object.delete"
    resource   = google_storage_bucket.gcs.name
  }
}