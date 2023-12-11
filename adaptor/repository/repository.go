package repository

import "github.com/ac-kurniawan/proxy/core"

type Repository struct {
	core.IDataSource
	core.IProxyDB
}

func NewRepository(module Repository) core.IRepository {
	return &module
}
