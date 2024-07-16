package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	mp "memory/genproto"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type MediasRepo struct {
	db *sql.DB
}

func NewMediasRepo(db *sql.DB) *MediasRepo {
	return &MediasRepo{db: db}
}

func (m *MediasRepo) Create(req *mp.MediasCreateReq) (*mp.Void, error) {
	id := uuid.New().String()
	void := mp.Void{}

	query := `
	INSERT INTO medias (
		id,
		memory_id,
		type,
		url
	) VALUES ($1, $2, $3, $4)
	`

	_, err := m.db.Exec(query, id, req.MemoryId, req.Type, req.Url)

	if err != nil {
		log.Println("Error while creating media: ", err)
		return nil, err
	}

	log.Println("Successfully created media")

	return &void, nil
}

func (m *MediasRepo) GetById(id *mp.ById) (*mp.MediasGetByIdRes, error) {
	media := mp.MediasGetByIdRes{
		Media: &mp.MediasRes{},
	}

	query := `
	SELECT 
		id,
		type,
		url,
		date,
		memory_id
	FROM 
		medias
	WHERE 
		id = $1
	AND 
		deleted_at = 0	
	`

	row := m.db.QueryRow(query, id.Id)

	err := row.Scan(
		&media.Media.Id,
		&media.Media.Type,
		&media.Media.Url,
		&media.Media.Date,
		&media.Media.MemoryId,
	)

	if err != nil {
		log.Println("Error while getting media by id: ", err)
		return nil, err
	}

	log.Println("Successfully got media")

	return &media, nil
}

func (m *MediasRepo) GetAll(req *mp.MediasGetAllReq) (*mp.MediasGetAllRes, error) {
	medias := mp.MediasGetAllRes{}

	query := `
	SELECT 
		id,
		type,
		url,
		date,
		memory_id
	FROM 
		medias
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.MemoryId != "" && req.MemoryId != "string" {
		conditions = append(conditions, " memory_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.MemoryId)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var limit int32
	var offset int32
	var count int32

	limit = 10
	offset = req.Filter.Page * limit

	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := m.db.Query(query, args...)

	if err != nil {
		log.Println("Error while retriving medias: ", err)
		return nil, err
	}

	for rows.Next() {
		media := mp.MediasRes{}

		err := rows.Scan(
			&media.Id,
			&media.Type,
			&media.Url,
			&media.Date,
			&media.MemoryId,
		)

		if err != nil {
			log.Println("Error while scanning all medias: ", err)
			return nil, err
		}

		count += 1
		medias.Medias = append(medias.Medias, &media)
	}

	medias.Count = count
	log.Println("Successfully fetched all medias")

	return &medias, nil
}

func (m *MediasRepo) Update(req *mp.MediasUpdateReq) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE
		medias
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Type != "" && req.Type != "string" {
		conditions = append(conditions, " type = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Type)
	}
	if req.Url != "" && req.Url != "string" {
		conditions = append(conditions, " url = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Url)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0 "

	args = append(args, req.Id)

	_, err := m.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating media: ", err)
		return nil, err
	}

	log.Println("Successfully updated media")

	return &void, nil
}

func (m *MediasRepo) Delete(id *mp.ById) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE 
		medias
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := m.db.Exec(query, id.Id)

	if err != nil {
		log.Println("Error while deleting media: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("media with id %s not found", id.Id)
	}

	log.Println("Successfully deleted media")

	return &void, nil
}
