package discovery_test

import (
	"regexp"
	"testing"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/stretchr/testify/assert"
)

func Test_MakePathValid(t *testing.T) {

	t.Run("WHEN given valid string SHOULD string should remain the same", func(t *testing.T) {
		t.Parallel()
		input := "/hello"

		registry := discovery.InMemoryRegistry{Services: make(map[string]discovery.ServiceInfo)}

		registry.MakePathValid(&input)

		assert.Equal(t, input, "/hello")

	})

	t.Run("WHEN given invalid string withough '/' as first character", func(t *testing.T) {
		t.Parallel()
		input := "hello/"

		registry := discovery.InMemoryRegistry{Services: make(map[string]discovery.ServiceInfo)}

		registry.MakePathValid(&input)

		assert.Equal(t, input, "/hello")
	})
}


func Test_SetServicePathRegex( t *testing.T){
	t.Run("WHEN there is no service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]discovery.ServiceInfo)
		registry := discovery.InMemoryRegistry{Services: services }


		registry.SetServicePathRegex()

		assert.Equal(t,registry.GetServicePathRegex(), "")
	})
	t.Run("WHEN there is one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]discovery.ServiceInfo)
		registry := discovery.InMemoryRegistry{Services: services }

		services["/hello"] = discovery.ServiceInfo{}

		registry.SetServicePathRegex()

		assert.Equal(t,registry.GetServicePathRegex(), "^(/hello)")
	})

	t.Run("WHEN there is more than one service SHOULD generate regex string when from service map", func(t *testing.T) {
		services := make(map[string]discovery.ServiceInfo)
		services["/req"] = discovery.ServiceInfo{}
		services["/hello"] = discovery.ServiceInfo{}
		registry := discovery.InMemoryRegistry{Services: services }

		registry.SetServicePathRegex()
		output := registry.GetServicePathRegex()

		t.Logf("%v\n hello", output)

		reg := regexp.MustCompile(`\^\((\/req|/hello)\|(\/req|/hello)\)`)
	
		assert.Equal(t,true, reg.MatchString(output))
	})
}

