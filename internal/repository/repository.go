package repository

import (
	"fmt"
	"log/slog"

	"news-rest-api/internal/dto"
	"news-rest-api/internal/models"
	e "news-rest-api/internal/pkg/errors"

	"github.com/lib/pq"
	"gopkg.in/reform.v1"
)

type Repository interface {
	UpdateNews(news dto.News, id uint64) error
	GetNewsList(pageSize, offset int) ([]dto.News, error)
}

type Repo struct {
	db     *reform.DB
	logger *slog.Logger
}

func NewRepository(db *reform.DB, logger *slog.Logger) *Repo {
	return &Repo{db: db, logger: logger}
}

const (
	queryGetNewsList = `SELECT "News"."Id", "News"."Title", "News"."Content", 
       COALESCE(array_agg("NewsCategories"."CategoryId"), '{}') AS "Categories"
	FROM "News" LEFT JOIN "NewsCategories" ON "News"."Id" = "NewsCategories"."NewsId"
	GROUP BY "News"."Id", "News"."Title", "News"."Content"
	LIMIT $1 OFFSET $2;`
)

func (r *Repo) UpdateNews(input dto.News, id uint64) error {
	r.logger.Info("Starting update news")
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	existingNews, err := tx.FindByPrimaryKeyFrom(models.NewsTable, id)
	if err != nil {
		r.logger.Error("Failed to find news", "news_id", id, "error", err)
		return e.ErrInvalidNewsId
	}

	news, ok := existingNews.(*models.News)
	if !ok {
		return fmt.Errorf("invalid news type")
	}

	if input.Title != "" {
		news.Title = input.Title
	}

	if input.Content != "" {
		news.Content = input.Content
	}

	err = tx.Update(news)
	if err != nil {
		r.logger.Error("Failed to update news", "news_id", news.ID, "error", err)
		return err
	}

	if input.Categories != nil {
		_, err = tx.DeleteFrom(models.NewsCategoryView, `WHERE "NewsId" = $1`, news.ID)
		if err != nil {
			r.logger.Error("Failed to delete categories news", "news_id", news.ID, "error", err)
			return err
		}

		for _, categoryID := range input.Categories {
			newsCategory := &models.NewsCategory{
				NewsID:     news.ID,
				CategoryID: uint64(categoryID),
			}
			err = tx.Insert(newsCategory)
			if err != nil {
				r.logger.Error("Failed to add categories news", "news_id", news.ID, "category_id", categoryID, "error", err)
				return err
			}
		}
	}

	r.logger.Info("News updated", "news_title", news.Title)
	return tx.Commit()
}

func (r *Repo) GetNewsList(pageSize, offset int) ([]dto.News, error) {
	r.logger.Info("Starting get news list")
	rows, err := r.db.Query(queryGetNewsList, pageSize, offset)
	if err != nil {
		r.logger.Error("Failed to execute query to get news list", "error", err)
		return nil, err
	}
	defer rows.Close()

	var newsList []dto.News
	for rows.Next() {
		var news dto.News
		var categories []int64

		if err = rows.Scan(&news.ID, &news.Title, &news.Content, pq.Array(&categories)); err != nil {
			r.logger.Error("Failed to parse row", "error", err)
			return nil, err
		}

		news.Categories = categories

		newsList = append(newsList, news)
	}

	if len(newsList) == 0 {
		r.logger.Error("No news found for the given page")
		return nil, e.ErrNotFoundNews
	}

	r.logger.Info("News list created")
	return newsList, nil
}
