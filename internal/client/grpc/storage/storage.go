package storage

import "github.com/igortoigildin/goph-keeper/internal/client/grpc/models"

type Saver interface {
	SaveText(id, info, text, etag string) error
	SaveCredentials(id, service, username, password, etag string) error
	SaveBankDetails(cardNumber, cvc, expDate, id, bankName, etag string) error
	SaveFile(id, filePath, etag string) error
}

type Downloader interface {
	GetText(id string) (models.Text, error)
	GetCredential(id string) (models.Credential, error)
	GetBankDetails(id string) (models.BankDetails, error)
	GetFile(id string) (models.File, error)
}

type SecretsReceiver interface {
	GetAllTexts() ([]models.Text, error)
	GetAllCredentials() ([]models.Credential, error)
	GetAllBankDetails() ([]models.BankDetails, error)
	ListAllFiles() ([]models.File, error)
}
