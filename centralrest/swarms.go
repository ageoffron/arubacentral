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

type SwarmStruct struct {
	FirmwareVersion string `json:"firmware_version"`
	GroupName       string `json:"group_name"`
	IPAddress       string `json:"ip_address"`
	Name            string `json:"name"`
	PublicIPAddress string `json:"public_ip_address"`
	Status          string `json:"status"`
	SwarmID         string `json:"swarm_id"`
}

type apiResultStruct struct {
	Count  int           `json:"count"`
	Swarms []SwarmStruct `json:"swarms"`
}

// Getswarms qweqe
func Getswarms(token TokenStruct, verbose bool) []SwarmStruct {
	//type output struct{}

	bearer := fmt.Sprintf("Bearer %v", token.AccessToken)
	url := fmt.Sprintf("%s/monitoring/v1/swarms", URIApi)
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
		if verbose {
			log.Printf("Device list Error %v", resp.Status)
		}
		panic(err)
	}

	var BodyMap apiResultStruct
	err = json.Unmarshal([]byte(body), &BodyMap)
	if err != nil {
		err = errors.New("Swarm Unmarshal error")
		panic(err)
	}
	if verbose {
		for _, v := range BodyMap.Swarms {
			log.Printf("%s, %s, %s, %s, %s\n", v.GroupName, v.FirmwareVersion, v.Name, v.Name, v.Status)
		}
	}
	return BodyMap.Swarms
}
