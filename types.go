package GormQueryFromJson

import (
	"fmt"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

type Results struct {
	AcceptedQueries []BaseQuery `json:"acceptedQueries"`
	// NotAcceptedQueries are queries that were not accepted by the schema
	// they may be valid for another DoQuery call, but not this one (subqueries .etc.)
	NotAcceptedQueries []BaseQuery `json:"notAcceptedQueries"`
}

type JSON struct {
	QueryType TypeIdentifier  `json:"type"`
	Query     json.RawMessage `json:"query"`
}

// MaxDepth is how deep a query will go before throwing.
const MaxDepth = 10

type ColumnUsageInfo struct {
	ColumnName string
	QueryType  TypeIdentifier
}

type BaseQuery interface {
	Query(database *gorm.DB) *gorm.DB
	// ColumnUsages returns a string of each column used in the query
	ColumnUsages() []ColumnUsageInfo
	// Depth returns the depth of the query or error if it is too deep
	Depth() (int, error)
	Type() TypeIdentifier
	// GetColumnName Returns false, "" if the query is an OrQuery
	GetColumnName() (bool, string)
}

type TypeIdentifier string

const (
	TypeNumber       TypeIdentifier = "number"
	TypeString       TypeIdentifier = "string"
	TypeStrictString TypeIdentifier = "strict_string"
	TypeOr           TypeIdentifier = "or"
	TypeAnd          TypeIdentifier = "and"
)

var ErrMaxDepth = fmt.Errorf("max depth exceeded")

type AcceptableQueryTypes struct {
	// ColumnName is the name of the column in the database
	ColumnName string `json:"columnName"`
	// QueryTypes are the types of queries that can be run on this column
	QueryTypes []TypeIdentifier `json:"queryTypes"`
	// Limit max number of queries that can be run on this column
	Limit int `json:"limit"`
}
