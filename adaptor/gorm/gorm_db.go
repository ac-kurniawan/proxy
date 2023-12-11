package gorm

import (
	"context"
	"time"

	"github.com/ac-kurniawan/proxy/adaptor/gorm/entity"
	"github.com/ac-kurniawan/proxy/core"
	"github.com/ac-kurniawan/proxy/core/model"
	"github.com/ac-kurniawan/proxy/library"
	gorm2 "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDB struct {
	Gorm  *gorm2.DB
	Trace library.AppTrace
}

// FindByUsername implements core.IProxyDB.
func (g *GormDB) FindByUsername(ctx context.Context, username string) (*model.UserModel, error) {
	ctx, span := g.Trace.StartTrace(ctx, "DATABASE - FindByUsername")
	defer g.Trace.EndTrace(span)
	var userEntity entity.UserEntity
	result := g.Gorm.WithContext(ctx).Where("username = ?", username).First(&userEntity)
	if result.Error != nil {
		g.Trace.TraceError(span, result.Error)
		return nil, result.Error
	}
	out := userEntity.ToModel()
	return &out, nil
}

// GetAlmostExpToken implements core.IProxyDB.
func (g *GormDB) GetAlmostExpToken(ctx context.Context, from time.Time, to time.Time) ([]model.UserModel, error) {
	ctx, span := g.Trace.StartTrace(ctx, "DATABASE - GetAlmostExpToken")
	defer g.Trace.EndTrace(span)
	var userEntities []entity.UserEntity
	result := g.Gorm.WithContext(ctx).Where("external_token_expired BETWEEN ? AND ?", from, to).Find(&userEntities)
	if result.Error != nil {
		g.Trace.TraceError(span, result.Error)
		return nil, result.Error
	}
	var out []model.UserModel
	for _, val := range userEntities {
		out = append(out, val.ToModel())
	}
	return out, nil
}

// Save implements core.IProxyDB.
func (g *GormDB) Save(ctx context.Context, userModel model.UserModel) error {
	ctx, span := g.Trace.StartTrace(ctx, "DATABASE - Save")
	defer g.Trace.EndTrace(span)
	var userEntity entity.UserEntity
	userEntity.FromModel(userModel)
	result := g.Gorm.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&userEntity)
	if result.Error != nil {
		g.Trace.TraceError(span, result.Error)
		return result.Error
	}
	return nil
}

func NewGormDB(module GormDB, enableAutoMigration bool) core.IProxyDB {
	if enableAutoMigration {
		module.Gorm.AutoMigrate(&entity.UserEntity{})
	}
	return &module
}
