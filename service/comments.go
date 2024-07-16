package service

import (
	"context"
	mp "memory/genproto"
	st "memory/storage/postgres"
)

type CommentsService struct {
	storage st.Storage
	mp.UnimplementedCommentsServiceServer
}

func NewCommentsService(storage *st.Storage) *CommentsService {
	return &CommentsService{storage: *storage}
}

func (s *CommentsService) Create(ctx context.Context, req *mp.CommentsCreateReq) (*mp.Void, error) {
	return s.storage.CommentS.Create(req)
}
func (s *CommentsService) GetById(ctx context.Context, id *mp.ById) (*mp.CommentsGetByIdRes, error) {
	return s.storage.CommentS.GetById(id)
}
func (s *CommentsService) GetAll(ctx context.Context, req *mp.CommentsGetAllReq) (*mp.CommentsGetAllRes, error) {
	return s.storage.CommentS.GetAll(req)
}
func (s *CommentsService) Update(ctx context.Context, req *mp.CommentsUpdateReq) (*mp.Void, error) {
	return s.storage.CommentS.Update(req)
}
func (s *CommentsService) Delete(ctx context.Context, id *mp.ById) (*mp.Void, error) {
	return s.storage.CommentS.Delete(id)
}
