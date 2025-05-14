package awsgw

import (
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/auth"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/aws"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/redis"
)

type AWSGatewayService struct {
	App     *aws.AWSApplication
	AuthAPI *auth.AuthAPI
	Redis   *redis.RedisHandler
}

func NewService(app *aws.AWSApplication, auth *auth.AuthAPI, redis *redis.RedisHandler) *AWSGatewayService {
	return &AWSGatewayService{
		App:     app,
		AuthAPI: auth,
		Redis:   redis,
	}
}
