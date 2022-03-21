package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spaolacci/murmur3"
	"github.com/zhaoyang1214/ginco/app/constant/cachekey"
	"github.com/zhaoyang1214/ginco/app/http/entity/shorturl"
	"github.com/zhaoyang1214/ginco/app/model"
	"github.com/zhaoyang1214/ginco/app/util"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/database"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

// Url Create
// @Summary Url Create
// @Accept json
// @Produce json
// @Param root body shorturl.UrlCreateRequest true "UrlCreate info"
// @Success 200 {object} entity.JSONResult{code=int,message=string,data=shorturl.UrlCreateResponse}
// @BasePath /api
// @Router /url [post]
func UrlCreate(app contract.Application) func(*gin.Context) {
	return func(c *gin.Context) {
		var ucr shorturl.UrlCreateRequest
		if err := c.ShouldBindJSON(&ucr); err != nil {
			util.JsonError(c, "Bind and validate params error, "+err.Error(), http.StatusBadRequest)
			return
		}
		su := model.ShortUrl{}
		db := app.GetI("db").(*database.Database)
		h32 := murmur3.New32WithSeed(0x00)
		if _, err := h32.Write([]byte(ucr.Url)); err != nil {
			util.JsonError(c, err.Error(), http.StatusInternalServerError)
			return
		}
		su.Hash = h32.Sum32()
		if err := db.First(&su).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			util.JsonError(c, err.Error(), http.StatusInternalServerError)
			return
		}
		su.Ttl = ucr.Ttl
		su.Url = ucr.Url
		su.Domain = strings.TrimRight(ucr.Domain, "/")

		result := db.Save(&su)
		if result.Error != nil {
			util.JsonError(c, "save error, "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		cache := app.GetI("cache").(contract.Cache)
		hash62 := util.FormatInt(uint64(su.Hash), 62)
		ttl := time.Duration(su.Ttl) * time.Second

		if err := cache.Set(context.Background(), fmt.Sprintf(cachekey.ShortUrlInfo, hash62), []byte(su.Url), ttl); err != nil {
			util.JsonError(c, "cache error, "+err.Error(), http.StatusInternalServerError)
			return
		}
		util.JsonResponse(c, shorturl.UrlCreateResponse{
			Url: su.Domain + "/" + hash62,
		}, http.StatusOK)
		return
	}
}
