package main

import (
	//"time"
	"os"
	"encoding/json"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Configuration struct {
	OagIP string `json:"oag_ip"`
	Community string `json:"oag_cs"`
}

type ConfigurationSet []Configuration


// func init() {
// 	log.SetLevel(log.InfoLevel)
// 	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
// }



func main() {

	// The API for setting attributes is a little different than the package level
	  // exported logger. See Godoc.
	  

  	log.Out = os.Stdout

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
	 log.Out = file
	} else {
	 log.Info("Failed to log to file, using default stderr")
	}

	log.Info("Create new cron")
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() { 
		// TriggerJob() 
		log.Info("[Job 1]Every minute job\n") })

	// Start cron with one scheduled job
	log.Info("Start cron")
	c.Start()
	printCronEntries(c.Entries())
	for {

	}
	// time.Sleep(2 * time.Minute)

	// // Funcs may also be added to a running Cron
	// log.Info("Add new job to a running cron")
	// entryID2, _ := c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	// printCronEntries(c.Entries())
	// time.Sleep(5 * time.Minute)

	// //Remove Job2 and add new Job2 that run every 1 minute
	// log.Info("Remove Job2 and add new Job2 with schedule run every minute")
	// c.Remove(entryID2)
	// c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
	// time.Sleep(5 * time.Minute)

}

func TriggerJob() {
	configSet := setConfig()

	for _, config := range configSet {
		log.Info("get Healthcheck data for OAG IP:", config.OagIP)
		SnmpPoller(&config)
}
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
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
			log.Info("error:", err)
		}
	
		return scConfig
	
	}


// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"github.com/robfig/cron/v3"
// )




// func main() {


// 	c := cron.New()
// 	fmt.Printf("Created New Cron Entry")
// 	c.AddFunc("0 1 * * * *", func() { fmt.Println("Every minute job")} )

// 	c.Start()

// 	fmt.Printf("cron started")
// 	// printCronEntries(c.Entries())
// 	// time.Sleep(2 * time.Minute)



// 	// fmt.Println("inside Main Method")



// 	// }


// }



// 


// func printCronEntries(cronEntries []cron.Entry) {
// 	fmt.Printf("Cron Info: %+v\n", cronEntries)
// }