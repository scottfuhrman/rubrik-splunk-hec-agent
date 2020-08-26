package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
)

// ManVolSummaryBody - return interface for managed volume summary
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

// GetManVolSummaryStats ...
func GetArchiveLocationUsageStats(rubrik *rubrikcdm.Credentials, clustername string) []string {
	response := []string{}

	archiveLocationUsage,err := rubrik.Get("internal","/stats/data_location/usage")
	if err != nil {
		log.Fatal(err)
	}
	for i := range archiveLocationUsage.(map[string]interface{})["data"].([]interface{}) {
		thisLocation := archiveLocationUsage.(map[string]interface{})["data"].([]interface{})[i].(map[string]interface{})
		thisEntry := ArchiveLocationUsageBody{
			ClusterName: 							clustername,
			LocationName:							GetLocationNameById(rubrik, thisLocation["locationId"].(string)),
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
			log.Fatal(err)
		}
		response = append(response,string(json))
	}
	return response
}

// returns an archive location name from it's ID
func GetLocationNameById(rubrik *rubrikcdm.Credentials, locationId string) string {
	archiveLocationSummary,err := rubrik.Get("internal","/archive/location")
	if err != nil {
		log.Fatal(err)
	}
	for j := range archiveLocationSummary.(map[string]interface{})["data"].([]interface{}) {
		thisLocation := archiveLocationSummary.(map[string]interface{})["data"].([]interface{})[j].(map[string]interface{})
		if thisLocation["id"] == locationId {
			return thisLocation["name"].(string)
		}
	}
	return ""
}