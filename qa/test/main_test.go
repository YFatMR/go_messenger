package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"testing"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

var envFile string

func init() {
	flag.StringVar(&envFile, "env-file", "", "configuration file")
}

func TestMain(m *testing.M) {
	flag.Parse()

	ctx, stopTests := signal.NotifyContext(context.Background(), os.Interrupt)

	// Setup tests
	config := cviper.New()
	config.SetConfigFile(envFile)
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(err)
	}

	grpcAuthorizationHeader = config.GetStringRequired("GRPC_AUTHORIZARION_HEADER")
	qaHost := config.GetStringRequired("QA_HOST")

	// restFrontServiceAddress := qaHost + ":" + config.GetStringRequired("PUBLIC_REST_FRONT_SERVICE_PORT")
	grpcFrontServiceAddress := qaHost + ":" + config.GetStringRequired("PUBLIC_GRPC_FRONT_SERVICE_PORT")
	userServiceAddress := qaHost + ":" + config.GetStringRequired("PUBLIC_USER_SERVICE_PORT")
	authServiceAddress := qaHost + ":" + config.GetStringRequired("PUBLIC_AUTH_SERVICE_PORT")

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

	stopTests()
	os.Exit(exitCode)
}
