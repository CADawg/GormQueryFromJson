package GormQueryFromJson

import "gorm.io/gorm"

type NumberQuery struct {
	ColumnName         string `json:"column"`
	Equals             *int   `json:"equals,omitempty"`
	NotEquals          *int   `json:"notEquals,omitempty"`
	GreaterThan        *int   `json:"greaterThan,omitempty"`
	GreaterThanOrEqual *int   `json:"greaterThanOrEqual,omitempty"`
	LessThan           *int   `json:"lessThan,omitempty"`
	LessThanOrEqual    *int   `json:"lessThanOrEqual,omitempty"`
}

func (query *NumberQuery) Query(database *gorm.DB) *gorm.DB {
	_, column := query.GetColumnName()

	if query.Equals != nil {
		database = database.Where(column+" = ?", query.Equals)
	}
	if query.NotEquals != nil {
		database = database.Where(column+" != ?", query.NotEquals)
	}
	if query.GreaterThan != nil {
		database = database.Where(column+" > ?", query.GreaterThan)
	}
	if query.GreaterThanOrEqual != nil {
		database = database.Where(column+" >= ?", query.GreaterThanOrEqual)
	}
	if query.LessThan != nil {
		database = database.Where(column+" < ?", query.LessThan)
	}
	if query.LessThanOrEqual != nil {
		database = database.Where(column+" <= ?", query.LessThanOrEqual)
	}

	return database
}

func (query *NumberQuery) ColumnUsages() []ColumnUsageInfo {
	return []ColumnUsageInfo{{query.ColumnName, query.Type()}}
}

func (query *NumberQuery) Depth() (int, error) {
	return 1, nil
}

func (query *NumberQuery) Type() TypeIdentifier {
	return TypeNumber
}

func (query *NumberQuery) GetColumnName() (bool, string) {
	return true, query.ColumnName
}
