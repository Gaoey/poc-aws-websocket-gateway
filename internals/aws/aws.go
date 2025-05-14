package aws

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithylog "github.com/aws/smithy-go/logging"
	"gopkg.in/square/go-jose.v2/json"
)

type AWSApplication struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
	Cfg       aws.Config
}

func checkCredentials(creds aws.CredentialsProvider, region string) bool {

	// Create an STS client
	stsClient := sts.New(sts.Options{
		Region:      region,
		Credentials: creds,
	})

	// Try to get the caller identity (this will throw an error if the credentials are invalid)
	resp, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("Error getting caller identity: %v", err)
		return false
	}

	// If successful, print the caller identity
	fmt.Printf("Successfully authenticated as %s with ARN %s\n", *resp.UserId, *resp.Arn)
	return true
}

func NewAWSApplication() *AWSApplication {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secreytKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_GATEWAY_ENDPOINT")

	creds := credentials.NewStaticCredentialsProvider(accessKey, secreytKey, "")
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody|aws.LogSigning), // logging
		config.WithLogger(smithylog.NewStandardLogger(os.Stderr)),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	checkCredentials(creds, region)

	return &AWSApplication{
		AccessKey: accessKey,
		SecretKey: secreytKey,
		Region:    region,
		Endpoint:  endpoint,
		Cfg:       awsCfg,
	}
}

func (a *AWSApplication) PostToConnection(ctx context.Context, connectionID string, data interface{}) error {
	var err error
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Prepare HTTP request
	u := &url.URL{
		Scheme: "https",
		Host:   "2gd7mfyzi2.execute-api.ap-southeast-1.amazonaws.com",
		Path:   "/dev/@connections/" + url.PathEscape(connectionID),
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	// Set required headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Amz-Date", time.Now().UTC().Format("20060102T150405Z"))

	// Compute payload hash for signing
	payloadHash := sha256.Sum256(payload)
	payloadHashHex := fmt.Sprintf("%x", payloadHash)

	// Prepare signer
	signer := v4.NewSigner()

	creds := aws.Credentials{
		AccessKeyID:     a.AccessKey,
		SecretAccessKey: a.SecretKey,
	}

	// Sign the request
	err = signer.SignHTTP(
		ctx,
		creds,
		req,
		payloadHashHex,
		"execute-api", // Service name
		a.Region,
		time.Now(),
	)
	if err != nil {
		log.Fatalf("Failed to sign request: %v", err)
	}

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Output response
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nResponse: %s\n", resp.StatusCode, string(body))

	return nil
}
