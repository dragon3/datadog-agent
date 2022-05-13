package netflow

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/epforwarder"
	"github.com/DataDog/datadog-agent/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestNewNetflowServer(t *testing.T) {
	// Setup NetFlow feature config
	port := uint16(52055)
	config.Datadog.SetConfigType("yaml")
	err := config.Datadog.MergeConfigOverride(strings.NewReader(fmt.Sprintf(`
network_devices:
  netflow:
    enabled: true
    aggregator_flush_interval: 1
    listeners:
      - flow_type: netflow5 # netflow, sflow, ipfix
        bind_host: 127.0.0.1
        port: %d # default 2055 for netflow
`, port)))
	require.NoError(t, err)

	// Setup NetFlow Server
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(1 * time.Millisecond)
	defer demux.Stop(false)

	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)

	// Send netflowV5Data twice to test aggregator
	// Flows will have 2x bytes/packets after aggregation
	err = sendUDPPacket(port, mockNetflowV5Data)
	require.NoError(t, err, "error sending udp packet")

	// Get Event Platform Events
	netflowEvents, err := demux.WaitEventPlatformEvents(epforwarder.EventTypeNetworkDevicesNetFlow, 6, 15*time.Second)
	require.NoError(t, err, "error waiting event platform events")
	assert.Equal(t, 6, len(netflowEvents))

	actualFlow, err := findEventBySourceDest(netflowEvents, "10.129.2.1", "10.128.2.119")
	assert.NoError(t, err)

	assert.Equal(t, "netflow5", actualFlow.FlowType)
	assert.Equal(t, uint64(0), actualFlow.SamplingRate)
	assert.Equal(t, "ingress", actualFlow.Direction)
	assert.Equal(t, uint64(1540209168), actualFlow.Start)
	assert.Equal(t, uint64(1540209169), actualFlow.End)
	assert.Equal(t, uint64(194), actualFlow.Bytes)
	assert.Equal(t, "2048", actualFlow.EtherType)
	assert.Equal(t, "6", actualFlow.IPProtocol)
	assert.Equal(t, uint32(0), actualFlow.Tos)
	assert.Equal(t, "127.0.0.1", actualFlow.Exporter.IP)
	assert.Equal(t, "10.129.2.1", actualFlow.Source.IP)
	assert.Equal(t, uint32(49452), actualFlow.Source.Port)
	assert.Equal(t, "00:00:00:00:00:00", actualFlow.Source.Mac)
	assert.Equal(t, "0.0.0.0/24", actualFlow.Source.Mask)
	assert.Equal(t, "10.128.2.119", actualFlow.Destination.IP)
	assert.Equal(t, uint32(8080), actualFlow.Destination.Port)
	assert.Equal(t, "", actualFlow.Destination.Mac)
	assert.Equal(t, "", actualFlow.Destination.Mask)
	assert.Equal(t, uint32(1), actualFlow.Ingress.Interface.Index)
	assert.Equal(t, uint32(7), actualFlow.Egress.Interface.Index)
	assert.Equal(t, "default", actualFlow.Namespace)
	hostname, _ := util.GetHostname(context.TODO())
	assert.Equal(t, hostname, actualFlow.Host)
	assert.ElementsMatch(t, []string{"SYN", "ACK"}, actualFlow.TCPFlags)
	assert.Equal(t, "0.0.0.0", actualFlow.NextHop.IP)
}

func TestStartServerAndStopServer(t *testing.T) {
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	err := StartServer(demux)
	require.NoError(t, err)
	require.NotNil(t, serverInstance)

	replaceWithDummyFlowProcessor(serverInstance, 123)

	StopServer()
	require.Nil(t, serverInstance)
}

func TestIsEnabled(t *testing.T) {
	saved := config.Datadog.Get("network_devices.netflow.enabled")
	defer config.Datadog.Set("network_devices.netflow.enabled", saved)

	config.Datadog.Set("network_devices.netflow.enabled", true)
	assert.Equal(t, true, IsEnabled())

	config.Datadog.Set("network_devices.netflow.enabled", false)
	assert.Equal(t, false, IsEnabled())
}

func TestServer_Stop(t *testing.T) {
	// Setup NetFlow config
	port := uint16(12056)
	config.Datadog.SetConfigType("yaml")
	err := config.Datadog.MergeConfigOverride(strings.NewReader(fmt.Sprintf(`
network_devices:
  netflow:
    enabled: true
    aggregator_flush_interval: 1
    listeners:
      - flow_type: netflow5 # netflow, sflow, ipfix
        bind_host: 0.0.0.0
        port: %d # default 2055 for netflow
`, port)))
	require.NoError(t, err)

	// Setup Netflow Server
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)

	flowProcessor := replaceWithDummyFlowProcessor(server, port)

	// Stops server
	server.stop()

	// Assert logs present
	assert.Equal(t, flowProcessor.stopped, true)
}
