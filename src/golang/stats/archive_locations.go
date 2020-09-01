package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/src/golang/utils"
	"encoding/json"
)

// ArchiveLocationUsageBody - return interface for archive location usage
type ArchiveLocationUsageBody struct {
	ClusterName 						string		`json:"clusterName"`
	LocationName						string		`json:"locationName"`
	LocationId							string		`json:"locationId"`
	DataDownloaded						float64		`json:"dataDownloaded"`
	DataArchived						float64		`json:"dataArchived"`
	NumVMsArchived						float64		`json:"numVMsArchived"`
	NumFilesetsArchived					float64		`json:"numFilesetsArchived"`
	NumLinuxFilesetsArchived			float64		`json:"numFilesetsArchived"`
	NumShareFilesetsArchived			float64		`json:"numShareFilesetsArchived"`
	NumMssqlDbsArchived					float64		`json:"numMssqlDbsArchived"`
	NumHypervVmsArchived				float64		`json:"numHypervVmsArchived"`
	NumNutanixVmsArchived				float64		`json:"numNutanixVmsArchived"`
	NumManagedVolumesArchived			float64		`json:"numManagedVolumesArchived"`
	NumStorageArrayVolumeGroupsArchived	float64		`json:"numStorageArrayVolumeGroupsArchived"`
	NumWindowsVolumeGroupsArchived		float64		`json:"numWindowsVolumeGroupsArchived"`
}

// ArchiveLocationBandwidthBody - return interface for archive location bandwidth
type ArchiveLocationBandwidthBody struct {
	ClusterName 						string		`json:"clusterName"`
	LocationName						string		`json:"locationName"`
	LocationId							string		`json:"locationId"`
	Time								int64		`json:"time"`
	Type								string		`json:"type"`
	Value								float64		`json:"value"`
}

// GetArchiveLocationUsageStats ...
func GetArchiveLocationUsageStats(rubrik *rubrikcdm.Credentials, clustername string) []string {
	response := []string{}

	archiveLocationUsage,err := rubrik.Get("internal","/stats/data_location/usage")
	if err != nil {
		log.Println("Error from stats.GetArchiveLocationUsageStats: ",err)
		return []string{}
	}
	for _, archiveLocation := range archiveLocationUsage.(map[string]interface{})["data"].([]interface{}) {
		thisLocation := archiveLocation.(map[string]interface{})
		thisEntry := ArchiveLocationUsageBody{
			ClusterName: 							clustername,
			LocationName:							utils.GetLocationNameById(rubrik, thisLocation["locationId"].(string)),
			LocationId:								thisLocation["locationId"].(string),
			DataDownloaded:							thisLocation["dataDownloaded"].(float64),
			DataArchived:							thisLocation["dataArchived"].(float64),
			NumVMsArchived:							thisLocation["numVMsArchived"].(float64),
			NumFilesetsArchived:					thisLocation["numFilesetsArchived"].(float64),
			NumLinuxFilesetsArchived:				thisLocation["numLinuxFilesetsArchived"].(float64),
			NumShareFilesetsArchived:				thisLocation["numShareFilesetsArchived"].(float64),
			NumMssqlDbsArchived:					thisLocation["numMssqlDbsArchived"].(float64),
			NumHypervVmsArchived:					thisLocation["numHypervVmsArchived"].(float64),
			NumNutanixVmsArchived:					thisLocation["numNutanixVmsArchived"].(float64),
			NumManagedVolumesArchived:				thisLocation["numManagedVolumesArchived"].(float64),
			NumStorageArrayVolumeGroupsArchived:	thisLocation["numStorageArrayVolumeGroupsArchived"].(float64),
			NumWindowsVolumeGroupsArchived:			thisLocation["numWindowsVolumeGroupsArchived"].(float64),
		}
		json, err := json.Marshal(thisEntry)
		if err != nil {
			log.Println("Error from stats.GetArchiveLocationUsageStats: ",err)
			return []string{}
		}
		response = append(response,string(json))
	}
	return response
}

// GetArchiveLocationBandwidthStats ...
func GetArchiveLocationBandwidthStats(rubrik *rubrikcdm.Credentials, clustername string) []string {
	response := []string{}

	archiveLocationSummary,err := rubrik.Get("internal","/archive/location")
	if err != nil {
		log.Println("Error from stats.GetArchiveLocationBandwidthStats: ",err)
		return []string{}
	}

	for _, archiveLocation := range archiveLocationSummary.(map[string]interface{})["data"].([]interface{}) {
		bwTypes := [...]string{"Incoming","Outgoing",}
		for _, bwType := range bwTypes {
			locationId := archiveLocation.(map[string]interface{})["id"].(string)
			bwTimeSeries,err := rubrik.Get("internal","/stats/archival/bandwidth/time_series?data_location_id="+locationId+"&range=-15min&bandwidth_type="+bwType)
			if err != nil {
				log.Println("Error from stats.GetArchiveLocationBandwidthStats: ",err)
				return []string{}
			}
			timeSeriesData := bwTimeSeries.([]interface{})
			thisEntry := ArchiveLocationBandwidthBody{
				ClusterName: 	clustername,
				LocationName:	archiveLocation.(map[string]interface{})["name"].(string),
				LocationId:		locationId,
				Type:			bwType,
				Time:			utils.ConvertRubrikTimeToUnixTime(timeSeriesData[0].(map[string]interface{})["time"].(string)),
				Value:			timeSeriesData[0].(map[string]interface{})["stat"].(float64),		
			}
			json, err := json.Marshal(thisEntry)
			if err != nil {
				log.Println("Error from stats.GetArchiveLocationBandwidthStats: ",err)
				return []string{}
			}
			response = append(response,string(json))
		}
	}
	return response
}
