package sandbox

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/mongodb"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ProgramDocument struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Language         string             `bson:"language,omitempty"`
	SourceCode       string             `bson:"sourceCode,omitempty"`
	CodeRunnerStdout string             `bson:"codeRunnerStdout,omitempty"`
	CodeRunnerStderr string             `bson:"codeRunnerStderr,omitempty"`
	LinterStdout     string             `bson:"linterStdout,omitempty"`
	LinterStderr     string             `bson:"linterStderr,omitempty"`
}

type repository struct {
	operationTimeout time.Duration
	collection       *mongo.Collection
	logger           *czap.Logger
}

func NewRepository(collection *mongo.Collection, operationTimeout time.Duration,
	logger *czap.Logger,
) apientity.SandboxRepository {
	return &repository{
		collection:       collection,
		operationTimeout: operationTimeout,
		logger:           logger,
	}
}

func RepositoryFromConfig(ctx context.Context, config *cviper.CustomViper, logger *czap.Logger) (
	apientity.SandboxRepository, error,
) {
	databaseSettings := cviper.NewDatabaseSettingsFromConfig(config)
	mongoCollection, err := mongodb.Connect(ctx, databaseSettings, logger)
	if err != nil {
		return nil, err
	}
	return NewRepository(mongoCollection, databaseSettings.GetOperationTimeout(), logger), nil
}

func (r *repository) GetProgramByID(ctx context.Context, programID *entity.ProgramID) (
	*entity.Program, error,
) {
	objectID, err := primitive.ObjectIDFromHex(programID.ID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return nil, ErrGetProgramByID
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	var document ProgramDocument
	err = r.collection.FindOne(mongoOperationCtx, bson.D{
		{Key: "_id", Value: objectID},
	}).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		r.logger.ErrorContext(ctx, "Program not found by id", zap.Error(err))
		return nil, ErrProgramNotFount
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, ErrProgramNotFount
	}

	program := programFromDocument(&document)
	return program, nil
}

func (r *repository) CreateProgram(ctx context.Context, programSource *entity.ProgramSource) (
	*entity.ProgramID, error,
) {
	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	document := programSourceToDocument(programSource)
	insertResult, err := r.collection.InsertOne(mongoOperationCtx, document)
	if err != nil {
		r.logger.ErrorContext(ctx, "Can't insert new program", zap.Error(err))
		return nil, ErrProgramCreation
	}

	programID := programIDFromInsertOneResult(insertResult)
	return programID, nil
}

func (r *repository) UpdateProgramSource(ctx context.Context, programID *entity.ProgramID,
	programSource *entity.ProgramSource,
) error {
	programMongoID, err := primitive.ObjectIDFromHex(programID.ID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return ErrUpdateProgramSource
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err = r.collection.UpdateOne(
		mongoOperationCtx,
		bson.M{"_id": programMongoID},
		bson.M{
			"$set": programSourceToDocument(programSource),
		},
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Update error", zap.Error(err), zap.String("programID", programID.ID))
		return ErrUpdateProgramSource
	}
	return nil
}

func (r *repository) UpdateCodeRunnerOutput(ctx context.Context, programID *entity.ProgramID,
	output *entity.ProgramOutput,
) error {
	programMongoID, err := primitive.ObjectIDFromHex(programID.ID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return ErrUpdateProgramRunnerOutput
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err = r.collection.UpdateOne(
		mongoOperationCtx,
		bson.M{"_id": programMongoID},
		bson.M{
			"$set": codeRunnerOutputToDocument(output),
		},
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Update error", zap.Error(err), zap.String("programID", programID.ID))
		return ErrUpdateProgramRunnerOutput
	}
	return nil
}

func (r *repository) UpdateLinterOutput(ctx context.Context, programID *entity.ProgramID,
	output *entity.ProgramOutput,
) error {
	programMongoID, err := primitive.ObjectIDFromHex(programID.ID)
	if err != nil {
		r.logger.ErrorContext(ctx, "Got wrong id format", zap.Error(err))
		return ErrUpdateLinterOutput
	}

	mongoOperationCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err = r.collection.UpdateOne(
		mongoOperationCtx,
		bson.M{"_id": programMongoID},
		bson.M{
			"$set": linterOutputToDocument(output),
		},
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Update error", zap.Error(err), zap.String("programID", programID.ID))
		return ErrUpdateLinterOutput
	}
	return nil
}
