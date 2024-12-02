package tenants

import (
	"context"
	"go_graphql/config"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"os"
	"reflect"
	"testing"
)

func TestTenantMutationResolver_UpdateTenant(t *testing.T) {
	os.Setenv("PERMIT_PROJECT", "dev")
	os.Setenv("PERMIT_ENV", "dev")
	os.Setenv("PERMIT_TOKEN", "permit_key_UeJJBnL50MovEqo4GZKbiQL0x4wKZY72hMO8nQmpR1z4gX9TOHk5STep5da7KdjJaa0utnczeLeNBqABvU02I9")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://api.permit.io")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "Harani@8500")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "iam_services")

	db := config.InitDB()
	type args struct {
		ctx   context.Context
		id    string
		input models.TenantInput
	}
	tests := []struct {
		name    string
		r       *TenantMutationResolver
		args    args
		want    *dto.Tenant
		wantErr bool
	}{
		{"test", &TenantMutationResolver{DB: db}, args{context.TODO(), "1", models.TenantInput{Name: "test"}}, nil, false},
		{"test2", &TenantMutationResolver{DB: db}, args{context.TODO(), "2", models.TenantInput{Name: "test2"}}, nil, false},
		{"test3", &TenantMutationResolver{DB: db}, args{context.TODO(), "3", models.TenantInput{Name: "test3"}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.UpdateTenant(tt.args.ctx, tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TenantMutationResolver.UpdateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(reflect.TypeOf(got))
		})
	}
}

func TestTenantMutationResolver_CreateTenant(t *testing.T) {
	// os.Setenv("PERMIT_PROJECT", "test_project")
	// os.Setenv("PERMIT_ENV", "test_env")
	// os.Setenv("PERMIT_TOKEN", "test_token")
	// os.Setenv("PERMIT_PDP_ENDPOINT", "http://test.endpoint")
	//db := config.InitDB()
	os.Setenv("PERMIT_PROJECT", "dev")
	os.Setenv("PERMIT_ENV", "dev")
	os.Setenv("PERMIT_TOKEN", "permit_key_UeJJBnL50MovEqo4GZKbiQL0x4wKZY72hMO8nQmpR1z4gX9TOHk5STep5da7KdjJaa0utnczeLeNBqABvU02I9")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://api.permit.io")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "Harani@8500")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "iam_services")

	db := config.InitDB()

	type args struct {
		ctx   context.Context
		input models.TenantInput
	}
	tests := []struct {
		name    string
		r       *TenantMutationResolver
		args    args
		want    *dto.Tenant
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			r: &TenantMutationResolver{
				DB: db,
			},
			args: args{
				ctx: context.Background(),
				input: models.TenantInput{
					Name: "test",
				},
			},
			want: &dto.Tenant{
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "test1",
			r:    &TenantMutationResolver{},
			args: args{
				ctx: context.Background(),
				input: models.TenantInput{
					Name: "test1",
				},
			},
			want: &dto.Tenant{
				Name: "test1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.CreateTenant(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TenantMutationResolver.CreateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TenantMutationResolver.CreateTenant() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestTenantMutationResolver_DeleteTenant(t *testing.T) {
	os.Setenv("PERMIT_PROJECT", "dev")
	os.Setenv("PERMIT_ENV", "dev")
	os.Setenv("PERMIT_TOKEN", "permit_key_UeJJBnL50MovEqo4GZKbiQL0x4wKZY72hMO8nQmpR1z4gX9TOHk5STep5da7KdjJaa0utnczeLeNBqABvU02I9")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://api.permit.io")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "Harani@8500")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "iam_services")

	db := config.InitDB()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		r       *TenantMutationResolver
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &TenantMutationResolver{DB: db}, args{context.TODO(), "1"}, true, false},
		{"test1", &TenantMutationResolver{DB: db}, args{context.TODO(), "2"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.DeleteTenant(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TenantMutationResolver.DeleteTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TenantMutationResolver.DeleteTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}
