// package upload

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/brianvoe/gofakeit/v6"
// 	authService "github.com/igortoigildin/goph-keeper/internal/client/grpc/service/auth"
// 	descAuth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
// 	"github.com/igortoigildin/goph-keeper/pkg/session"
// 	descUpload "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func TestClientService_SendPassword(t *testing.T) {
// 	conn, err := grpc.NewClient(fmt.Sprintf(":9000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	type args struct {
// 		loginStr string
// 		passStr  string
// 		id       string
// 	}
// 	tests := []struct {
// 		name         string
// 		clientAuth   descAuth.AuthV1Client
// 		clientUpload descUpload.UploadV1Client
// 		args         args
// 		wantErr      bool
// 	}{
// 		{
// 			name:         "happy_test",
// 			clientAuth:   descAuth.NewAuthV1Client(conn),
// 			clientUpload: descUpload.NewUploadV1Client(conn),
// 			args: args{
// 				loginStr: gofakeit.Email(),
// 				passStr:  "fake_pass",
// 				id:       gofakeit.UUID(),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			auth := authService.New(fmt.Sprintf(":9000"))

// 			err := auth.RegisterNewUser(context.Background(), tt.args.loginStr, tt.args.passStr)
// 			assert.NoError(t, err)

// 			s := New()

// 			var token string
// 			if token, err = auth.Login(context.Background(), tt.args.loginStr, tt.args.passStr); (err != nil) != tt.wantErr {
// 				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			assert.NotEmpty(t, token)

// 			sessionData := &session.Session{
// 				Login:     tt.args.loginStr,
// 				Token:     token,
// 				ExpiresAt: time.Now().Add(time.Hour),
// 			}

// 			_ = session.SaveSession(sessionData)

// 			if err := s.SendPassword(fmt.Sprintf(":9000"), tt.args.loginStr, tt.args.passStr, tt.args.id); (err != nil) != tt.wantErr {
// 				t.Errorf("ClientService.SendPassword() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestClientService_SendBankDetails(t *testing.T) {
// 	conn, err := grpc.NewClient(fmt.Sprintf(":9000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	type args struct {
// 		addr       string
// 		cardNumber string
// 		cvc        string
// 		expDate    string
// 		loginStr   string
// 		passStr    string
// 		id         string
// 	}
// 	tests := []struct {
// 		name         string
// 		clientAuth   descAuth.AuthV1Client
// 		clientUpload descUpload.UploadV1Client
// 		args         args
// 		wantErr      bool
// 	}{
// 		{
// 			name:         "happy_test",
// 			clientAuth:   descAuth.NewAuthV1Client(conn),
// 			clientUpload: descUpload.NewUploadV1Client(conn),
// 			args: args{
// 				cardNumber: gofakeit.CreditCardNumber(nil),
// 				cvc:        gofakeit.CreditCardCvv(),
// 				expDate:    gofakeit.CreditCardExp(),
// 				id:         gofakeit.UUID(),
// 				loginStr:   gofakeit.Email(),
// 				passStr:    "fake_pass",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			auth := authService.New(fmt.Sprintf(":9000"))

// 			err := auth.RegisterNewUser(context.Background(), tt.args.loginStr, tt.args.passStr)
// 			assert.NoError(t, err)

// 			s := New()

// 			var token string
// 			if token, err = auth.Login(context.Background(), tt.args.loginStr, tt.args.passStr); (err != nil) != tt.wantErr {
// 				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			assert.NotEmpty(t, token)

// 			sessionData := &session.Session{
// 				Login:     tt.args.loginStr,
// 				Token:     token,
// 				ExpiresAt: time.Now().Add(time.Hour),
// 			}

// 			_ = session.SaveSession(sessionData)

// 			if err := s.SendBankDetails(fmt.Sprintf(":9000"), tt.args.cardNumber, tt.args.cvc, tt.args.expDate, tt.args.id); (err != nil) != tt.wantErr {
// 				t.Errorf("ClientService.SendBankDetails() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestClientService_SendText(t *testing.T) {
// 	conn, err := grpc.NewClient(fmt.Sprintf(":9000"), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	type args struct {
// 		text     string
// 		id       string
// 		loginStr string
// 		passStr  string
// 	}
// 	tests := []struct {
// 		name         string
// 		clientAuth   descAuth.AuthV1Client
// 		clientUpload descUpload.UploadV1Client
// 		args         args
// 		wantErr      bool
// 	}{
// 		{
// 			name:         "happy_test",
// 			clientAuth:   descAuth.NewAuthV1Client(conn),
// 			clientUpload: descUpload.NewUploadV1Client(conn),
// 			args: args{
// 				text:     gofakeit.AdverbPlace(),
// 				id:       gofakeit.UUID(),
// 				loginStr: gofakeit.Email(),
// 				passStr:  "fake_pass",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			auth := authService.New(fmt.Sprintf(":9000"))

// 			err := auth.RegisterNewUser(context.Background(), tt.args.loginStr, tt.args.passStr)
// 			assert.NoError(t, err)

// 			s := New()

// 			var token string
// 			if token, err = auth.Login(context.Background(), tt.args.loginStr, tt.args.passStr); (err != nil) != tt.wantErr {
// 				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			assert.NotEmpty(t, token)

// 			sessionData := &session.Session{
// 				Login:     tt.args.loginStr,
// 				Token:     token,
// 				ExpiresAt: time.Now().Add(time.Hour),
// 			}

// 			_ = session.SaveSession(sessionData)

//				if err := s.SendText(fmt.Sprintf(":9000"), tt.args.text, tt.args.id); (err != nil) != tt.wantErr {
//					t.Errorf("ClientService.SendText() error = %v, wantErr %v", err, tt.wantErr)
//				}
//			})
//		}
//	}
package upload
