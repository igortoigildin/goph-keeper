package tests

import (
	"context"
	"strconv"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/pkg/download_v1"
	"github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"github.com/igortoigildin/goph-keeper/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestSaveText_Happy(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	resp, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := resp.GetToken()

	text := gofakeit.Adverb()
	id := gofakeit.Digit()

	md := metadata.Pairs("login", login, "id", id, "authorization", "Bearer "+token)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = st.UploadClient.UploadText(ctx, &upload_v1.UploadTextRequest{
		Text: text,
	})
	require.NoError(t, err)
}

func TestDownloadText_Happy(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	resp, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := resp.GetToken()

	text := gofakeit.Adverb()
	id := gofakeit.Number(23, 1948)
	idUpd := strconv.Itoa(id)

	md := metadata.Pairs("login", login, "id", idUpd, "authorization", "Bearer "+token)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = st.UploadClient.UploadText(ctx, &upload_v1.UploadTextRequest{
		Text: text,
	})
	require.NoError(t, err)

	_, err = st.DownloadClient.DownloadText(ctx, &download_v1.DownloadTextRequest{
		Uuid: idUpd,
	})

	md = metadata.Pairs("login", login, "id", idUpd)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	require.NoError(t, err)
}

func TestSaveText_Empty_Login(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	_, err = st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	text := gofakeit.Adverb()
	id := gofakeit.Digit()

	md := metadata.Pairs("login", "", "id", id)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = st.UploadClient.UploadText(ctx, &upload_v1.UploadTextRequest{
		Text: text,
	})
	require.Error(t, err)
}

func TestSaveText_Empty_Id(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	_, err = st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	text := gofakeit.Adverb()
	id := gofakeit.Digit()

	md := metadata.Pairs("login", "", "id", id)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	_, err = st.UploadClient.UploadText(ctx, &upload_v1.UploadTextRequest{
		Text: text,
	})
	require.Error(t, err)
}

func TestSaveBankDetails_Happy(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()
	id := strconv.Itoa(gofakeit.Number(100, 1000))

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	resp, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := resp.GetToken()

	md := metadata.Pairs("login", login, "id", id, "authorization", "Bearer "+token)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	data := make(map[string]string, 3)
	data["card_number"] = gofakeit.CreditCardNumber(nil)
	data["CVC"] = gofakeit.CreditCardCvv()
	data["expiration_date"] = gofakeit.CreditCardExp()

	_, err = st.UploadClient.UploadBankData(ctx, &upload_v1.UploadBankDataRequest{
		Data: data,
	})
	require.NoError(t, err)
}

func TestDownloadBankDetails_Happy(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()
	id := strconv.Itoa(gofakeit.Number(100, 1000))

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	resp, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := resp.GetToken()

	md := metadata.Pairs("login", login, "id", id, "authorization", "Bearer "+token)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	data := make(map[string]string, 3)
	data["card_number"] = gofakeit.CreditCardNumber(nil)
	data["CVC"] = gofakeit.CreditCardCvv()
	data["expiration_date"] = gofakeit.CreditCardExp()

	_, err = st.UploadClient.UploadBankData(ctx, &upload_v1.UploadBankDataRequest{
		Data: data,
	})
	require.NoError(t, err)

	_, err = st.DownloadClient.DownloadBankData(ctx, &download_v1.DownloadBankDataRequest{
		Uuid: id,
	})
	require.NoError(t, err)
}

func TestSaveBankDetails_Emty_Login(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()
	id := strconv.Itoa(gofakeit.Number(100, 1000))

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	_, err = st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	md := metadata.Pairs("login", "", "id", id)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	data := make(map[string]string, 3)
	data["card_number"] = gofakeit.CreditCardNumber(nil)
	data["CVC"] = gofakeit.CreditCardCvv()
	data["expiration_date"] = gofakeit.CreditCardExp()

	_, err = st.UploadClient.UploadBankData(ctx, &upload_v1.UploadBankDataRequest{
		Data: data,
	})
	require.Error(t, err)
}

func TestSaveBankDetails_Emty_Bank_details(t *testing.T) {
	ctx, st := suite.New(t)
	login := gofakeit.Email()
	pass := randomFakePassword()
	id := strconv.Itoa(gofakeit.Number(100, 1000))

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	_, err = st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	md := metadata.Pairs("login", login, "id", id)

	ctx = metadata.NewOutgoingContext(context.Background(), md)

	var data map[string]string

	_, err = st.UploadClient.UploadBankData(ctx, &upload_v1.UploadBankDataRequest{
		Data: data,
	})
	require.Error(t, err)
}
