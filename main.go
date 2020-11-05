package main

import (
	"io/ioutil"
	"net/http"

	//"time"
	"encoding/json"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	// b "math/big"
)

var log = logrus.New()

type Configuration struct {
	NodeName  string `json:"node_name"`
	OagIP     string `json:"oag_ip"`
	Community string `json:"oag_cs"`
}

type ConfigurationSet []Configuration

type Nodes struct {
	Nodes []Configuration `json:"nodes"`
}

type OidResult struct {
	OagNode  string
	OidName  string
	Oid      string
	Response string
}

type OidResultSet struct {
	oidResult OidResult
}

type HttpResponse struct {
	OidResultSet []string
}

// func init() {
// 	log.SetLevel(log.InfoLevel)
// 	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
// }

func main() {
	// Prometheus: Histogram to collect required metrics
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "greeting_seconds",
		Help:    "Time take to greet someone",
		Buckets: []float64{1, 2, 5, 6, 10}, //defining small buckets as this app should not take more than 1 sec to respond
	}, []string{"code"}) // this will be partitioned by the HTTP code.

	router := mux.NewRouter()
	router.Handle("/metrics", Sayhello(histogram))
	// router.Handle("/metrics", prometheus.Handler())

	//Registering the defined metric with Prometheus
	prometheus.Register(histogram)

	log.Fatal(http.ListenAndServe(":9443", router))
}

func TriggerJob() []string {
	configSet := setConfig()

	//var OidResultSet []string
	var httpResponse []OidResultSet

	for _, config := range configSet.Nodes {
		log.Info("get Healthcheck data for OAG IP:", config.OagIP)
		OidResultSet1 := SnmpPoller(&config)
		OidResultSet = append(OidResultSet, OidResultSet1)
	}
	return OidResultSet1
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}

func setConfig() Nodes {

	configFile, err := os.Open("config/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	var nodes Nodes

	json.Unmarshal(byteValue, &nodes)

	return nodes
}
