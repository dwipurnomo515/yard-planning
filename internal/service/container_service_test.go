package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerService_ValidateContainerSpec(t *testing.T) {
	service := &ContainerService{}

	tests := []struct {
		name          string
		size          int
		height        float64
		containerType string
		wantErr       bool
	}{
		{
			name:          "valid 20ft DRY",
			size:          20,
			height:        8.6,
			containerType: "DRY",
			wantErr:       false,
		},
		{
			name:          "valid 40ft REEFER",
			size:          40,
			height:        9.6,
			containerType: "REEFER",
			wantErr:       false,
		},
		{
			name:          "invalid size",
			size:          30,
			height:        8.6,
			containerType: "DRY",
			wantErr:       true,
		},
		{
			name:          "invalid height",
			size:          20,
			height:        10.0,
			containerType: "DRY",
			wantErr:       true,
		},
		{
			name:          "invalid type",
			size:          20,
			height:        8.6,
			containerType: "INVALID",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateContainerSpec(tt.size, tt.height, tt.containerType)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
