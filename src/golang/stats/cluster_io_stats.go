package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
	"fmt"
)

// ClusterIOStatsBody - return interface for cluster IO stats
type ClusterIOStatsBody struct {
	ClusterName 			string		`json:"clusterName"`
	ReadsPerSecond			float64		`json:"readsPerSecond"`
	WritesPerSecond			float64		`json:"writesPerSecond"`
	ReadBytePerSecond		float64		`json:"readBytePerSecond"`
	WriteBytePerSecond		float64		`json:"writeBytePerSecond"`
}

// GetClusterIOStats ...
func GetClusterIOStats(rubrik *rubrikcdm.Credentials, clustername string) (string,string) {
	clusterIOStats,err := rubrik.Get("internal","/cluster/me/io_stats?range=-10min")
	if err != nil {
		log.Fatal(err)
	}
	iopsData := clusterIOStats.(map[string]interface{})["iops"]
	ioThroughputData := clusterIOStats.(map[string]interface{})["ioThroughput"]
	readPerSecData := iopsData.(map[string]interface{})["readsPerSecond"].([]interface{})
	writePerSecData := iopsData.(map[string]interface{})["writesPerSecond"].([]interface{})
	readBpsData := ioThroughputData.(map[string]interface{})["readBytePerSecond"].([]interface{})
	writeBpsData := ioThroughputData.(map[string]interface{})["writeBytePerSecond"].([]interface{})
	// if one of these slices is empty we will return an empty string
	if len(readPerSecData) == 0 || len(writePerSecData) == 0 || len(readBpsData) == 0 || len(writeBpsData) == 0 {
		fmt.Println(clusterIOStats)
		return "",""
	}
	timeStamp := readPerSecData[0].(map[string]interface{})["time"].(string)
	response := ClusterIOStatsBody{
		ClusterName: 			clustername,
		ReadsPerSecond:  		readPerSecData[0].(map[string]interface{})["stat"].(float64),
		WritesPerSecond:  		writePerSecData[0].(map[string]interface{})["stat"].(float64),
		ReadBytePerSecond:  	readBpsData[0].(map[string]interface{})["stat"].(float64),
		WriteBytePerSecond:  	writeBpsData[0].(map[string]interface{})["stat"].(float64),
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	return string(json),timeStamp
}
