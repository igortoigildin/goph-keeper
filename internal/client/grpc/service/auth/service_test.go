package register

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	desc "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAuthService_RegisterNewUser(t *testing.T) {
	conn, err := grpc.NewClient(fmt.Sprintf(":9000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	type args struct {
		ctx   context.Context
		login string
		pass  string
	}
	tests := []struct {
		name    string
		client  desc.AuthV1Client
		args    args
		wantErr bool
	}{
		{
			name:   "happy_test",
			client: desc.NewAuthV1Client(conn),
			args: args{
				ctx:   context.Background(),
				login: gofakeit.Email(),
				pass:  "temp_pass",
			},
			wantErr: false,
		},
		{
			name:   "fail_empty_login",
			client: desc.NewAuthV1Client(conn),
			args: args{
				ctx:   context.Background(),
				login: "",
				pass:  "temp_pass",
			},
			wantErr: true,
		},
		{
			name:   "fail_empty_pass",
			client: desc.NewAuthV1Client(conn),
			args: args{
				ctx:   context.Background(),
				login: gofakeit.Email(),
				pass:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := New(fmt.Sprintf(":9000"))

			if err := auth.RegisterNewUser(tt.args.ctx, tt.args.login, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	conn, err := grpc.NewClient(fmt.Sprintf(":9000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	type args struct {
		ctx   context.Context
		login string
		pass  string
	}
	tests := []struct {
		name    string
		client  desc.AuthV1Client
		args    args
		wantErr bool
	}{
		{
			name:   "happy_test",
			client: desc.NewAuthV1Client(conn),
			args: args{
				ctx:   context.Background(),
				login: gofakeit.Email(),
				pass:  "temp_pass",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := New(fmt.Sprintf(":9000"))

			err := auth.RegisterNewUser(context.Background(), tt.args.login, tt.args.pass)

			if _, err = auth.Login(tt.args.ctx, tt.args.login, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
