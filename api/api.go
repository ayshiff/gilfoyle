//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go --parseDependency
package api

import (
	"errors"
	"github.com/dreamvo/gilfoyle"
	_ "github.com/dreamvo/gilfoyle/api/docs"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	defaultItemsPerPage = 50
	maxItemsPerPage     = 100
)

var (
	ErrInvalidUUID      = errors.New("invalid UUID provided")
	ErrResourceNotFound = errors.New("resource not found")
)

// @title Gilfoyle server
// @description Cloud-native media hosting & streaming server for businesses.
// @version v1
// @host demo-v1.gilfoyle.dreamvo.com
// @BasePath /
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/dreamvo/gilfoyle/blob/master/LICENSE

func NewServer() *gin.Engine {
	r := gin.New()
	RegisterMiddlewares(r)
	RegisterRoutes(r)
	return r
}

// RegisterMiddlewares adds middlewares to a given router instance
func RegisterMiddlewares(r *gin.Engine) *gin.Engine {
	r.Use(gin.Recovery())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	// TODO(sundowndev): update Gin to enable this feature. See https://github.com/gin-gonic/gin/commits/master/recovery.go
	//r.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
	//	if err, ok := recovered.(string); ok {
	//		util.NewError(ctx, http.StatusInternalServerError, errors.New(err))
	//	}
	//	util.NewError(ctx, http.StatusInternalServerError, errors.New("an unexpected error occurred"))
	//}))

	r.Use(func(ctx *gin.Context) {
		ctx.Next()

		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		errorMsg := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log := gilfoyle.Logger.With(
			zap.String("Method", ctx.Request.Method),
			zap.String("Path", path),
			zap.Int("StatusCode", ctx.Writer.Status()),
			zap.Int("BodySize", ctx.Writer.Size()),
			zap.String("UserAgent", ctx.Request.UserAgent()),
		)

		if errorMsg != "" {
			log.Error("Incoming HTTP Request",
				zap.String("ErrorMessage", errorMsg),
			)
			return
		}

		log.Info("Incoming HTTP Request")
	})

	return r
}

// RegisterRoutes adds routes to a given router instance
func RegisterRoutes(r *gin.Engine) *gin.Engine {
	r.GET("/healthz", healthCheckHandler)
	r.GET("/metricsz", getMetrics)

	medias := r.Group("/medias")
	{
		medias.GET("", paginateHandler, getAllMedias)
		medias.GET(":id", getMedia)
		medias.DELETE(":id", deleteMedia)
		medias.POST("", createMedia)
		medias.PATCH(":id", updateMedia)
		medias.POST(":id/upload/video", uploadVideoFile)
		medias.POST(":id/upload/audio", uploadAudioFile)
		medias.GET(":id/attachments", getMediaAttachments)
		medias.POST(":id/attachments", addMediaAttachments)
		medias.DELETE(":id/attachments/:attachment_id", deleteMediaAttachments)
		medias.GET(":id/stream/:preset", streamMedia)
	}

	if gilfoyle.Config.Settings.ExposeSwaggerUI {
		// Register swagger docs handler
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Use(func(ctx *gin.Context) {
		util.NewError(ctx, http.StatusNotFound, errors.New("resource not found"))
	})

	return r
}

func paginateHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil || limitInt > maxItemsPerPage {
		limitInt = defaultItemsPerPage
	}

	offset := ctx.Query("offset")
	offsetInt, err := strconv.ParseInt(offset, 10, 64)

	if err != nil {
		offsetInt = 0
	}

	ctx.Set("limit", int(limitInt))
	ctx.Set("offset", int(offsetInt))
	ctx.Next()
}

func rollbackWithError(ctx *gin.Context, tx *ent.Tx, statusCode int, err error) {
	if txErr := tx.Rollback(); txErr != nil {
		util.NewError(ctx, statusCode, txErr)
	}

	util.NewError(ctx, statusCode, err)
}
