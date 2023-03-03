package user

import (
	"github.com/YFatMR/go_messenger/user_service/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func userDocumentFromEntities(user *entity.User, credential *entity.Credential) *userDocument {
	return &userDocument{
		Login:          credential.Login,
		HashedPassword: credential.HashedPassword,
		UserRole:       credential.Role.Name,
		Nickname:       user.Nickname,
		Name:           user.Name,
		Surname:        user.Surname,
	}
}

func insertOneResultToUserID(result *mongo.InsertOneResult) *entity.UserID {
	return &entity.UserID{ID: result.InsertedID.(primitive.ObjectID).Hex()}
}

func userDocumentToUser(document *userDocument) *entity.User {
	return &entity.User{
		Nickname: document.Nickname,
		Name:     document.Name,
		Surname:  document.Surname,
	}
}

func userDocumentToAccount(document *userDocument) (
	*entity.Account, error,
) {
	role, err := entity.UserRoleFromString(document.UserRole)
	if err != nil {
		return nil, err
	}
	return &entity.Account{
		UserID:         document.ID.Hex(),
		Login:          document.Login,
		HashedPassword: document.HashedPassword,
		Role:           role,
		Nickname:       document.Nickname,
		Name:           document.Name,
		Surname:        document.Surname,
	}, nil
}

// func programIDFromInsertOneResult(result *mongo.InsertOneResult) *entity.ProgramID {
// 	return &entity.ProgramID{ID: result.InsertedID.(primitive.ObjectID).Hex()}
// }

// func programSourceFromDocument(document *ProgramDocument) *entity.ProgramSource {
// 	return &entity.ProgramSource{
// 		SourceCode: document.SourceCode,
// 		Language:   document.Language,
// 	}
// }

// func programSourceToDocument(programSource *entity.ProgramSource) *ProgramDocument {
// 	return &ProgramDocument{
// 		Language:   programSource.Language,
// 		SourceCode: programSource.SourceCode,
// 	}
// }

// func codeRunnerOutputFromDocument(document *ProgramDocument) *entity.ProgramOutput {
// 	return &entity.ProgramOutput{
// 		Stdout: document.CodeRunnerStdout,
// 		Stderr: document.CodeRunnerStderr,
// 	}
// }

// func codeRunnerOutputToDocument(programOutput *entity.ProgramOutput) *ProgramDocument {
// 	return &ProgramDocument{
// 		CodeRunnerStdout: programOutput.Stdout,
// 		CodeRunnerStderr: programOutput.Stderr,
// 	}
// }

// func linterOutputFromDocument(document *ProgramDocument) *entity.ProgramOutput {
// 	return &entity.ProgramOutput{
// 		Stdout: document.LinterStdout,
// 		Stderr: document.LinterStderr,
// 	}
// }

// func linterOutputToDocument(programOutput *entity.ProgramOutput) *ProgramDocument {
// 	return &ProgramDocument{
// 		LinterStdout: programOutput.Stdout,
// 		LinterStderr: programOutput.Stderr,
// 	}
// }

// func programFromDocument(document *ProgramDocument) *entity.Program {
// 	return &entity.Program{
// 		ID:               entity.ProgramID{ID: document.ID.Hex()},
// 		Source:           *programSourceFromDocument(document),
// 		CodeRunnerOutput: *codeRunnerOutputFromDocument(document),
// 		LinterOutput:     *linterOutputFromDocument(document),
// 	}
// }
