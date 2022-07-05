package db

import (
	"strings"

	"github.com/lethe3000/go-common/pkg/log"

	"github.com/lethe3000/go-common/pkg/cecontext"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	Database
	logger *log.Logger
}

func NewRepository(db Database, logger *log.Logger) Repository {
	return Repository{db, logger}
}

func (r Repository) db(ctx *cecontext.Context) *gorm.DB {
	if ctx == nil {
		ctx = cecontext.NewContext()
	}
	if ctx.DbTrx == nil {
		ctx.DbTrx = r.DB
	}
	return ctx.DbTrx
}

func (r Repository) Create(ctx *cecontext.Context, obj interface{}) error {
	return r.db(ctx).Omit(clause.Associations).Create(obj).Error
}

func (r Repository) CreateWithAssociations(ctx *cecontext.Context, obj interface{}) error {
	return r.db(ctx).Create(obj).Error
}

func (r Repository) CreateInBatches(ctx *cecontext.Context, objs interface{}) error {
	return r.db(ctx).CreateInBatches(objs, 100).Error
}

// GetByUid 根据单个uid查询，结果写入obj
func (r Repository) GetByUid(ctx *cecontext.Context, uid string, obj interface{}, preloads ...string) error {
	return withPreloads(r.db(ctx), preloads).Where("uid = ?", uid).First(obj).Error
}

func (r Repository) GetByUids(ctx *cecontext.Context, uids []string, result interface{}, preloads ...string) error {
	return byUids(r.db(ctx), uids, preloads, result)
}

type LimitOffsetPaginator interface {
	Offset() int
	Limit() int
}

// List 默认分页查询
func (r Repository) List(ctx *cecontext.Context, obj interface{}, paginator LimitOffsetPaginator, preloads ...string) (total int64, err error) {
	err = r.db(ctx).Count(&total).Error
	if err != nil {
		return -1, err
	}

	err = withPreloads(r.db(ctx), preloads).Limit(paginator.Limit()).Offset(paginator.Offset()).Find(obj).Error

	return total, err
}

func prepareLike(content string) string {
	specials := []string{"%", "_", "'", "\\"}
	for _, special := range specials {
		content = strings.Replace(content, special, "\\"+special, -1)
	}
	return content
}

// 查看是否有重复名称
func withPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	if preloads != nil && len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	return db
}

// 查询单个默认配置，结果写入result
func byDefault(db *gorm.DB, preloads []string, result interface{}) (err error) {
	return withPreloads(db, preloads).Where("system_default = true").First(result).Error
}

// 根据多个uid查询，结果列表写入result
func byUids(db *gorm.DB, uids []string, preloads []string, result interface{}) (err error) {
	return withPreloads(db, preloads).Where("uid in ?", uids).Find(result).Error
}

// 获取分页列表的通用方法，结果列表写入result
func list(db *gorm.DB, search LimitOffsetPaginator, preloads []string, result interface{}) (total int64, err error) {
	err = db.Count(&total).Error
	if err != nil {
		return -1, err
	}

	err = withPreloads(db, preloads).Limit(search.Limit()).Offset(search.Offset()).Find(result).Error

	return total, err
}

func (r Repository) Migrate(ctx *cecontext.Context, domains ...interface{}) error {
	if err := r.DB.AutoMigrate(domains...); err != nil {
		r.logger.WithContext(ctx).Errorf("自动迁移数据库失败 err=%v", err)
		return err
	}
	return nil
}
