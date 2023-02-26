package container

import (
	"fmt"
	"sync"

	"github.com/sarulabs/di/v2"
)

var (
	mu        sync.Mutex
	container di.Container
	builders  []buildFn
	lock      sync.Mutex
)

type (
	Def       = di.Def
	Key       = di.ContainerKey
	Builder   = di.Builder
	Container = di.Container

	buildFn func(builder *Builder, params map[string]interface{}) error
)

const (
	App        = di.App
	Request    = di.Request
	SubRequest = di.SubRequest
)

// Register definition builder
func Register(fn buildFn) {
	mu.Lock()
	defer mu.Unlock()

	builders = append(builders, fn)
}

// Instance return container
func Instance(scopes []string, params map[string]interface{}) (di.Container, error) {
	lock.Lock()
	defer lock.Unlock()
	if container != nil {
		return container, nil
	}

	var (
		builder *di.Builder
		err     error
	)
	if builder, err = di.NewBuilder(scopes...); err != nil {
		return nil, fmt.Errorf("can't create container builder: %v", err)
	}

	for _, fn := range builders {
		if err := fn(builder, params); err != nil {
			return nil, err
		}
	}

	return builder.Build(), nil
}
