package rmq

import (
	"context"
	"sync"

	rmqClDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/client/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server"
)

type (
	Manager interface {
		StartAll(ctx context.Context)
	}

	manager struct {
		connections []connection.Manager
		clients     []rmqClDef.Manager
		servers     []server.Manager
		stopChan    chan struct{}
	}
)

func NewManager(
	connections []connection.Manager,
	clients []rmqClDef.Manager,
	servers []server.Manager,
) *manager {
	return &manager{
		connections: connections,
		clients:     clients,
		servers:     servers,
		stopChan:    make(chan struct{}),
	}
}

func (m *manager) StartAll(ctx context.Context) {
	go m.pendingClosure(ctx)

	m.startServers()

	<-m.stopChan
	m.closeConnections()
}

func (m *manager) pendingClosure(ctx context.Context) {
	<-ctx.Done()

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() { defer wg.Done(); m.closeServers() }()
	go func() { defer wg.Done(); m.closeClients() }()

	wg.Wait()

	close(m.stopChan)
}

func (m *manager) closeClients() {
	wg := &sync.WaitGroup{}

	for _, cl := range m.clients {
		wg.Add(1)
		go func(cl rmqClDef.Manager) {
			defer wg.Done()
			cl.Close()
		}(cl)
	}

	wg.Wait()
}

func (m *manager) closeServers() {
	wg := &sync.WaitGroup{}

	for _, srv := range m.servers {
		wg.Add(1)
		go func(srv server.Manager) {
			srv.CloseAll()
			wg.Done()
		}(srv)
	}

	wg.Wait()
}

func (m *manager) closeConnections() {
	wg := &sync.WaitGroup{}

	for _, conn := range m.connections {
		wg.Add(1)
		go func(conn connection.Manager) {
			conn.Close()
			wg.Done()
		}(conn)
	}

	wg.Wait()
}

func (m *manager) startServers() {
	wg := &sync.WaitGroup{}
	for _, srv := range m.servers {
		wg.Add(1)
		go func(srv server.Manager) {
			srv.StartAll()
			wg.Done()
		}(srv)
	}

	wg.Wait()
}
