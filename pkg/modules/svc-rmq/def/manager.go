package def

import (
	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	rmq "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq"
	rmqClDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/client/def"
	rmqConnDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection/def"
	rmqSrvDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/def"
)

type (
	Config struct {
		Connections rmqConnDef.Configs `mapstructure:"connections"`
		Clients     rmqClDef.Configs   `mapstructure:"clients"`
		Servers     rmqSrvDef.Configs  `mapstructure:"servers"`
	}

	Manager = rmq.Manager
)

const (
	DIManager = "rmq.manager"
)

func (c *Config) GetConnections(cont container.Container) ([]rmqConnDef.Manager, error) {
	connections := make([]rmqConnDef.Manager, 0, len(c.Connections))
	for _, connection := range c.Connections {
		var conn rmqConnDef.Manager
		if err := cont.Fill(rmqConnDef.DIManagerPrefix+connection.Name, &conn); err != nil {
			return nil, err
		}

		connections = append(connections, conn)
	}

	return connections, nil
}

func (c *Config) GetClients(cont container.Container) ([]rmqClDef.Manager, error) {
	clients := make([]rmqClDef.Manager, 0, len(c.Servers))
	for _, server := range c.Clients {
		var rmqCl rmqClDef.Manager
		if err := cont.Fill(rmqClDef.DIManagerPrefix+server.Name, &rmqCl); err != nil {
			return nil, err
		}

		clients = append(clients, rmqCl)
	}

	return clients, nil
}

func (c *Config) GetServers(cont container.Container) ([]rmqSrvDef.Manager, error) {
	servers := make([]rmqSrvDef.Manager, 0, len(c.Servers))
	for _, server := range c.Servers {
		var srv rmqSrvDef.Manager
		if err := cont.Fill(rmqSrvDef.DIManagerPrefix+server.Name, &srv); err != nil {
			return nil, err
		}

		servers = append(servers, srv)
	}

	return servers, nil
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		cont := builder.Build()

		var cfg cfgDef.Wrapper
		if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
			return err
		}

		var kCfg Config
		if err := cfg.UnmarshalKey("rmq", &kCfg); err != nil {
			return err
		}

		if err := kCfg.Connections.AddDefs(builder); err != nil {
			return err
		}

		if err := kCfg.Clients.AddDefs(builder); err != nil {
			return err
		}

		return kCfg.Servers.AddDefs(builder)
	})

	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIManager,
			Build: func(cont container.Container) (interface{}, error) {
				var cfg cfgDef.Wrapper
				if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
					return nil, err
				}

				var kCfg Config
				if err := cfg.UnmarshalKey("rmq", &kCfg); err != nil {
					return nil, err
				}

				var (
					connections []rmqConnDef.Manager
					err         error
				)
				if connections, err = kCfg.GetConnections(cont); err != nil {
					return nil, err
				}

				var clients []rmqClDef.Manager
				if clients, err = kCfg.GetClients(cont); err != nil {
					return nil, err
				}

				var servers []rmqSrvDef.Manager
				if servers, err = kCfg.GetServers(cont); err != nil {
					return nil, err
				}

				return rmq.NewManager(connections, clients, servers), nil
			},
		})
	})
}
