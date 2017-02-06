package vshard

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"strconv"
	"sync"
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

// HashKeyStrategy defines the signature for the key hashing function
type HashKeyStrategy func(key string) string

// Close closes connections in a pool
func (r VitessResource) Close() {
	r.Connection.Close()
}

// Pool defines the pool
type Pool struct {
	Servers         []string
	ServerStrategy  ServerStrategy
	HashKeyStrategy HashKeyStrategy
	numServers      int
	pool            []*pools.ResourcePool
	sync.RWMutex
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
		Servers:         servers,
		numServers:      numServers,
		pool:            []*pools.ResourcePool{},
		ServerStrategy:  ShardedServerStrategyFarmhash,
		HashKeyStrategy: HashKeyStrategyFarmhash,
	}

	for i, server := range servers {
		func(_pool *Pool, _server string) {
			_pool.Lock()
			_pool.pool = append(_pool.pool, pools.NewResourcePool(func() (pools.Resource, error) {
				c, err := memcache.Connect(_server, time.Minute)
				return VitessResource{c}, err
			}, capacity, maxCap, idleTimeout))
			_pool.Unlock()

			conn, err := pool.GetPoolConnection(i)
			defer pool.ReturnConnection(i, conn)
			if err != nil {
				log.Fatalf("Can't connect to memcached: %s", err)
			}
		}(pool, server)
	}

	return pool, nil
}

// ShardedServerStrategyMD5 uses md5+jump to pick a server
func ShardedServerStrategyMD5(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	hash := md5.Sum([]byte(key))
	hashInt := big.NewInt(0)
	hashInt.SetString(hex.EncodeToString(hash[:]), 16)

	return int(jump.Hash(hashInt.Uint64(), numServers))
}

// ShardedServerStrategyFarmhash uses farmhash+jump to pick a server
func ShardedServerStrategyFarmhash(key string, numServers int) int {
	if numServers == 1 {
		return 0
	}

	return int(jump.Hash(farm.Fingerprint64([]byte(key)), numServers))
}

// HashKeyStrategyFarmhash uses farmhash to normalize key names for storage
func HashKeyStrategyFarmhash(key string) string {
	return strconv.FormatUint(farm.Fingerprint64([]byte(key)), 10)
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

	v.RLock()
	resource, err := v.pool[poolNum].Get(ctx)
	v.RUnlock()
	if err != nil {
		return nil, err
	}

	connection := resource.(VitessResource)

	return &connection, nil
}

// ReturnConnection returns a connection to the pool
func (v *Pool) ReturnConnection(poolNum int, resource *VitessResource) {
	if poolNum > v.numServers || poolNum < 0 {
		log.Fatalf("error: invalid server %d (of total %d)", poolNum, v.numServers)
	}

	v.RLock()
	v.pool[poolNum].Put(*resource)
	v.RUnlock()
}

// GetKeyMapping returns a mapping of server to a list of keys, useful for Gets()
func (v *Pool) GetKeyMapping(keys ...string) map[int][]string {
	mapping := make(map[int][]string)

	for i := 0; i < v.numServers; i++ {
		mapping[i] = []string{}
	}

	for _, key := range keys {
		poolNum := v.ServerStrategy(key, v.numServers)
		mapping[poolNum] = append(mapping[poolNum], v.HashKeyStrategy(key))
	}

	return mapping
}
