package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country_name"`
	Region   string `json:"region_name"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

func GetLocationFromIP(ip string) (*LocationInfo, error) {

	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info LocationInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	info.IP = ip
	return &info, nil
}
