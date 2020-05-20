package dao

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/sagapi/models/tables"
	"github.com/ontio/sagapi/sagaconfig"
	"time"
)

var DefSagaApiDB *SagaApiDB

type SagaApiDB struct {
	DB *sqlx.DB
}

func NewSagaApiDB(dbConfig *sagaconfig.DBConfig) (*SagaApiDB, error) {
	dbx, dberr := sqlx.Open("mysql",
		dbConfig.ProjectDBUser+
			":"+dbConfig.ProjectDBPassword+
			"@tcp("+dbConfig.ProjectDBUrl+
			")/"+dbConfig.ProjectDBName+
			"?charset=utf8&parseTime=true")
	if dberr != nil {
		return nil, dberr
	}

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	err := dbx.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	dbx.SetMaxIdleConns(256)

	return &SagaApiDB{
		DB: dbx,
	}, nil
}

func (this *SagaApiDB) Exec(tx *sqlx.Tx, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		_, err = tx.Exec(query, args...)
	} else {
		_, err = this.DB.Exec(query, args...)
	}
	return err
}

func (this *SagaApiDB) Select(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		err = tx.Select(dest, query, args...)
	} else {
		err = this.DB.Select(dest, query, args...)
	}
	return err
}

func (this *SagaApiDB) Get(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	var err error
	if tx != nil {
		err = tx.Get(dest, query, args...)
	} else {
		err = this.DB.Get(dest, query, args...)
	}
	return err
}

/////////////////////
func (this *SagaApiDB) InsertApiTag(tx *sqlx.Tx, apiTag *tables.ApiTag) error {
	sqlStr := `insert into tbl_api_tag (ApiId, TagId, State) values (?,?,?)`
	err := this.Exec(tx, sqlStr, apiTag.ApiId, apiTag.TagId, apiTag.State)
	return err
}

func (this *SagaApiDB) InsertTag(tx *sqlx.Tx, Tag *tables.Tag) error {
	sqlStr := `insert into tbl_tag (Name, CategoryId, State) values (?,?,?)`
	err := this.Exec(tx, sqlStr, Tag.Name, Tag.CategoryId, Tag.State)
	return err
}

func (this *SagaApiDB) QueryTagByNameId(tx *sqlx.Tx, categoryId uint32, name string) (*tables.Tag, error) {
	var res tables.Tag
	strSql := `select * from tbl_tag where CategoryId=? and Name=?`
	err := this.Get(tx, &res, strSql, categoryId, name)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *SagaApiDB) InsertCategory(tx *sqlx.Tx, cat *tables.Category) error {
	sqlStr := `insert into tbl_category (NameZh, NameEn,Icon,State) values (?,?,?,?)`
	err := this.Exec(tx, sqlStr, cat.NameZh, cat.NameEn, cat.Icon, cat.State)
	return err
}

func (this *SagaApiDB) QueryCategoryByName(tx *sqlx.Tx, NameEn string) (*tables.Category, error) {
	var res tables.Category
	sqlStr := `select * from tbl_category where NameEn=?`
	err := this.Get(tx, &res, sqlStr, NameEn)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *SagaApiDB) QueryCategoryById(tx *sqlx.Tx, categoryId uint32) (*tables.Category, error) {
	var res tables.Category
	sqlStr := `select * from tbl_category where Id=?`
	err := this.Get(tx, &res, sqlStr, categoryId)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (this *SagaApiDB) Close() error {
	return this.DB.Close()
}
