package helpers

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func IfElementExists(arr []string, e string) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

func ReturningColumnsCalculator(db *gorm.DB, _DataEntity interface{}, _SelectFields, _OmitFields []string, args ...interface{}) clause.Returning {
	var _OverrideOmit bool = false

	if len(args) != 0 {
		_OverrideOmit = args[0].(bool)
	}

	returningColumnNames := []string{}

	if len(_SelectFields) == 0 {
		columns, _ := db.Migrator().ColumnTypes(_DataEntity)
		for _, column := range columns {
			if !IfElementExists(_OmitFields, column.Name()) {
				returningColumnNames = append(returningColumnNames, column.Name())
			}
		}
	} else if _OverrideOmit {
		returningColumnNames = _SelectFields
	} else {
		for _, columnName := range _SelectFields {
			if !IfElementExists(_OmitFields, columnName) {
				returningColumnNames = append(returningColumnNames, columnName)
			}
		}
	}

	returningColumns := []clause.Column{}
	for _, columnName := range returningColumnNames {
		returningColumns = append(returningColumns, clause.Column{Name: columnName})
	}
	return clause.Returning{Columns: returningColumns}
}
