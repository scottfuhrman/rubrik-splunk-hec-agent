package main

import (
	"log"
	"time"
	"os"
	"net/http"
	"crypto/tls"
	"encoding/json"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/src/golang/stats"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/src/golang/events"
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
			log.Panic("The `",envVar,"` environment variable is not present")
		}
	}
	// set our Splunk variables
	splunkToken, _ := os.LookupEnv("SPLUNK_HEC_TOKEN")
	splunkURL, _ := os.LookupEnv("SPLUNK_URL")
	splunkIndex, _ := os.LookupEnv("SPLUNK_INDEX")
	// create Rubrik client
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Panic(err)
	}
	// get cluster name, also tests connection before we go any further
	clusterDetails,err := rubrik.Get("v1","/cluster/me")
	if err != nil {
		log.Panic(err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"].(string)
	log.Printf("Cluster name: %s",clusterName)
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
				log.Panic(err)
			}
			log.Printf("Posted rubrik:storagesummary event.")
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()
	// get our cluster IO stats
	go func() {
		for {
			clusterIoStats := stats.GetClusterIOStats(rubrik,clusterName)
			if len(clusterIoStats) > 0 {
				var statsDetails map[string]interface{}
				json.Unmarshal([]byte(clusterIoStats), &statsDetails)		
				err := splunkClient.LogEvent(&splunk.Event{
					int64(statsDetails["time"].(float64)),
					clusterName,
					"rubrikhec",
					"rubrik:clusteriostats",
					splunkIndex,
					clusterIoStats,
				})
				if err != nil {
					log.Panic(err)
				}
			}
			log.Printf("Posted rubrik:clusteriostats event.")
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
			if err != nil {
				log.Panic(err)
			}
			log.Printf("Posted rubrik:runwayremaining event.")
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()
	// go get event feed
	go func() {
		for {
			eventList := events.GetEventFeed(rubrik,clusterName)
			if len(eventList) > 0 {
				for event := range eventList {
					var eventDetails map[string]interface{}
					json.Unmarshal([]byte(eventList[event]), &eventDetails)		
					err := splunkClient.LogEvent(&splunk.Event{
						int64(eventDetails["time"].(float64)),
						clusterName,
						"rubrikhec",
						"rubrik:eventfeed",
						splunkIndex,
						eventList[event],
					})
					if err != nil {
						log.Panic(err)
					}
				}
			}
			log.Printf("Posted %d rubrik:eventfeed events.",len(eventList))
			time.Sleep(time.Duration(20) * time.Minute)
		}
	}()
	// go get org capacity report stats
	go func() {
		for {
			reportEntryList := stats.GetOrgCapacityReports(rubrik,clusterName)
			if len(reportEntryList) > 0 {
				for reportEntry := range reportEntryList {
					err := splunkClient.LogEvent(&splunk.Event{
						time.Now().Unix(),
						clusterName,
						"rubrikhec",
						"rubrik:orgcapacityreport",
						splunkIndex,
						reportEntryList[reportEntry],
					})
					if err != nil {
						log.Panic(err)
					}
				}
			}
			log.Printf("Posted %d rubrik:orgcapacityreport events.",len(reportEntryList))
			time.Sleep(time.Duration(4) * time.Hour)
		}
	}()
	// go get man vol summary stats
	go func() {
		for {
			mvSummary := stats.GetManVolSummaryStats(rubrik,clusterName)
			if len(mvSummary) > 0 {
				err := splunkClient.LogEvent(&splunk.Event{
					time.Now().Unix(),
					clusterName,
					"rubrikhec",
					"rubrik:manvolsummary",
					splunkIndex,
					mvSummary,
				})
				if err != nil {
					log.Panic(err)
				}
				log.Printf("Posted rubrik:manvolsummary event.")
			} else {
				log.Printf("No data received for rubrik:manvolsummary.")
			}
			time.Sleep(time.Duration(4) * time.Hour)
		}
	}()
	// go get archive location usage stats
	go func() {
		for {
			archiveLocationList := stats.GetArchiveLocationUsageStats(rubrik,clusterName)
			if len(archiveLocationList) > 0 {
				for archiveEntry := range archiveLocationList {
					err := splunkClient.LogEvent(&splunk.Event{
						time.Now().Unix(),
						clusterName,
						"rubrikhec",
						"rubrik:archivelocationusage",
						splunkIndex,
						archiveLocationList[archiveEntry],
					})
					if err != nil {
						log.Panic(err)
					}
				}
				log.Printf("Posted %d rubrik:archivelocationusage events.",len(archiveLocationList))
			} else {
				log.Printf("No data received for rubrik:archivelocationusage.")
			}
			time.Sleep(time.Duration(4) * time.Hour)
		}
	}()
	// go get archive location bandwidth stats
	go func() {
		for {
			archiveBandwidthStats := stats.GetArchiveLocationBandwidthStats(rubrik,clusterName)
			if len(archiveBandwidthStats) > 0 {
				for _, statEntry := range archiveBandwidthStats {
					var statDetails map[string]interface{}
					json.Unmarshal([]byte(statEntry), &statDetails)	
					err := splunkClient.LogEvent(&splunk.Event{
						int64(statDetails["time"].(float64)),
						clusterName,
						"rubrikhec",
						"rubrik:archivelocationbandwidth",
						splunkIndex,
						statDetails,
					})
					if err != nil {
						log.Panic(err)
					}
				}
			}
			log.Printf("Posted %d rubrik:archivelocationbandwidth event.", len(archiveBandwidthStats))
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()
	// keep application open until terminated
	for {
		time.Sleep(time.Duration(1) * time.Hour)
	}
}
