package awsilvanus

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type FunctionWrapper struct {
	LambdaClient   *lambda.Client
	cloudWatchLogs *cloudwatchlogs.Client
}

func (wrapper FunctionWrapper) GetLogs(functionName string) string {
	// Create a CloudWatch Logs client
	svc := wrapper.InitCloudWatchClient()
	logGroupName := "/aws/lambda/" + functionName
	// Call DescribeLogStreams to get the log stream details

	describeStreamsOutput, err := svc.DescribeLogStreams(context.TODO(), &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		// log.Fatalf("failed to describe log streams, %v", err)
	}

	logStreamName := aws.ToString(describeStreamsOutput.LogStreams[0].LogStreamName)

	// Call GetLogEvents to retrieve log events
	getEventsOutput, err := svc.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})
	if err != nil {
		// log.Fatalf("failed to get log events, %v", err)
	}

	var sb strings.Builder

	// Print log events
	for _, event := range getEventsOutput.Events {
		sb.WriteString(aws.ToString(event.Message) + "\n")
		// fmt.Printf("Timestamp: %d, Message: %s\n", aws.ToInt64(event.Timestamp), aws.ToString(event.Message))
	}

	return sb.String()
}

// ListFunctions lists up to maxItems functions for the account. This function uses a
// lambda.ListFunctionsPaginator to paginate the results.
func (wrapper FunctionWrapper) ListFunctions(maxItems int) []types.FunctionConfiguration {
	var functions []types.FunctionConfiguration
	paginator := lambda.NewListFunctionsPaginator(wrapper.LambdaClient, &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(int32(maxItems)),
	})
	for paginator.HasMorePages() && len(functions) < maxItems {
		pageOutput, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Panicf("Couldn't list functions for your account. Here's why: %v\n", err)
		}
		functions = append(functions, pageOutput.Functions...)
	}
	return functions
}

func (wrapper FunctionWrapper) InitCloudWatchClient() *cloudwatchlogs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("endi"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a CloudWatch Logs client
	return cloudwatchlogs.NewFromConfig(cfg)
}

func (wrapper FunctionWrapper) InitLambdaClient() *lambda.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("endi"))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return lambda.NewFromConfig(cfg)
}
