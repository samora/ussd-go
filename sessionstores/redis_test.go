package sessionstores

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRedis(t *testing.T) {
	store := NewRedis("localhost:6379")
	suite.Run(t, NewStoreSuite(store))
}
