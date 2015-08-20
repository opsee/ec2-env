package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

var metadataURI = "http://169.254.169.254/latest/meta-data"
var identityURI = "http://169.254.169.254/latest/dynamic/instance-identity/document/"

// {
//   "instanceId" : "i-01debbc2",
//   "billingProducts" : null,
//   "imageId" : "ami-076e6542",
//   "architecture" : "x86_64",
//   "pendingTime" : "2015-01-17T02:17:13Z",
//   "instanceType" : "t2.micro",
//   "accountId" : "933693344490",
//   "kernelId" : null,
//   "ramdiskId" : null,
//   "region" : "us-west-1",
//   "version" : "2010-08-31",
//   "availabilityZone" : "us-west-1c",
//   "privateIp" : "172.31.10.200",
//   "devpayProductCodes" : null
// }

type awsCredentials struct {
	Code            string `json:"Code"`
	LastUpdated     string `json:"LastUpdated"`
	Type            string `json:"Type"`
	AccessKeyID     string `json:"AccessKeyId" shell:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"SecretAccessKey" shell:"AWS_SECRET_ACCESS_KEY"`
	Token           string `json:"Token" shell:"AWS_SESSION_TOKEN"`
	Expiration      string `json:"Expiration"`
}

type instanceData struct {
	InstanceID string `json:"instanceId" shell:"AWS_INSTANCE_ID"`
	ImageID    string `json:"imageId" shell:"AWS_IMAGE_ID"`
	AccountID  string `json:"accountId" shell:"AWS_ACCOUNT_ID"`
	Region     string `json:"region" shell:"AWS_DEFAULT_REGION"`
}

func buildURL(s string) string {
	return metadataURI + "/" + s
}

func toShellVar(s string) string {
	path := strings.Split(s, "/")
	s = path[len(path)-1]
	s = strings.Replace(s, "-", "_", -1)
	s = strings.ToUpper(s)
	return s
}

var client http.Client

func makeHTTPRequest(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getInstanceData() (*instanceData, error) {
	var d = new(instanceData)

	jsonBody, err := makeHTTPRequest(identityURI)
	if err != nil {
		return d, err
	}

	err = json.Unmarshal(jsonBody, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}

func getAwsCredentials() (*awsCredentials, error) {
	var creds = new(awsCredentials)

	credentialsURI := buildURL("iam/security-credentials/CoreOS_Cluster_Role")

	jsonBody, err := makeHTTPRequest(credentialsURI)
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(jsonBody, &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}

func shellEncode(i interface{}) ([]byte, error) {
	var b bytes.Buffer

	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		v := val.Field(i)
		tag := f.Tag.Get("shell")
		if tag != "" {
			_, err := b.WriteString(fmt.Sprintf("export %s=%s\n", tag, v.String()))
			if err != nil {
				return nil, err
			}
		}
	}

	return b.Bytes(), nil
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(2*time.Second))
}

func main() {
	transport := http.Transport{
		Dial: dialTimeout,
	}

	client = http.Client{
		Transport: &transport,
	}

	instanceData, err := getInstanceData()
	if err != nil {
		fmt.Println("ERROR getting instance identity: ", err)
		os.Exit(255)
	}

	credentials, err := getAwsCredentials()
	if err != nil {
		fmt.Println("ERROR getting instance identity: ", err)
		os.Exit(255)
	}

	encoded, err := shellEncode(*instanceData)
	if err != nil {
		fmt.Println("ERROR encoding to shell variables: ", err)
	}
	fmt.Print(string(encoded))

	encoded, err = shellEncode(*credentials)
	if err != nil {
		fmt.Println("ERROR encoding to shell variables: ", err)
	}

	fmt.Print(string(encoded))
}
