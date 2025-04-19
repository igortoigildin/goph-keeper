package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	models "github.com/igortoigildin/goph-keeper/internal/client/grpc/models"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(path string) (*ClientRepository, error) {
	db, err := InitDB(path)
	if err != nil {
		return nil, err
	}

	c := ClientRepository{
		db: db,
	}

	return &c, nil
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS credentials (
		id TEXT PRIMARY KEY,
		service TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS texts (
		id TEXT PRIMARY KEY,
		info TEXT NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS bank_data (
		id TEXT PRIMARY KEY,
		bank_name TEXT,
		card_number TEXT,
		expiry TEXT,
		cvc TEXT,
		created_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		filename TEXT,
		data BLOB,
		updated_at DATETIME
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (rep *ClientRepository) SaveText(id, info, text string) error {
	_, err := rep.db.Exec(`
		INSERT INTO texts (id, info, text, created_at)
		VALUES (?, ?, ?, ?)`,
		id, info, text, time.Now())

	return err
}

func (rep *ClientRepository) GetAllTexts() ([]models.Text, error) {
	rows, err := rep.db.Query("SELECT id, info, text, created_at FROM texts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []models.Text
	for rows.Next() {
		var t models.Text

		err := rows.Scan(&t.ID, &t.Info, &t.Text, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		texts = append(texts, t)
	}

	return texts, nil
}

func (rep *ClientRepository) GetText(id string) (models.Text, error) {
	var t models.Text

	err := rep.db.QueryRow(`
		SELECT id, info, text, created_at
		FROM texts
		WHERE id = ?
	`, id).Scan(&t.ID, &t.Info, &t.Text, &t.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Text{}, fmt.Errorf("данные с id '%s' не найдены", id)
		}
		return models.Text{}, fmt.Errorf("ошибка при получении данных: %w", err)
	}

	return t, nil
}

func (rep *ClientRepository) SaveCredentials(id, service, username, password string) error {
	_, err := rep.db.Exec(`
		INSERT INTO credentials (id, service, username, password, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		id, service, username, password, time.Now())

	return err
}

func (rep *ClientRepository) GetAllCredentials() ([]models.Credential, error) {
	rows, err := rep.db.Query("SELECT id, service, username, password, created_at FROM credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []models.Credential
	for rows.Next() {
		var c models.Credential

		err := rows.Scan(&c.ID, &c.Service, &c.Username, &c.Password, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}

	return creds, nil
}

func (rep *ClientRepository) GetCredential(id string) (models.Credential, error) {
	var c models.Credential

	err := rep.db.QueryRow(`
		SELECT id, service, username, password, created_at
		FROM credentials
		WHERE id = ?
	`, id).Scan(&c.ID, &c.Service, &c.Username, &c.Password, &c.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Credential{}, fmt.Errorf("данные с id '%s' не найдены", id)
		}
		return models.Credential{}, fmt.Errorf("ошибка при получении данных: %w", err)
	}

	return c, nil
}

func (rep *ClientRepository) SaveBankDetails(cardNumber, cvc, expDate string, id, bankName string) error {
	_, err := rep.db.Exec(`
		INSERT INTO bank_data (id, bank_name, card_number, expiry, cvc, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		id, bankName, cardNumber, expDate, cvc, time.Now())

	return err
}

func (rep *ClientRepository) GetAllBankDetails() ([]models.BankDetails, error) {
	rows, err := rep.db.Query("SELECT id, service, username, password, created_at FROM bank_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.BankDetails
	for rows.Next() {
		var c models.BankDetails

		err := rows.Scan(&c.ID, &c.Info, &c.CardNumber, &c.Cvc, &c.ExpDate, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	return cards, nil
}

func (rep *ClientRepository) GetBankDetails(id string) (models.BankDetails, error) {
	var b models.BankDetails

	err := rep.db.QueryRow(`
		SELECT id, bank_name, card_number, expiry, cvc, created_at
		FROM bank_data
		WHERE id = ?
	`, id).Scan(&b.ID, &b.Info, &b.CardNumber, b.ExpDate, &b.Cvc, &b.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BankDetails{}, fmt.Errorf("данные с id '%s' не найдены", id)
		}
		return models.BankDetails{}, fmt.Errorf("ошибка при получении данных: %w", err)
	}

	return b, nil
}

func (rep *ClientRepository) SaveFile(id, filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("error while reading file", zap.Error(err))
		return err
	}

	f := models.File{
		ID:        id,
		Filename:  filePath,
		Data:      fileData,
		UpdatedAt: time.Now(),
	}

	_, err = rep.db.Exec("INSERT OR REPLACE INTO files (id, filename, data, updated_at) VALUES (?, ?, ?, ?)",
		f.ID, f.Filename, f.Data, f.UpdatedAt)
	return err
}

func (rep *ClientRepository) ListAllFiles(db *sql.DB) ([]models.File, error) {
	rows, err := db.Query("SELECT id, filename, data, updated_at FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		var ignored []byte
		err = rows.Scan(&f.ID, &f.Filename, &ignored, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (rep *ClientRepository) GetFile(id string) (models.File, error) {
	var f models.File

	err := rep.db.QueryRow(`
		SELECT id, filename, data, updated_at
		FROM files
		WHERE id = ?
	`, id).Scan(&f.ID, &f.Filename, &f.Data, &f.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.File{}, fmt.Errorf("файл с id '%s' не найден или уже синхронизирован", id)
		}
		return models.File{}, fmt.Errorf("ошибка при получении файла: %w", err)
	}
	return f, nil
}
