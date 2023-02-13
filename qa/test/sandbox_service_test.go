package test

import (
	"context"
	"testing"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"github.com/stretchr/testify/suite"
)

type SandboxTestSuite struct {
	userManager UserManager
	suite.Suite
}

func TestSandboxTestSuite(t *testing.T) {
	suite.Run(t, new(SandboxTestSuite))
}

func (s *SandboxTestSuite) TestHelloWorldExecution() {
	ctx := context.Background()
	require := s.Require()

	_, token, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)

	ctx = s.userManager.NewContextWithToken(ctx, token)
	program, err := frontServicegRPCClient.Execute(ctx, &proto.Program{
		Language:   "go",
		SourceCode: GetHelloWorldProgramSource(),
	})
	require.NoError(err)
	require.NotNil(program)

	require.Equal(program.GetStderr(), "")
	require.Equal(program.GetStdout(), "Hello world!\n")
}
