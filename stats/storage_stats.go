package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
	"fmt"
)

// RunwayBody - return interface for runway remaining
type RunwayBody struct {
	clusterName string
	days		float64
}

// StorageSummaryBody - return interface for storage summary
type StorageSummaryBody struct {
	clusterName 	string
	available		float64
	liveMount		float64
	miscellaneous	float64
	snapshot		float64
	total			float64
	used			float64
}

// GetRunwayRemaining ...
func GetRunwayRemaining(rubrik *rubrikcdm.Credentials, clustername string) string {
	runwayRemaining,err := rubrik.Get("internal","/stats/runway_remaining")
	if err != nil {
		log.Fatal(err)
	}
	response := &RunwayBody{
		clusterName: 	clustername,
		days:  			runwayRemaining.(map[string]interface{})["days"].(float64),
	}
	fmt.Println(response)
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json)
	return string(json)
}

// GetStorageSummary ...
func GetStorageSummary(rubrik *rubrikcdm.Credentials, clustername string) string {
	storageStats,err := rubrik.Get("internal","/stats/system_storage")
	if err != nil {
		log.Fatal(err)
	}
	response := &StorageSummaryBody{
		clusterName: 	clustername,
		available:		storageStats.(map[string]interface{})["available"].(float64),
		liveMount:		storageStats.(map[string]interface{})["liveMount"].(float64),
		miscellaneous:	storageStats.(map[string]interface{})["miscellaneous"].(float64),
		snapshot:		storageStats.(map[string]interface{})["snapshot"].(float64),
		total:			storageStats.(map[string]interface{})["total"].(float64),
		used:			storageStats.(map[string]interface{})["used"].(float64),
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json)
	return string(json)
}