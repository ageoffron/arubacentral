package centralrest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//type output struct{}

// DeviceStruct device struct
type DeviceStruct struct {
	ArubaPartNo  string `json:"aruba_part_no"`
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	DeviceType   string `json:"device_type"`
	Macaddr      string `json:"macaddr"`
	Model        string `json:"model"`
	Serial       string `json:"serial"`
}

type apiDeviceResultStruct struct {
	Devices []DeviceStruct `json:"devices"`
	Total   int            `json:"total"`
}

// Getdevices get a list of devices
func Getdevices(token TokenStruct, verbose bool) []DeviceStruct {
	bearer := fmt.Sprintf("Bearer %v", token.AccessToken)
	url := fmt.Sprintf("%s/device_inventory/v2/devices?limit=200&offset=0&sku_type=IAP", URIApi)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		panic(err)
	}
	var BodyMap apiDeviceResultStruct
	err = json.Unmarshal(body, &BodyMap)
	if err != nil {
		panic(err.Error())
	}
	if verbose {
		for _, v := range BodyMap.Devices {
			log.Printf("%s, %s, %s, %s\n", v.DeviceType, v.ArubaPartNo, v.Macaddr, v.Model)
		}
	}
	return BodyMap.Devices

}
