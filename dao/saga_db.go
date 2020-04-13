package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/saga/config"
	"github.com/ontio/saga/models/tables"
)

type SagaDB struct {
	db *gorm.DB
}

var DefDB *SagaDB

func NewDB() (*SagaDB, error) {
	db, dberr := gorm.Open("mysql", config.DefConfig.ProjectDBUser+
		":"+config.DefConfig.ProjectDBPassword+
		"@tcp("+config.DefConfig.ProjectDBUrl+
		")/"+config.DefConfig.ProjectDBName+
		"?charset=utf8")
	if dberr != nil {
		return nil, fmt.Errorf("[NewSagaDB] open db error: %s", dberr)
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return &SagaDB{
		db: db,
	}, nil
}

func initTables() []interface{} {
	return []interface{}{
		&tables.BuyRecord{},
		&tables.APIInfo{},
		&tables.APIKey{},
		&tables.APIInterfaceInfo{},
		&tables.ApiTestRecord{},
	}
}
func (this *SagaDB) Init() error {
	tabs := initTables()
	for _, table := range tabs {
		if !this.db.HasTable(table) {
			db := this.db.CreateTable(table)
			if db.Error != nil {
				return db.Error
			}
			this.db = db
		}
	}
	return nil
}

func (this *SagaDB) DeleteTable() {
	tabs := initTables()
	for _, table := range tabs {
		this.db.DropTableIfExists(table)
	}
}

func (this *SagaDB) InsertApiInfo(apiInfo *tables.APIInfo) error {
	db := this.db.Create(apiInfo)
	if db.Error != nil {
		return db.Error
	}
	this.db = db
	return nil
}

func (this *SagaDB) QueryPriceByApiId(ApiId int) (string, error) {
	info := &tables.APIInfo{}
	db := this.db.Table("api_infos").Find(info, "api_id=?", ApiId)
	if db.Error != nil {
		return "", db.Error
	}
	return info.ApiPrice, nil
}

func (this *SagaDB) InsertBuyRecord(buyRecord *tables.BuyRecord) error {
	db := this.db.Create(buyRecord)
	if db.Error != nil {
		return db.Error
	}
	this.db = db
	return nil
}

func (this *SagaDB) InsertApiKey(apiKey *tables.APIKey) error {
	db := this.db.Create(apiKey)
	if db.Error != nil {
		return db.Error
	}
	this.db = db
	return nil
}

func (this *SagaDB) QueryRequestNum(apiKey string) (int, error) {
	key := &tables.APIKey{}
	db := this.db.Table("api_keys").Find(key, "api_key=?", apiKey)
	if db.Error != nil {
		return 0, db.Error
	}
	this.db = db
	return key.UsedNum, nil
}

func (this *SagaDB) QueryApiInfoByPage(start, pageSize int) (infos []tables.APIInfo, err error) {
	db := this.db.Table("api_infos").Limit(pageSize).Find(&infos, "api_id>=?", start)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}

func (this *SagaDB) QueryApiInfoByApiId(apiId uint) (*tables.APIInfo, error) {
	info := tables.APIInfo{}
	db := this.db.Table("api_infos").Find(&info, "api_id=?", apiId)
	if db.Error != nil && db.Error.Error() != "record not found" {
		return nil, db.Error
	}
	return &info, nil
}

func (this *SagaDB) SearchApi(key string) ([]tables.APIInfo, error) {
	var info []tables.APIInfo
	k := "%" + key + "%"
	db := this.db.Table("api_infos").Where("api_desc like ?", k).Find(&info)
	if db.Error != nil {
		return nil, db.Error
	}
	return info, nil
}

func (this *SagaDB) VerifyApiKey(apiKey string) error {
	key := &tables.APIKey{}
	db := this.db.Table("api_keys").Find(key, "api_key=?", apiKey)
	if db.Error != nil {
		return db.Error
	}
	if key == nil {
		return fmt.Errorf("invalid api key: %s", apiKey)
	}
	if key.UsedNum >= key.Limit {
		return fmt.Errorf("Available times:%d, has used times: %d", key.Limit, key.UsedNum)
	}
	return nil
}

func QueryTestRecord() {

}

func (this *SagaDB) Close() error {
	return this.db.Close()
}
