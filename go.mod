module github.com/CADawg/GormQueryFromJson

go 1.22

// v0.1.1 had an issue with parsing "and" queries, and so is not recommended
retract v0.1.1

require (
	github.com/goccy/go-json v0.10.3
	gorm.io/gorm v1.25.11
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.17.0 // indirect
)
