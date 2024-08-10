package GormQueryFromJson

import (
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

type AndQuery struct {
	Queries []BaseQuery `json:"queries"`
}

func (query *AndQuery) UnmarshalJSON(data []byte) error {
	var dataJson []JSON

	err := json.Unmarshal(data, &dataJson)

	if err != nil {
		return err
	}

	var queries []BaseQuery

	for _, d := range dataJson {
		q := GetQueryType(d.QueryType)

		err := json.Unmarshal(d.Query, q)

		if err != nil {
			return err
		}

		queries = append(queries, q)
	}

	query.Queries = queries

	return nil
}

func (query *AndQuery) Query(database *gorm.DB) *gorm.DB {
	if len(query.Queries) == 0 {
		return database.Session(&gorm.Session{})
	}

	return database.Where(query.buildAndCondition(database))
}

func (query *AndQuery) buildAndCondition(database *gorm.DB) *gorm.DB {
	var conditions []*gorm.DB

	for _, q := range query.Queries {
		session := database.Session(&gorm.Session{})

		conditions = append(conditions, q.Query(session))
	}

	andCondition := conditions[0] //nolint:nilaway

	for i := 1; i < len(conditions); i++ {
		andCondition.Where(conditions[i])
	}

	return andCondition
}

func (query *AndQuery) Type() TypeIdentifier {
	return TypeAnd
}

func (query *AndQuery) ColumnUsages() []ColumnUsageInfo {
	var columnUsages []ColumnUsageInfo

	for _, q := range query.Queries {
		columnUsages = append(columnUsages, q.ColumnUsages()...)
	}

	return columnUsages
}

func (query *AndQuery) Depth() (int, error) {
	maxDepth := 0

	for _, q := range query.Queries {
		qDepth, err := q.Depth()

		if err != nil {
			return 0, err
		}

		if qDepth >= maxDepth {
			return 0, ErrMaxDepth
		}
	}

	// Remember this is another layer
	return maxDepth + 1, nil
}

func (query *AndQuery) GetColumnName() (bool, string) {
	return false, ""
}
