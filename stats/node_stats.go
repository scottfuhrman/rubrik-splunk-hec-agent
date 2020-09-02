package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/src/golang/utils"
	"encoding/json"
)

// NodeIOStatsBody - return interface for node IO stats
type NodeStatsBody struct {
	ClusterName 			string		`json:"clusterName"`
	NodeId					string		`json:"id"`
	Time					int64		`json:"time"`
	Status					string		`json:"status"`
	IPAddress				string		`json:"ipAddress"`
	IsTunnelEnabled			bool		`json:"isTunnelEnabled"`
	CPUStat					float64		`json:"cpuStat"`
	ReadsPerSecond			float64		`json:"readsPerSecond"`
	WritesPerSecond			float64		`json:"writesPerSecond"`
	ReadBytePerSecond		float64		`json:"readBytePerSecond"`
	WriteBytePerSecond		float64		`json:"writeBytePerSecond"`
	BytesTransmitted		float64		`json:"bytesTransmitted"`
	BytesReceived			float64		`json:"bytesReceived"`
}

// GetNodeStats ...
func GetNodeStats(rubrik *rubrikcdm.Credentials, clustername string) []string {
	response := []string{}
	nodeList,err := rubrik.Get("internal","/node")
	if err != nil {
		log.Println("Error from stats.GetNodeStats: ",err)
		return []string{}
	}
	for _, node := range nodeList.(map[string]interface{})["data"].([]interface{}) {	
		thisNodeStats,err := rubrik.Get("internal","/node/"+node.(map[string]interface{})["id"].(string)+"/stats")
		if err != nil {
			log.Println("Error from stats.GetNodeStats: ",err)
			return []string{}
		}
		thisNodeNetworkStats := thisNodeStats.(map[string]interface{})["networkStat"]
		iopsData := thisNodeStats.(map[string]interface{})["iops"]
		ioThroughputData := thisNodeStats.(map[string]interface{})["ioThroughput"]
		thisNodeCpuStats := thisNodeStats.(map[string]interface{})["cpuStat"].([]interface{})
		readPerSecData := iopsData.(map[string]interface{})["readsPerSecond"].([]interface{})
		writePerSecData := iopsData.(map[string]interface{})["writesPerSecond"].([]interface{})
		readBpsData := ioThroughputData.(map[string]interface{})["readBytePerSecond"].([]interface{})
		writeBpsData := ioThroughputData.(map[string]interface{})["writeBytePerSecond"].([]interface{})
		bytesReceivedData := thisNodeNetworkStats.(map[string]interface{})["bytesReceived"].([]interface{})
		bytesTransmittedData := thisNodeNetworkStats.(map[string]interface{})["bytesTransmitted"].([]interface{})
		// if one of these slices is empty we will return an empty string
		if len(readPerSecData) == 0 || len(writePerSecData) == 0 || len(readBpsData) == 0 || len(writeBpsData) == 0 {
			continue
		}
		thisEntry := NodeStatsBody{
			ClusterName: 			clustername,
			NodeId:					node.(map[string]interface{})["id"].(string),
			Time:					utils.ConvertRubrikTimeToUnixTime(readPerSecData[0].(map[string]interface{})["time"].(string)),
			ReadsPerSecond:  		readPerSecData[len(readPerSecData) - 1].(map[string]interface{})["stat"].(float64),
			WritesPerSecond:  		writePerSecData[len(writePerSecData) - 1].(map[string]interface{})["stat"].(float64),
			ReadBytePerSecond:  	readBpsData[len(readBpsData) - 1].(map[string]interface{})["stat"].(float64),
			WriteBytePerSecond:  	writeBpsData[len(writeBpsData) - 1].(map[string]interface{})["stat"].(float64),
			Status:					thisNodeStats.(map[string]interface{})["status"].(string),
			IPAddress:				thisNodeStats.(map[string]interface{})["ipAddress"].(string),
			IsTunnelEnabled:		thisNodeStats.(map[string]interface{})["supportTunnel"].(map[string]interface{})["isTunnelEnabled"].(bool),
			CPUStat:				thisNodeCpuStats[len(thisNodeCpuStats) - 1].(map[string]interface{})["stat"].(float64),
			BytesTransmitted:		bytesReceivedData[len(bytesReceivedData) - 1].(map[string]interface{})["stat"].(float64),
			BytesReceived:			bytesTransmittedData[len(bytesTransmittedData) - 1].(map[string]interface{})["stat"].(float64),
		}
		json, err := json.Marshal(thisEntry)
		if err != nil {
			log.Println("Error from stats.GetNodeStats: ",err)
			return []string{}
		}
		response = append(response,string(json))
	}
	return response
}
