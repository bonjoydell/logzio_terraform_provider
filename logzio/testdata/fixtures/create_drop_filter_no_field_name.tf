resource "logzio_drop_filter" "%s" {
  log_type = "some_type"

  field_conditions {
    value = "some_string_value"
  }
}