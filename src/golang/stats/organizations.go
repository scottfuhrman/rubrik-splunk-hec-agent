package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/rubrikinc/rubrik-splunk-hec-agent/src/golang/utils"
	"encoding/json"
	"strconv"
	"time"
)

type OrgCapacityReportBody struct {
	ClusterName							string	`json:"clusterName"`
	Month								string	`json:"month"`
    OrganizationId						string	`json:"organizationId"`
    Organization						string	`json:"organization"`
	OrganizationState					string	`json:"organizationState"`
    LocalStorage						int64	`json:"localStorage"`
    LocalLogicalBytes					int64	`json:"localLogicalStorage"`
    ArchiveStorage						int64	`json:"archiveStorage"`
    LocalDataReductionPercent			float64	`json:"localDataReductionPercent"`
    LocalLogicalDataReductionPercent	float64	`json:"localLogicalDataReductionPercent"`
    ReplicaStorage						int64	`json:"replicaStorage"`
    LocalStorageGrowth					int64	`json:"localStorageGrowth"`
    ReplicaStorageGrowth				int64	`json:"replicaStorageGrowth"`
    ArchiveStorageGrowth				int64	`json:"archiveStorageGrowth"`
}

// GetRunwayRemaining ...
func GetOrgCapacityReports(rubrik *rubrikcdm.Credentials, clustername string) []string {
	response := []string{}
	reportList,err := rubrik.Get("internal","/report?report_template=CapacityOverTime&report_type=Canned")
	if err != nil {
		log.Println("Error from stats.GetOrgCapacityReports: ",err)
		return []string{}
	}
	reportId := reportList.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	body := map[string]interface{}{
		"limit": 100,
		"sortBy": "Month",
		"sortOrder": "desc",
	}
	reportDetails,err := rubrik.Post("internal","/report/"+reportId+"/table",body)
	thisMonth := utils.LeftPad(strconv.Itoa(int(time.Now().Month())),"0",2)
	thisYear := strconv.Itoa(int(time.Now().Year()))
	currentMonthStr := thisYear + "-" + thisMonth
	reportDataGrid := reportDetails.(map[string]interface{})["dataGrid"].([]interface{})
	reportColumns := reportDetails.(map[string]interface{})["columns"].([]interface{})
	for reportEntry := range reportDataGrid {
		thisReportEntryMapped := map[string]interface{}{ }
		for i := range reportColumns {
			thisReportEntryMapped[reportColumns[i].(string)] = reportDataGrid[reportEntry].([]interface{})[i]
		}
		if thisReportEntryMapped["Month"] == currentMonthStr {
			thisEntry := OrgCapacityReportBody{
				ClusterName: 						clustername,
				Month:								thisReportEntryMapped["Month"].(string),
				OrganizationId:						thisReportEntryMapped["OrganizationId"].(string),
				Organization:						thisReportEntryMapped["Organization"].(string),
				OrganizationState:					thisReportEntryMapped["OrganizationState"].(string),
				LocalStorage:						utils.ConvertToInt64(thisReportEntryMapped["LocalStorage"].(string)),
				LocalLogicalBytes:					utils.ConvertToInt64(thisReportEntryMapped["LocalLogicalBytes"].(string)),
				ArchiveStorage:						utils.ConvertToInt64(thisReportEntryMapped["ArchiveStorage"].(string)),
				LocalDataReductionPercent:			utils.ConvertToFloat64(thisReportEntryMapped["LocalDataReductionPercent"].(string)),
				LocalLogicalDataReductionPercent:	utils.ConvertToFloat64(thisReportEntryMapped["LocalLogicalDataReductionPercent"].(string)),
				ReplicaStorage:						utils.ConvertToInt64(thisReportEntryMapped["ReplicaStorage"].(string)),
				LocalStorageGrowth:					utils.ConvertToInt64(thisReportEntryMapped["LocalStorageGrowth"].(string)),
				ReplicaStorageGrowth:				utils.ConvertToInt64(thisReportEntryMapped["ReplicaStorageGrowth"].(string)),
				ArchiveStorageGrowth:				utils.ConvertToInt64(thisReportEntryMapped["ArchiveStorageGrowth"].(string)),
			}
			json, err := json.Marshal(thisEntry)
			if err != nil {
				log.Println("Error from stats.GetOrgCapacityReports: ",err)
				return []string{}
			}
			response = append(response,string(json))
		}
	}
	return response
}
