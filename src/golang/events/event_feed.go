package events

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"encoding/json"
	"time"
	"strings"
	"strconv"
)

// EventBody - return interface for events
type EventBody struct {
	ClusterName 	string		`json:"clusterName"`
	Id				string		`json:"id"`
	ObjectId		string		`json:"objectId"`
	ObjectName		string		`json:"objectName"`
	Message			string		`json:"message"`
	EventStatus		string		`json:"eventStatus"`
	EventType		string		`json:"eventType"`
	LocationName	string		`json:"locationName"`
	Username		string		`json:"username"`
	OrgName			string		`json:"orgName"`
	OrgId			string		`json:"orgId"`
	Hostname		string		`json:"hostname"`
	Time			int64		`json:"time"`
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
		eventList,err := rubrik.Get("internal","/event?limit="+eventLimit+"&after_date="+fAfterDate,30)
		if err != nil {
			log.Fatal(err)
		}
		eventDataArray := eventList.(map[string]interface{})["data"].([]interface{})
		for event := range eventDataArray {
			eventStatus := eventDataArray[event].(map[string]interface{})["eventStatus"].(string)
			// extract and marshall the JSON from the eventInfo field
			eventInfoJson := eventDataArray[event].(map[string]interface{})["eventInfo"].(string)
			var eventInfo map[string]interface{}
			json.Unmarshal([]byte(eventInfoJson), &eventInfo)
			// if the event status is of a type we want then build the event
			if IsStatusValid(eventStatus) {
				eventType := eventDataArray[event].(map[string]interface{})["eventType"].(string)
				thisEvent := EventBody{
					ClusterName: 	clustername,
					Id:				eventDataArray[event].(map[string]interface{})["id"].(string),
					ObjectId:		eventDataArray[event].(map[string]interface{})["objectId"].(string),
					Message:		eventInfo["message"].(string),
					EventStatus:	eventStatus,
					EventType:		eventType,
					Time:			convertRubrikTimeToUnixTime(eventDataArray[event].(map[string]interface{})["time"].(string)),
				}
				if _, ok := eventDataArray[event].(map[string]interface{})["objectName"]; ok {
					thisEvent.ObjectName = eventDataArray[event].(map[string]interface{})["objectName"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${locationName}"]; ok {
					thisEvent.LocationName = eventInfo["params"].(map[string]interface{})["${locationName}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${username}"]; ok {
					thisEvent.Username = eventInfo["params"].(map[string]interface{})["${username}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${orgName}"]; ok {
					thisEvent.OrgName = eventInfo["params"].(map[string]interface{})["${orgName}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${orgId}"]; ok {
					thisEvent.OrgId = eventInfo["params"].(map[string]interface{})["${orgId}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${hostname}"]; ok {
					thisEvent.LocationName = eventInfo["params"].(map[string]interface{})["${hostname}"].(string)
				}
				if eventType == "Recovery" {
					eventSeriesId := eventDataArray[event].(map[string]interface{})["eventSeriesId"].(string)
					eventSeriesData,err := rubrik.Get("internal","/event_series/"+eventSeriesId)
					if err != nil {
						log.Fatal(err)
					}
					if _, ok := eventSeriesData.(map[string]interface{})["username"]; ok {
						thisEvent.Username = eventSeriesData.(map[string]interface{})["username"].(string)
					}	
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
		eventDataArray := eventList.(map[string]interface{})["data"].([]interface{})
		for event := range eventDataArray {
			eventStatus := eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["eventStatus"].(string)
			// extract and marshall the JSON from the eventInfo field
			eventInfoJson := eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["eventInfo"].(string)
			var eventInfo map[string]interface{}
			json.Unmarshal([]byte(eventInfoJson), &eventInfo)
			if IsStatusValid(eventStatus) {
				eventType := eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["eventType"].(string)
				thisEvent := EventBody{
					ClusterName: 	clustername,
					Id:				eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["id"].(string),
					ObjectId:		eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["objectId"].(string),
					Message:		eventInfo["message"].(string),
					EventStatus:	eventStatus,
					EventType:		eventType,
					Time:			convertRubrikTimeToUnixTime(eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["time"].(string)),
				}
				if _, ok := eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["objectName"]; ok {
					thisEvent.ObjectName = eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["objectName"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${locationName}"]; ok {
					thisEvent.LocationName = eventInfo["params"].(map[string]interface{})["${locationName}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${username}"]; ok {
					thisEvent.Username = eventInfo["params"].(map[string]interface{})["${username}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${orgName}"]; ok {
					thisEvent.OrgName = eventInfo["params"].(map[string]interface{})["${orgName}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${orgId}"]; ok {
					thisEvent.OrgId = eventInfo["params"].(map[string]interface{})["${orgId}"].(string)
				}
				if _, ok := eventInfo["params"].(map[string]interface{})["${hostname}"]; ok {
					thisEvent.LocationName = eventInfo["params"].(map[string]interface{})["${hostname}"].(string)
				}
				if eventType == "Recovery" {
					eventSeriesId := eventDataArray[event].(map[string]interface{})["latestEvent"].(map[string]interface{})["eventSeriesId"].(string)
					eventSeriesData,err := rubrik.Get("v1","/event_series/"+eventSeriesId)
					if err != nil {
						log.Fatal(err)
					}
					if _, ok := eventSeriesData.(map[string]interface{})["username"]; ok {
						thisEvent.Username = eventSeriesData.(map[string]interface{})["username"].(string)
					}	
				}

				json, err := json.Marshal(thisEvent)
				if err != nil {
					log.Fatal(err)
				}
				response = append(response,string(json))
			}
		}
	}
	return response
}

// checks if an event status is something we should process
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

// converts rubrik timestamp (RFC3339 format) to an int64 epoch time
func convertRubrikTimeToUnixTime(RubrikTime string) int64 {
	parsedTime, e := time.Parse(time.RFC3339, RubrikTime)
	if e != nil {
		log.Fatal(e)
	}
	return parsedTime.Unix()
}