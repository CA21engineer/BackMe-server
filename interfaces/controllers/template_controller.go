package controllers

import (
	"ca-zoooom/entity"
	"ca-zoooom/interfaces/db"
	"ca-zoooom/usecase"
	"math"
	"strconv"
	"time"
)

type TemplateController struct {
	Interactor usecase.TemplateInteractor
}

type templateRequest struct {
	DesignPatternId    int      `json:"design_pattern_id"`
	IsPrivate          bool     `json:"is_private"`
	BackGroundUrl      string   `json:"background_url"`
	GeneratedSampleUrl string   `json:"generated_sample_url"`
	Tags               []string `json:"tags"`
}

type templateResponse struct {
	Id                 int       `json:"id"`
	Uid                string    `json:"uid"`
	DesignPatternId    int       `json:"design_pattern_id"`
	IsPrivate          bool      `json:"is_private"`
	BackGroundUrl      string    `json:"background_url"`
	GeneratedSampleUrl string    `json:"generated_sample_url"`
	Tags               []string  `json:"tags"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
}

func NewTemplateController(sqlHandler db.SqlHandler) *TemplateController {
	return &TemplateController{
		Interactor: usecase.TemplateInteractor{
			TemplateRepository: &db.TemplateRepository{
				SqlHandler: sqlHandler,
			},
			TagRepository: &db.TagRepository{
				SqlHandler: sqlHandler,
			},
			TemplateTagRepository: &db.TemplateTagRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *TemplateController) Index(c Context) {
	// ページネーション処理
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pages", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	keyword := c.DefaultQuery("keyword", "")
	offset := limit * (pageNumber - 1)

	templates, total, err := controller.Interactor.ListTemplates(limit, offset, keyword)
	if err != nil {
		c.JSON(controller.Interactor.StatusCode, NewError(err))
		return
	}

	// 全部のページ数。Goの言語仕様上、小数点切り上げを行うためにfloatで計算する必要がある
	totalPageCount := math.Ceil(float64(total) / float64(limit))

	c.JSON(controller.Interactor.StatusCode, H{"templates": templates, "pagination": Pagination{pageNumber, limit, int(totalPageCount)}})
}

func (controller *TemplateController) Show(c Context) {
	uid := c.Param("uid")
	template, tags, err := controller.Interactor.GetByUniqueId(uid)
	if err != nil {
		c.JSON(controller.Interactor.StatusCode, NewError(err))
		return
	}

	c.JSON(controller.Interactor.StatusCode, responseBuilder(template, tags))
}

func (controller *TemplateController) Create(c Context) {
	t := &templateRequest{}
	_ = c.Bind(&t)

	tp, tg := requestConverter(t)

	template, tags, err := controller.Interactor.Add(tp, tg)
	if err != nil {
		c.JSON(controller.Interactor.StatusCode, NewError(err))
		return
	}

	c.JSON(controller.Interactor.StatusCode, responseBuilder(template, tags))
}

func requestConverter(t *templateRequest) (tp *entity.Template, tg []entity.Tag) {
	tp = &entity.Template{
		DesignPatternId:    t.DesignPatternId,
		IsPrivate:          t.IsPrivate,
		BackGroundUrl:      t.BackGroundUrl,
		GeneratedSampleUrl: t.GeneratedSampleUrl,
	}

	for _, title := range t.Tags {
		tg = append(tg, entity.Tag{Title: title})
	}

	return
}

func responseBuilder(tp entity.Template, tg []entity.Tag) (t templateResponse) {
	t.Id = tp.Id
	t.Uid = tp.Uid
	t.DesignPatternId = tp.DesignPatternId
	t.IsPrivate = tp.IsPrivate
	t.GeneratedSampleUrl = tp.GeneratedSampleUrl
	t.BackGroundUrl = tp.BackGroundUrl
	t.UpdatedAt = tp.UpdatedAt
	t.CreatedAt = tp.CreatedAt
	for _, tag := range tg {
		t.Tags = append(t.Tags, tag.Title)
	}

	return
}
