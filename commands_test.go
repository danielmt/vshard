package vshard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type VShardCommandsTestSuite struct {
	suite.Suite
	Pool *Pool
}

func (suite *VShardCommandsTestSuite) SetupSuite() {
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

func (suite *VShardCommandsTestSuite) TearDownTest() {
	errs := suite.Pool.FlushAll()
	for _, err := range errs {
		if err != nil {
			suite.FailNow("Failure on FlushAll", err)
		}
	}
}

func (suite *VShardCommandsTestSuite) TestSetGet() {
	key := "x"
	expected := "hello-test"
	ok, err := suite.Pool.Set(key, 0, 0, []byte(expected))
	suite.True(ok)
	suite.NoError(err)

	value, err := suite.Pool.Get(key)
	suite.NoError(err)
	suite.Equal(expected, string(value))
}

func (suite *VShardCommandsTestSuite) TestGetInexistentKey() {
	key := "key-do-not-exist"
	value, err := suite.Pool.Get(key)
	suite.Equal(ErrKeyNotFound, err)
	suite.Empty(value)
}

func TestVShardCommandsTestSuite(t *testing.T) {
	suite.Run(t, new(VShardCommandsTestSuite))
}
