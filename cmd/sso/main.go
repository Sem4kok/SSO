package main

import (
	"SSO/internal/config"
	"fmt"
)

// project layer's schema
// --------------------------------
// 1) Transport layer.
// Transport layer gets request and
// communicate with Service layer
//
// 2) Service layer. (auth, permission, userinfo)
// Service layer implement business-logic
// communicate with Data layer
//
// 3) Data layer.
// Data layer communicate with data (include storage)
// return response's to Service layer

func main() {
	cfg := config.Load()

	fmt.Println(cfg)

	// TODO: initialize logger

	// TODO: Initialize application(app)

	// TODO: start gRPC-server app
}
