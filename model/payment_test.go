package model

import (
	"testing"
)

func TestPayment_Valid(t *testing.T) {
	type fields struct {
		Type           string
		ID             PaymentID
		Version        uint
		OrganisationID string
		Attributes     Attributes
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Fully empty",
			want: false,
		},
		{
			name: "Wrong Type",
			fields: fields{
				Type: "something",
			},
			want: false,
		},
		{
			name: "Correct entry",
			fields: fields{
				Type:           "Payment",
				ID:             PaymentID("343423423"),
				Version:        1,
				OrganisationID: "87847584385",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				Version:        tt.fields.Version,
				OrganisationID: tt.fields.OrganisationID,
				Attributes:     tt.fields.Attributes,
			}
			if got := p.Valid(); got != tt.want {
				t.Errorf("Payment.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
