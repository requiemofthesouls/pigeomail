package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type eventImport struct {
	path string
	name string
}

func newEventImport(path string, num int) eventImport {
	return eventImport{
		path: strings.Replace(path, "/proto", "/pb", -1),
		name: fmt.Sprintf("events%d", num),
	}
}

func listImports(file *protogen.File) map[string]eventImport {
	imports, numEvents := make(map[string]eventImport, 0), 0
	for _, service := range file.Services {
		for _, method := range service.Methods {
			path := method.Input.GoIdent.GoImportPath.String()
			if _, ok := imports[path]; !ok {
				numEvents++
				imports[path] = newEventImport(path, numEvents)
			}
		}
	}

	return imports
}

func getEventStructureName(listEventsImports map[string]eventImport, msg *protogen.Message) string {
	return fmt.Sprintf("%s.%s", listEventsImports[msg.GoIdent.GoImportPath.String()].name, msg.GoIdent.GoName)
}
