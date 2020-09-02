package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
)

// RunwayBody - return interface for runway remaining
type RunwayBody struct {
	ClusterName string		`json:"clusterName"`
	Days		float64		`json:"days"`
}

// StorageSummaryBody - return interface for storage summary
type StorageSummaryBody struct {
	ClusterName 	string		`json:"clusterName"`
	Available		float64		`json:"available"`
	LiveMount		float64		`json:"liveMount"`
	Miscellaneous	float64		`json:"miscellaneous"`
	Snapshot		float64		`json:"snapshot"`
	Total			float64		`json:"total"`
	Used			float64		`json:"used"`
}

// GetRunwayRemaining ...
func GetRunwayRemaining(rubrik *rubrikcdm.Credentials, clustername string) string {
	runwayRemaining,err := rubrik.Get("internal","/stats/runway_remaining")
	if err != nil {
		log.Println("Error from stats.GetRunwayRemaining: ",err)
		return ""
	}
	response := RunwayBody{
		ClusterName: 	clustername,
		Days:  			runwayRemaining.(map[string]interface{})["days"].(float64),
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Println("Error from stats.GetRunwayRemaining: ",err)
		return ""
	}
	return string(json)
}

// GetStorageSummary ...
func GetStorageSummary(rubrik *rubrikcdm.Credentials, clustername string) string {
	storageStats,err := rubrik.Get("internal","/stats/system_storage")
	if err != nil {
		log.Println("Error from stats.GetStorageSummary: ",err)
		return ""
	}
	response := StorageSummaryBody{
		ClusterName: 	clustername,
		Available:		storageStats.(map[string]interface{})["available"].(float64),
		LiveMount:		storageStats.(map[string]interface{})["liveMount"].(float64),
		Miscellaneous:	storageStats.(map[string]interface{})["miscellaneous"].(float64),
		Snapshot:		storageStats.(map[string]interface{})["snapshot"].(float64),
		Total:			storageStats.(map[string]interface{})["total"].(float64),
		Used:			storageStats.(map[string]interface{})["used"].(float64),
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Println("Error from stats.GetStorageSummary: ",err)
		return ""
	}
	return string(json)
}