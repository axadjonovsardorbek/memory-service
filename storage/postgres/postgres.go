package postgres

import (
	"database/sql"
	"fmt"

	"memory/config"
	"memory/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db        *sql.DB
	MemoryS    storage.MemoriesI
	MediaS storage.MediasI
	SharedMemoryS    storage.SharedMemoriesI
	CommentS storage.CommentsI
}

func NewPostgresStorage(config config.Config) (*Storage, error) {
	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.DB_HOST, config.DB_USER, config.DB_NAME, config.DB_PASSWORD, config.DB_PORT)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	memory := NewMemoriesRepo(db)
	media := NewMediasRepo(db)
	comment := NewCommentsRepo(db)
	shared := NewSharedMemoriesRepo(db)

	return &Storage{
		Db:        db,
		MemoryS:    memory,
		MediaS: media,
		SharedMemoryS:    shared,
		CommentS: comment,
	}, nil
}
