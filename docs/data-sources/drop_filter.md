# Drop Filter Resource

Provides a Logz.io drop filter resource. This can be used to create and manage Logz.io drop filters.

* Learn more about drop filters in the [Logz.io Docs](https://docs.logz.io/api/#tag/Drop-filters).

## Argument Reference
* `field_conditions` - Filters for an exact match of a field:value pair. **Note** that changing this field after creation will cause the resource to be destroyed and re-created. See below for **nested schema**.
* `log_type` - Filters for the [log type](https://docs.logz.io/user-guide/log-shipping/built-in-log-types.html). Emit or leave empty if you want this filter to apply to all types. **Note** that changing this field after creation will cause the resource to be destroyed and re-created. 
* `drop_filter_id` - Drop filter ID in the Logz.io database.

## Attribute Reference
* `active` - If true, the drop filter is active and logs that match the filter are dropped before indexing. If false, the drop filter is disabled. **Note** this argument can only be changed after the creation of the filter. Each filter is created with the `active` argument set to true.
