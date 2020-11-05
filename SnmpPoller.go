package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	b "math/big"
	"os"
	"time"

	g "github.com/gosnmp/gosnmp"
)

type Nodes struct {
	Nodes []Configuration `json:"nodes"`
}

type Configuration struct {
	NodeName  string `json:"node_name"`
	OagIP     string `json:"oag_ip"`
	Community string `json:"oag_cs"`
}

type OidResult struct {
	OagNode  string
	OidName  string
	Oid      string
	Response *b.Int
}

func SnmpPoller(config *Configuration) {

	var OidResultSet []OidResult

	mibs := map[string]string{
		".1.3.6.1.4.1.2021.4.5.0":  "memory_total_installed",
		".1.3.6.1.4.1.2021.4.6.0":  "memory_total_used",
		".1.3.6.1.4.1.2021.4.11.0": "memory_total_free",
		".1.3.6.1.4.1.2021.4.13.0": "memory_total_shared",
		".1.3.6.1.4.1.2021.4.14.0": "memory_total_buffered",
		".1.3.6.1.4.1.2021.4.15.0": "memory_total_cached",
		".1.3.6.1.4.1.2021.4.3.0":  "swap_mem_total_size",
		".1.3.6.1.4.1.2021.4.4.0":  "swap_mem_available",
		//	".1.3.6.1.2.1.25.1.1.0":    "system_uptime",
		//".1.3.6.1.4.1.2021.2.*.1":    "session_cache_service",
		".1.3.6.1.2.1.31.1.1.1.6.3":  "network_interface_in",
		".1.3.6.1.2.1.31.1.1.1.10.3": "network_interface_out",
		//".1.3.6.1.4.1.2021.9.1.2.1":  "disk_path",
		".1.3.6.1.4.1.2021.9.1.6.1": "disk_total_size",
		".1.3.6.1.4.1.2021.9.1.7.1": "disk_total_available",
		".1.3.6.1.4.1.2021.9.1.9.1": "disk_percentage_used",
		//	".1.3.6.1.2.1.1.1.0":         "sys_obj_oag_appliance",
		//	".1.3.6.1.2.1.1.4.0":         "sys_obj_oag_support",
		//	".1.3.6.1.2.1.1.5.0":         "sys_obj_dev",
	}

	oids := make([]string, 0, len(mibs))
	for key, _ := range mibs {
		oids = append(oids, key)
	}

	config := setConfig()
	for i := 0; i < len(config.Nodes); i++ {
		fmt.Println("oag_ip: " + config.Nodes[i].NodeName)
		fmt.Println("oag_ip: " + config.Nodes[i].OagIP)
		fmt.Println("oag_cs: " + config.Nodes[i].Community)
		//initialize
		oidResult := OidResult{}
		var err error
		var err2 error
		//node_name := config.Nodes[i].NodeName
		g.Default.Target = config.Nodes[0].OagIP
		g.Default.Community = config.Nodes[0].Community
		g.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking

		err = g.Default.Connect()
		if err != nil {
			log.Fatalf("Connect() err: %v", err)
		}
		defer g.Default.Conn.Close()
		//var result g.SnmpPacket
		result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err2 != nil {
			log.Fatalf("Get() err: %v", err2)
		}
		oidResult.OagNode = config.Nodes[i].NodeName
		for _, variable := range result.Variables {
			//fmt.Printf("%d: oid: %s ", i, variable.Name)
			//fmt.Println("i ", i)
			//fmt.Println("variableName", variable.Name)

			oidResult.Oid = variable.Name
			//fmt.Println("oidResult.Oid", oidResult.Oid)
			oidResult.OidName = mibs[variable.Name]
			//fmt.Println("oidResult.OidName", oidResult.OidName)
			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {
			case g.OctetString:
				bytes := variable.Value.([]byte)

				fmt.Printf("string: %s\n", string(bytes))

			default:
				// ... or often you're just interested in numeric values.
				// ToBigInt() will return the Value as a BigInt, for plugging
				// into your calculations.
				oidResult.Response = g.ToBigInt(variable.Value)
				OidResultSet = append(OidResultSet, oidResult)
				//fmt.Printf("number: %d\n", oidResult.Response)

			}
			//fmt.Printf("OidResultSet: %d\n", OidResultSet[i].Response)
		}
	}
	fmt.Println("length of ResultSet", len(OidResultSet))
	// Marshal the map into a JSON string.
	nodeData, err := json.Marshal(OidResultSet)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonStr := string(nodeData)
	fmt.Println("JSON data is:")
	fmt.Println(jsonStr)
}

func setConfig() *Nodes {

	configFile, err := os.Open("output/config/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	var nodes Nodes

	json.Unmarshal(byteValue, &nodes)

	return &nodes
}
