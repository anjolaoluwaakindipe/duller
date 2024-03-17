package discovery_test

import (
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Discommand_Init(t *testing.T) {
	t.Run("Test that default values are used when no flags are passed", func(t *testing.T) {
		command := discovery.NewDiscCommand()

		err := command.Init([]string{})

		assert.Nil(t, err)

		assert.Equal(t, utils.DISCOVERY_KEY, command.DiscoveryKey)
		assert.Equal(t, utils.DISCOVERY_PORT, command.DiscoveryPort)
		assert.Equal(t, utils.DISCOVERY_SERVICE_PATH, command.DiscoveryServicePath)
		assert.Equal(t, utils.HEARTBEAT_INTERVAL, command.DiscoveryHeartbeatInterval)
	})

	t.Run("Test that arguement values are used when flags are passed", func(t *testing.T) {
		heartbeat := (10 * time.Second)
		argMap := map[string]string{
			utils.DISCOVERY_KEY_FLAG:          "hello",
			utils.DISCOVERY_PORT_FLAG:         "9999",
			utils.DISCOVERY_SERVICE_PATH_FLAG: "/new-service-path",
			utils.HEARTBEAT_INTERVAL_FLAG:     heartbeat.String(),
		}

		args := make([]string, 0)

		for k, v := range argMap {
			args = append(args, "--"+k+"="+v)
		}

		command := discovery.NewDiscCommand()

		err := command.Init(args)

		assert.Nil(t, err)

		assert.Equal(t, argMap[utils.DISCOVERY_KEY_FLAG], command.DiscoveryKey)
		assert.Equal(t, argMap[utils.DISCOVERY_PORT_FLAG], command.DiscoveryPort)
		assert.Equal(t, argMap[utils.DISCOVERY_SERVICE_PATH_FLAG], command.DiscoveryServicePath)
		assert.Equal(t, heartbeat, command.DiscoveryHeartbeatInterval)
	})
}
