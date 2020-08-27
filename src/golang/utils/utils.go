package utils

import (
	"time"
	"strconv"
	"log"
	"regexp"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

// returns an archive location name from it's ID
func GetLocationNameById(rubrik *rubrikcdm.Credentials, locationId string) string {
	archiveLocationSummary,err := rubrik.Get("internal","/archive/location")
	if err != nil {
		log.Panic(err)
	}
	for _, archiveLocation := range archiveLocationSummary.(map[string]interface{})["data"].([]interface{}) {
		thisLocation := archiveLocation.(map[string]interface{})
		if thisLocation["id"] == locationId {
			return thisLocation["name"].(string)
		}
	}
	return ""
}

// converts rubrik timestamp to an int64 epoch time
func ConvertRubrikTimeToUnixTime(RubrikTime string) int64 {
	/*
		Time will be in one of two formats:
		Thu Aug 27 08:36:47 UTC 2020
		2020-08-27T08:48:33.852Z
	*/
	rfc3339type := regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}Z)")
	stringtype := regexp.MustCompile("([A-Z][a-z]{2} [A-Z][a-z]{2} [0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} UTC [0-9]{4})")
	if rfc3339type.MatchString(RubrikTime) {
		parsedTime, e := time.Parse(time.RFC3339, RubrikTime)
		if e != nil {
			log.Panic(e)
		}
		return parsedTime.Unix()
	} else if stringtype.MatchString(RubrikTime) {
		const layout = "Mon Jan 2 15:04:05 MST 2006"
		parsedTime, e := time.Parse(layout, RubrikTime)
		if e != nil {
			log.Panic(e)
		}
		return parsedTime.Unix()
	}
	return 0
}

// pads a string with `pad` to length `plength`
func LeftPad(s string, pad string, plength int) string {
    for i := len(s); i < plength; i++ {
        s = pad + s
    }
    return s
}

// converts a string to a float64
func ConvertToFloat64(s string) float64 {
	c,err := strconv.ParseFloat(s,64)
	if (err != nil) {
		log.Panic(err)
	}
	return c
}

// converts a string to int64
func ConvertToInt64(s string) int64 {
	c,err := strconv.ParseInt(s,10,64)
	if (err != nil) {
		log.Panic(err)
	}
	return c
}