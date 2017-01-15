package vshard

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	farm "github.com/dgryski/go-farm"
	jump "github.com/dgryski/go-jump"
	"github.com/youtube/vitess/go/memcache"
	"github.com/youtube/vitess/go/pools"
)

var (
	// ErrKeyNotFound defines the error mensage when key is not found on memcached
	ErrKeyNotFound = errors.New("error: key not found")
)

// VitessResource implements the expected interface for vitess internal pool
type VitessResource struct {
	*memcache.Connection
}

// ServerStrategy defines the signature for the sharding function
type ServerStrategy func(key string, numServers int) int

// Close closes connections in a pool
func (r VitessResource) Close() {
	r.Connection.Close()
}

// Pool defines the pool
type Pool struct {
	Servers        []string
	ServerStrategy ServerStrategy
	numServers     int
	pool           []*pools.ResourcePool
}

// PoolStats defines all stats vitess memcached driver exposes
type PoolStats struct {
	Slot        int
	Server      string
	Capacity    int64
	Available   int64
	MaxCap      int64
	WaitCount   int64
	WaitTime    time.Duration
	IdleTimeout time.Duration
}

// NewPool returns a new VitessPool
func NewPool(servers []string, capacity, maxCap int, idleTimeout time.Duration) (*Pool, error) {
	numServers := len(servers)

	pool := &Pool{
		Servers:        servers,
		numServers:     numServers,
		pool:           make([]*pools.ResourcePool, numServers),
		ServerStrategy: ShardedServerStrategyMD5,
	}

	for i, server := range servers {
		func(_server string) {
			pool.pool[i] = pools.NewResourcePool(func() (pools.Resource, error) {
				c, err := memcache.Connect(_server, time.Minute)
				return VitessResource{c}, err
			}, capacity, maxCap, idleTimeout)
		}(server)
	}

	return pool, nil
}

// ShardedServerStrategyMD5 uses md5+jump to pick a server
func ShardedServerStrategyMD5(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	hash := md5.Sum([]byte(key))
	hashHex := hex.EncodeToString(hash[:])

	hashInt := big.NewInt(0)
	hashInt.SetString(hashHex, 16)

	server := int(jump.Hash(hashInt.Uint64(), numServers))

	return server
}

// ShardedServerStrategyFarmhash uses farmhash+jump to pick a server
func ShardedServerStrategyFarmhash(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	return int(jump.Hash(farm.Fingerprint64([]byte(key)), numServers))
}

// GetConnection returns a connection from the sharding pool, based on the key
func (v *Pool) GetConnection(key string) (*VitessResource, int, error) {
	poolNum := v.ServerStrategy(key, v.numServers)

	connection, err := v.GetPoolConnection(poolNum)
	if err != nil {
		return nil, -1, err
	}

	return connection, poolNum, nil
}

// GetPoolConnection returns a connection from a specific pool number
func (v *Pool) GetPoolConnection(poolNum int) (*VitessResource, error) {
	ctx := context.Background()

	resource, err := v.pool[poolNum].Get(ctx)
	if err != nil {
		return nil, err
	}

	connection := resource.(VitessResource)

	return &connection, nil
}

// ReturnConnection returns a connection to the pool
func (v *Pool) ReturnConnection(poolNum int, resource *VitessResource) {
	v.pool[poolNum].Put(*resource)
}

// GetKeyMapping returns a mapping of server to a list of keys, useful for Gets()
func (v *Pool) GetKeyMapping(keys ...string) map[int][]string {
	mapping := make(map[int][]string)

	for i := 0; i < v.numServers; i++ {
		mapping[i] = []string{}
	}

	for _, key := range keys {
		poolNum := v.ServerStrategy(key, v.numServers)
		mapping[poolNum] = append(mapping[poolNum], key)
	}

	return mapping
}

func (v *Pool) GetKeyMapping2(keys ...string) map[int][]string {
	mapping := make(map[int][]string)

	for _, key := range keys {
		poolNum := v.ServerStrategy(key, v.numServers)
		if _, ok := mapping[poolNum]; !ok {
			mapping[poolNum] = []string{}
		}
		mapping[poolNum] = append(mapping[poolNum], key)
	}

	return mapping
}
