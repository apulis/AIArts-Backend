package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupSavedImage(r *gin.Engine) {
	group := r.Group("/ai_arts/api/saved_imgs")
	group.Use(Auth())
	group.GET("/", wrapper(lsSavedImages))
	group.GET("/:id", wrapper(getSavedImage))
	group.POST("/", wrapper(createSavedImage))
	group.POST("/:id", wrapper(updateSavedImage))
	group.DELETE("/:id", wrapper(deleteSavedImage))
}

type savedImageId struct {
	ID int `uri:"id" binding:"required"`
}

type lsSavedImagesReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize,default=10"`
	Name     string `form:"name"`
	OrderBy  string `form:"orderBy,default=created_at"`
	Order    string `form:"order,default=desc"`
}

type createSavedImageReq struct {
	Name        string `json:"name" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Description string `json:"description"`
	JobId       string `json:"jobId"`
	IsPrivate   bool   `json:"isPrivate"`
}

type updateSavedImageReq struct {
	Description string `json:"description" binding:"required"`
}

type GetSavedImageResp struct {
	SavedImage models.SavedImage `json:"savedImages"`
}

type GetSavedImagesResp struct {
	SavedImages []models.SavedImage `json:"savedImages"`
	Total       int                 `json:"total"`
	TotalPage   int                 `json:"totalPage"`
	PageNum     int                 `json:"pageNum"`
	PageSize    int                 `json:"pageSize"`
}

// @Summary get saved_images by id
// @Produce  json
// @Param query query lsSavedImagesReq true "query"
// @Success 200 {object} GetSavedImagesResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/saved_imgs [get]
func lsSavedImages(c *gin.Context) error {
	var req lsSavedImagesReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	var savedImages []models.SavedImage
	var total int
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}
	savedImages, total, err = services.ListSavedImages(req.PageNum, req.PageSize, req.OrderBy, req.Order, req.Name, username)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	data := GetSavedImagesResp{
		SavedImages: savedImages,
		Total:       total,
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
		TotalPage:   total/req.PageSize + 1,
	}
	return SuccessResp(c, data)
}

// @Summary get saved image by id
// @Produce  json
// @Param id path int true "saved image id"
// @Success 200 {object} GetSavedImageResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/saved_imgs/:id [get]
func getSavedImage(c *gin.Context) error {
	var id savedImageId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	savedImage, err := services.GetSavedImage(id.ID)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	data := GetSavedImageResp{SavedImage: savedImage}
	return SuccessResp(c, data)
}

// @Summary create saved_image
// @Produce  json
// @Param body body createSavedImageReq true "json body"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/saved_imgs [post]
func createSavedImage(c *gin.Context) error {
	var req createSavedImageReq
	err := c.ShouldBind(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	t, err := services.CreateSavedImage(req.Name, req.Version, req.Description, req.JobId, username, req.IsPrivate)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	data := gin.H{"duration" : t}
	return SuccessResp(c, data)
}

// @Summary update saved_image
// @Produce  json
// @Param description path string true "saved_image description"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/saved_imgs/:id [post]
func updateSavedImage(c *gin.Context) error {
	var id savedImageId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req updateSavedImageReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateSavedImage(id.ID, req.Description)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary delete saved_image by id
// @Produce  json
// @Param id path int true "saved image id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/saved_images/:id [delete]
func deleteSavedImage(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.DeleteSavedImage(id.ID)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
