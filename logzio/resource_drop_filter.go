package logzio

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/drop_filters"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	dropFilterIdField         = "drop_filter_id"
	dropFilterActive          = "active"
	dropFilterLogType         = "log_type"
	dropFilterFieldConditions = "field_conditions"
	dropFilterFieldName       = "field_name"
	dropFilterValue           = "value"
)

// Returns the drop filters client with the api token from the provider
func dropFilterClient(m interface{}) *drop_filters.DropFiltersClient {
	var client *drop_filters.DropFiltersClient
	client, _ = drop_filters.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceDropFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceDropFilterCreate,
		Read:   resourceDropFilterRead,
		Update: resourceDropFilterUpdate,
		Delete: resourceDropFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			dropFilterIdField: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			dropFilterActive: {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			dropFilterLogType: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			dropFilterFieldConditions: {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dropFilterFieldName: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						dropFilterValue: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

// Creates a new drop filter in logzio
func resourceDropFilterCreate(d *schema.ResourceData, m interface{}) error {
	createDropFilter := createCreatDropFilterFromSchema(d)
	active, isSet := d.GetOkExists(dropFilterActive)
	if isSet {
		log.Printf("active attribute is set to %t, note that this field is ignored for creation. A drop filter will always be active after creation.\n", active)
	}
	dropFilter, err := dropFilterClient(m).CreateDropFilter(createDropFilter)
	if err != nil {
		return err
	}

	d.SetId(dropFilter.Id)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceDropFilterRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "could not find drop filter with id") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

// Gets drop filter by id
func resourceDropFilterRead(d *schema.ResourceData, m interface{}) error {
	dropFilters, err := dropFilterClient(m).RetrieveDropFilters()
	if err != nil {
		return err
	}

	dropFilter := findDropFilterById(d.Id(), dropFilters)
	if dropFilter == nil {
		return fmt.Errorf("could not find drop filter with id: %s", d.Id())

	}

	setDropFilter(d, dropFilter)
	return nil
}

// Updates drop field by id - activate or deactivate
func resourceDropFilterUpdate(d *schema.ResourceData, m interface{}) error {
	activate := d.Get(dropFilterActive).(bool)
	var err error
	if activate {
		_, err = dropFilterClient(m).ActivateDropFilter(d.Id())
	} else {
		_, err = dropFilterClient(m).DeactivateDropFilter(d.Id())
	}

	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceDropFilterRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "could not find drop filter with id") {
				return resource.RetryableError(err)
			}
		}

		dropFilterFromSchema := createDropFilterFromSchema(d)
		if activate != dropFilterFromSchema.Active {
			return resource.RetryableError(fmt.Errorf("drop filter %s was not updated yet", d.Id()))
		}

		return resource.NonRetryableError(err)
	})
}

// Deletes drop filter by id
func resourceDropFilterDelete(d *schema.ResourceData, m interface{}) error {
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := dropFilterClient(m).DeleteDropFilter(d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func setDropFilter(d *schema.ResourceData, dropFilter *drop_filters.DropFilter) {
	d.Set(dropFilterIdField, dropFilter.Id)
	d.Set(dropFilterActive, dropFilter.Active)
	d.Set(dropFilterLogType, dropFilter.LogType)

	fieldConditions := getFieldConditionsMapping(dropFilter.FieldCondition)
	d.Set(dropFilterFieldConditions, fieldConditions)
}

func getFieldConditionsMapping(conditions []drop_filters.FieldConditionObject) []map[string]interface{} {
	var conditionsMapping []map[string]interface{}

	for _, condition := range conditions {
		mapping := map[string]interface{}{
			dropFilterFieldName: condition.FieldName,
			dropFilterValue:     convertObjectToString(condition.Value),
		}

		conditionsMapping = append(conditionsMapping, mapping)
	}

	return conditionsMapping
}

func createCreatDropFilterFromSchema(d *schema.ResourceData) drop_filters.CreateDropFilter {
	fieldConditionsFromSchema := d.Get(dropFilterFieldConditions).([]interface{})
	fieldConditions := getFieldConditionsList(fieldConditionsFromSchema)

	return drop_filters.CreateDropFilter{
		LogType:         d.Get(dropFilterLogType).(string),
		FieldConditions: fieldConditions,
	}
}

func getFieldConditionsList(conditionsFromSchemas []interface{}) []drop_filters.FieldConditionObject {
	var fieldConditions []drop_filters.FieldConditionObject
	var conditionToAppend drop_filters.FieldConditionObject
	for _, element := range conditionsFromSchemas {
		condition := element.(map[string]interface{})
		conditionToAppend.FieldName = condition[dropFilterFieldName].(string)
		conditionToAppend.Value = getValueObjectByType(condition[dropFilterValue].(string))
		fieldConditions = append(fieldConditions, conditionToAppend)
	}

	return fieldConditions
}

func getValueObjectByType(value string) interface{} {
	// will try to parse in this order: json, float, int, bool, string
	var returnObject map[string]interface{}
	err := json.Unmarshal([]byte(value), &returnObject)
	if err == nil {
		return returnObject
	}

	returnFloat, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return returnFloat
	}

	returnInt, err := strconv.Atoi(value)
	if err == nil {
		return returnInt
	}

	returnBool, err := strconv.ParseBool(value)
	if err == nil {
		return returnBool
	}

	return value
}

func convertObjectToString(value interface{}) string {
	switch value.(type) {
	case map[string]interface{}:
		byteArray, _ := json.Marshal(value)
		return string(byteArray)
	case string:
		return value.(string)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func createDropFilterFromSchema(d *schema.ResourceData) drop_filters.DropFilter {
	fieldConditionsFromSchema := d.Get(dropFilterFieldConditions).([]interface{})
	fieldConditions := getFieldConditionsList(fieldConditionsFromSchema)
	id, ok := d.GetOk(dropFilterIdField)
	if !ok {
		id = d.Id()
	}

	return drop_filters.DropFilter{
		Id:             id.(string),
		Active:         d.Get(dropFilterActive).(bool),
		LogType:        d.Get(dropFilterLogType).(string),
		FieldCondition: fieldConditions,
	}
}

func findDropFilterById(dropFilterId string, dropFilters []drop_filters.DropFilter) *drop_filters.DropFilter {
	for _, filter := range dropFilters {
		if filter.Id == dropFilterId {
			return &filter
		}
	}

	return nil
}
