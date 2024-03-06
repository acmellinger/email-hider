terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.19.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = "us-central1"
}

resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "bucket" {
  name                        = "${random_id.bucket_prefix.hex}-gcf-source" # Every bucket name must be globally unique
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function.zip"
  bucket = google_storage_bucket.bucket.name
  source = "function.zip" # Add path to the zipped function source code
}

resource "google_cloudfunctions2_function" "function" {
  name        = "email-hider"
  location    = "us-central1"
  description = "email-hider"

  build_config {
    runtime     = "go120"
    entry_point = "handleRequest" #
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    environment_variables = var.env
  }

  lifecycle {
    replace_triggered_by = [
      google_storage_bucket_object.object
    ]
  }
}

output "function_uri" {
  value = google_cloudfunctions2_function.function.service_config[0].uri
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloudfunctions2_function.function.location
  service  = google_cloudfunctions2_function.function.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
  depends_on = [google_cloudfunctions2_function.function]
}
