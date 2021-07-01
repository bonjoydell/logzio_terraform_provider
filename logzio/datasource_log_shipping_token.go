package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/log_shipping_tokens"
	"strconv"
)

func dataSourceLogShippingToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogShippingTokenRead,
		Schema: map[string]*schema.Schema{
			logShippingTokenId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenEnabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			logShippingTokenToken: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenUpdatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenCreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLogShippingTokenRead(d *schema.ResourceData, m interface{}) error {
	client, _ := log_shipping_tokens.New(m.(Config).apiToken, m.(Config).baseUrl)
	tokenIdString, ok := d.GetOk(logShippingTokenId)

	if ok {
		id, err := strconv.Atoi(tokenIdString.(string))
		if err != nil {
			return err
		}

		token, err := client.GetLogShippingToken(int32(id))
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("%d", id))
		setLogShippingToken(d, token)

		return nil
	}

	// If for some reason we couldn't find the token by id,
	// looking for the token by it's name
	tokenName, ok := d.GetOk(logShippingTokenEnabled)
	if ok {
		tokenEnabled, ok := d.GetOk(logShippingTokenEnabled)
		if ok {
			token, err := findLogShippingTokenByName(tokenName.(string), tokenEnabled.(bool), client)
			if err != nil {
				return err
			}

			if token != nil {
				d.SetId(fmt.Sprintf("%d", token.Id))
				setLogShippingToken(d, token)
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find log shipping token with specified attributes")
}

func findLogShippingTokenByName(name string, enabled bool, client *log_shipping_tokens.LogShippingTokensClient) (*log_shipping_tokens.LogShippingToken, error) {
	retrieveRequest := log_shipping_tokens.RetrieveLogShippingTokensRequest{
		Filter: log_shipping_tokens.ShippingTokensFilterRequest{Enabled: strconv.FormatBool(enabled)},
		Pagination: log_shipping_tokens.ShippingTokensPaginationRequest{
			PageNumber: 1,
			PageSize:   25,
		},
	}

	tokenFound, total, totalRetrieved, err := findToken(name, client, retrieveRequest)
	if err != nil {
		return nil, err
	}

	if tokenFound != nil {
		return tokenFound, nil
	}

	// Pagination
	for total > totalRetrieved {
		retrieveRequest.Pagination.PageNumber += 1
		tokenFound, _, currentlyRetrieved, err := findToken(name, client, retrieveRequest)
		if err != nil {
			return nil, err
		}

		if tokenFound != nil {
			return tokenFound, nil
		}

		totalRetrieved += currentlyRetrieved
	}

	return nil, fmt.Errorf("couldn't find log shipping token with specified attributes")
}

func findTokenInResultsListByName(name string, tokens []log_shipping_tokens.LogShippingToken) *log_shipping_tokens.LogShippingToken {
	for _, token := range tokens {
		if token.Name == name {
			return &token
		}
	}

	return nil
}

func findToken(name string, client *log_shipping_tokens.LogShippingTokensClient, request log_shipping_tokens.RetrieveLogShippingTokensRequest) (*log_shipping_tokens.LogShippingToken, int, int, error) {
	tokens, err := client.RetrieveLogShippingTokens(request)
	if err != nil {
		return nil, 0, 0, err
	}

	tokenFound := findTokenInResultsListByName(name, tokens.Results)

	return tokenFound, int(tokens.Total), len(tokens.Results), nil
}
