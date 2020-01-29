package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"net/http"
	"crypto/tls"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/stats"
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk"
)

func main() {
	// check we got env variables
	envVarsList := [...]string{
		"SPLUNK_HEC_TOKEN",
		"SPLUNK_URL",
		"SPLUNK_INDEX",
	}
	for _, envVar := range envVarsList {
		_, ok := os.LookupEnv(envVar)
		if ok != true {
			log.Fatal("The `",envVar,"` environment variable is not present")
		}
	}
	// set our Splunk variables
	splunkToken, _ := os.LookupEnv("SPLUNK_HEC_TOKEN")
	splunkURL, _ := os.LookupEnv("SPLUNK_URL")
	splunkIndex, _ := os.LookupEnv("SPLUNK_INDEX")
	// create Rubrik client
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}
	// get cluster name, also tests connection before we go any further
	clusterDetails,err := rubrik.Get("v1","/cluster/me")
	if err != nil {
		log.Fatal(err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"].(string)
	fmt.Println("Cluster name: "+clusterName)
	// create HTTP client (change InsecureSkipVerify to false if not using self-signed certs)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var httpClient *http.Client
	httpClient = &http.Client{Timeout: time.Second * 20, Transport: tr}
	// create HEC client
	splunkClient := splunk.NewClient(
		httpClient,
		splunkURL,
		splunkToken,
		"rubrikhec",
		"rubrikhec:default",
		splunkIndex,
	)
	// get our storage summary stats
	go func() {
		for {
			err := splunkClient.LogEvent(&splunk.Event{
				time.Now().Unix(),
				clusterName,
				"rubrikhec",
				"rubrik:storagesummary",
				splunkIndex,
				stats.GetStorageSummary(rubrik,clusterName),
			})
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()
	// get our runway remaining stats
	go func() {
		for {
			err := splunkClient.LogEvent(&splunk.Event{
				time.Now().Unix(),
				clusterName,
				"rubrikhec",
				"rubrik:runwayremaining",
				splunkIndex,
				stats.GetRunwayRemaining(rubrik,clusterName),
			})
			//err := splunkClient.Log(stats.GetRunwayRemaining(rubrik))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()
	// keep application open until terminated
	for {
		time.Sleep(time.Duration(1) * time.Hour)
	}
}