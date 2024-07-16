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

type SharedMemoriesRepo struct {
	db *sql.DB
}

func NewSharedMemoriesRepo(db *sql.DB) *SharedMemoriesRepo {
	return &SharedMemoriesRepo{db: db}
}

func (m *SharedMemoriesRepo) Create(req *mp.SharedMemoriesCreateReq) (*mp.Void, error) {
	id := uuid.New().String()
	void := mp.Void{}

	query := `
	INSERT INTO shared_memories (
		id,
		memory_id,
		shared_id,
		recipient_id,
		message,
	) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := m.db.Exec(query, id, req.MemoryId, req.SharedId, req.RecipientId, req.Message)

	if err != nil {
		log.Println("Error while creating shared_memory: ", err)
		return nil, err
	}

	log.Println("Successfully created shared_memory")

	return &void, nil
}

func (m *SharedMemoriesRepo) GetById(id *mp.ById) (*mp.SharedMemoriesGetByIdRes, error) {
	memory := mp.SharedMemoriesGetByIdRes{
		Memory: &mp.SharedMemoriesRes{},
	}
	// Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Message     string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// SharedAt    string `protobuf:"bytes,3,opt,name=shared_at,json=sharedAt,proto3" json:"shared_at,omitempty"`
	// MemoryId    string `protobuf:"bytes,4,opt,name=memory_id,json=memoryId,proto3" json:"memory_id,omitempty"`
	// SharedId    string `protobuf:"bytes,5,opt,name=shared_id,json=sharedId,proto3" json:"shared_id,omitempty"`
	// RecipientId
	query := `
	SELECT 
		id,
		message,
		created_at as shared_at,
		memory_id,
		shared_id,
		recipient_id
	FROM 
		shared_memories
	WHERE 
		id = $1
	AND 
		deleted_at = 0	
	`

	row := m.db.QueryRow(query, id.Id)

	err := row.Scan(
		&memory.Memory.Id,
		&memory.Memory.Message,
		&memory.Memory.SharedAt,
		&memory.Memory.MemoryId,
		&memory.Memory.SharedId,
		&memory.Memory.RecipientId,
	)

	if err != nil {
		log.Println("Error while getting shared memory by id: ", err)
		return nil, err
	}

	log.Println("Successfully got shared memory")

	return &memory, nil
}

func (m *SharedMemoriesRepo) GetAll(req *mp.SharedMemoriesGetAllReq) (*mp.SharedMemoriesGetAllRes, error) {
	memories := mp.SharedMemoriesGetAllRes{}

	query := `
	SELECT 
		id,
		message,
		created_at as shared_at,
		memory_id,
		shared_id,
		recipient_id
	FROM 
		shared_memories
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.SharedId != "" && req.SharedId != "string" {
		conditions = append(conditions, " shared_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.SharedId)
	}
	if req.RecipientId != "" && req.RecipientId != "string" {
		conditions = append(conditions, " recipient_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.RecipientId)
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
		log.Println("Error while retriving shared memories: ", err)
		return nil, err
	}

	for rows.Next() {
		memory := mp.SharedMemoriesRes{}

		err := rows.Scan(
			&memory.Id,
			&memory.Message,
			&memory.SharedAt,
			&memory.MemoryId,
			&memory.SharedId,
			&memory.RecipientId,
		)

		if err != nil {
			log.Println("Error while scanning all shared memories: ", err)
			return nil, err
		}

		count += 1
		memories.Memories = append(memories.Memories, &memory)
	}

	memories.Count = count
	log.Println("Successfully fetched all shared memories")

	return &memories, nil
}

func (m *SharedMemoriesRepo) Update(req *mp.SharedMemoriesUpdateReq) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE
		shared_memories
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Message != "" && req.Message != "string" {
		conditions = append(conditions, " message = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Message)
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
		log.Println("Error while updating shared memory: ", err)
		return nil, err
	}

	log.Println("Successfully updated shared memory")

	return &void, nil
}

func (m *SharedMemoriesRepo) Delete(id *mp.ById) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE 
		shared_memories
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := m.db.Exec(query, id.Id)

	if err != nil {
		log.Println("Error while deleting shared_memory: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("shared_memory with id %s not found", id.Id)
	}

	log.Println("Successfully deleted shared_memory")

	return &void, nil
}
