package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strings"
)

func NewConsulClient(address string) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = address
	return api.NewClient(config)
}

func RegisterService(
	client *api.Client,
	serviceName, serviceID, address, serviceType string,
	tags []string,
) error {
	port := getPortFromAddress(address)

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: getHostFromAddress(address),
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s/healthz", address),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	if serviceType == "grpc" {
		registration.Check.GRPC = fmt.Sprintf("%s/grpc.health.v1.Health/Check", address)
		registration.Check.Interval = "15s"
		registration.Check.HTTP = ""
	}

	return client.Agent().ServiceRegister(registration)
}
func DeregisterService(client *api.Client, serviceID string) error {
	return client.Agent().ServiceDeregister(serviceID)
}

func getPortFromAddress(address string) int {
	var port int
	_, err := fmt.Sscanf(address, ":%d", &port)
	if err != nil {

		return 0
	}
	return port
}
func getHostFromAddress(address string) string {
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		return address
	}
	return strings.Join(parts[:len(parts)-1], ":")
}
