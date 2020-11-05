package main


import (
	// "strconv"
	// "io/ioutil"
	"fmt"
	"net/http"
	"time"
	// "os"
	"github.com/prometheus/client_golang/prometheus"
)


func Sayhello(histogram *prometheus.HistogramVec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var returnString string
		//monitoring how long it takes to respond
		start := time.Now()
		defer r.Body.Close()
		code := 500

		defer func() {
			httpDuration := time.Since(start)
			histogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(httpDuration.Seconds())
		}()

		code = http.StatusBadRequest // if req is not GET
		if r.Method == "GET" {
			code = http.StatusOK

			// event, err := os.Open("config/data.txt")

			httpResp := TriggerJob()

			for _,oneNodeResult := range httpResp {
				returnString += oneNodeResult
			}

			// if err != nil {
			// 	log.Fatal(err)
			// }
		
			// defer event.Close()
		
			// byteValue, _ := ioutil.ReadAll(ioutil.Read)
			// var nodes Nodes

			b := []byte(returnString)
			w.Write(b)
		} else {
			w.WriteHeader(code)
		}
	}
}