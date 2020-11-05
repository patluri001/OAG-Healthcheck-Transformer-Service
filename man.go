package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)



type Configuration struct {
	OagIP string `json:"oag_ip"`
	Community string `json:"oag_cs"`
}

type ConfigurationSet []Configuration

func main() {


	fmt.Println("inside Main Method")


	configSet := setConfig()

	for _, config := range configSet {
		fmt.Println("get Healthcheck data for OAG IP:", config.OagIP)
		SnmpPoller(&config)
	}


}


func setConfig() (ConfigurationSet){

	configFile, err := os.Open("config/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	scConfig := ConfigurationSet{}
	err = decoder.Decode(&scConfig)
	if err != nil {
		fmt.Println("error:", err)
	}

	return scConfig

}
