package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"{{if index .Modules "postgres"}}
	"go.uber.org/mock/gomock"{{end}}

	"{{index .App "git"}}/internal/modules/dummy"{{if index .Modules "postgres"}}
	"{{index .App "git"}}/internal/modules/dummy/dummymocks"{{end}}
)

func Test_Create(t *testing.T) {
	ctx := context.Background(){{if index .Modules "postgres"}}
	ctrl := gomock.NewController(t)
	mockRepo := dummymocks.NewMockRepository(ctrl){{end}}
	service := New({{if index .Modules "postgres"}}
		mockRepo,{{end}}
	)

	expected := &dummy.Dummy{
		ID:   100,
		Name: "Some name",
	}
{{if index .Modules "postgres"}}
	mockRepo.EXPECT().Create(ctx, expected.Name).Return(expected, nil){{end}}
	actual, err := service.Create(ctx, expected.Name)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestService_List(t *testing.T) {
	ctx := context.Background(){{if index .Modules "postgres"}}
	ctrl := gomock.NewController(t)
	mockRepo := dummymocks.NewMockRepository(ctrl){{end}}
	service := New({{if index .Modules "postgres"}}
		mockRepo,{{end}}
	)

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
{{if index .Modules "postgres"}}
	mockRepo.EXPECT().List(ctx, dummy.ListRequest{}).Return(expected, nil){{end}}
	actual, err := service.List(ctx, dummy.ListRequest{})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
