package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"log"
)

func CreateTenant(req *types.Page) *types.TenantCreationResponse {
	lmsServiceUrl := GetEnv("LMS_SERVICE_URL", "")
	restyClient := resty.New()

	var tenantCreated *types.TenantCreationResponse
	post, err := restyClient.R().
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		SetBody(req).
		Post(fmt.Sprintf("%s%s", lmsServiceUrl, "/tenants/"))
	if err != nil {
		return nil
	}

	err = json.Unmarshal(post.Body(), &tenantCreated)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return tenantCreated
}
