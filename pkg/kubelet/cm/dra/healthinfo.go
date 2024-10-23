/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dra

import (
	"sync"

	"k8s.io/kubernetes/pkg/kubelet/cm/dra/state"
)

// healthInfoCache is a cache of known device health.
type healthInfoCache struct {
	sync.RWMutex
	HealthInfo *state.DevicesHealthMap
}

// newClaimInfoCache creates a new claim info cache object, pre-populated from a checkpoint (if present).
func newHealthInfoCache() (*healthInfoCache, error) {
	cache := &healthInfoCache{
		HealthInfo: &state.DevicesHealthMap{},
	}

	return cache, nil
}

// withLock runs a function while holding the claimInfoCache lock.
func (cache *healthInfoCache) withLock(f func() error) error {
	cache.Lock()
	defer cache.Unlock()
	return f()
}

// withRLock runs a function while holding the claimInfoCache rlock.
func (cache *healthInfoCache) withRLock(f func() error) error {
	cache.RLock()
	defer cache.RUnlock()
	return f()
}

// getHealthInfo returns the current health info.
func (cache *healthInfoCache) getHealthInfo(driverName string, poolName string, deviceName string) state.DeviceHealthString {
	res := state.DeviceHealthString("unknown")

	cache.withRLock(func() error {
		if driver, ok := (*cache.HealthInfo)[driverName]; ok {

			for _, device := range driver.Devices {
				if device.PoolName == poolName && device.DeviceName == deviceName {
					res = device.Health
				}
			}

			// if we reached here - device health was never reported

		} else {
			// The driver does not implement the health API
		}
		return nil
	})

	return res
}

// updateHealthInfo updates the health info for a device. Returns true if the health info was updated.
func (cache *healthInfoCache) updateHealthInfo(resourceName, poolName, deviceName string, health state.DeviceHealthString) (state.DeviceHealth, bool) {
	var d state.DeviceHealth
	res := false
	cache.withLock(func() error {
		if _, ok := (*cache.HealthInfo)[resourceName]; !ok {
			(*cache.HealthInfo)[resourceName] = state.DriverHealthState{
				Devices: make([]state.DeviceHealth, 8),
			}
		}

		for i, dH := range (*cache.HealthInfo)[resourceName].Devices {
			if dH.PoolName == poolName && dH.DeviceName == deviceName {
				if dH.Health == health {
					// do nothing if the health is the same
					d = dH
					res = false
					return nil
				} else {
					(*cache.HealthInfo)[resourceName].Devices[i].Health = health
					d = (*cache.HealthInfo)[resourceName].Devices[i]
					res = true
					return nil
				}
			}
		}

		d = state.DeviceHealth{
			PoolName:   poolName,
			DeviceName: deviceName,
			Health:     health,
		}
		res = true
		driverHealthState := (*cache.HealthInfo)[resourceName]
		driverHealthState.Devices = append(driverHealthState.Devices, d)
		(*cache.HealthInfo)[resourceName] = driverHealthState

		return nil

	})
	return d, res
}
