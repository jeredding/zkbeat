package mntr

import (
	"bufio"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/helper"
	"github.com/jeredding/zkbeat/module/zookeeper"
	"io"
	"net"
	"regexp"
	"strconv"
	"time"
)

var (
	zkVersionRegex                 = regexp.MustCompile("zk_version\\s+(.*$)")
	zkServerStateRegex             = regexp.MustCompile("zk_server_state\\s+(.*$)")
	zkAvgLatencyRegex              = regexp.MustCompile("zk_avg_latency\\s+(.*$)")
	zkMaxLatencyRegex              = regexp.MustCompile("zk_max_latency\\s+(.*$)")
	zkMinLatencyRegex              = regexp.MustCompile("zk_min_latency\\s+(.*$)")
	zkPacketsReceivedRegex         = regexp.MustCompile("zk_packets_received\\s+(.*$)")
	zkPacketsSentRegex             = regexp.MustCompile("zk_packets_sent\\s+(.*$)")
	zkNumAliveConnectionsRegex     = regexp.MustCompile("zk_num_alive_connections\\s+(.*$)")
	zkOutstandingRequestsRegex     = regexp.MustCompile("zk_outstanding_requests\\s+(.*$)")
	zkNodeCountRegex               = regexp.MustCompile("zk_znode_count\\s+(.*$)")
	zkWatchCountRegex              = regexp.MustCompile("zk_watch_count\\s+(.*$)")
	zkEphemeralsCountRegex         = regexp.MustCompile("zk_ephemerals_count\\s+(.*$)")
	zkApproximateDataSizeRegex     = regexp.MustCompile("zk_approximate_data_size\\s+(.*$)")
	zkOpenFileDescriptorCountRegex = regexp.MustCompile("zk_open_file_descriptor_count\\s+(.*$)")
	zkMaxFileDescriptorCountRegex  = regexp.MustCompile("zk_max_file_descriptor_count\\s+(.*$)")
	zkFollowersRegex               = regexp.MustCompile("zk_followers\\s+(.*$)")
	zkSyncedFollowersRegex         = regexp.MustCompile("zk_synced_followers\\s+(.*$)")
	zkPendingSyncsRegex            = regexp.MustCompile("zk_pending_syncs\\s+(.*$)")
)

const command = "mntr"

func init() {
	helper.Registry.AddMetricSeter("zookeeper", "mntr", New)
}

func New() helper.MetricSeter {
	return &MetricSeter{}
}

type MetricSeter struct {
	Hostname string
	Port     string
	Timeout  time.Duration
}

func (m *MetricSeter) Setup(ms *helper.MetricSet) error {

	// Additional configuration options
	config := struct {
		Hostname string `config:"hostname"`
		Port     string `config:"port"`
		Timeout  string `config:"timeout"`
	}{
		Hostname: "127.0.0.1",
		Port:     "2181",
		Timeout:  "60s",
	}

	if err := ms.Module.ProcessConfig(&config); err != nil {
		return err
	}

	m.Hostname = config.Hostname
	m.Port = config.Port
	m.Timeout, _ = time.ParseDuration(config.Timeout)
	return nil
}

func (m *MetricSeter) Fetch(ms *helper.MetricSet, host string) (event common.MapStr, err error) {
	connectionString := net.JoinHostPort(m.Hostname, m.Port)
	outputReader, err := zookeeper.RunCommand(command, connectionString, m.Timeout)
	if err != nil {
		logp.Err("Error running four-letter command %s on %s: %v", command, connectionString, err)
		return nil, err
	}
	return mntrEventMapping(outputReader), nil
}

