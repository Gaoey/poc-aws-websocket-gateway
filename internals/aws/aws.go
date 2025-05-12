package aws

import (
	"context"
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"gopkg.in/square/go-jose.v2/json"
)

type AWSApplication struct {
	Client *apigatewaymanagementapi.Client
}

func NewAWSApplication() *AWSApplication {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secreytKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	credentialsProvider := credentials.NewStaticCredentialsProvider(accessKey, secreytKey, "")
	credentials := aws.NewCredentialsCache(credentialsProvider)
	endpoint := os.Getenv("AWS_ENDPOINT") // Add endpoint from environment variable
	client := apigatewaymanagementapi.New(apigatewaymanagementapi.Options{
		Region:           region, // e.g., "us-east-1"
		Credentials:      credentials,
		EndpointResolver: apigatewaymanagementapi.EndpointResolverFromURL(endpoint),
	})

	return &AWSApplication{
		Client: client,
	}
}

func (a *AWSApplication) PostToConnection(ctx context.Context, connectionID string, data interface{}) error {
	// Implement the logic to post to connection
	// Prepare payload
	var err error
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Call PostToConnection
	_, err = a.Client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         payload,
	})
	if err != nil {
		if ok := strings.Contains(err.Error(), "GoneException"); ok || (err != nil && strings.Contains(err.Error(), "410")) {
			log.Printf("Connection is gone: %s", connectionID)
			// clear connection
			return fmt.Errorf("connection is gone: %w", err)
		} else {
			log.Printf("Failed to post to connection: %s, error: %v", connectionID, err)
			return fmt.Errorf("failed to post to connection: %w", err)
		}
	}

	fmt.Println("Message sent successfully!")
	return nil
}
