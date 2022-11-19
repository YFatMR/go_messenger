package main

import (
	. "core/pkg/utils"
	. "front_server/internal"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func main() {
	// Init vars
	frontUserRESTServerAddress := GetFullServiceAddress("FRONT")
	userServerAddress := GetFullServiceAddress("USER")
	mux := runtime.NewServeMux()

	// Register all services
	RegisterRestUserServer(mux, userServerAddress)

	// Start REST service
	if err := http.ListenAndServe(frontUserRESTServerAddress, mux); err != nil {
		panic(err)
	}
}
