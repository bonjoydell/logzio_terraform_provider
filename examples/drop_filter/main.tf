variable "api_token" {
  type = string
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_drop_filter" "test_filter" {
  log_type = "some_type"

  field_conditions {
    field_name = "some_field"
    value = "string_value"
  }
  field_conditions {
    field_name = "other_field_int"
    value = 200
  }
}
