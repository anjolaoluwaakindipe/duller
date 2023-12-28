package discovery_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/stretchr/testify/assert"
)

func Test_SetServicePathRegex(t *testing.T) {
	t.Run("WHEN there is no service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]*discovery.ServiceInfo)
		registry := discovery.InMemoryRegistry{Services: services, Clock: &FakeTime{}}

		registry.SetServicePathRegex()

		assert.Equal(t, registry.GetServicePathRegex(), "")
	})
	t.Run("WHEN there is one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]*discovery.ServiceInfo)
		registry := discovery.InMemoryRegistry{Services: services, Clock: &FakeTime{}}
		services["/hello"] = &discovery.ServiceInfo{}
		registry.SetServicePathRegex()
		assert.Equal(t, registry.GetServicePathRegex(), "^(/hello)")
	})

	t.Run("WHEN there is more than one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]*discovery.ServiceInfo)
		services["/req"] = &discovery.ServiceInfo{}
		services["/hello"] = &discovery.ServiceInfo{}
		registry := discovery.InMemoryRegistry{Services: services, Clock: &FakeTime{}}
		registry.SetServicePathRegex()
		output := registry.GetServicePathRegex()
		reg := regexp.MustCompile(`\^\((\/req|/hello)\|(\/req|/hello)\)`)
		assert.Equal(t, true, reg.MatchString(output))
	})
}

func Test_RegisterService(t *testing.T) {
	t.Run("SHOULD create new add it to the registry service WHEN given a valid RegisterServiceMessage with a service that does not exist ", func(t *testing.T) {
		newMessage := discovery.RegisterServiceMessage{Path: "/hello", Address: "http://localhost:3000", ServiceName: "server_1"}

		registry := discovery.InMemoryRegistry{Services: make(map[string]*discovery.ServiceInfo), Clock: &FakeTime{}}

		if err := registry.RegisterService(newMessage); err != nil {
			t.Error("Error while Registering Service")
			return
		}

		val, ok := registry.Services[newMessage.Path]

		assert.True(t, ok)

		if ok {
			assert.Equal(t, newMessage.Address, val.Address)
			assert.Equal(t, newMessage.Path, val.Path)
			assert.Equal(t, newMessage.ServiceName, val.ServiceId)
		}
	})

	t.Run("SHOULD update a service WHEN a valid RegisterServiceMessage with a service that already exists is given", func(t *testing.T) {
		// setup up two calls to RegisterService
		newMessage := discovery.RegisterServiceMessage{Path: "/hello", Address: "http://localhost:3000", ServiceName: "server_1"}

		stubTime := &FakeTime{time.Now()}
		registry := discovery.InMemoryRegistry{Services: make(map[string]*discovery.ServiceInfo), Clock: stubTime}
		if err := registry.RegisterService(newMessage); err != nil {
			t.Error("Error while Registering Service")
			return
		}

		service := registry.Services[newMessage.Path]

		stubTime.CurrentTime = stubTime.CurrentTime.Add(5 * time.Second)

		createdAt := service.LastHeartbeat

		if err := registry.RegisterService(newMessage); err != nil {
			t.Error("Error while Registering Service")
			return
		}

		service = registry.Services[newMessage.Path]
		updatedAt := service.LastHeartbeat

		assert.NotEqual(t, createdAt, updatedAt)
	})
}

// Mocks

type FakeTime struct {
	CurrentTime time.Time
}

func (ft *FakeTime) Now() time.Time {
	return ft.CurrentTime
}
