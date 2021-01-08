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

// ApRadioStruct radio struct
type ApRadioStruct struct {
	Band    float64 `json:"band"`
	Index   float64 `json:"index"`
	MacAddr string  `json:"macaddr"`
	Status  string  `json:"status"`
}

// ApStruct access point struct
type ApStruct struct {
	ApDeploymentMode string          `json:"ap_deployment_mode"`
	ClientCount      float64         `json:"client_count"`
	APGroup          string          `json:"ap_group"`
	ClusterID        string          `json:"cluster_id"`
	Labels           []string        `json:"labels"`
	Model            string          `json:"model"`
	Radios           []ApRadioStruct `json:"radios"`
	Serial           string          `json:"serial"`
	FirmwareVersion  string          `json:"firmware_version"`
	IPAddress        string          `json:"ip_address"`
	LastModified     float64         `json:"last_modified"`
	MeshRole         string          `json:"mesh_role"`
	Name             string          `json:"name"`
	Status           string          `json:"status"`
	Macaddr          string          `json:"macaddr"`
	Notes            string          `json:"notes"`
	PublicIPAddress  string          `json:"public_ip_address"`
	SubnetMask       string          `json:"subnet_mask"`
	GroupName        string          `json:"group_name"`
	Site             string          `json:"site"`
	SwarmID          string          `json:"swarm_id"`
	SwarmMaster      bool            `json:"swarm_master"`
}

type apiApResultStruct struct {
	Aps   []ApStruct `json:"aps"`
	Count int        `json:"count"`
}

// Getaps get a list of devices
func Getaps(token TokenStruct, verbose bool) []ApStruct {
	bearer := fmt.Sprintf("Bearer %v", token.AccessToken)
	url := fmt.Sprintf("%s/monitoring/v1/aps?calculate_client_count=true", URIApi)
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
	var BodyMap apiApResultStruct
	//var BodyMap interface{}
	err = json.Unmarshal(body, &BodyMap)
	if err != nil {
		panic(err.Error())
	}
	if verbose {
		for _, v := range BodyMap.Aps {
			log.Printf("%s, %s, %s, %s\n", v.APGroup, v.ClientCount, v.Name, v.Site)
		}
	}
	//var BodyMapfake []ApStruct
	return BodyMap.Aps

}
