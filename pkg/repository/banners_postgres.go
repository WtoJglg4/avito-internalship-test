package repository

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github/avito/entities"
	"sort"
	"strconv"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	featuresTable       = "features"
	bannersTable        = "banners"
	bannersTagsTable    = "banner_tags"
	tagsTable           = "tags"
	bannerVersionsTable = "banner_versions"
)

type BannersPostgres struct {
	db *sqlx.DB
}

func NewBannersPostgres(db *sqlx.DB) *BannersPostgres {
	return &BannersPostgres{db: db}
}

func (r *BannersPostgres) CreateBanner(banner entities.Banner) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error(err.Error())
		return 0, tx.Rollback()
	}

	//check exist all tags and features
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 LIMIT 1);", featuresTable)
	if err := r.db.Get(&exists, query, banner.FeatureID); err != nil {
		logrus.Error(err.Error())
		return 0, tx.Rollback()
	}
	if !exists {
		err := errors.New("feature with this ID not exists")
		tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 LIMIT 1);", tagsTable)
	for _, tagId := range banner.TagIDs {
		if err := r.db.Get(&exists, query, tagId); err != nil {
			logrus.Error(err.Error())
			return 0, tx.Rollback()
		}
		if !exists {
			err := fmt.Errorf("tag with ID %d not exists", tagId)
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}
	}

	//sort tags in asc order
	sort.Ints(banner.TagIDs)

	//get tags_hash
	hashBytes := hashInts(banner.TagIDs)
	banner.TagsHash = fmt.Sprintf("%x", hashBytes)

	query = fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE tags_hash = $1 AND feature_id = $2)", bannersTable)
	if err := r.db.Get(&exists, query, banner.TagsHash, banner.FeatureID); err != nil {
		logrus.Error(err.Error())
		return 0, tx.Rollback()
	}

	if !exists {
		query = fmt.Sprintf("INSERT INTO %s (feature_id, content, is_active, created_at, updated_at, version, tags_hash) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", bannersTable)
		content, err := json.Marshal(banner.Content)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}

		err = r.db.QueryRow(query, banner.FeatureID, content, banner.IsActive, banner.CreatedAt, banner.UpdatedAt, 1, banner.TagsHash).Scan(&banner.BannerID)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}

		query = fmt.Sprintf("INSERT INTO %s (banner_id, tag_id) VALUES($1, $2)", bannersTagsTable)
		for _, tag := range banner.TagIDs {
			_, err = r.db.Exec(query, banner.BannerID, tag)
			if err != nil {
				logrus.Error(err.Error())
				tx.Rollback()
				return 0, err
			}
		}

		return banner.BannerID, tx.Commit()
	} else {
		var lastVersionBanner struct {
			entities.Banner
			JsonContent string `db:"content"`
		}

		query = fmt.Sprintf("SELECT * FROM %s WHERE tags_hash = $1 AND feature_id = $2", bannersTable)
		if err := r.db.Get(&lastVersionBanner, query, banner.TagsHash, banner.FeatureID); err != nil {
			logrus.Error(err.Error())
			return 0, tx.Rollback()
		}

		// if exist: copy update updated_at, update version, update IsActive, return last version banner index
		query = fmt.Sprintf("INSERT INTO %s (banner_id, content, version, updated_at, is_active) VALUES($1, $2, $3, $4, $5)", bannerVersionsTable)
		_, err = r.db.Exec(query, lastVersionBanner.BannerID, lastVersionBanner.JsonContent, lastVersionBanner.Version, lastVersionBanner.UpdatedAt, false)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}

		content, err := json.Marshal(banner.Content)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}

		query = fmt.Sprintf("UPDATE %s SET content = $1, version = $2, updated_at = $3, is_active = $4 WHERE id = $5", bannersTable)
		_, err = r.db.Exec(query, content, lastVersionBanner.Version+1, banner.UpdatedAt, banner.IsActive, lastVersionBanner.BannerID)
		if err != nil {
			logrus.Error(err.Error())
			tx.Rollback()
			return 0, err
		}

		return lastVersionBanner.BannerID, tx.Commit()
	}

}

func hashInts(ints []int) []byte {
	hash := sha256.New()

	buf := make([]byte, 8) // Буфер для представления int64 (максимально возможный размер int в Go)
	for _, n := range ints {
		binary.LittleEndian.PutUint64(buf, uint64(n)) // Преобразование int в []byte
		hash.Write(buf)                               // Добавление в хеш
	}

	return hash.Sum(nil)
}

