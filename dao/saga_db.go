package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ontio/sagapi/config"
	"github.com/ontio/sagapi/models/tables"
)

type SagaDB struct {
	db *gorm.DB
}

var DefDB *SagaDB

func NewDB() (*SagaDB, error) {
	db, dberr := gorm.Open("mysql", config.DefConfig.DbConfig.ProjectDBUser+
		":"+config.DefConfig.DbConfig.ProjectDBPassword+
		"@tcp("+config.DefConfig.DbConfig.ProjectDBUrl+
		")/"+config.DefConfig.DbConfig.ProjectDBName+
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
		&tables.Order{},
		&tables.ApiBasicInfo{},
		&tables.APIKey{},
		&tables.ApiDetailInfo{},
		&tables.QrCode{},
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

func (this *SagaDB) InsertApiBasicInfo(apiInfo *tables.ApiBasicInfo) error {
	db := this.db.Create(apiInfo)
	if db.Error != nil {
		return db.Error
	}
	this.db = db
	return nil
}

func (this *SagaDB) QueryPriceByApiId(ApiId int) (string, error) {
	info := &tables.ApiBasicInfo{}
	db := this.db.Table("api_basic_infos").Find(info, "api_id=?", ApiId)
	if db.Error != nil {
		return "", db.Error
	}
	return info.ApiPrice, nil
}

func (this *SagaDB) InsertOrder(buyRecord *tables.Order) error {
	db := this.db.Create(buyRecord)
	if db.Error != nil {
		return db.Error
	}
	this.db = db
	return nil
}

func (this *SagaDB) UpdateTxInfoByOrderId(orderId string, txHash string, status config.OrderStatus) error {
	return this.db.Table("orders").Where("order_id=?", orderId).Update("tx_hash=?,order_status=?", txHash, status).Error
}

func (this *SagaDB) QueryOrderStatusByOrderId(orderId string) (config.OrderStatus, error) {
	order := &tables.Order{}
	err := this.db.Table("orders").Find(order, "order_id=?", orderId).Error
	if err != nil {
		return config.Processing, err
	}
	return order.OrderStatus, nil
}

func (this *SagaDB) InsertQrCode(code *tables.QrCode) error {
	return this.db.Create(code).Error
}

func (this *SagaDB) QueryOrderIdByQrCodeId(qrCodeId string) (string, error) {
	code, err := this.QueryQrCodeById(qrCodeId)
	if err != nil {
		return "", err
	}
	return code.OrderId, nil
}

func (this *SagaDB) QueryQrCodeById(id string) (*tables.QrCode, error) {
	code := &tables.QrCode{}
	err := this.db.Table("qr_codes").Find(code, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return code, nil
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

func (this *SagaDB) QueryApiInfoByPage(start, pageSize int) (infos []tables.ApiBasicInfo, err error) {
	db := this.db.Table("api_basic_infos").Limit(pageSize).Find(&infos, "api_id>=?", start)
	if db.Error != nil {
		return nil, db.Error
	}
	return
}

func (this *SagaDB) QueryApiBasicInfoByApiId(apiId uint) (*tables.ApiBasicInfo, error) {
	info := tables.ApiBasicInfo{}
	db := this.db.Table("api_basic_infos").Find(&info, "api_id=?", apiId)
	if db.Error != nil && db.Error.Error() != "record not found" {
		return nil, db.Error
	}
	return &info, nil
}

func (this *SagaDB) SearchApi(key string) ([]tables.ApiBasicInfo, error) {
	var info []tables.ApiBasicInfo
	k := "%" + key + "%"
	db := this.db.Table("api_basic_infos").Where("api_desc like ?", k).Find(&info)
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
