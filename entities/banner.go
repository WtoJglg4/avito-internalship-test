package entities

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Banner содержит информацию о баннере
type Banner struct {
	BannerID  int       `json:"banner_id" db:"id"`
	TagIDs    []int     `json:"tag_ids" db:"tag_ids"`
	FeatureID int       `json:"feature_id" db:"feature_id"`
	Content   Content   `json:"content" db:"-"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Version   int       `json:"version" db:"version"`
	TagsHash  string    `json:"-" db:"tags_hash"`
}

// Content представляет содержимое баннера
type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

// BannerPostRequest содержит поля, необходимые для создания баннера
type BannerPostRequest struct {
	TagIDs    []int   `json:"tag_ids"`
	FeatureID int     `json:"feature_id"`
	Content   Content `json:"content"`
	IsActive  bool    `json:"is_active"`
}

// BannerDeleteQueryParams содержит параметры запроса для удаления баннеров
type BannerDeleteQueryParams struct {
	FeatureID int `json:"feature_id"`
	TagID     int `json:"tag_id"`
}

func (b *BannerPostRequest) Validate() error {

	if len(b.TagIDs) == 0 {
		return errors.New("tag_ids is empty")
	}
	if err := b.Content.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *Content) Validate() error {
	if len(c.Text) == 0 || len(c.Title) == 0 || len(c.URL) == 0 {
		return errors.New("ivalid content field")
	}
	return nil
}

type QueryFilters struct {
	Feature_id        int
	Tags_id           int
	Limit             int
	Offset            int
	use_last_revision bool
}

func GetAllQueryParams(c *gin.Context) (QueryFilters, error) {
	f := QueryFilters{
		Feature_id: -1,
		Tags_id:    -1,
		Limit:      -1,
		Offset:     -1,
	}

	param := c.Query("feature_id")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Feature_id = num
	}
	param = c.Query("tag_id")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Tags_id = num
	}
	param = c.Query("limit")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Limit = num
	}
	param = c.Query("offset")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Offset = num
	}

	return f, nil
}

func UserGetQueryParams(c *gin.Context) (QueryFilters, error) {
	f := QueryFilters{
		Feature_id:        -1,
		Tags_id:           -1,
		use_last_revision: false,
	}

	param := c.Query("feature_id")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Feature_id = num
	}
	param = c.Query("tag_id")
	if param != "" {
		num, err := strconv.Atoi(param)
		if err != nil {
			return f, err
		}
		if num < 0 {
			return f, errors.New("ivalid query parameter")
		}
		f.Tags_id = num
	}
	param = c.Query("use_last_revision")
	if param == "true" {
		f.use_last_revision = true
	}

	return f, nil
}
