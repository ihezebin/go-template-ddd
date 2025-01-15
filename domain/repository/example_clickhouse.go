package repository

import (
	"context"
	"fmt"

	"github.com/ihezebin/go-template-ddd/component/storage"
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type exampleClickhouseRepository struct{}

func NewExampleClickhouseRepository() ExampleRepository {
	return &exampleClickhouseRepository{}
}

var _ ExampleRepository = (*exampleClickhouseRepository)(nil)

func (e *exampleClickhouseRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	query := `
		INSERT INTO examples (id, username, password, email, salt)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := storage.ClickhouseStorageConn().ExecContext(ctx, query,
		example.Id,
		example.Username,
		example.Password,
		example.Email,
		example.Salt,
	)
	return err
}

func (e *exampleClickhouseRepository) FindByUsername(ctx context.Context, username string) (*entity.Example, error) {
	query := `
		SELECT id, username, password, email, salt
		FROM examples 
		WHERE username = ?
		LIMIT 1
	`

	example := &entity.Example{}
	err := storage.ClickhouseStorageConn().QueryRowContext(ctx, query, username).Scan(
		&example.Id,
		&example.Username,
		&example.Password,
		&example.Email,
		&example.Salt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan example error: %v", err)
	}

	return example, nil
}

func (e *exampleClickhouseRepository) FindByEmail(ctx context.Context, email string) (*entity.Example, error) {
	query := `
		SELECT id, username, password, email, salt
		FROM examples 
		WHERE email = ?
		LIMIT 1
	`

	example := &entity.Example{}
	err := storage.ClickhouseStorageConn().QueryRowContext(ctx, query, email).Scan(
		&example.Id,
		&example.Username,
		&example.Password,
		&example.Email,
		&example.Salt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan example error: %v", err)
	}

	return example, nil
}
