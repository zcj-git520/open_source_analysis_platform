package data

import (
	"collect_open_source_data/internal/conf"
	"collect_open_source_data/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewOpenSourceRepo)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db, rdb}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名是否加 s
		},
	})

	if err != nil {
		log.Errorf("failed opening connection to mysql: %v", err)
		//panic("failed to connect database")
	}
	_ = db.AutoMigrate(&domain.RepoInfo{}, &domain.Language{}, &domain.Owner{}, &domain.RepoCategory{},
		domain.RepoCategoryId{}, &domain.RepoMetrics{}, &domain.RepoFav{}, &domain.MessageType{}, &domain.Message{})
	return db
}

func NewRedis(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return rdb
}
