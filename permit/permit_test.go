package permit

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/permitio/permit-golang/pkg/models"
)

func TestCreateTenant(t *testing.T) {
	os.Setenv("PERMIT_PROJECT", "dev")
	os.Setenv("PERMIT_ENV", "dev")
	os.Setenv("PERMIT_TOKEN", "permit_key_mXEiRTTrdC78tEWyjy2dazPGFruDrhJswiNthVQ0YYqjT7VWH1dn9Wiw3vtVFUl25tweXwICha30RsPe0XRwwU")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://api.permit.io")
	type args struct {
		tenantname string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.TenantRead
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{"test"}, nil, false},
		{"test2", args{"test2"}, nil, false},
		{"test3", args{""}, nil, false},
		{"test4", args{"$&^@*^()"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTenant(tt.args.tenantname)
			if (err != nil) != tt.wantErr {
				fmt.Println(err)
				t.Errorf("CreateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(reflect.TypeOf(got))
		})
	}
}

func TestDeleteTenant(t *testing.T) {
	type args struct {
		tenantid string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{"1"}, true, false},
		{"test2", args{"2"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteTenant(tt.args.tenantid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTenant(t *testing.T) {
	type args struct {
		tenantid   string
		tenantname string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.TenantRead
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{"1", "test"}, nil, false},
		{"test2", args{"2", "test2"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateTenant(tt.args.tenantid, tt.args.tenantname)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}
