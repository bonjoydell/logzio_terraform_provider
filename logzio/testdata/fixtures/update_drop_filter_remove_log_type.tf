resource "logzio_drop_filter" "%s" {
  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
