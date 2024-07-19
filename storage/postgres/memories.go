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
	"github.com/lib/pq"
)

type MemoriesRepo struct {
	db *sql.DB
}

func NewMemoriesRepo(db *sql.DB) *MemoriesRepo {
	return &MemoriesRepo{db: db}
}

func (m *MemoriesRepo) Create(req *mp.MemoriesCreateReq) (*mp.Void, error) {
	id := uuid.New().String()
	void := mp.Void{}

	query := `
	INSERT INTO memories (
		id,
		user_id,
		title,
		description,
		date,
		tags,
		location,
		place_name,
		privacy
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := m.db.Exec(query, id, req.UserId, req.Title, req.Description, req.Date, pq.Array(req.Tags), req.Location, req.PlaceName, req.Privacy)

	if err != nil {
		log.Println("Error while creating memory: ", err)
		return nil, err
	}

	log.Println("Successfully created memory")

	return &void, nil
}

func (m *MemoriesRepo) GetById(id *mp.ById) (*mp.MemoriesGetByIdRes, error) {
	memory := mp.MemoriesGetByIdRes{
		Memory: &mp.MemoriesRes{},
	}

	query := `
	SELECT 
		id,
		title,
		description,
		date,
		tags,
		location,
		place_name,
		privacy,
		user_id
	FROM 
		memories
	WHERE 
		id = $1
	AND 
		deleted_at = 0	
	`

	row := m.db.QueryRow(query, id.Id)

	err := row.Scan(
		&memory.Memory.Id,
		&memory.Memory.Title,
		&memory.Memory.Description,
		&memory.Memory.Date,
		&memory.Memory.Tags,
		&memory.Memory.Location,
		&memory.Memory.PlaceName,
		&memory.Memory.Privacy,
		&memory.Memory.UserId,
	)

	if err != nil {
		log.Println("Error while getting memory by id: ", err)
		return nil, err
	}

	log.Println("Successfully got memory")

	return &memory, nil
}

func (m *MemoriesRepo) GetAll(req *mp.MemoriesGetAllReq) (*mp.MemoriesGetAllRes, error) {
	memories := mp.MemoriesGetAllRes{}

	query := `
	SELECT 
		id,
		title,
		description,
		date,
		tags,
		location,
		place_name,
		privacy,
		user_id
	FROM 
		memories
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.UserId != "" && req.UserId != "string" {
		conditions = append(conditions, " user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.UserId)
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
		log.Println("Error while retriving memories: ", err)
		return nil, err
	}

	for rows.Next() {
		memory := mp.MemoriesRes{}

		err := rows.Scan(
			&memory.Id,
			&memory.Title,
			&memory.Description,
			&memory.Date,
			&memory.Tags,
			&memory.Location,
			&memory.PlaceName,
			&memory.Privacy,
			&memory.UserId,
		)

		if err != nil {
			log.Println("Error while scanning all memories: ", err)
			return nil, err
		}

		count += 1
		memories.Memories = append(memories.Memories, &memory)
	}

	memories.Count = count
	log.Println("Successfully fetched all memories")

	return &memories, nil
}

func (m *MemoriesRepo) Update(req *mp.MemoriesUpdateReq) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE
		memories
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Title != "" && req.Title != "string" {
		conditions = append(conditions, " title = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Title)
	}
	if req.Description != "" && req.Description != "string" {
		conditions = append(conditions, " description = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Description)
	}
	if req.Privacy != "" && req.Privacy != "string" {
		conditions = append(conditions, " privacy = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Privacy)
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
		log.Println("Error while updating memory: ", err)
		return nil, err
	}

	log.Println("Successfully updated memory")

	return &void, nil
}

func (m *MemoriesRepo) Delete(id *mp.ById) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE 
		memories
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := m.db.Exec(query, id.Id)

	if err != nil {
		log.Println("Error while deleting memory: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("memory with id %s not found", id.Id)
	}

	log.Println("Successfully deleted memory")

	return &void, nil
}
