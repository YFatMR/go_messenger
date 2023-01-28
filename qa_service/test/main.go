package test

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

func TestMain(m *testing.M) {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	grpcAuthorizationHeader = config.GetStringRequired("GRPC_AUTHORIZARION_HEADER")

	// restFrontServiceAddress := config.GetStringRequired("REST_FRONT_SERVICE_ADDRESS")
	grpcFrontServiceAddress := config.GetStringRequired("GRPC_FRONT_SERVICE_ADDRESS")
	userServiceAddress := config.GetStringRequired("USER_SERVICE_ADDRESS")
	authServiceAddress := config.GetStringRequired("AUTH_SERVICE_ADDRESS")

	testResponseTimeout := config.GetMillisecondsDurationRequired("TEST_RESPONSE_TIMEOUT_MILLISECONDS")
	testSetupTimeout := config.GetMillisecondsDurationRequired("TEST_SETUP_TIMEOUT_MILLISECONDS")

	testServiceSetupRetryCount := config.GetIntRequired("TEST_SERVICE_SETUP_RETRY_COUNT")
	testServiceSetupRetryInterval := config.GetMillisecondsDurationRequired(
		"TEST_SERVICE_SETUP_RETRY_INTERVAL_MILLISECONDS",
	)

	// Setup auth service
	ctxSetup, cancelSetup := context.WithTimeout(ctx, testSetupTimeout)
	authClient, err := newProtobufAuthClient(ctxSetup, authServiceAddress, testResponseTimeout)
	if err != nil {
		panic(err)
	}
	pingAuthServiceCallback := func(ctx context.Context) (*proto.Pong, error) {
		return authClient.Ping(ctx, &proto.Void{})
	}
	err = pingService(ctx, pingAuthServiceCallback, testServiceSetupRetryCount, testServiceSetupRetryInterval)
	if err != nil {
		panic(err)
	}

	// Setup user service
	userClient, err := newProtobufUserClient(ctxSetup, userServiceAddress, testResponseTimeout)
	if err != nil {
		panic(err)
	}
	pingUserServiceCallback := func(ctx context.Context) (*proto.Pong, error) {
		return userClient.Ping(ctx, &proto.Void{})
	}
	err = pingService(ctx, pingUserServiceCallback, testServiceSetupRetryCount, testServiceSetupRetryInterval)
	if err != nil {
		panic(err)
	}

	// Setup front service
	frontServicegRPCClient, err = newProtobufFrontClient(ctxSetup, grpcFrontServiceAddress, testResponseTimeout)
	if err != nil {
		panic(err)
	}
	pingFrontServiceCallback := func(ctx context.Context) (*proto.Pong, error) {
		return frontServicegRPCClient.Ping(ctx, &proto.Void{})
	}
	err = pingService(ctx, pingFrontServiceCallback, testServiceSetupRetryCount, testServiceSetupRetryInterval)
	if err != nil {
		panic(err)
	}

	cancelSetup()

	exitCode := m.Run()

	os.Exit(exitCode)
}
