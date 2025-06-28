package gormrepository

import "gorm.io/gorm"

type Repository interface {
	FindByRaw(target interface{}, query string) error

	FindAll(target interface{}, preloads ...string) error
	FindAllAndLimit(target interface{}, limit int, preloads ...string) error
	FindBatch(target interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error)

	FindWhere(target interface{}, condition string, preloads ...string) error
	FindWhereBatch(target interface{}, condition string, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error)

	FindByField(target interface{}, field string, value interface{}, preloads ...string) error
	FindByFields(target interface{}, fields map[string]interface{}, preloads ...string) error
	FindByFieldsOrder(target interface{}, fields map[string]interface{}, order string, preloads ...string) error

	FindByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error)
	FindByFieldsBatch(target interface{}, fields map[string]interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error)

	FindOneByField(target interface{}, field string, value interface{}, preloads ...string) error
	FindOneLastByField(target interface{}, field string, value interface{}, preloads ...string) error
	FindOneByFields(target interface{}, fields map[string]interface{}, preloads ...string) error
	FindOneByCondition(target interface{}, condition string, preloads ...string) error
	FindOneLastByCondition(target interface{}, condition string, preloads ...string) error
	FindOneByFieldOrder(target interface{}, field string, value interface{}, order string, preloads ...string) error

	FindOneByID(target interface{}, id interface{}, preloads ...string) error
	FindByConditionThenLock(target interface{}, condition string, preloads ...string) error

	Create(target interface{}) error
	Save(target interface{}) error
	Delete(target interface{}) error
	UpdateNoTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}) error
	DeleteByCondition(target interface{}, condition string) error
	BatchInsertTx(target interface{}, perBatch int, tx *gorm.DB) error

	DB() *gorm.DB
	DBWithPreloads(preloads []string) *gorm.DB
	HandleError(res *gorm.DB) error
	HandleOneError(res *gorm.DB) error

	ToSQL(query *gorm.DB, target interface{}) (string, error)
}

type TransactionRepository interface {
	Repository
	CreateTx(target interface{}, tx *gorm.DB) error
	SaveTx(target interface{}, tx *gorm.DB) error
	UpdateTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) error
	UpdateOrCreateTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) error
	UpdateOrCreateTxV2(target interface{}, attributes map[string]interface{}, tx *gorm.DB) error
	UpdateOrCreateTxReturn(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) (oldData, newData interface{}, action string, err error)
	DeleteTx(target interface{}, tx *gorm.DB) error
	DeleteTxByCondition(target interface{}, condition string, tx *gorm.DB) error
	UpdateByConditionTx(Txtarget interface{}, condition string, values map[string]interface{}, tx *gorm.DB) error
	FindWhereTx(target interface{}, condition string, tx *gorm.DB) error
}
