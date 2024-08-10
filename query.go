package GormQueryFromJson

import (
	"encoding/json"
	"gorm.io/gorm"
)

func GetQueryType(id TypeIdentifier) BaseQuery {
	switch id {
	case TypeNumber:
		return &NumberQuery{}
	case TypeString:
		return &StringQuery{}
	case TypeStrictString:
		return &StrictStringQuery{}
	case TypeOr:
		return &OrQuery{}
	case TypeAnd:
		return &AndQuery{}
	default:
		return nil
	}
}

func DoQuery(queries []JSON, database *gorm.DB, schema []AcceptableQueryTypes) (*gorm.DB, Results, error) {
	var goodQueries []BaseQuery
	var badQueries []BaseQuery

	var columnUsageCounts = make(map[string]int)

	for _, query := range queries {
		thisQuery := GetQueryType(query.QueryType)

		if thisQuery == nil {
			continue
		}

		// deserialise
		err := json.Unmarshal(query.Query, thisQuery)

		if err != nil {
			return database, Results{}, err
		}

		// check we're happy with the query
		usages := thisQuery.ColumnUsages()

		var passed = true

		for _, usage := range usages {
			if _, ok := columnUsageCounts[usage.ColumnName]; !ok {
				columnUsageCounts[usage.ColumnName] = 0
			}

			var haveMatched = false

			for _, col := range schema {
				if col.ColumnName == usage.ColumnName {
					haveMatched = true

					if columnUsageCounts[usage.ColumnName]+1 > col.Limit {
						passed = false
					}
					break
				}
			}

			if !haveMatched {
				passed = false
				break
			}

			columnUsageCounts[usage.ColumnName] += 1
		}

		if !passed {
			badQueries = append(badQueries, thisQuery)
			continue
		}

		goodQueries = append(goodQueries, thisQuery)
	}

	// if we don't keep this separate, it messes up the query generation with or queries
	var dbQueries []*gorm.DB

	for _, query := range goodQueries {
		session := database.Session(&gorm.Session{})
		dbQueries = append(dbQueries, query.Query(session))

		if database.Error != nil {
			return database, Results{}, database.Error
		}
	}

	for _, dbQuery := range dbQueries {
		database.Where(dbQuery)
	}

	return database, Results{
		AcceptedQueries:    goodQueries,
		NotAcceptedQueries: badQueries,
	}, nil
}
