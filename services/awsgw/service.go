package awsgw

import "github.com/Gaoey/poc-aws-websocket-gateway.git/internals/aws"

type AWSGatewayService struct {
	// Add any necessary fields here, such as configuration or dependencies
	App *aws.AWSApplication
}

func NewService(app *aws.AWSApplication) *AWSGatewayService {
	// Initialize the AWS application or any other dependencies here
	return &AWSGatewayService{
		// Initialize any necessary fields here
		App: app,
	}
}
