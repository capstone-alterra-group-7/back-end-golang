package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/models"
	"back-end-golang/usecases"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleController interface {
	GetAllArticles(c echo.Context) error
	GetArticleByID(c echo.Context) error
	CreateArticle(c echo.Context) error
	UpdateArticle(c echo.Context) error
	DeleteArticle(c echo.Context) error
}

type articleController struct {
	articleUsecase usecases.ArticleUsecase
}

func NewArticleController(articleUsecase usecases.ArticleUsecase) ArticleController {
	return &articleController{articleUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *articleController) GetAllArticles(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	articles, count, err := c.articleUsecase.GetAllArticles(page, limit)
	if err != nil {

		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching articles",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all article",
			articles,
			page,
			limit,
			count,
		),
	)
}

func (c *articleController) GetArticleByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	article, err := c.articleUsecase.GetArticleByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get article by id",
			article,
		),
	)

}

func (c *articleController) CreateArticle(ctx echo.Context) error {
	var articleInput dtos.ArticleInput
	if err := ctx.Bind(&articleInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding article",
				helpers.GetErrorData(err),
			),
		)
	}

	if articleInput.Image == "" {
		formHeader, err := ctx.FormFile("file")
		if err != nil {
			if err != nil {
				return ctx.JSON(
					http.StatusInternalServerError,
					helpers.NewErrorResponse(
						http.StatusInternalServerError,
						"Error uploading photo",
						helpers.GetErrorData(err),
					),
				)
			}
		}

		//get file from header
		formFile, err := formHeader.Open()
		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)

		if !re.MatchString(formHeader.Filename) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(err),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().FileUpload(models.File{File: formFile})

		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}
		articleInput.Image = uploadUrl
	} else {
		var url models.Url
		url.Url = articleInput.Image

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)
		if !re.MatchString(articleInput.Image) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(nil),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().RemoteUpload(url)
		if uploadUrl == "" || err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		articleInput.Image = uploadUrl
	}

	article, err := c.articleUsecase.CreateArticle(&articleInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a article",
			article,
		),
	)
}

func (c *articleController) UpdateArticle(ctx echo.Context) error {

	var articleInput dtos.ArticleInput
	if err := ctx.Bind(&articleInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding article",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	article, err := c.articleUsecase.GetArticleByID(uint(id))
	if article.ArticleID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article by id",
				helpers.GetErrorData(err),
			),
		)
	}

	if articleInput.Image == "" {
		formHeader, err := ctx.FormFile("file")
		if err != nil {
			if err != nil {
				return ctx.JSON(
					http.StatusInternalServerError,
					helpers.NewErrorResponse(
						http.StatusInternalServerError,
						"Error uploading photo",
						helpers.GetErrorData(err),
					),
				)
			}
		}

		//get file from header
		formFile, err := formHeader.Open()
		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)

		if !re.MatchString(formHeader.Filename) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(err),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().FileUpload(models.File{File: formFile})

		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}
		articleInput.Image = uploadUrl
	} else {
		var url models.Url
		url.Url = articleInput.Image

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)
		if !re.MatchString(articleInput.Image) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(nil),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().RemoteUpload(url)
		if uploadUrl == "" || err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		articleInput.Image = uploadUrl
	}

	articleResp, err := c.articleUsecase.UpdateArticle(uint(id), articleInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated article",
			articleResp,
		),
	)
}

func (c *articleController) DeleteArticle(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.articleUsecase.DeleteArticle(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete article",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted article",
			nil,
		),
	)
}
