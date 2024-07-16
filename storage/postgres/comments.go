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

type CommentsRepo struct {
	db *sql.DB
}

func NewCommentsRepo(db *sql.DB) *CommentsRepo {
	return &CommentsRepo{db: db}
}

func (c *CommentsRepo) Create(req *mp.CommentsCreateReq) (*mp.Void, error) {
	id := uuid.New().String()
	void := mp.Void{}

	query := `
	INSERT INTO comments (
		id,
		memory_id,
		user_id,
		content
	) VALUES ($1, $2, $3, $4)
	`

	_, err := c.db.Exec(query, id, req.MemoryId, req.UserId, req.Content)

	if err != nil {
		log.Println("Error while creating comment: ", err)
		return nil, err
	}

	log.Println("Successfully created comment")

	return &void, nil
}

func (c *CommentsRepo) GetById(id *mp.ById) (*mp.CommentsGetByIdRes, error) {
	comment := mp.CommentsGetByIdRes{
		Comment: &mp.CommentsRes{},
	}

	query := `
	SELECT 
		id,
		content,
		memory_id,
		user_id
	FROM 
		comments
	WHERE 
		id = $1
	AND 
		deleted_at = 0	
	`

	row := c.db.QueryRow(query, id.Id)

	err := row.Scan(
		&comment.Comment.Id,
		&comment.Comment.Content,
		&comment.Comment.MemoryId,
		&comment.Comment.UserId,
	)

	if err != nil {
		log.Println("Error while getting comment by id: ", err)
		return nil, err
	}

	log.Println("Successfully got comment")

	return &comment, nil
}

func (c *CommentsRepo) GetAll(req *mp.CommentsGetAllReq) (*mp.CommentsGetAllRes, error) {
	comments := mp.CommentsGetAllRes{}

	query := `
	SELECT 
		id,
		content,
		memory_id,
		user_id
	FROM 
		comments
	WHERE 
		deleted_at = 0	
	`

	var args []interface{}
	var conditions []string

	if req.MemoryId != "" && req.MemoryId != "string" {
		conditions = append(conditions, " memory_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.MemoryId)
	}
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

	rows, err := c.db.Query(query, args...)

	if err != nil {
		log.Println("Error while retriving commments: ", err)
		return nil, err
	}

	for rows.Next() {
		comment := mp.CommentsRes{}

		err := rows.Scan(
			&comment.Id,
			&comment.Content,
			&comment.MemoryId,
			&comment.UserId,
		)

		if err != nil {
			log.Println("Error while scanning all comments: ", err)
			return nil, err
		}
		count += 1
		comments.Comments = append(comments.Comments, &comment)
	}

	comments.Count = count
	log.Println("Successfully fetched all comments")

	return &comments, nil
}

func (c *CommentsRepo) Update(req *mp.CommentsUpdateReq) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE
		comments
	SET 
	`

	var conditions []string
	var args []interface{}

	if req.Content != "" && req.Content != "string" {
		conditions = append(conditions, " content = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Content)
	}

	if len(conditions) == 0 {
		return nil, errors.New("nothing to update")
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0 "

	args = append(args, req.Id)

	_, err := c.db.Exec(query, args...)

	if err != nil {
		log.Println("Error while updating comment: ", err)
		return nil, err
	}

	log.Println("Successfully updated comment")

	return &void, nil
}

func (c *CommentsRepo) Delete(id *mp.ById) (*mp.Void, error) {
	void := mp.Void{}

	query := `
	UPDATE 
		comments
	SET 
		deleted_at = EXTRACT(EPOCH FROM NOW())
	WHERE 
		id = $1
	AND 
		deleted_at = 0
	`

	res, err := c.db.Exec(query, id.Id)

	if err != nil {
		log.Println("Error while deleting comment: ", err)
		return nil, err
	}

	if r, err := res.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("comment with id %s not found", id.Id)
	}

	log.Println("Successfully deleted comment")

	return &void, nil
}
