package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

type awsCredentials struct {
	AccessKeyID     string `shell:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `shell:"AWS_SECRET_ACCESS_KEY"`
	SessionToken    string `shell:"AWS_SESSION_TOKEN"`
}

type instanceData struct {
	InstanceID string `shell:"AWS_INSTANCE_ID"`
	Region     string `shell:"AWS_DEFAULT_REGION"`
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
			_, err := b.WriteString(fmt.Sprintf("%s=%s\n", tag, v.String()))
			if err != nil {
				return nil, err
			}
		}
	}

	return b.Bytes(), nil
}

func main() {
	metadataClient := ec2metadata.New(&ec2metadata.Config{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	})
	region, err := metadataClient.Region()
	if err != nil {
		log.Fatalf(err.Error())
	}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
		})

	v, err := creds.Get()
	if err != nil {
		log.Fatalf(err.Error())
	}
	awsCreds := awsCredentials{
		v.AccessKeyID,
		v.SecretAccessKey,
		v.SessionToken,
	}

	instanceID, err := metadataClient.GetMetadata("instance-id")
	if err != nil {
		log.Fatalf(err.Error())
	}

	instanceData := instanceData{
		instanceID,
		region,
	}

	encoded, err := shellEncode(instanceData)
	if err != nil {
		fmt.Println("ERROR encoding to shell variables: ", err)
	}
	fmt.Print(string(encoded))

	encoded, err = shellEncode(awsCreds)
	if err != nil {
		fmt.Println("ERROR encoding to shell variables: ", err)
	}

	fmt.Print(string(encoded))
}
