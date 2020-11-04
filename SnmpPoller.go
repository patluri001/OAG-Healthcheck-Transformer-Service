
package main
import (
	"fmt"
	//"log"
	"os"
    "time"
    g "github.com/gosnmp/gosnmp"
)
func main() {

	g.Default.Target = "192.168.1.12"
	g.Default.Community = "Ro4OAG4tw4yM0n1t0r1ng"
	g.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	/*err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()*/
	//var m map[string]string
	mibs := map[string]string{
		"sys_description": ".1.3.6.1.2.1.1.1.0",
		"sys_contact":   ".1.3.6.1.2.1.1.4.0",
		"sys_name":		".1.3.6.1.2.1.1.5.0",
		"sys_location":	".1.3.6.1.2.1.1.6.0",
		"memory_total": ".1.3.6.1.4.1.2021.4.5.0",
		"memory_total_used": ".1.3.6.1.4.1.2021.4.6.0",
	}

	var oids []string
	for key, value := range mibs {
		fmt.Println("Key:", key, "Value:", value)
		oids = append(oids,value)
	}

	/*topmibs := map[string]string{
		"objects": ".1.3.6.1.4.1.2021.10",
		"memory":   ".1.3.6.1.4.1.2021.4",
		"network":	".1.3.6.1.2.1.31.1",
		"disk":	".1.3.6.1.4.1.2021.9",
	}*/

	err := g.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}

	defer g.Default.Conn.Close()

	//var topoids []string
	for key, value := range mibs {
		fmt.Println("Key:", key, "Value:", value)
		//topoids = append(topoids,value)
		err = g.Default.BulkWalk(value, printValue)
		if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
	}

	/*err := g.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer g.Default.Conn.Close()

	err = g.Default.BulkWalk("", printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}*/



	//oids := []string{".1.3.6.1.4.1.2021.4.5.0", ".1.3.6.1.4.1.2021.4.6.0",".1.3.6.1.4.1.2021.9.1.2",".1.3.6.1.4.1.2021.4.3.0",".1.3.6.1.4.1.2021.9.1.6"}
	//oids := []string{".1.3.6.1.4.1.2021.4", ".1.3.6.1.4.1.2021.9",".1.3.6.1.2.1.31.1",".1.3.6.1.4.1.2021.10"}
	/*result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)
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
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}*/
}


func printValue(pdu g.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)
	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
	}
	return nil
}
