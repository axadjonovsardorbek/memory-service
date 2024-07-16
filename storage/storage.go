package storage

import (
	mp "memory/genproto"
)

type CommentsI interface {
	Create(*mp.CommentsCreateReq) (*mp.Void, error)
	GetById(*mp.ById) (*mp.CommentsGetByIdRes, error)
	GetAll(*mp.CommentsGetAllReq) (*mp.CommentsGetAllRes, error)
	Update(*mp.CommentsUpdateReq) (*mp.Void, error)
	Delete(*mp.ById) (*mp.Void, error)
}

type MemoriesI interface {
	Create(*mp.MemoriesCreateReq) (*mp.Void, error)
	GetById(*mp.ById) (*mp.MemoriesGetByIdRes, error)
	GetAll(*mp.MemoriesGetAllReq) (*mp.MemoriesGetAllRes, error)
	Update(*mp.MemoriesUpdateReq) (*mp.Void, error)
	Delete(*mp.ById) (*mp.Void, error)
}

type MediasI interface {
	Create(*mp.MediasCreateReq) (*mp.Void, error)
	GetById(*mp.ById) (*mp.MediasGetByIdRes, error)
	GetAll(*mp.MediasGetAllReq) (*mp.MediasGetAllRes, error)
	Update(*mp.MediasUpdateReq) (*mp.Void, error)
	Delete(*mp.ById) (*mp.Void, error)
}

type SharedMemoriesI interface {
	Create(*mp.SharedMemoriesCreateReq) (*mp.Void, error)
	GetById(*mp.ById) (*mp.SharedMemoriesGetByIdRes, error)
	GetAll(*mp.SharedMemoriesGetAllReq) (*mp.SharedMemoriesGetAllRes, error)
	Update(*mp.SharedMemoriesUpdateReq) (*mp.Void, error)
	Delete(*mp.ById) (*mp.Void, error)
}
