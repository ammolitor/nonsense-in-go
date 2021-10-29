package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"encoding/json"
	"log"
	"os"
)

func main() {
	const (
		functionName = "ebs_dashboard"
		region       = "us-west-2"
		namespace    = "aaronm"
	)

	type lambdaPayload struct {
		Namespace string `json:"namespace"`
	}

	request := lambdaPayload{Namespace: namespace}

	payload, err := json.Marshal(request)
	if err != nil {
		log.Println("Error marshaling request")
		os.Exit(1)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
			LogLevel:                      aws.LogLevel(aws.LogDebug),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	_, sessionError := sess.Config.Credentials.Get()
	if sessionError != nil {
		log.Print(sessionError)
	}
	client := lambda.New(sess)

	_, err = client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName), Payload: payload})
	if err != nil {
		log.Printf("Error calling %s\n error=%s", functionName, err.Error())
	}
	var dashboardURL = "https://" + region + ".console.aws.amazon.com/cloudwatch/home?region=" + region + "#dashboards:name=" + namespace + ";start=PT1H"
	fmt.Printf(dashboardURL)
}
