package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAllArticles(page, limit int) ([]models.Article, int, error)
	GetArticleByID(id uint) (models.Article, error)
	CreateArticle(article models.Article) (models.Article, error)
	UpdateArticle(article models.Article) (models.Article, error)
	DeleteArticle(id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *articleRepository) GetAllArticles(page, limit int) ([]models.Article, int, error) {
	var (
		articles []models.Article
		count    int64
	)
	err := r.db.Find(&articles).Count(&count).Error
	if err != nil {
		return articles, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&articles).Error

	return articles, int(count), err
}

func (r *articleRepository) GetArticleByID(id uint) (models.Article, error) {
	var article models.Article
	err := r.db.Where("id = ?", id).First(&article).Error
	return article, err
}

func (r *articleRepository) CreateArticle(article models.Article) (models.Article, error) {
	err := r.db.Create(&article).Error
	return article, err
}

func (r *articleRepository) UpdateArticle(article models.Article) (models.Article, error) {
	err := r.db.Save(&article).Error
	return article, err
}

func (r *articleRepository) DeleteArticle(id uint) error {
	var article models.Article
	err := r.db.Where("id = ?", id).Delete(&article).Error
	return err
}
