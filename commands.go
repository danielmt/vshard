package vshard

import "github.com/youtube/vitess/go/cacheservice"

// Status returns all statistics exposed by the memcached driver
func (v *Pool) Status() []*PoolStats {
	stats := make([]*PoolStats, v.numServers)

	for i, pool := range v.pool {
		capacity, available, maxCap, waitCount, waitTime, idleTimeout := pool.Stats()
		status := &PoolStats{
			Slot:        i,
			Server:      v.Servers[i],
			Capacity:    capacity,
			Available:   available,
			MaxCap:      maxCap,
			WaitCount:   waitCount,
			WaitTime:    waitTime,
			IdleTimeout: idleTimeout,
		}

		stats[i] = status
	}

	return stats
}

// Get returns a key from the memcached server
func (v *Pool) Get(key string) ([]byte, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return nil, err
	}

	result, err := resource.Get(key)
	if err != nil {
		return nil, err
	}

	if len(result) < 1 {
		return nil, ErrKeyNotFound
	}

	return result[0].Value, nil
}

// Gets returns cached data for given keys, it is an alternative Get api
// for using with CAS. Gets returns a CAS identifier with the item. If
// the item's CAS value has changed since you Gets'ed it, it will not be stored.
func (v *Pool) Gets(keys ...string) ([]cacheservice.Result, error) {
	mapping := v.GetKeyMapping(keys...)
	results := []cacheservice.Result{}

	for poolNum, keys := range mapping {
		if len(keys) > 0 {
			connection, err := v.GetPoolConnection(poolNum)
			defer v.ReturnConnection(poolNum, connection)
			if err != nil {
				return nil, err
			}

			result, err := connection.Gets(keys...)
			if err != nil {
				return nil, err
			}

			results = append(results, result...)
		}
	}

	return results, nil
}

// Set set the value with specified cache key.
func (v *Pool) Set(key string, flags uint16, timeout uint64, value []byte) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Set(key, flags, timeout, value)
}

// Add store the value only if it does not already exist.
func (v *Pool) Add(key string, flags uint16, timeout uint64, value []byte) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Add(key, flags, timeout, value)
}

// Replace replaces the value, only if the value already exists,
// for the specified cache key.
func (v *Pool) Replace(key string, flags uint16, timeout uint64, value []byte) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Replace(key, flags, timeout, value)
}

// Append appends the value after the last bytes in an existing item.
func (v *Pool) Append(key string, flags uint16, timeout uint64, value []byte) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Append(key, flags, timeout, value)
}

// Prepend prepends the value before existing value.
func (v *Pool) Prepend(key string, flags uint16, timeout uint64, value []byte) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Prepend(key, flags, timeout, value)
}

// Cas stores the value only if no one else has updated the data since you read it last.
func (v *Pool) Cas(key string, flags uint16, timeout uint64, value []byte, cas uint64) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Cas(key, flags, timeout, value, cas)
}

// Delete delete the value for the specified cache key.
func (v *Pool) Delete(key string) (bool, error) {
	resource, poolNum, err := v.GetConnection(key)
	defer v.ReturnConnection(poolNum, resource)
	if err != nil {
		return false, err
	}

	return resource.Delete(key)
}

// FlushAll purges the entire cache on all servers.
func (v *Pool) FlushAll() []error {
	errs := []error{}

	for poolNum := range v.pool {
		resource, err := v.GetPoolConnection(poolNum)
		defer v.ReturnConnection(poolNum, resource)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = resource.FlushAll()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
