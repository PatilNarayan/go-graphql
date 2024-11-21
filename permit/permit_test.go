package permit

import (
	"reflect"
	"testing"

	"github.com/permitio/permit-golang/pkg/models"
)

func TestPermitClient_CreateTenant(t *testing.T) {
	type args struct {
		name        *string
		description *string
		attributes  map[string]interface{}
	}
	tests := []struct {
		name    string
		pc      *PermitClient
		args    args
		want    *models.TenantRead
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pc.CreateTenant(tt.args.name, tt.args.description, tt.args.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PermitClient.CreateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PermitClient.CreateTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermitClient_UpdateTenant(t *testing.T) {
	type args struct {
		tenantKey   string
		name        *string
		description *string
		attributes  map[string]interface{}
	}
	tests := []struct {
		name    string
		pc      *PermitClient
		args    args
		want    *models.TenantRead
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pc.UpdateTenant(tt.args.tenantKey, tt.args.name, tt.args.description, tt.args.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PermitClient.UpdateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PermitClient.UpdateTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermitClient_DeleteTenant(t *testing.T) {
	type args struct {
		tenantKey string
	}
	tests := []struct {
		name    string
		pc      *PermitClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pc.DeleteTenant(tt.args.tenantKey); (err != nil) != tt.wantErr {
				t.Errorf("PermitClient.DeleteTenant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
