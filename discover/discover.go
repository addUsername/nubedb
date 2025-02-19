// Package discover is responsible for handling the discovery of nubedb nodes.
package discover

import (
	"errors"
	"github.com/narvikd/errorskit"
	"github.com/narvikd/mdns"
	"log"
	"net"
	"nubedb/cluster"
	"nubedb/internal/config"
	"strings"
	"sync"
	"time"
)

const (
	// The service name identifier used for the discovery.
	serviceName       = "_nubedb._tcp"
	ErrLeaderNotFound = "couldn't find a leader"
)

// ServeAndBlock creates a new discovery service with the given node ID and port, blocks indefinitely.
func ServeAndBlock(nodeID string, port int) {
	const errGen = "Discover serve and block: "
	info := []string{"nubedb Discover"}

	ip, errGetIP := getIP(nodeID)
	if errGetIP != nil {
		errorskit.FatalWrap(errGetIP, errGen)
	}

	// Create a new mDNS service for the node.
	service, errService := mdns.NewMDNSService(nodeID, serviceName, "", "", port, []net.IP{ip}, info)
	if errService != nil {
		errorskit.FatalWrap(errService, errGen+"discover service")
	}

	// Create a new mDNS server for the service.
	server, errServer := mdns.NewServer(&mdns.Config{Zone: service})
	if errServer != nil {
		errorskit.FatalWrap(errService, errGen+"discover server")
	}

	// Shut down the server when the function returns. (Which shouldn't)
	defer func(server *mdns.Server) {
		_ = server.Shutdown()
	}(server)

	// Block indefinitely.
	select {}
}

func getIP(nodeID string) (net.IP, error) {
	hosts, errLookup := net.LookupHost(nodeID)
	if errLookup != nil {
		return nil, errorskit.Wrap(errLookup, "couldn't lookup host")
	}

	return net.ParseIP(hosts[0]), nil
}

// SearchNodes returns a list of all discovered nodes, excluding the one passed as a parameter.
func SearchNodes(currentNode string) ([]string, error) {
	// map to store the discovered nodes.
	hosts := make(map[string]bool)
	var lastError error

	// Try to discover nodes 3 times to add any missing nodes in the first scan.
	for i := 0; i < 3; i++ {
		hostsQuery, err := query()
		if err != nil {
			log.Println(err)
			lastError = err
			continue
		}

		for _, host := range hostsQuery {
			// In some linux versions it reports "$name." (name and a dot)
			host = strings.ReplaceAll(host, ".", "")
			hosts[host] = true
		}
		// Wait for 100 milliseconds before trying again to not spam/have some space between requests.
		time.Sleep(100 * time.Millisecond) // TODO: Try to refactor this
	}

	// Convert the map to a slice of strings and exclude the current node.
	result := make([]string, 0, len(hosts))
	for host := range hosts {
		if currentNode == host {
			continue
		}
		result = append(result, host)
	}

	return result, lastError
}

// query sends an mDNS query to discover nubedb nodes and returns a list of their hosts.
func query() ([]string, error) {
	var mu sync.Mutex
	var hosts []string
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			mu.Lock()
			hosts = append(hosts, entry.Host)
			mu.Unlock()
		}
	}()

	params := mdns.DefaultParams(serviceName)
	params.DisableIPv6 = true
	params.Entries = entriesCh

	defer close(entriesCh)
	err := mdns.Query(params)
	if err != nil {
		return nil, errorskit.Wrap(err, "discover search")
	}

	mu.Lock()
	defer mu.Unlock()
	return hosts, nil
}

// SearchLeader will return an error if a leader is not found,
// since it skips the current node and this could be a leader.
//
// If the current node is as leader, it will still return an error
func SearchLeader(currentNode string) (string, error) {
	nodes, errNodes := SearchNodes(currentNode)
	if errNodes != nil {
		return "", errNodes
	}

	for _, node := range nodes {
		leader, err := cluster.IsLeader(config.MakeGrpcAddress(node))
		if err != nil {
			errorskit.LogWrap(err, "couldn't contact node while searching for leaders")
			continue
		}
		if leader {
			return node, nil
		}
	}

	return "", errors.New(ErrLeaderNotFound)
}
