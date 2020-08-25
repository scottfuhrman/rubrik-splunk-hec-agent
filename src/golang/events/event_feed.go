package events

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	"fmt"
)

// EventBody - return interface for events
type EventBody struct {
	ClusterName string		`json:"clusterName"`
	Id			string		`json:"id"`
	ObjectId	string		`json:"objectId"`
	Message		string		`json:"message"`
}


// GetEventFeed ...
func GetEventFeed(rubrik *rubrikcdm.Credentials, clustername string) []string {

	// get cluster version
	clusterVersion,err := rubrik.ClusterVersion()
	if err != nil {
		log.Fatal(err)
	}
	clusterMajorVersion,err := strconv.ParseInt(strings.Split(clusterVersion,".")[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	clusterMinorVersion,err := strconv.ParseInt(strings.Split(clusterVersion,".")[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	response := []string{}
	// get the timestamp for 20 minutes ago
	now := time.Now()
	afterDate := now.UTC().Add(time.Minute * -20)
	fAfterDate := afterDate.Format(time.RFC3339)
	// set the max number of events to retrieve
	eventLimit := "9999"
	if (clusterMajorVersion == 5 && clusterMinorVersion < 2) || clusterMajorVersion < 5 { // cluster version is older than 5.1
		eventList,err := rubrik.Get("internal","/event?limit="+eventLimit+"&after_date="+fAfterDate)
		if err != nil {
			log.Fatal(err)
		}
		eventDataArray := eventList.(map[string]interface{})["data"].([]interface{})
		for event := range eventDataArray {
			fmt.Println(eventDataArray[event])
			eventStatus := eventDataArray[event].(map[string]interface{})["eventStatus"].(string)
			if IsStatusValid(eventStatus) {
				thisEvent := EventBody{
					ClusterName: 	clustername,
					Id:				eventDataArray[event].(map[string]interface{})["id"].(string),
					ObjectId:		eventDataArray[event].(map[string]interface{})["objectId"].(string),
					Message:		eventDataArray[event].(map[string]interface{})["eventInfo"].(map[string]interface{})["message"].(string),
				}
				json, err := json.Marshal(thisEvent)
				if err != nil {
					log.Fatal(err)
				}
				response = append(response,string(json))
			}
		}
	} else { // cluster version is 5.2 or newer
		eventList,err := rubrik.Get("v1","/event/latest?limit="+eventLimit+"&before_date="+fAfterDate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(eventList)
	}
	fmt.Println(response)
	return response
}

func IsStatusValid(status string) bool {
	switch status {
		case
			"Failure",
			"Warning",
			"Success",
			"Canceled":
			return true
	}
    return false
}