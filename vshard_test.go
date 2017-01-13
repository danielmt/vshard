package vshard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type VShardTestSuite struct {
	suite.Suite
	Pool *Pool
}

func (suite *VShardTestSuite) SetupSuite() {
	servers := []string{
		"127.0.0.1:21210",
		"127.0.0.1:21211",
		"127.0.0.1:21212",
		"127.0.0.1:21213",
		"127.0.0.1:21214",
		"127.0.0.1:21215",
		"127.0.0.1:21216",
		"127.0.0.1:21217",
		"127.0.0.1:21218",
		"127.0.0.1:21219",
	}

	var err error
	suite.Pool, err = NewPool(servers, 10, 10, time.Second*5)
	if err != nil {
		suite.FailNow("Failure bringing pool up.", err)
	}
}

func (suite *VShardTestSuite) TearDownTest() {
	errs := suite.Pool.FlushAll()
	for _, err := range errs {
		if err != nil {
			suite.FailNow("Failure on FlushAll", err)
		}
	}
}

func (suite *VShardTestSuite) TestNumberOfServers() {
	suite.Equal(10, suite.Pool.numServers)
	suite.Len(suite.Pool.Servers, 10)
	suite.Len(suite.Pool.pool, 10)
}

func (suite *VShardTestSuite) TestStatus() {
	capacity := int64(10)
	available := int64(10)
	maxCap := int64(10)
	idleTimeout := time.Second * 5

	statusList := suite.Pool.Status()
	suite.Len(statusList, 10)

	for _, status := range statusList {
		suite.Equal(capacity, status.Capacity)
		suite.Equal(available, status.Available)
		suite.Equal(maxCap, status.MaxCap)
		suite.Equal(idleTimeout, status.IdleTimeout)
	}
}

func (suite *VShardTestSuite) testShardingDistribution(key, value string, poolNum int) {
	ok, err := suite.Pool.Set(key, 0, 0, []byte(value))

	suite.True(ok)
	suite.NoError(err)

	resource, err := suite.Pool.GetPoolConnection(poolNum)
	defer suite.Pool.ReturnConnection(poolNum, resource)
	if err != nil {
		suite.FailNow("Failure getting specific connection from pool", err)
	}

	result, err := resource.Get(key)
	if err != nil {
		suite.FailNow("Failure getting key", err)
	}

	suite.Equal(value, string(result[0].Value))
}

func (suite *VShardTestSuite) TestShardingDistribution() {
	suite.testShardingDistribution("f", "test-server-1", 0)
	suite.testShardingDistribution("o", "test-server-2", 1)
	suite.testShardingDistribution("d", "test-server-3", 2)
	suite.testShardingDistribution("e", "test-server-4", 3)
	suite.testShardingDistribution("b", "test-server-5", 4)
	suite.testShardingDistribution("g", "test-server-6", 5)
	suite.testShardingDistribution("p", "test-server-7", 6)
	suite.testShardingDistribution("c", "test-server-8", 7)
	suite.testShardingDistribution("l", "test-server-9", 8)
	suite.testShardingDistribution("a", "test-server-10", 9)
}

func TestVShardTestSuite(t *testing.T) {
	suite.Run(t, new(VShardTestSuite))
}
