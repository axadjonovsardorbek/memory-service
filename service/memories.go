package service

import (
	"context"
	mp "memory/genproto"
	st "memory/storage/postgres"
)

type MemoriesService struct {
	storage st.Storage
	mp.UnimplementedMemoriesServiceServer
}

func NewMemoriesService(storage *st.Storage) *MemoriesService {
	return &MemoriesService{storage: *storage}
}

func (s *MemoriesService) Create(ctx context.Context, req *mp.MemoriesCreateReq) (*mp.Void, error) {
	return s.storage.MemoryS.Create(req)
}
func (s *MemoriesService) GetById(ctx context.Context, id *mp.ById) (*mp.MemoriesGetByIdRes, error) {
	return s.storage.MemoryS.GetById(id)
}
func (s *MemoriesService) GetAll(ctx context.Context, req *mp.MemoriesGetAllReq) (*mp.MemoriesGetAllRes, error) {
	return s.storage.MemoryS.GetAll(req)
}
func (s *MemoriesService) Update(ctx context.Context, req *mp.MemoriesUpdateReq) (*mp.Void, error) {
	return s.storage.MemoryS.Update(req)
}
func (s *MemoriesService) Delete(ctx context.Context, id *mp.ById) (*mp.Void, error) {
	return s.storage.MemoryS.Delete(id)
}
