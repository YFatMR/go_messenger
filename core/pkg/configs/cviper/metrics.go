package cviper

import "time"

type MetricServiceSettings struct {
	address           string
	listingSuffix     string
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
}

func NewMetricMetricServiceSettingsFromConfig(config *CustomViper) *MetricServiceSettings {
	enableMetricService := config.GetBoolRequired("ENABLE_METRIC_SERVICE")
	if !enableMetricService {
		return nil
	}
	return &MetricServiceSettings{
		address:           config.GetStringRequired("METRICS_SERVICE_ADDRESS"),
		listingSuffix:     config.GetStringRequired("METRICS_SERVICE_LISTING_SUFFIX"),
		readTimeout:       config.GetSecondsDurationRequired("METRICS_SERVICE_READ_TIMEOUT_SECONDS"),
		writeTimeout:      config.GetSecondsDurationRequired("METRICS_SERVICE_WRITE_TIMEOUT_SECONDS"),
		idleTimeout:       config.GetSecondsDurationRequired("METRICS_SERVICE_IDLE_TIMEOUT_SECONDS"),
		readHeaderTimeout: config.GetSecondsDurationRequired("METRICS_SERVICE_READ_HEADER_TIMEOUT_SECONDS"),
	}
}

func (s *MetricServiceSettings) GetAddress() string {
	return s.address
}

func (s *MetricServiceSettings) GetListingSuffix() string {
	return s.listingSuffix
}

func (s *MetricServiceSettings) GetReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *MetricServiceSettings) GetWriteTimeout() time.Duration {
	return s.writeTimeout
}

func (s *MetricServiceSettings) GetIdleTimeout() time.Duration {
	return s.idleTimeout
}

func (s *MetricServiceSettings) GetReadHeaderTimeout() time.Duration {
	return s.readHeaderTimeout
}
