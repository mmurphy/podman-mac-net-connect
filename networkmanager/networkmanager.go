package networkmanager

import (
	"bytes"
	"fmt"
	"os/exec"

	containerTypes "github.com/containers/common/libnetwork/types"
)

type NetworkManager struct {
	PodmanNetworks map[string]containerTypes.Network
}

func New() NetworkManager {
	return NetworkManager{
		PodmanNetworks: map[string]containerTypes.Network{},
	}
}

// Set the point-to-point IP address configuration on a network interface.
func (manager *NetworkManager) SetInterfaceAddress(ip string, peerIp string, iface string) (string, string, error) {

	cmd := exec.Command("ifconfig", iface, "inet", ip+"/32", peerIp)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

// Add a route to the macOS routing table.
func (manager *NetworkManager) AddRoute(net string, iface string) (string, string, error) {

	cmd := exec.Command("route", "-q", "-n", "add", "-inet", net, "-interface", iface)
	fmt.Printf("cmd: %v\n", cmd.String())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

// Delete a route from the macOS routing table.
func (manager *NetworkManager) DeleteRoute(net string) (string, string, error) {

	cmd := exec.Command("route", "-q", "-n", "delete", "-inet", net)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func (manager *NetworkManager) ProcessPodmanNetworkCreate(network containerTypes.Network, iface string) {
	manager.PodmanNetworks[network.ID] = network

	for _, subnet := range network.Subnets {
		fmt.Printf("network interface: %v\n", network.NetworkInterface)
		if network.Driver == "bridge" {
			fmt.Printf("Adding route for %s -> %s (%s)\n", subnet.Subnet, iface, network.Name)

			_, stderr, err := manager.AddRoute(subnet.Subnet.String(), iface)

			if err != nil {
				fmt.Errorf("Failed to add route: %v. %v\n", err, stderr)
			}
		}
	}
}
