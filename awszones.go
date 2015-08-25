package main

import (
	"fmt"
	"log"

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

func main() {

	// TODO(mleone896): add distinction between public and private zones

	c := New()

	log.Printf("querying all route53 zones")
	zones := c.HostedZones()

	for id, _ := range zones {

		z, err := c.r53.ListResourceRecordSets(id, nil)

		if err != nil {
			log.Fatal(err)
		}

		for _, record := range z.Records {
			fmt.Printf("record: %s       type: %s\n", record.Name, record.Type)

		}
	}

}
