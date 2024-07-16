package service

import (
	"context"
	mp "memory/genproto"
	st "memory/storage/postgres"
)

type MediasService struct {
	storage st.Storage
	mp.UnimplementedMediasServiceServer
}

func NewMediasService(storage *st.Storage) *MediasService {
	return &MediasService{storage: *storage}
}

func (s *MediasService) Create(ctx context.Context, req *mp.MediasCreateReq) (*mp.Void, error) {
	return s.storage.MediaS.Create(req)
}
func (s *MediasService) GetById(ctx context.Context, id *mp.ById) (*mp.MediasGetByIdRes, error) {
	return s.storage.MediaS.GetById(id)
}
func (s *MediasService) GetAll(ctx context.Context, req *mp.MediasGetAllReq) (*mp.MediasGetAllRes, error) {
	return s.storage.MediaS.GetAll(req)
}
func (s *MediasService) Update(ctx context.Context, req *mp.MediasUpdateReq) (*mp.Void, error) {
	return s.storage.MediaS.Update(req)
}
func (s *MediasService) Delete(ctx context.Context, id *mp.ById) (*mp.Void, error) {
	return s.storage.MediaS.Delete(id)
}
