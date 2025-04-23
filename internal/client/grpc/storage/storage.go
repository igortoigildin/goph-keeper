package storage

import "github.com/igortoigildin/goph-keeper/internal/client/grpc/models"

type Saver interface {
	SaveText(id, info, text string) error
	SaveCredentials(id, service, username, password string) error
	SaveBankDetails(cardNumber, cvc, expDate string, id, bankName string) error
	SaveFile(id, filePath string) error
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
