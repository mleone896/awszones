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

func (c *Conn) HostedZones() []string {

	var Zones []string

	zones, err := c.r53.ListHostedZones("", 50)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range zones.HostedZones {
		Zones = append(Zones, val.Name)
	}

	return Zones
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

	c := New()

	log.Printf("querying all route53 zones")
	zones := c.HostedZones()

	for _, val := range zones {
		fmt.Println(val)
	}

}
