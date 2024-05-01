package suite

import (
	ssov1 "SSO/contract/gen/go/sso"
	"SSO/internal/config"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}
