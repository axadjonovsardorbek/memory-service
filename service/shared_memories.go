package service

import (
	"context"
	mp "memory/genproto"
	st "memory/storage/postgres"
)

type SharedMemoriesService struct {
	storage st.Storage
	mp.UnimplementedSharedMemoriesServiceServer
}

func NewSharedMemoriesService(storage *st.Storage) *SharedMemoriesService {
	return &SharedMemoriesService{storage: *storage}
}

func (s *SharedMemoriesService) Create(ctx context.Context, req *mp.SharedMemoriesCreateReq) (*mp.Void, error) {
	return s.storage.SharedMemoryS.Create(req)
}
func (s *SharedMemoriesService) GetById(ctx context.Context, id *mp.ById) (*mp.SharedMemoriesGetByIdRes, error) {
	return s.storage.SharedMemoryS.GetById(id)
}
func (s *SharedMemoriesService) GetAll(ctx context.Context, req *mp.SharedMemoriesGetAllReq) (*mp.SharedMemoriesGetAllRes, error) {
	return s.storage.SharedMemoryS.GetAll(req)
}
func (s *SharedMemoriesService) Update(ctx context.Context, req *mp.SharedMemoriesUpdateReq) (*mp.Void, error) {
	return s.storage.SharedMemoryS.Update(req)
}
func (s *SharedMemoriesService) Delete(ctx context.Context, id *mp.ById) (*mp.Void, error) {
	return s.storage.SharedMemoryS.Delete(id)
}
