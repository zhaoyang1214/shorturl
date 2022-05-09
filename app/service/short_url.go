package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/zhaoyang1214/ginco/app/constant/cachekey"
	"github.com/zhaoyang1214/ginco/app/entity"
	"github.com/zhaoyang1214/ginco/app/model"
	"github.com/zhaoyang1214/ginco/app/util"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/database"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type ShortUrl struct {
	*Service
}

func NewShortUrl(app contract.Application) *ShortUrl {
	return &ShortUrl{
		&Service{
			app: app,
		},
	}
}

func (s ShortUrl) Create(ucr entity.ShortUrlCreateRequest) (code int, data entity.ShortUrlCreateResponse, err error) {
	code = http.StatusInternalServerError
	data = entity.ShortUrlCreateResponse{}
	err = nil

	su := model.ShortUrl{}
	db := s.app.GetI("db").(*database.Database)

	hash := murmur3.Sum32([]byte(ucr.Url))
	if err = db.Where("hash=?", hash).First(&su).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			su.Hash = hash
		} else {
			return
		}
	}

	su.Ttl = ucr.Ttl
	su.Url = ucr.Url
	su.Domain = strings.TrimRight(ucr.Domain, "/")

	result := db.Save(&su)
	if result.Error != nil {
		err = errors.New("save error, " + result.Error.Error())
		return
	}

	cache := s.app.GetI("cache").(contract.Cache)
	hash62 := util.FormatInt(uint64(su.Hash), 62)
	ttl := time.Duration(su.Ttl) * time.Second

	if err = cache.Set(context.Background(), fmt.Sprintf(cachekey.ShortUrlInfo, hash62), []byte(su.Url), ttl); err != nil {
		err = errors.New("cache error, " + err.Error())
		return
	}

	return http.StatusOK, entity.ShortUrlCreateResponse{Url: su.Domain + "/" + hash62}, nil
}

func (s ShortUrl) GetUrlByHash(hash string) (int, string, error) {
	cache := s.app.GetI("cache").(contract.Cache)
	urlByte, err := cache.Get(context.Background(), fmt.Sprintf(cachekey.ShortUrlInfo, hash))
	if err == nil {
		url := string(urlByte)
		if url == "" {
			return http.StatusNotFound, "", errors.New("not found")
		}
		return http.StatusFound, url, nil
	}

	hint64, err := util.ParseUint(hash, 62)
	if err != nil {
		return http.StatusNotFound, "", errors.New("not found")
	}

	su := model.ShortUrl{}
	db := s.app.GetI("db").(*database.Database)

	err = db.Where("hash=?", hint64).First(&su).Error
	_ = cache.Set(context.Background(), fmt.Sprintf(cachekey.ShortUrlInfo, hash), []byte(su.Url), time.Duration(0))
	if err != nil {
		return http.StatusNotFound, "", errors.New("not found")
	}
	return http.StatusFound, su.Url, nil
}

func (s ShortUrl) List(ulr entity.ShortUrlListRequest) (data entity.ShortUrlListResponse, err error) {
	page, size := ulr.GetPage(), ulr.GetSize()
	db := s.app.GetI("db").(*database.Database)
	var total int64
	var urls []model.ShortUrl
	result := db.Limit(size).Offset((page - 1) * size).Find(&urls).Count(&total)
	err = result.Error
	for _, v := range urls {
		data.List = append(data.List, entity.ShortUrlListResponseWithList{
			Hash:      util.FormatInt(uint64(v.Hash), 62),
			Url:       v.Url,
			Ttl:       v.Ttl,
			Domain:    v.Domain,
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: v.UpdatedAt.String(),
		})
	}
	return
}
