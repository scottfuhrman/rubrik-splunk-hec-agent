package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
)

// ClusterIOStatsBody - return interface for cluster IO stats
type ManVolSummaryBody struct {
	ClusterName 			string		`json:"clusterName"`
	Count					float64		`json:"count"`
	Exported				int			`json:"exported"`
	Writable				int			`json:"writable"`
	Nfs						int			`json:"nfs"`
	Smb						int			`json:"smb"`
}

// GetClusterIOStats ...
func GetManVolSummaryStats(rubrik *rubrikcdm.Credentials, clustername string) string {
	mvSummary,err := rubrik.Get("internal","/managed_volume?is_relic=false&primary_cluster_id=local")
	if err != nil {
		log.Fatal(err)
	}
	countExported := 0
	countWritable := 0
	countNfs := 0
	countSmb := 0
	for i := range mvSummary.(map[string]interface{})["data"].([]interface{}) {
		thisMv := mvSummary.(map[string]interface{})["data"].([]interface{})[i].(map[string]interface{})
		if thisMv["state"] == "Exported" {
			countExported++
		}
		if thisMv["isWritable"].(bool) {
			countWritable++
		}
		if thisMv["shareType"] == "NFS" {
			countNfs++
		} else {
			countSmb++
		}
	}
	response := ManVolSummaryBody{
		ClusterName: 			clustername,
		Count:  				mvSummary.(map[string]interface{})["total"].(float64),
		Exported:				countExported,
		Writable:				countWritable,
		Nfs:					countNfs,
		Smb:					countSmb,
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}
