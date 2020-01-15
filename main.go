package main

import (
	"fmt"
	"log"
	//"net/http"
	//"time"
	//"strconv"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func main() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}
	clusterDetails,err := rubrik.Get("v1","/cluster/me")
	if err != nil {
		log.Fatal(err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"]
	fmt.Println("Cluster name: "+clusterName.(string))
	fmt.Println(GetRunwayRemaining(rubrik))
	fmt.Println(GetStorageSummary(rubrik))
}