package registry_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/registry"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_SetServicePathRegex(t *testing.T) {
	t.Run("WHEN there is no service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string][]*service.ServiceInfo)
		registry := registry.InMemoryRegistry{PathTable: services, Clock: &FakeTime{}}

		registry.SetServicePathRegex()

		assert.Equal(t, registry.GetServicePathRegex(), "")
	})
	t.Run("WHEN there is one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string][]*service.ServiceInfo)
		registry := registry.InMemoryRegistry{PathTable: services, Clock: &FakeTime{}}
		services["/hello"] = []*service.ServiceInfo{{Path: "/hello"}}
		registry.SetServicePathRegex()
		assert.Equal(t, registry.GetServicePathRegex(), "^(/hello)")
	})

	t.Run("WHEN there is more than one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string][]*service.ServiceInfo)
		services["/req"] = []*service.ServiceInfo{}
		services["/hello"] = []*service.ServiceInfo{}
		registry := registry.InMemoryRegistry{PathTable: services, Clock: &FakeTime{}}
		registry.SetServicePathRegex()
		output := registry.GetServicePathRegex()
		reg := regexp.MustCompile(`\^\((\/req|/hello)\|(\/req|/hello)\)`)
		assert.Equal(t, true, reg.MatchString(output))
	})
}

func Test_RegisterService(t *testing.T) {
	t.Run("SHOULD create new add it to the registry service WHEN given a valid RegisterServiceMessage with a service that does not exist ", func(t *testing.T) {
		newMessage := service.ServiceInfo{Path: "/hello", IP: "http://localhost", Port: "3000", ServiceId: "server_1"}

		registry := registry.InMemoryRegistry{PathTable: make(map[string][]*service.ServiceInfo), Clock: &FakeTime{}, ServiceIdTable: make(map[string]*service.ServiceInfo)}

		if err := registry.RegisterService(&newMessage); err != nil {
			t.Error("Error while Registering Service")
			return
		}

		services, err := registry.GetServicesByPath(newMessage.Path)

		assert.Nil(t, err)

		if err != nil {
			assert.Len(t, services, 1)
			assert.Equal(t, newMessage.IP, services[0].IP)
			assert.Equal(t, newMessage.Path, services[0].Path)
			assert.Equal(t, newMessage.ServiceId, services[0].ServiceId)
			assert.Equal(t, newMessage.Port, services[0].Port)
		}
	})

	t.Run("SHOULD update a service WHEN a valid RegisterServiceMessage with a service that already exists is given", func(t *testing.T) {
		newMessage := service.ServiceInfo{Path: "/hello", IP: "127.0.0.1", Port: "9990", ServiceId: "server_1"}

		stubTime := &FakeTime{time.Now()}
		registry := registry.InMemoryRegistry{PathTable: make(map[string][]*service.ServiceInfo), Clock: stubTime, ServiceIdTable: make(map[string]*service.ServiceInfo)}
		if err := registry.RegisterService(&newMessage); err != nil {
			t.Error("Error while Registering Service ")
			return
		}

		service, err := registry.GetServiceById(newMessage.ServiceId)
		if err != nil {
			t.Errorf("Error while getting Service %v: %v", newMessage.ServiceId, err)
			return
		}

		stubTime.CurrentTime = stubTime.CurrentTime.Add(5 * time.Second)

		createdAt := service.LastHeartbeat

		if err := registry.RegisterService(&newMessage); err != nil {
			t.Error("Error while Updating Service")
			return
		}

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