func (r *BannersPostgres) GetAllBanners(filters entities.QueryFilters) ([]entities.Banner, error) {
	// SELECT b.id, b.feature_id, b.content, b.is_active, b.version, b.created_at, b.updated_at, array_agg(bt.tag_id) AS tag_ids FROM  banners b LEFT JOIN banner_tags bt ON b.id = bt.banner_id GROUP BY b.id, b.feature_id, b.content, b.is_active, b.version, b.created_at, b.updated_at;
	banners := make([]entities.Banner, 0)
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().Select("b.id, b.feature_id, b.content, b.is_active, b.version, b.created_at, b.updated_at, array_agg(bt.tag_id) AS tag_ids").
		From(fmt.Sprintf("%s b", bannersTable)).
		SQL("LEFT").Join("banner_tags bt ON b.id = bt.banner_id")

	where := ""
	if filters.Feature_id != -1 {
		where = fmt.Sprintf("WHERE %s feature_id = %d ", where, filters.Feature_id)
	}
	if filters.Tags_id != -1 {
		if filters.Feature_id != -1 {
			where += "AND "
		} else {
			where = "WHERE "
		}
		where = fmt.Sprintf("%s  b.id IN (SELECT banner_id FROM banner_tags WHERE tag_id = %d) ", where, filters.Tags_id)
	}
	sb.SQL(where).
		GroupBy("b.id, b.feature_id, b.content, b.is_active, b.version, b.created_at, b.updated_at")

	if filters.Limit != -1 {
		sb.Limit(filters.Limit)
	}
	if filters.Offset != -1 {
		sb.Offset(filters.Offset)
	}
	query, _ := sb.Build()
	// logrus.Println(query, args)

	rows, err := r.db.Query(query)
	if err != nil {
		return banners, err
	}
	defer rows.Close()

	for rows.Next() {
		bannerWithContent := struct {
			entities.Banner
			JsonContent string `json:"content"`
		}{}
		var tagIDs IntArray
		if err := rows.Scan(
			&bannerWithContent.BannerID,
			&bannerWithContent.FeatureID,
			&bannerWithContent.JsonContent,
			&bannerWithContent.IsActive,
			&bannerWithContent.Version,
			&bannerWithContent.CreatedAt,
			&bannerWithContent.UpdatedAt,
			&tagIDs); err != nil {
			// pq.Array(&bannerWithContent.TagIDs)); err != nil {

			return banners, err
		}
		for _, v := range tagIDs {
			bannerWithContent.TagIDs = append(bannerWithContent.TagIDs, v)
		}
		if err := json.Unmarshal([]byte(bannerWithContent.JsonContent), &bannerWithContent.Content); err != nil {
			return banners, err
		}
		// logrus.Println(bannerWithContent)
		banners = append(banners, bannerWithContent.Banner)
	}

	return banners, nil
}

type IntArray []int

// Scan - метод, который реализует интерфейс sql.Scanner.
func (a *IntArray) Scan(src interface{}) error {
	// Преобразуем значение в строку.
	s, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("IntArray.Scan: expected []byte, got %T", src)
	}

	// Удаляем фигурные скобки.
	s = bytes.Trim(s, "{}")

	// Разделяем строку по запятым.
	parts := bytes.Split(s, []byte(","))

	// Преобразуем каждую часть в int и добавляем в срез.
	for _, part := range parts {
		val, err := strconv.Atoi(string(part))
		if err != nil {
			return err
		}
		*a = append(*a, val)
	}

	return nil
}

var ErrNoRowsDeleted = errors.New("no rows deleted")

func (r *BannersPostgres) DeleteBanners(filters entities.QueryFilters) error {
	where := ""
	if filters.Feature_id != -1 {
		where += fmt.Sprintf("AND b.feature_id = %d", filters.Feature_id)
	}
	if filters.Tags_id != -1 {
		where += fmt.Sprintf("AND bt.tag_id = %d", filters.Tags_id)
	}

	query := fmt.Sprintf("DELETE FROM %s b USING %s bt WHERE b.id = bt.banner_id %s", bannersTable, bannersTagsTable, where)
	logrus.Println(query)

	res, err := r.db.Exec(query)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.Errorf("rows affected: %v, %v", rowsAffected, err.Error())
		return err
	}

	if rowsAffected == 0 {
		logrus.Errorf("rows affected: %v, %v", rowsAffected, ErrNoRowsDeleted)
		return ErrNoRowsDeleted
	}

	return nil
}

var ErrNoRowsSelected = errors.New("banner not found")

func (r *BannersPostgres) UserBanner(filters entities.QueryFilters) (entities.Content, error) {
	where := "WHERE b.is_active = true "
	if filters.Feature_id != -1 {
		where += fmt.Sprintf("AND b.feature_id = %d ", filters.Feature_id)
	}
	if filters.Tags_id != -1 {
		where += fmt.Sprintf("AND bt.tag_id = %d", filters.Tags_id)
	}
	query := "SELECT b.content FROM banners b join banner_tags bt on b.id = bt.banner_id " + where

	logrus.Println(query)

	var jsonContent []byte
	var bannerContent entities.Content

	row, err := r.db.Query(query)
	if err != nil {
		logrus.Error(err.Error())
		return bannerContent, err
	}
	defer row.Close()

	count := 0
	for row.Next() {
		count++
		if err := row.Scan(&jsonContent); err != nil {
			logrus.Error(err.Error())
			return bannerContent, err
		}
	}
	if count == 0 {
		logrus.Error(ErrNoRowsSelected)
		return bannerContent, ErrNoRowsSelected
	}
	if err := json.Unmarshal(jsonContent, &bannerContent); err != nil {
		logrus.Error(err.Error())
		return bannerContent, err
	}

	return bannerContent, nil
}

func (r *BannersPostgres) InsertTestTags() {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ('test')", tagsTable)
	for i := 0; i < 100; i++ {
		r.db.Exec(query)
	}
}

func (r *BannersPostgres) InsertTestFeatures() {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ('test')", featuresTable)
	for i := 0; i < 100; i++ {
		r.db.Exec(query)
	}
}
