package server

import (
	"review_service/internal/conf"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v3"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewGRPCServer, NewHTTPServer)

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	cfg := api.DefaultConfig()

	if conf.GetConsul().GetAddr() != "" {
		cfg.Address = conf.GetConsul().GetAddr()
	}

	if conf.GetConsul().GetScheme() != "" {
		cfg.Scheme = conf.GetConsul().GetScheme()
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return consul.New(
		client,
		consul.WithHealthCheck(true),
	)
}
