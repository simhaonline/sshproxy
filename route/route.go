package route

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"sshproxy/utils"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("sshproxy/route")

type selectDestinationFunc func([]string, bool) (string, string, error)

var (
	// default algorithm to find route
	DefaultAlgorithm = "ordered"
	// keyword for default route
	DefaultRouteKeyword = "default"
)

var (
	routeSelecters = map[string]selectDestinationFunc{
		"ordered": selectDestinationOrdered,
		"random":  selectDestinationRandom,
	}
)

// CanConnect tests if a connection to host:port can be made (with a 1s timeout).
func CanConnect(host, port string) bool {
	c, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 1*time.Second)
	if err != nil {
		log.Info("cannot connect to %s:%s: %s", host, port, err)
		return false
	}
	c.Close()
	return true
}

// selectDestinationOrdered selects the first reachable destination from a list
// of destinations. It returns its host and port.
func selectDestinationOrdered(destinations []string, check_host bool) (string, string, error) {
	for i, dst := range destinations {
		host, port, err := utils.SplitHostPort(dst)
		if err != nil {
			return "", "", err
		}

		// always return the last destination without trying to connect
		if i == len(destinations)-1 {
			return host, port, nil
		}
		if !check_host || CanConnect(host, port) {
			return host, port, nil
		}
	}
	return "", "", fmt.Errorf("no valid destination found")
}

// selectDestinationRandom randomizes the order of the provided list of
// destinations and selects the first reachable one. It returns its host and
// port.
func selectDestinationRandom(destinations []string, check_host bool) (string, string, error) {
	rand.Seed(time.Now().UnixNano())
	rdestinations := make([]string, len(destinations))
	perm := rand.Perm(len(destinations))
	for i, v := range perm {
		rdestinations[i] = destinations[v]
	}
	log.Debug("randomized destinations: %v", rdestinations)
	return selectDestinationOrdered(rdestinations, check_host)
}

func Select(route_select string, destinations []string, check_host bool) (string, string, error) {
	return routeSelecters[route_select](destinations, check_host)
}

func IsAlgorithm(algo string) bool {
	_, ok := routeSelecters[algo]
	return ok
}