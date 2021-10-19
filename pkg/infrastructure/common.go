package infrastructure

import (
	"math"
	"strings"

	"github.com/takemo101/dc-scheduler/pkg/application"
	"gorm.io/gorm"
)

// --- GormPaging ---

// GormPagingParameter
type GormPaging struct {
	db      *gorm.DB
	page    int
	limit   int
	orderBy []string
}

// NewGormPaging コンストラクタ
func NewGormPaging(
	db *gorm.DB,
	page int,
	limit int,
	orderBy []string,
) GormPaging {
	return GormPaging{
		db,
		page,
		limit,
		orderBy,
	}
}

// GormPaging Gormを利用してページネーターを生成する
func (p *GormPaging) Paging(anyType interface{}) (paginator application.Paginator, err error) {
	db := p.db

	if p.page < 1 {
		p.page = 1
	}
	if p.limit == 0 {
		p.limit = 10
	}
	if len(p.orderBy) > 0 {
		for _, o := range p.orderBy {
			db = db.Order(o)
		}
	}

	var offset int

	// record count
	count, err := p.countRecord(anyType)
	if err != nil {
		return paginator, err
	}

	if p.page <= 1 {
		offset = 0
	} else {
		offset = (p.page - 1) * p.limit
	}

	// レコードを取得
	if err := db.Limit(p.limit).Offset(offset).Find(anyType).Error; err != nil {
		return paginator, err
	}

	// Paginatorに値をセット
	paginator.TotalCount = count
	paginator.CurrentPage = p.page

	paginator.Offset = offset
	paginator.Limit = p.limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.limit)))

	if p.page > 1 {
		paginator.PrevPage = p.page - 1
	} else {
		paginator.PrevPage = p.page
	}

	if p.page == paginator.TotalPage {
		paginator.NextPage = p.page
	} else {
		paginator.NextPage = p.page + 1
	}

	paginator.FirstCount = offset
	nextCount := p.page * p.limit
	if count < nextCount {
		paginator.LastCount = count
	} else {
		paginator.LastCount = nextCount
	}

	return paginator, err
}

// countRecord 指定したGormModelからレコードをカウントする
func (p GormPaging) countRecord(anyType interface{}) (int, error) {
	var count int64
	err := p.db.Model(anyType).Count(&count).Error
	return int(count), err
}

// --- GetNextIdentitySelectSQL ---

// GetNextIdentitySelectSQL Gormで次のIDを生成するSelect文を出力する
func GetNextIdentitySelectSQL(dbType string) string {
	lowerType := strings.ToLower(dbType)
	// DBタイプがSQLiteの場合とその他を分ける
	if lowerType == "sqlite" {
		return "LAST_INSERT_ROWID()"
	}
	return "LAST_INSERT_ID()"
}
