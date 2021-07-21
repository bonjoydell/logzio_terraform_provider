package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/drop_filters"
	"log"
	"reflect"
)

func dataSourceDropFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDropFilterRead,
		Schema: map[string]*schema.Schema{
			dropFilterIdField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			dropFilterActive: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			dropFilterLogType: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			dropFilterFieldConditions: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dropFilterFieldName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						dropFilterValue: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDropFilterRead(d *schema.ResourceData, m interface{}) error {
	client, _ := drop_filters.New(m.(Config).apiToken, m.(Config).baseUrl)
	dropFilterId, ok := d.GetOk(dropFilterIdField)
	dropFilters, err := client.RetrieveDropFilters()
	if err != nil {
		return err
	}

	if len(dropFilters) == 0 {
		log.Println("no drop filters in account")
		return nil
	}

	if ok {
		dropFilter := findDropFilterById(dropFilterId.(string), dropFilters)
		if dropFilter != nil {
			d.SetId(dropFilterId.(string))
			setDropFilter(d, dropFilter)
			return nil
		}
	}

	// If could not find id, search for drop filter with the specified attributes
	dropFilterToSearch := createDropFilterFromSchema(d)
	for _, filter := range dropFilters {
		if isSameDropFilter(dropFilterToSearch, filter) {
			d.SetId(filter.Id)
			setDropFilter(d, &filter)
			return nil
		}
	}

	return fmt.Errorf("couldn't find drop filter with specified attributes")
}

func isSameDropFilter(dropFilterToSearch drop_filters.DropFilter, dropFilter drop_filters.DropFilter) bool {
	return reflect.DeepEqual(dropFilterToSearch.FieldCondition, dropFilter.FieldCondition) &&
		reflect.DeepEqual(dropFilterToSearch.LogType, dropFilter.LogType)
}
