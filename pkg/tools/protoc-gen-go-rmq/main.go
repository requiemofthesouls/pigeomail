package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	GenTypeServer = "server"
	GenTypeClient = "client"
)

func main() {
	var flags flag.FlagSet
	genType := flags.String("gen_type", "", "generation type: server or client")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			switch *genType {
			case GenTypeServer:
				generateServers(gen, f)
			case GenTypeClient:
				generateClients(gen, f)
			default:
				return fmt.Errorf("unknown generation type %s", *genType)
			}
		}

		return nil
	})
}
