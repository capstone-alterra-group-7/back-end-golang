package repositories

import (
	"back-end-golang/configs"
	"back-end-golang/models"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestArticleRepository_GetAllArticles(t *testing.T) {
	db, err := configs.ConnectDBTest()
	assert.NoError(t, err)

	repo := NewArticleRepository(db)

	// Artikel yang akan dibuat
	newArticles := []models.Article{
		{Title: "Article 1", Image: "image.jpg", Description: "Content 1", Label: "train"},
		{Title: "Article 2", Image: "image.jpg", Description: "Content 2", Label: "train"},
		{Title: "Article 3", Image: "image.jpg", Description: "Content 3", Label: "train"},
		// Add the remaining expected articles here
	}

	// Membuat artikel baru
	for _, newArticle := range newArticles {
		createdArticle, err := repo.CreateArticle(newArticle)
		if err != nil {
			t.Errorf("Failed to create article: %v", err)
			return
		}

		// Assersikan artikel yang berhasil dibuat
		assert.NotZero(t, createdArticle.ID)
		assert.Equal(t, newArticle.Title, createdArticle.Title)
		assert.Equal(t, newArticle.Image, createdArticle.Image)
		assert.Equal(t, newArticle.Description, createdArticle.Description)
		assert.Equal(t, newArticle.Label, createdArticle.Label)
	}

	articles, _, err := repo.GetAllArticles(1, 3)
	if err != nil {
		t.Errorf("Error getting all articles: %v", err)
		return
	}

	expectedNumArticles := 3
	if len(articles) != expectedNumArticles {
		t.Errorf("Expected %d articles, but got %d", expectedNumArticles, len(articles))
		return
	}

	expectedArticles := []models.Article{
		{Title: "Article 1", Image: "image.jpg", Description: "Content 1", Label: "train"},
		{Title: "Article 2", Image: "image.jpg", Description: "Content 2", Label: "train"},
		{Title: "Article 3", Image: "image.jpg", Description: "Content 3", Label: "train"},
		// Add the remaining expected articles here
	}

	for i := 0; i < expectedNumArticles; i++ {
		if !reflect.DeepEqual(articles[i].Title, expectedArticles[i].Title) ||
			!reflect.DeepEqual(articles[i].Image, expectedArticles[i].Image) ||
			!reflect.DeepEqual(articles[i].Description, expectedArticles[i].Description) ||
			!reflect.DeepEqual(articles[i].Label, expectedArticles[i].Label) {
			t.Errorf("Mismatch in article at index %d\nExpected: %+v\nActual: %+v",
				i, expectedArticles[i], articles[i])
		}
	}

	// Menghapus artikel yang dibuat
	for _, createdArticle := range articles {
		err = repo.ForceDeleteArticle(createdArticle.ID)
		if err != nil {
			t.Errorf("Failed to delete article: %v", err)
		}
	}
}

func TestArticleRepository_GetArticleByID(t *testing.T) {
	db, err := configs.ConnectDBTest()
	assert.NoError(t, err)

	repo := NewArticleRepository(db)

	// Simulasikan artikel dengan ID 1 di database
	expectedArticle := models.Article{Model: gorm.Model{ID: 1, CreatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), UpdatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false}}, Title: "Article 1", Image: "image.jpg", Description: "Content 1", Label: "train"}

	// Menambahkan artikel ke dalam database mock
	_, err = repo.CreateArticle(expectedArticle)
	assert.NoError(t, err)

	// Mendapatkan artikel dengan ID 1
	article, err := repo.GetArticleByID(expectedArticle.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedArticle, article)

	// Hapus artikel setelah pengujian selesai
	defer func() {
		err := repo.ForceDeleteArticle(article.ID)
		if err != nil {
			t.Errorf("Failed to delete article: %v", err)
		}
	}()
}
func TestArticleRepository_CreateArticle(t *testing.T) {
	db, err := configs.ConnectDBTest()
	assert.NoError(t, err)

	repo := NewArticleRepository(db)

	// Artikel yang akan dibuat
	newArticle := models.Article{Title: "New Article", Image: "image.jpg", Description: "New Content", Label: "train"}

	// Membuat artikel baru
	createdArticle, err := repo.CreateArticle(newArticle)
	if err != nil {
		t.Errorf("Failed to create article: %v", err)
		return
	}

	// Verifikasi artikel yang dibuat
	assert.NotZero(t, createdArticle.ID)
	assert.Equal(t, newArticle.Title, createdArticle.Title)
	assert.Equal(t, newArticle.Image, createdArticle.Image)
	assert.Equal(t, newArticle.Description, createdArticle.Description)
	assert.Equal(t, newArticle.Label, createdArticle.Label)

	// Hapus artikel setelah pengujian selesai
	defer func() {
		err := repo.ForceDeleteArticle(createdArticle.ID)
		if err != nil {
			t.Errorf("Failed to delete article: %v", err)
		}
	}()
}

func TestArticleRepository_UpdateArticle(t *testing.T) {
	db, err := configs.ConnectDBTest()
	assert.NoError(t, err)

	repo := NewArticleRepository(db)

	// Simulasikan artikel dengan ID 1 di database
	article := models.Article{Model: gorm.Model{ID: 1, CreatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), UpdatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false}}, Title: "Article 1", Image: "image.jpg", Description: "Content 1", Label: "train"}

	// Menambahkan artikel ke dalam database mock
	article, err = repo.CreateArticle(article)
	assert.NoError(t, err)

	// Mengubah judul artikel
	updatedTitle := "Updated Article 1"
	article.Title = updatedTitle

	// Mengupdate artikel
	updatedArticle, err := repo.UpdateArticle(article)

	assert.NoError(t, err)
	assert.Equal(t, updatedTitle, updatedArticle.Title)

	// Hapus artikel setelah pengujian selesai
	defer func() {
		err := repo.ForceDeleteArticle(article.ID)
		if err != nil {
			t.Errorf("Failed to delete article: %v", err)
		}
	}()
}

func TestArticleRepository_DeleteArticle(t *testing.T) {
	db, err := configs.ConnectDBTest()
	assert.NoError(t, err)

	repo := NewArticleRepository(db)

	// Simulasikan artikel dengan ID 1 di database
	article := models.Article{Model: gorm.Model{ID: 1, CreatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), UpdatedAt: time.Date(2023, time.June, 21, 9, 2, 53, 603000000, time.Local), DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false}}, Title: "Article 1", Image: "image.jpg", Description: "Content 1", Label: "train"}

	// Menambahkan artikel ke dalam database mock
	_, err = repo.CreateArticle(article)
	assert.NoError(t, err)

	// Menghapus artikel dengan ID 1
	err = repo.ForceDeleteArticle(article.ID)

	assert.NoError(t, err)
}
