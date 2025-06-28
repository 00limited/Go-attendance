package gormrepository

import (
	"errors"
	"fmt"

	"github.com/yourname/payslip-system/internal/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	db           *gorm.DB
	defaultJoins []string
}

func NewGormRepository(db *gorm.DB, defaultJoins ...string) *gormRepository {
	return &gormRepository{db: db, defaultJoins: defaultJoins}
}

func (r *gormRepository) DB() *gorm.DB {
	return r.DBWithPreloads(nil)
}

func (r *gormRepository) DBWithPreloads(preloads []string) *gorm.DB {
	conn := r.db

	for _, join := range r.defaultJoins {
		conn = conn.Joins(join)
	}

	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	return conn
}

func (r *gormRepository) FindByRaw(target interface{}, query string) error {
	res := r.DB().Raw(query).
		Scan(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindAll(target interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindAllAndLimit(target interface{}, limit int, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Limit(limit).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindBatch(target interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error) {
	query := r.DBWithPreloads(preloads).Model(target)

	if search != "" {
		query = query.Where(search)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if groupBy != "" {
		query = query.Group(groupBy)
	}

	query.Count(&count)

	res := query.
		Limit(limit).
		Offset(offset).
		Find(target)

	return count, r.HandleError(res)
}

func (r *gormRepository) FindWhere(target interface{}, condition string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindWhereBatch(target interface{}, condition string, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error) {
	query := r.DBWithPreloads(preloads).Model(target)

	if search != "" {
		query = query.Where(search)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if groupBy != "" {
		query = query.Group(groupBy)
	}

	query = query.Where(condition)

	query.Count(&count)

	res := query.Limit(limit).
		Offset(offset).
		Find(target)

	return count, r.HandleError(res)
}

func (r *gormRepository) FindByField(target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%s = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindOneByFieldOrder(target interface{}, field string, value interface{}, order string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%s = ?", field), value).
		Order(order).
		First(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindByFields(target interface{}, fields map[string]interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fields).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindByFieldsOrder(target interface{}, fields map[string]interface{}, order string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fields).
		Order(order).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error) {
	query := r.DBWithPreloads(preloads).Model(target)

	if search != "" {
		query = query.Where(search)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if groupBy != "" {
		query = query.Group(groupBy)
	}

	query = query.Where(fmt.Sprintf("%s = ?", field), value)

	query.Count(&count)

	res := query.
		Limit(limit).
		Offset(offset).
		Find(target)

	return count, r.HandleError(res)
}

func (r *gormRepository) FindByFieldsBatch(target interface{}, fields map[string]interface{}, limit, offset int, search, orderBy, groupBy string, preloads ...string) (count int64, err error) {
	query := r.DBWithPreloads(preloads).Model(target)

	if search != "" {
		query = query.Where(search)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if groupBy != "" {
		query = query.Group(groupBy)
	}

	query = query.Where(fields)
	query.Count(&count)

	res := query.
		Limit(limit).
		Offset(offset).
		Find(target)

	return count, r.HandleError(res)
}

func (r *gormRepository) FindOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%s = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) FindOneLastByField(target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%s = ?", field), value).
		Last(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) FindOneByFields(target interface{}, fields map[string]interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fields).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) FindOneByID(target interface{}, id interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) FindOneByCondition(target interface{}, condition string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(condition).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) FindOneLastByCondition(target interface{}, condition string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(condition).
		Order("id desc").
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) Create(target interface{}) error {
	res := r.DB().Create(target)

	return r.HandleError(res)
}

func (r *gormRepository) CreateTx(target interface{}, tx *gorm.DB) error {
	res := tx.Create(target)

	return r.HandleError(res)
}

func (r *gormRepository) Save(target interface{}) error {
	res := r.DB().Save(target)

	return r.HandleError(res)
}

func (r *gormRepository) SaveTx(target interface{}, tx *gorm.DB) error {
	res := tx.Save(target)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) error {
	res := tx.Model(target).Where(attributes).Updates(values)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateNoTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}) error {
	res := r.DB().Model(target).Where(attributes).Updates(values)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateOrCreateTx(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) error {
	res := tx.Where(attributes).Assign(values).FirstOrCreate(target)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateOrCreateTxV2(target interface{}, attributes map[string]interface{}, tx *gorm.DB) error {
	res := tx.Where(attributes).Assign(attributes).FirstOrCreate(target)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateOrCreateTxReturn(target interface{}, attributes map[string]interface{}, values map[string]interface{}, tx *gorm.DB) (oldData, newData interface{}, action string, err error) {
	// find first
	res := tx.Where(attributes).First(target)

	// if not found, create new
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		res = tx.Where(attributes).Assign(values).FirstOrCreate(target)
		return nil, target, "create", r.HandleError(res)
	}

	// if found, update
	oldData = helper.InterfaceToMapStringInterface(target)

	res = tx.Model(target).Where(attributes).Updates(values)
	newData = target

	return oldData, newData, "update", r.HandleError(res)
}

func (r *gormRepository) Delete(target interface{}) error {
	res := r.DB().Delete(target)

	return r.HandleError(res)
}

func (r *gormRepository) DeleteByCondition(target interface{}, condition string) error {
	res := r.DB().Where(condition).Delete(target)

	return r.HandleError(res)
}

func (r *gormRepository) DeleteTx(target interface{}, tx *gorm.DB) error {
	res := tx.Delete(target)

	return r.HandleError(res)
}

func (r *gormRepository) DeleteTxByCondition(target interface{}, condition string, tx *gorm.DB) error {
	res := tx.Where(condition).Delete(target)

	return r.HandleError(res)
}

func (r *gormRepository) BatchInsertTx(target interface{}, perBatch int, tx *gorm.DB) error {
	if perBatch <= 0 {
		perBatch = 100
	}

	res := tx.CreateInBatches(target, perBatch)

	return r.HandleError(res)
}

// handle error
func (r *gormRepository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("error: %w", res.Error)

		return err
	}

	return nil
}

func (r *gormRepository) HandleOneError(res *gorm.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return ErrRecordNotFound
	}

	return nil
}

func (r *gormRepository) ToSQL(query *gorm.DB, target interface{}) (string, error) {
	sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scan(target)
	})

	return sql, nil
}

// add function update by condition
func (r *gormRepository) UpdateByConditionTx(target interface{}, condition string, values map[string]interface{}, tx *gorm.DB) error {
	res := tx.Model(target).Where(condition).Updates(values)

	return r.HandleError(res)
}

// func for lock find by condition
func (r *gormRepository) FindByConditionThenLock(target interface{}, condition string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).Clauses(clause.Locking{Strength: "UPDATE"}).Where(condition).Find(target)

	return r.HandleError(res)
}

// func find where using transaction
func (r *gormRepository) FindWhereTx(target interface{}, condition string, tx *gorm.DB) error {
	res := tx.Where(condition).Find(target)

	return r.HandleError(res)
}
