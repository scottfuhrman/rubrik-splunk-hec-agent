package main

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
    "encoding/json"
)

// GetRunwayRemaining ...
func GetRunwayRemaining(rubrik *rubrikcdm.Credentials) string {
	runwayRemaining,err := rubrik.Get("internal","/stats/runway_remaining")
	if err != nil {
		log.Fatal(err)
	}
	json, err := json.Marshal(runwayRemaining)
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}

// GetStorageSummary ...
func GetStorageSummary(rubrik *rubrikcdm.Credentials) string {
	storageStats,err := rubrik.Get("internal","/stats/system_storage")
	if err != nil {
		log.Fatal(err)
	}
	json, err := json.Marshal(storageStats)
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}