package GormQueryFromJson

import (
	"fmt"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

type OrQuery struct {
	Queries []BaseQuery `json:"queries"`
}

func (query *OrQuery) UnmarshalJSON(data []byte) error {
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

func (query *OrQuery) Query(database *gorm.DB) *gorm.DB {
	if len(query.Queries) == 0 {
		return database.Session(&gorm.Session{})
	}

	return database.Where(query.buildOrCondition(database))
}

func (query *OrQuery) buildOrCondition(database *gorm.DB) *gorm.DB {
	var conditions []*gorm.DB

	for _, q := range query.Queries {
		session := database.Session(&gorm.Session{})

		conditions = append(conditions, q.Query(session))
	}

	orCondition := conditions[0] //nolint:nilaway
	for i := 1; i < len(conditions); i++ {
		orCondition = orCondition.Or(conditions[i])
	}

	return orCondition
}

func (query *OrQuery) Type() TypeIdentifier {
	return TypeOr
}

func (query *OrQuery) ColumnUsages() []ColumnUsageInfo {
	var columnUsages []ColumnUsageInfo

	for _, q := range query.Queries {
		columnUsages = append(columnUsages, q.ColumnUsages()...)
	}

	return columnUsages
}

var ErrMaxDepth = fmt.Errorf("max depth exceeded")

func (query *OrQuery) Depth() (int, error) {
	maxDepth := 0

	for _, q := range query.Queries {
		qDepth, err := q.Depth()
		if err != nil {
			return 0, err
		}

		if qDepth > maxDepth {
			maxDepth = qDepth
		}

		if maxDepth > MaxDepth {
			return 0, ErrMaxDepth
		}
	}

	return maxDepth, nil
}

func (query *OrQuery) GetColumnName() (bool, string) {
	return false, ""
}
