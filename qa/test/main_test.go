package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"testing"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

var (
	dockerComposeFile string
	envFile           string
)

func init() {
	flag.StringVar(&dockerComposeFile, "docker-compose-file", "", "docker compose file path")
	flag.StringVar(&envFile, "env-file", "", "docker compose --env-file flag and viper config")
}

func TestMain(m *testing.M) {
	flag.Parse()

	ctx, stopDocker := signal.NotifyContext(context.Background(), os.Interrupt)

	// Run docker-compose
	command := exec.CommandContext(
		ctx, "docker-compose", "--file", dockerComposeFile, "--env-file", envFile, "up",
	)
	err := command.Start()
	if err != nil {
		panic(err)
	}

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

	stopDocker()
	os.Exit(exitCode)
}
