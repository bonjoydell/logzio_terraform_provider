package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceDropFilter(t *testing.T) {
	filterNameResource := "test_create_drop_filter"
	resourceName := "logzio_drop_filter." + filterNameResource
	dataSourceName := "data.logzio_drop_filter.my_drop_filter_datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterNameResource, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceTestDropFilter(filterNameResource, dropFilterResourceCreateDropFilter) +
					testAccKubernetesDataSourceDropFilterById(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(dataSourceName, "field_conditions.#", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceDropFilterByAttributes(t *testing.T) {
	filterNameResource := "test_create_drop_filter"
	resourceName := "logzio_drop_filter." + filterNameResource
	dataSourceName := "data.logzio_drop_filter.my_drop_filter_datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterNameResource, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceTestDropFilter(filterNameResource, dropFilterResourceCreateDropFilter) +
					testAccKubernetesDataSourceDropFilterByAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(dataSourceName, "field_conditions.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "drop_filter_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccKubernetesDataSourceDropFilterById() string {
	return fmt.Sprintf(`data "logzio_drop_filter" "my_drop_filter_datasource" {
drop_filter_id = "${logzio_drop_filter.test_create_drop_filter.id}"
}
`)
}

func testAccKubernetesDataSourceDropFilterByAttributes() string {
	return fmt.Sprintf(`data "logzio_drop_filter" "my_drop_filter_datasource" {
  log_type = "some_type"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
`)
}
