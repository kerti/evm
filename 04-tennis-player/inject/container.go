package inject

import (
	"github.com/facebookgo/inject"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

// Service is the service interface
type Service interface {
	Startup()
	Shutdown()
}

// ServiceWithPrepareShutdown is the service interface with prepare shutdown method
type ServiceWithPrepareShutdown interface {
	Startup()
	PrepareShutdown()
	Shutdown()
}

// ContainerAware is an interface for container-aware components
type ContainerAware interface {
	SetContainer(c ServiceContainer)
}

// ServiceContainer is the service container interface
type ServiceContainer interface {
	Ready() error
	GetService(id string) (interface{}, bool)
	RegisterService(id string, svc interface{})
	RegisterServices(services map[string]interface{})
	PrepareShutdown()
	Shutdown()
}

// ServiceRegistry is a struct for keeping tabs on registered services
type ServiceRegistry struct {
	graph    inject.Graph
	services map[string]interface{}
	order    map[int]string
	ready    bool
}

// GetService fetches a service by its ID
func (reg *ServiceRegistry) GetService(id string) (interface{}, bool) {
	svc, ok := reg.services[id]
	return svc, ok
}

// RegisterService registers a service
func (reg *ServiceRegistry) RegisterService(id string, svc interface{}) {
	err := reg.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		panic(err.Error())
	}
	reg.order[len(reg.order)] = id
	reg.services[id] = svc
}

// RegisterServices registers multiple services
func (reg *ServiceRegistry) RegisterServices(services map[string]interface{}) {
	for id, svc := range reg.services {
		reg.RegisterService(id, svc)
	}
}

// Ready starts up the service graph and returns error if it's not ready
func (reg *ServiceRegistry) Ready() error {
	if reg.ready {
		return nil
	}
	err := reg.graph.Populate()
	if err != nil {
		return err
	}
	for i := 0; i < len(reg.order); i++ {
		k := reg.order[i]
		obj := reg.services[k]
		if s, ok := obj.(Service); ok {
			containerAware, okc := s.(ContainerAware)
			if okc {
				containerAware.SetContainer(reg)
			}
			s.Startup()
		}
	}
	reg.ready = err == nil
	return err
}

// PrepareShutdown prepares all services for eventual shutdown
func (reg *ServiceRegistry) PrepareShutdown() {
	for _, name := range reg.order {
		if service, ok := reg.services[name]; ok {
			if s, ok := service.(ServiceWithPrepareShutdown); ok {
				s.PrepareShutdown()
			}
		}
	}
}

// Shutdown shuts down all services
func (reg *ServiceRegistry) Shutdown() {
	for _, name := range reg.order {
		if service, ok := reg.services[name]; ok {
			if s, ok := service.(Service); ok {
				s.Shutdown()
			}
		}
	}
}

// NewContainer creates a new service container
func NewContainer() ServiceContainer {
	logger.Debug("Initializing service container...")
	return &ServiceRegistry{services: make(map[string]interface{}), order: make(map[int]string), ready: false}
}
