package model

import "testing"

func TestProductPayloadValidate(t *testing.T) {
	tests := []struct {
		name    string
		payload ProductPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: ProductPayload{
				SKU:    "SKU-001",
				Name:   "Keyboard",
				Price:  100,
				Status: "active",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			payload: ProductPayload{
				SKU:    "SKU-001",
				Name:   "Keyboard",
				Price:  100,
				Status: "draft",
			},
			wantErr: true,
		},
		{
			name: "negative price",
			payload: ProductPayload{
				SKU:    "SKU-001",
				Name:   "Keyboard",
				Price:  -10,
				Status: "active",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.payload.Validate()
			if tt.wantErr && len(errs) == 0 {
				t.Fatalf("expected errors, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Fatalf("expected no errors, got %v", errs)
			}
		})
	}
}
