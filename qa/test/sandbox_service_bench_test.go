//go:build bench
// +build bench

package test

import (
	"context"
	"strings"
	"testing"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

func BenchmarkCreateProgram(b *testing.B) {
	ctx := context.Background()
	userManager := UserManager{}
	_, token, err := userManager.NewAuthorizedUser(ctx)
	if err != nil {
		b.Fail()
	}
	program := NewHugeHelloWorldWithComment(strings.Repeat("a", 3000))
	ctx = userManager.NewContextWithToken(ctx, token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := frontServicegRPCClient.CreateProgram(ctx, &proto.ProgramSource{
			SourceCode: program.SourceCode,
			Language:   "go",
		})
		if result == nil || err != nil {
			b.Fail()
		}
	}
}

func BenchmarkGetProgramByID(b *testing.B) {
	ctx := context.Background()
	userManager := UserManager{}
	_, token, err := userManager.NewAuthorizedUser(ctx)
	if err != nil {
		b.Fail()
	}
	program := NewHugeHelloWorldWithComment(strings.Repeat("a", 3000))
	ctx = userManager.NewContextWithToken(ctx, token)
	programID, err := frontServicegRPCClient.CreateProgram(ctx, &proto.ProgramSource{
		SourceCode: program.SourceCode,
		Language:   "go",
	})
	if programID == nil || err != nil {
		b.Fail()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := frontServicegRPCClient.GetProgramByID(ctx, programID)
		if result == nil || err != nil {
			b.Fail()
		}
	}
}

func BenchmarkUpdateProgramSource(b *testing.B) {
	ctx := context.Background()
	userManager := UserManager{}
	_, token, err := userManager.NewAuthorizedUser(ctx)
	if err != nil {
		b.Fail()
	}

	startProgram := NewHelloWorldProgram()

	ctx = userManager.NewContextWithToken(ctx, token)
	programID, err := frontServicegRPCClient.CreateProgram(ctx, &proto.ProgramSource{
		SourceCode: startProgram.SourceCode,
		Language:   "go",
	})
	if programID == nil || err != nil {
		b.Fail()
	}

	updateProgram1 := NewHugeHelloWorldWithComment(strings.Repeat("a", 3000))
	updateProgram2 := NewHugeHelloWorldWithComment(strings.Repeat("b", 2987))

	request1 := &proto.UpdateProgramSourceRequest{
		ProgramID: programID,
		ProgramSource: &proto.ProgramSource{
			Language:   "go",
			SourceCode: updateProgram1.SourceCode,
		},
	}
	request2 := &proto.UpdateProgramSourceRequest{
		ProgramID: programID,
		ProgramSource: &proto.ProgramSource{
			Language:   "go",
			SourceCode: updateProgram2.SourceCode,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use defferent requests to prevent database optimizations
		request := request1
		if i%2 == 0 {
			request = request2
		}
		result, err := frontServicegRPCClient.UpdateProgramSource(ctx, request)
		if result == nil || err != nil {
			b.Fail()
		}
	}
}

// TODO: implement dry run to prevent "no free worker" error
// func BenchmarkRunProgram(b *testing.B) {
// 	ctx := context.Background()
// 	userManager := UserManager{}
// 	_, token, err := userManager.NewAuthorizedUser(ctx)
// 	if err != nil {
// 		b.Error(err)
// 	}
// 	program := NewHelloWorldProgram()
// 	ctx = userManager.NewContextWithToken(ctx, token)
// 	programID, err := frontServicegRPCClient.CreateProgram(ctx, &proto.ProgramSource{
// 		SourceCode: program.SourceCode,
// 		Language:   "go",
// 	})
// 	if programID == nil || err != nil {
// 		b.Error(err)
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		result, err := frontServicegRPCClient.RunProgram(ctx, programID)
// 		if result == nil || err != nil {
// 			b.Error(err)
// 		}
// 	}
// }
