package cache

import (
	"context"
	"errors"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"gorm.io/gorm"
	"time"
)

type DatabaseDriver struct {
	db     *gorm.DB
	table  string
	prefix string
}

type cacheData struct {
	Key        string `gorm:"primaryKey"`
	Value      string
	Expiration time.Time
}

var _ contract.Cache = (*DatabaseDriver)(nil)

func NewDatabaseDriver(db *gorm.DB, table, prefix string) *DatabaseDriver {
	return &DatabaseDriver{
		db:     db,
		table:  table,
		prefix: prefix,
	}
}

func (d *DatabaseDriver) Get(ctx context.Context, key string) ([]byte, error) {
	db := d.db.WithContext(ctx)
	var data cacheData
	if err := db.Table(d.table).First(&data, "`key` = ?", d.prefix+key).Error; err != nil {
		return nil, err
	}
	if data.Expiration.Before(time.Now()) {
		db.Table(d.table).Delete(&data)
		return nil, errors.New("database: expiration")
	}

	return []byte(data.Value), nil
}

func (d *DatabaseDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	db := d.db.WithContext(ctx)
	data := cacheData{
		d.prefix + key,
		string(value),
		time.Now().Add(ttl),
	}

	return db.Table(d.table).Save(&data).Error
}

func (d *DatabaseDriver) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	db := d.db.WithContext(ctx)
	for i, v := range keys {
		keys[i] = d.prefix + v
	}

	return db.Table(d.table).Delete(cacheData{}, "`key` IN ? ", keys).Error
}

func (d *DatabaseDriver) Has(ctx context.Context, key string) bool {
	db := d.db.WithContext(ctx)
	var Found bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM "+d.table+" WHERE `key` = ?) AS found", d.prefix+key).Scan(&Found)
	return Found
}

func (d *DatabaseDriver) ClearPrefix(ctx context.Context, prefix string) error {
	db := d.db.WithContext(ctx)
	return db.Table(d.table).Delete(cacheData{}, "`key` LIKE ? ", d.prefix+prefix+"%").Error
}

func (d *DatabaseDriver) Clear(ctx context.Context) error {
	db := d.db.WithContext(ctx)
	return db.Table(d.table).Where("1=1").Delete(cacheData{}).Error
}
