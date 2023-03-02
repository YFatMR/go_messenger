package sandbox

import (
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func programIDFromInsertOneResult(result *mongo.InsertOneResult) *entity.ProgramID {
	return &entity.ProgramID{ID: result.InsertedID.(primitive.ObjectID).Hex()}
}

func programSourceFromDocument(document *ProgramDocument) *entity.ProgramSource {
	return &entity.ProgramSource{
		SourceCode: document.SourceCode,
		Language:   document.Language,
	}
}

func programSourceToDocument(programSource *entity.ProgramSource) *ProgramDocument {
	return &ProgramDocument{
		Language:   programSource.Language,
		SourceCode: programSource.SourceCode,
	}
}

func codeRunnerOutputFromDocument(document *ProgramDocument) *entity.ProgramOutput {
	return &entity.ProgramOutput{
		Stdout: document.CodeRunnerStdout,
		Stderr: document.CodeRunnerStderr,
	}
}

func codeRunnerOutputToDocument(programOutput *entity.ProgramOutput) *ProgramDocument {
	return &ProgramDocument{
		CodeRunnerStdout: programOutput.Stdout,
		CodeRunnerStderr: programOutput.Stderr,
	}
}

func linterOutputFromDocument(document *ProgramDocument) *entity.ProgramOutput {
	return &entity.ProgramOutput{
		Stdout: document.LinterStdout,
		Stderr: document.LinterStderr,
	}
}

func linterOutputToDocument(programOutput *entity.ProgramOutput) *ProgramDocument {
	return &ProgramDocument{
		LinterStdout: programOutput.Stdout,
		LinterStderr: programOutput.Stderr,
	}
}

func programFromDocument(document *ProgramDocument) *entity.Program {
	return &entity.Program{
		ID:               entity.ProgramID{ID: document.ID.Hex()},
		Source:           *programSourceFromDocument(document),
		CodeRunnerOutput: *codeRunnerOutputFromDocument(document),
		LinterOutput:     *linterOutputFromDocument(document),
	}
}
