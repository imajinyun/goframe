package orm

import (
	"context"

	"github.com/jianfengye/collection"
	"gorm.io/gorm"

	"github.com/imajinyun/goframe/contract"
)

func (s *OrmService) GetTables(ctx context.Context, db *gorm.DB) ([]string, error) {
	return db.Migrator().GetTables()
}

func (s *OrmService) HasTable(ctx context.Context, db *gorm.DB, table string) (bool, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return false, err
	}
	collects := collection.NewStrCollection(tables)

	return collects.Contains(table), nil
}

func (s *OrmService) GetTableColumns(ctx context.Context, db *gorm.DB, table string) ([]contract.TableColumn, error) {
	var columns []contract.TableColumn
	result := db.Raw("SHOW COLUMNS FROM " + table).Scan(&columns)
	if result.Error != nil {
		return nil, result.Error
	}

	return columns, nil
}
