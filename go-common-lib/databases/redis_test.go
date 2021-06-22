package databases

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedis(t *testing.T) {
	tests := []struct {
		context          context.Context
		connectionString string
		err              error
	}{
		{context: context.Background(), connectionString: "127.0.0.1:6379", err: nil},
		{context: context.Background(), connectionString: "127.0.0.1:2379", err: errors.New("dial tcp 127.0.0.1:6379: connect: connection refused")},
	}
	for _, test := range tests {
		_, err := NewRedis(test.context, test.connectionString, "")
		if err != nil {
			assert.Error(t, test.err, err)
			continue
		}
		assert.IsType(t, err, nil)

	}
}
