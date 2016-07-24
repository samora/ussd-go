package sessionstores

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/appengine/aetest"
)

func TestMemcacheGAE(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	store := NewMemcacheGAE(ctx)
	suite.Run(t, NewStoreSuite(store))
}
