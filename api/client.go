package api

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// ServiceType represents the type of OpenStack service.
type ServiceType string

const (
	// BareMetalV1 identifies the OpenStack Baremetal V1 service (BareMetal).
	BareMetalV1 ServiceType = "openstack_baremetal_v1"
	// IdentityV3 identifies the OpenStack Identity V3 service (Identity).
	IdentityV3 = "openstack_identity_v3"
	// Compute identifies the penStack Compute V2 service (Compute).
	ComputeV2 = "openstack_compute_v2"
	// NetworkingV2 identifies the OpenStack Network V2 service (Networking).
	NetworkingV2 = "openstack_networking_v2"
	// BlockStorageV3 identifies the OpenStack Block Storage V3 service (BlockStorage).
	BlockStorageV3 = "openstack_blockstorage_v3"
	// ImageV2 identifies the OpenStack Image Service V2 service (Image).
	ImageV2 = "openstack_image_v2"
)

type ServiceInfo struct {
	constructor  func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	microversion string
}

type ServiceMap map[ServiceType]ServiceInfo

var serviceConfigMap = map[ServiceType]ServiceInfo{
	BareMetalV1: {
		constructor:  openstack.NewBareMetalV1,
		microversion: "1.58",
	},
	IdentityV3: {
		constructor:  openstack.NewIdentityV3,
		microversion: "3.13",
	},
	ComputeV2: {
		constructor:  openstack.NewComputeV2,
		microversion: "2.79",
	},
	NetworkingV2: {
		constructor:  openstack.NewNetworkV2,
		microversion: "2.35",
	},
	BlockStorageV3: {
		constructor:  openstack.NewBlockStorageV3,
		microversion: "3.59",
	},
	ImageV2: {
		constructor:  openstack.NewImageServiceV2,
		microversion: "2.9",
	},
}

const Region = "RegionOne"

type APIClient struct {
	Client   *gophercloud.ProviderClient
	mutex    sync.RWMutex
	services map[ServiceType]*gophercloud.ServiceClient
}

func (c *APIClient) GetServiceClient(key ServiceType) (*gophercloud.ServiceClient, error) {
	c.mutex.RLock()
	if service, ok := c.services[key]; ok && service != nil {
		slog.Debug("returning existing service client", "type", string(key))
		c.mutex.RUnlock()
		return service, nil
	}
	c.mutex.RUnlock()
	return c.initServiceClient(key)
}

func (c *APIClient) initServiceClient(key ServiceType) (*gophercloud.ServiceClient, error) {

	// no existing service client, need to initialise one
	if _, ok := serviceConfigMap[key]; !ok {
		slog.Error("invalid service client type", "type", string(key))
		return nil, fmt.Errorf("invalid service type: %q", string(key))
	}

	slog.Debug("creating new service client", "type", string(key))

	client, err := serviceConfigMap[key].constructor(c.Client, gophercloud.EndpointOpts{Region: Region})

	if err != nil {
		slog.Error("error creating service client", "type", string(key), "error", err)
		return nil, err
	}
	client.Microversion = serviceConfigMap[key].microversion

	// save to object
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.services[key] = client

	slog.Info("new service client ready", "type", string(key))

	return client, nil
}
