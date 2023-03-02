package prometheus

import (
	"net/http"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServeMetrcirService(config *cviper.MetricServiceSettings, logger *czap.Logger) {
	server := &http.Server{
		Addr:              config.GetAddress(),
		ReadTimeout:       config.GetReadTimeout(),
		WriteTimeout:      config.GetWriteTimeout(),
		IdleTimeout:       config.GetIdleTimeout(),
		ReadHeaderTimeout: config.GetReadHeaderTimeout(),
		Handler:           nil,
	}
	http.Handle(config.GetListingSuffix(), promhttp.Handler())
	//#nosec G114: Use of net/http serve function that has no support for setting timeouts
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Can't up metrics server with endpoint" + config.GetAddress() +
			". Operation finished with error: " + err.Error())
		panic(err)
	}
}
