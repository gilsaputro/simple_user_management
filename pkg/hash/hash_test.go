package hash

import (
	"reflect"
	"testing"
)

func TestNewHashMethod(t *testing.T) {
	type args struct {
		cost int
	}
	tests := []struct {
		name string
		args args
		want HashMethod
	}{
		{
			name: "success flow",
			args: args{
				cost: 10,
			},
			want: &HashConfig{
				cost: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHashMethod(tt.args.cost); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHashMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashConfig_HashValue(t *testing.T) {
	type args struct {
		password string
		check    string
	}
	tests := []struct {
		name      string
		h         *HashConfig
		args      args
		sameValue bool
		wantErr   bool
	}{
		{
			name: "success flow",
			h: &HashConfig{
				cost: 10,
			},
			args: args{
				password: "password",
				check:    "password",
			},
			sameValue: true,
			wantErr:   false,
		},
		{
			name: "error compare flow",
			h: &HashConfig{
				cost: 10,
			},
			args: args{
				password: "password",
				check:    "password#1",
			},
			sameValue: false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.HashValue(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashConfig.HashValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			validation := tt.h.CompareValue(string(got), tt.args.check)
			if validation != tt.sameValue {
				t.Errorf("HashConfig.CompareValue() validation = %v, sameValue %v", validation, tt.sameValue)
				return
			}

		})
	}
}
