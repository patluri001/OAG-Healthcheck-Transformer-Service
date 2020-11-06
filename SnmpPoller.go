package main

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "log"
	"time"

	g "github.com/gosnmp/gosnmp"
)

func SnmpPoller(config *Configuration, OidResultSet *[]string) {

	//initialize
	oidResult := OidResult{}

	mibs := map[string]string{
		".1.3.6.1.4.1.2021.4.5.0":      "memory_total_installed",
		".1.3.6.1.4.1.2021.4.6.0":      "memory_total_used",
		".1.3.6.1.4.1.2021.4.11.0":     "memory_total_free",
		".1.3.6.1.4.1.2021.4.13.0":     "memory_total_shared",
		".1.3.6.1.4.1.2021.4.14.0":     "memory_total_buffered",
		".1.3.6.1.4.1.2021.4.15.0":     "memory_total_cached",
		".1.3.6.1.4.1.2021.4.3.0":      "swap_mem_total_size",
		".1.3.6.1.4.1.2021.4.4.0":      "swap_mem_available",
		".1.3.6.1.2.1.25.1.1.0":        "system_uptime",
		".1.3.6.1.2.1.1.1.0":           "system_description",
		".1.3.6.1.2.1.1.4.0":           "system_contact",
		".1.3.6.1.2.1.1.5.0":           "system_Location",
		".1.3.6.1.2.1.1.6.0":           "system_Name",
		".1.3.6.1.4.1.2021.10.1.3.1":   "sys_onemin_average",
		".1.3.6.1.4.1.2021.10.1.3.2":   "sys_fivemin_average",
		".1.3.6.1.4.1.2021.10.1.3.3":   "sys_fifteenmin_average",
		".1.3.6.1.2.1.31.1.1.1.6.1":    "network_interface_in_1",
		".1.3.6.1.2.1.31.1.1.1.6.2":    "network_interface_in_2",
		".1.3.6.1.2.1.31.1.1.1.6.3":    "network_interface_in_3",
		".1.3.6.1.2.1.31.1.1.1.10.1":   "network_interface_out_1",
		".1.3.6.1.2.1.31.1.1.1.10.2":   "network_interface_out_2",
		".1.3.6.1.2.1.31.1.1.1.10.3":   "network_interface_out_3",
		".1.3.6.1.4.1.2021.9.1.2.1":    "disk_path",
		".1.3.6.1.4.1.2021.9.1.5.1":    "disk_min_percentage",
		".1.3.6.1.4.1.2021.9.1.6.1":    "disk_total_size",
		".1.3.6.1.4.1.2021.9.1.7.1":    "disk_total_available",
		".1.3.6.1.4.1.2021.9.1.9.1":    "disk_percentage_used",
		".1.3.6.1.4.1.2021.16.2.1.1.1": "poll_session_logwatch_1",
		".1.3.6.1.4.1.2021.16.2.1.1.2": "poll_session_logwatch_2",
		".1.3.6.1.4.1.2021.16.2.1.1.3": "poll_session_logwatch_3",
		".1.3.6.1.4.1.2021.16.2.1.2.1": "poll_session_logwatch_session_db_connection",
		".1.3.6.1.4.1.2021.16.2.1.2.2": "poll_session_logwatch_session_db_storing",
		".1.3.6.1.4.1.2021.16.2.1.2.3": "poll_session_logwatch_session_db_get",
		".1.3.6.1.4.1.2021.2.1.5.3":    "poll_process_objects_HA",
		".1.3.6.1.4.1.2021.2.1.4.3":    "poll_process_objects_Time_Svc",
		".1.3.6.1.4.1.2021.2.1.3.3":    "poll_process_objects_web_proc_svc",
	}

	oids := make([]string, 0, len(mibs))
	for key, _ := range mibs {
		oids = append(oids, key)
	}

	fmt.Println("oag_ip: " + config.NodeName)
	fmt.Println("oag_ip: " + config.OagIP)
	fmt.Println("oag_cs: " + config.Community)

	var err error
	var err2 error
	//node_name := config.Nodes[i].NodeName
	g.Default.Target = config.OagIP
	g.Default.Community = config.Community
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
	oidResult.OagNode = config.NodeName
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
		//fmt.Println("variable.Type ", variable.Type)
		//fmt.Println("oidResult.Oid ", oidResult.Oid)

		switch variable.Type {
		case g.OctetString:
			//bytes := variable.Value.([]byte)
			value := variable.Value.(string)
			//fmt.Printf("string: %s\n", string(bytes))
			//fmt.Println("string: %s\n", value)
			*OidResultSet = append(*OidResultSet, oidResult.OagNode+"_"+oidResult.OidName+" "+value)
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			// oidResult.Response = g.ToBigInt(variable.Value)

			// oidResult.Response = variable.Value.(string)

			log.Info("Pinting the Value::", fmt.Sprint(g.ToBigInt(variable.Value)))
			*OidResultSet = append(*OidResultSet, oidResult.OagNode+"_"+oidResult.OidName+" "+fmt.Sprint(g.ToBigInt(variable.Value)))

			//log.Info("Length of ResultSet::",)
			//OidResultSet = append(OidResultSet, oidResult)
			//fmt.Printf("number: %d\n", oidResult.Response)

		}
		//fmt.Printf("OidResultSet: %d\n", OidResultSet[i].Response)
	}

	// fmt.Println("length of ResultSet", len(*OidResultSet))

}
