package GormQueryFromJson

import "gorm.io/gorm"

type StringQuery struct {
	ColumnName string  `json:"column"`
	Equals     *string `json:"equals,omitempty"`
	Contains   *string `json:"contains,omitempty"`
}

func (query *StringQuery) Query(database *gorm.DB) *gorm.DB {
	column := query.ColumnName

	if query.Equals != nil {
		database = database.Where(column+" = ?", query.Equals)
	}
	if query.Contains != nil {
		database = database.Where(column+" LIKE ?", "%"+*query.Contains+"%")
	}

	return database
}

func (query *StringQuery) Type() TypeIdentifier {
	return TypeString
}

func (query *StringQuery) ColumnUsages() []ColumnUsageInfo {
	return []ColumnUsageInfo{{query.ColumnName, query.Type()}}
}

func (query *StringQuery) Depth() (int, error) {
	return 1, nil
}

func (query *StringQuery) GetColumnName() (bool, string) {
	return true, query.ColumnName
}
