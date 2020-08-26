package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
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
		log.Fatal(err)
	}
	reportId := reportList.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	body := map[string]interface{}{
		"limit": 100,
		"sortBy": "Month",
		"sortOrder": "desc",
	}
	reportDetails,err := rubrik.Post("internal","/report/"+reportId+"/table",body)
	thisMonth := LeftPad(strconv.Itoa(int(time.Now().Month())),"0",2)
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
				LocalStorage:						ConvertToInt64(thisReportEntryMapped["LocalStorage"].(string)),
				LocalLogicalBytes:					ConvertToInt64(thisReportEntryMapped["LocalLogicalBytes"].(string)),
				ArchiveStorage:						ConvertToInt64(thisReportEntryMapped["ArchiveStorage"].(string)),
				LocalDataReductionPercent:			ConvertToFloat64(thisReportEntryMapped["LocalDataReductionPercent"].(string)),
				LocalLogicalDataReductionPercent:	ConvertToFloat64(thisReportEntryMapped["LocalLogicalDataReductionPercent"].(string)),
				ReplicaStorage:						ConvertToInt64(thisReportEntryMapped["ReplicaStorage"].(string)),
				LocalStorageGrowth:					ConvertToInt64(thisReportEntryMapped["LocalStorageGrowth"].(string)),
				ReplicaStorageGrowth:				ConvertToInt64(thisReportEntryMapped["ReplicaStorageGrowth"].(string)),
				ArchiveStorageGrowth:				ConvertToInt64(thisReportEntryMapped["ArchiveStorageGrowth"].(string)),
			}
			json, err := json.Marshal(thisEntry)
			if err != nil {
				log.Fatal(err)
			}
			response = append(response,string(json))
		}
	}
	return response
}

// pads a string with `pad` to length `plength`
func LeftPad(s string, pad string, plength int) string {
    for i := len(s); i < plength; i++ {
        s = pad + s
    }
    return s
}

func ConvertToFloat64(s string) float64 {
	c,err := strconv.ParseFloat(s,64)
	if (err != nil) {
		log.Fatal(err)
	}
	return c
}

func ConvertToInt64(s string) int64 {
	c,err := strconv.ParseInt(s,10,64)
	if (err != nil) {
		log.Fatal(err)
	}
	return c
}