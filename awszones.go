package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/route53"
)

type Conn struct {
	r53 *route53.Route53
}

func (c *Conn) HostedZones() (ZoneMap map[string]string) {

	ZoneMap = make(map[string]string)

	zones, err := c.r53.ListHostedZones("", 50)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range zones.HostedZones {
		ZoneMap[route53.CleanZoneID(val.ID)] = val.Name
	}

	return ZoneMap
}

func (c *Conn) RecordTypeMap(zones map[string]string) (recordMap map[string]string) {
	recordMap = make(map[string]string)
	for id, _ := range zones {

		z, err := c.r53.ListResourceRecordSets(id, nil)

		if err != nil {
			log.Fatal(err)
		}

		for _, record := range z.Records {
			recordMap[record.Name] = record.Type
		}
	}
	return recordMap
}

func New() *Conn {

	c := new(Conn)

	// this is looking for keys in env
	auth, err := aws.EnvAuth() // TODO(mleone896): maybe make a switch to use from config ?
	if err != nil {
		log.Fatal(err)
	}
	c.r53 = route53.New(auth, aws.USWest)
	return c

}

func (c *Conn) Records(w http.ResponseWriter, r *http.Request) {
	zones := c.HostedZones()

	records := c.RecordTypeMap(zones)

	out, _ := json.Marshal(records)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))

}

func main() {
	c := New()

	// TODO(mleone896): add distinction between public and private zones
	mux := http.NewServeMux()

	mux.HandleFunc("/records", c.Records)

	http.ListenAndServe(":9999", mux)

}