func mntrEventMapping(response io.Reader) common.MapStr {

	var (
		versionString           string
		serverState             string
		avgLatency              int
		minLatency              int
		maxLatency              int
		packetsReceived         int
		packetsSent             int
		numAliveConnections     int
		outstandingRequests     int
		znodeCount              int
		watchCount              int
		ephemeralsCount         int
		approximateDataSize     int
		openFileDescriptorCount int
		maxFileDescriptorCount  int
		followers               int
		syncedFollowers         int
		pendingSyncs            int
	)

	//  zk_version      3.5.1-alpha-1693007, built on 07/28/2015 07:19 GMT
	//  zk_avg_latency  0
	//  zk_max_latency  1789
	//  zk_min_latency  0
	//  zk_packets_received     22152032
	//  zk_packets_sent 30959914
	//  zk_num_alive_connections        1033
	//  zk_outstanding_requests 0
	//  zk_server_state leader
	//  zk_znode_count  242609
	//  zk_watch_count  940522
	//  zk_ephemerals_count     8565
	//  zk_approximate_data_size        372143564
	//  zk_open_file_descriptor_count   1083
	//  zk_max_file_descriptor_count    1048576
	//  zk_followers    5
	//  zk_synced_followers     2
	//  zk_pending_syncs        0

	scanner := bufio.NewScanner(response)
	for scanner.Scan() {

		if matches := zkVersionRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			versionString = matches[1]
		}

		if matches := zkServerStateRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			serverState = matches[1]
		}

		if matches := zkAvgLatencyRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			avgLatency, _ = strconv.Atoi(matches[1])
		}

		if matches := zkMaxLatencyRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			maxLatency, _ = strconv.Atoi(matches[1])
		}

		if matches := zkMinLatencyRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			minLatency, _ = strconv.Atoi(matches[1])
		}

		if matches := zkPacketsReceivedRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			packetsReceived, _ = strconv.Atoi(matches[1])
		}

		if matches := zkPacketsSentRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			packetsSent, _ = strconv.Atoi(matches[1])
		}

		if matches := zkNumAliveConnectionsRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			numAliveConnections, _ = strconv.Atoi(matches[1])
		}

		if matches := zkOutstandingRequestsRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			outstandingRequests, _ = strconv.Atoi(matches[1])
		}

		if matches := zkNodeCountRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			znodeCount, _ = strconv.Atoi(matches[1])
		}

		if matches := zkWatchCountRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			watchCount, _ = strconv.Atoi(matches[1])
		}

		if matches := zkEphemeralsCountRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			ephemeralsCount, _ = strconv.Atoi(matches[1])
		}

		if matches := zkApproximateDataSizeRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			approximateDataSize, _ = strconv.Atoi(matches[1])
		}

		if matches := zkOpenFileDescriptorCountRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			openFileDescriptorCount, _ = strconv.Atoi(matches[1])
		}

		if matches := zkMaxFileDescriptorCountRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			maxFileDescriptorCount, _ = strconv.Atoi(matches[1])
		}

		if matches := zkFollowersRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			followers, _ = strconv.Atoi(matches[1])
		}

		if matches := zkSyncedFollowersRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			syncedFollowers, _ = strconv.Atoi(matches[1])
		}

		if matches := zkPendingSyncsRegex.FindStringSubmatch(scanner.Text()); matches != nil {
			pendingSyncs, _ = strconv.Atoi(matches[1])
		}
	}

	event := common.MapStr{
		"version_string":             versionString,
		"avg_latency":                avgLatency,
		"min_latency":                minLatency,
		"max_latency":                maxLatency,
		"packets_received":           packetsReceived,
		"packets_sent":               packetsSent,
		"num_alive_connections":      numAliveConnections,
		"outstanding_requests":       outstandingRequests,
		"server_state":               serverState,
		"znode_count":                znodeCount,
		"watch_count":                watchCount,
		"ephemerals_count":           ephemeralsCount,
		"approximate_data_size":      approximateDataSize,
		"open_file_descriptor_count": openFileDescriptorCount,
		"max_file_descriptor_count":  maxFileDescriptorCount,
		"followers":                  followers,
		"synced_followers":           syncedFollowers,
		"pending_syncs":              pendingSyncs,
	}
	return event
}
