package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var metadataURI = "http://169.254.169.254/latest/meta-data"

func getAvailabilityZone() string {
	return getMetadata("placement/availability-zone")
}

// map of functions that return (string,error)
func optionFunctionMap() map[string]func() string {
	m := make(map[string]func() string)

	m["availability_zone"] = getAvailabilityZone
	return m
}

func buildURL(s string) string {
	return metadataURI + s
}

func toShellVar(s string) string {
	path := strings.Split(s, "/")
	s = path[len(path)-1]
	s = strings.Replace(s, "-", "_", -1)
	s = strings.ToUpper(s)
	return s
}

func getMetadata(key string) string {
	resp, err := http.Get(buildURL(key))
	if err != nil {
		fmt.Println("# ERROR getting ", key, ": ", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("# ERROR getting ", key, ": ", err)
		return ""
	}

	return string(body)
}

var (
	amiID      bool
	hostname   bool
	instanceID bool
	localIPV4  bool
)

func init() {
	flag.BoolVar(&amiID, "ami_id", false, "Get instance AMI ID")
	flag.BoolVar(&hostname, "hostname", false, "Get instance hostname")
	flag.BoolVar(&instanceID, "instance_id", false, "Get instance ID")
	flag.BoolVar(&localIPV4, "local_ipv4", false, "Get local IPV4 address")
}

func main() {

	fMap := optionFunctionMap()

	defaultMetadata := []string{
		"availability_zone",
	}

	var additionalMetadata = []string{}

	flag.Parse()

	metadata := append(defaultMetadata, additionalMetadata...)

	for _, element := range metadata {
		value := (fMap[element])()
		fmt.Println("export ", toShellVar(element), "=", value)
	}

}
