package GormQueryFromJson

import "gorm.io/gorm"

type StrictStringQuery struct {
	ColumnName string `json:"column"`
	Equals     string `json:"equals"`
}

func (query *StrictStringQuery) Type() TypeIdentifier {
	return TypeStrictString
}

func (query *StrictStringQuery) Query(database *gorm.DB) *gorm.DB {
	column := query.ColumnName

	database = database.Where(column+" = ?", query.Equals)

	return database
}

func (query *StrictStringQuery) ColumnUsages() []ColumnUsageInfo {
	return []ColumnUsageInfo{{query.ColumnName, query.Type()}}
}

func (query *StrictStringQuery) Depth() (int, error) {
	return 1, nil
}

func (query *StrictStringQuery) GetColumnName() (bool, string) {
	return true, query.ColumnName
}
