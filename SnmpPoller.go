package main

import (
	"fmt"
	"encoding/json"
	"log"
	"os"
	"time"
    g "github.com/gosnmp/gosnmp"
)

type Configuration struct {
	OagIP string `json:"oag_ip"`
	Community string `json:"oag_cs"`
}

type OidResult struct {
	OidName string
	Oid 	string
	Response  string
}





func main() {

	oidResult := OidResult{}

	var OidResultSet []OidResult

	oids := []string{".1.3.6.1.4.1.2021.4.5.0", ".1.3.6.1.4.1.2021.4.6.0"}

	// mibs := map[string]string{
	// 	".1.3.6.1.4.1.2021.4.5.0": "memory_total",
	// 	".1.3.6.1.4.1.2021.4.6.0": "memory_total_used",
	// }

	config := setConfig()

	g.Default.Target = config.OagIP
	g.Default.Community = config.Community
	g.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking

	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		oidResult.Oid = variable.Name

		// oidResult.OidName = mibs(variable.Name)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.OctetString:
			bytes := variable.Value.([]byte)
			fmt.Printf("string: %s\n", string(bytes))

			oidResult.Response = string(bytes)

		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			OidResultSet = append(OidResultSet, oidResult)

			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
}

func setConfig() (*Configuration){

	configFile, err := os.Open("config/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	scConfig := Configuration{}
	err = decoder.Decode(&scConfig)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &scConfig

}