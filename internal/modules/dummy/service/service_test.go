package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"{{index .App "git"}}/internal/modules/dummy"
)

func Test_Create(t *testing.T) {
	ctx := context.Background()
	expected := &dummy.Dummy{
		ID:   100,
		Name: "Some name",
	}
	service := New()

	actual, err := service.Create(ctx, expected.Name)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestService_List(t *testing.T) {
	ctx := context.Background()
	expected := []dummy.Dummy{
		{
			ID:   1,
			Name: "First Object",
		},
		{
			ID:   2,
			Name: "Second Object",
		},
	}
	service := New()

	actual, err := service.List(ctx, dummy.ListRequest{})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
