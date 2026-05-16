package v1_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"tracking/internal/ent"
	v1 "tracking/internal/modules/sites/v1"
	"tracking/internal/modules/sites/v1/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateSite(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name        string
		request     v1.CreateRequest
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid site creation",
			request: v1.CreateRequest{
				Name:   "Test Site",
				Domain: "example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid domain",
			request: v1.CreateRequest{
				Name:   "Test Site",
				Domain: "invalid domain with spaces",
			},
			wantErr:     true,
			expectedErr: v1.ErrInvalidDomain,
		},
		{
			name: "empty name",
			request: v1.CreateRequest{
				Name:   "",
				Domain: "example.com",
			},
			wantErr:     true,
			expectedErr: v1.ErrInvalidName,
		},
		{
			name: "name too long",
			request: v1.CreateRequest{
				Name:   string(make([]byte, 256)),
				Domain: "example.com",
			},
			wantErr:     true,
			expectedErr: v1.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockRepository(t)

			if !tt.wantErr {
				mockRepo.On("GetByDomain", mock.Anything, mock.Anything).
					Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).
					Return(&ent.Site{}, nil)
			}

			service := v1.NewService(mockRepo, logger)
			_, err := service.CreateSite(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
