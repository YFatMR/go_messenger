package cviper

import (
	"time"
)

type DatabaseSettings struct {
	uri               string
	databaseName      string
	collectionName    string
	operationTimeout  time.Duration
	connectionTimeout time.Duration

	startupReconnectionCount    int
	startupReconnectionInterval time.Duration
}

func NewDatabaseSettingsFromConfig(config *CustomViper) *DatabaseSettings {
	return &DatabaseSettings{
		uri:                      config.GetStringRequired("DATABASE_URI"),
		databaseName:             config.GetStringRequired("DATABASE_NAME"),
		collectionName:           config.GetStringRequired("DATABASE_COLLECTION_NAME"),
		operationTimeout:         config.GetMillisecondsDurationRequired("DATABASE_OPERATION_TIMEOUT_MILLISECONDS"),
		connectionTimeout:        config.GetMillisecondsDurationRequired("DATABASE_CONNECTION_TIMEOUT_MILLISECONDS"),
		startupReconnectionCount: config.GetIntRequired("DATABASE_STARTUP_RECONNECTION_COUNT"),
		startupReconnectionInterval: config.GetMillisecondsDurationRequired(
			"DATABASE_STURTUP_RECONNECTIONION_INTERVAL_MILLISECONDS",
		),
	}
}

func (s *DatabaseSettings) GetURI() string {
	return s.uri
}

func (s *DatabaseSettings) GetDatabaseName() string {
	return s.databaseName
}

func (s *DatabaseSettings) GetCollectionName() string {
	return s.collectionName
}

func (s *DatabaseSettings) GetOperationTimeout() time.Duration {
	return s.operationTimeout
}

func (s *DatabaseSettings) GetConnectionTimeout() time.Duration {
	return s.connectionTimeout
}

func (s *DatabaseSettings) GetStartupReconnectionCount() int {
	return s.startupReconnectionCount
}

func (s *DatabaseSettings) GetSturtupReconnectionInterval() time.Duration {
	return s.startupReconnectionInterval
}
