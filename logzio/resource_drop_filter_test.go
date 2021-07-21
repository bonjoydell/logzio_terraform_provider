package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

const (
	dropFilterResourceCreateDropFilter                  = "create_drop_filter"
	dropFilterResourceCreateDropFilterNoFieldConditions = "create_drop_filter_no_field_conditions"
	dropFilterResourceCreateDropFilterNoFieldName       = "create_drop_filter_no_field_name"
	dropFilterResourceCreateDropFilterNoValue           = "create_drop_filter_no_value"
	dropFilterResourceCreateDropFilterEmptyLogType      = "create_drop_filter_empty_log_type"
	dropFilterResourceUpdateDropFilter                  = "update_drop_filter"
	dropFilterResourceUpdateDropFilterChangeLogType     = "update_drop_filter_change_log_type"
	dropFilterResourceUpdateDropFilterRemoveLogType     = "update_drop_filter_remove_log_type"
)

func TestAccLogzioDropFilter_CreateDropFilter(t *testing.T) {
	filterName := "test_create_drop_filter"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config:            resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioDropFilter_CreateDropEmptyLogType(t *testing.T) {
	filterName := "test_create_drop_filter_empty_log_type"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilterEmptyLogType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, ""),
				),
			},
			{
				Config:            resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioDropFilter_UpdateDropFilter(t *testing.T) {
	filterName := "test_update_drop_filter"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropFilterActive, "true"),
				),
			},
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceUpdateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropFilterActive, "false"),
				),
			},
		},
	})
}

func TestAccLogzioDropFilter_UpdateDropFilterChangeLogType(t *testing.T) {
	filterName := "test_update_drop_filter_change_log_type"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceUpdateDropFilterChangeLogType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "other_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
		},
	})
}

func TestAccLogzioDropFilter_UpdateDropFilterRemoveLogType(t *testing.T) {
	filterName := "test_update_drop_filter_remove_log_type"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceUpdateDropFilterRemoveLogType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, ""),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
		},
	})
}

func TestAccLogzioDropFilter_CreateDropFilterNoFieldConditions(t *testing.T) {
	filterName := "test_create_drop_filter_no_field_conditions"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilterNoFieldConditions),
				ExpectError: regexp.MustCompile("required field is not set"),
			},
		},
	})
}

func TestAccLogzioDropFilter_CreateDropFilterNoFieldName(t *testing.T) {
	filterName := "test_create_drop_filter_no_field_conditions"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilterNoFieldName),
				ExpectError: regexp.MustCompile("no definition was found"),
			},
		},
	})
}

func TestAccLogzioDropFilter_CreateDropFilterNoValue(t *testing.T) {
	filterName := "test_create_drop_filter_no_field_conditions"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestDropFilter(filterName, dropFilterResourceCreateDropFilterNoValue),
				ExpectError: regexp.MustCompile("no definition was found"),
			},
		},
	})
}

func resourceTestDropFilter(name string, path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}
