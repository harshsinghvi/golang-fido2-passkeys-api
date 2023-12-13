package helpers

import (
	"log"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/iancoleman/strcase"
)

func ReturningColumnsCalculator(db *gorm.DB, _DataEntity interface{}, config models.Config) clause.Returning {
	returningColumns := []clause.Column{}

	returningColumnNames := GetReturningColumnNames(db, _DataEntity, config)
	for _, columnName := range returningColumnNames {
		returningColumns = append(returningColumns, clause.Column{Name: columnName})
	}

	return clause.Returning{Columns: returningColumns}
}

func GetReturningColumnNames(db *gorm.DB, _DataEntity interface{}, config models.Config) []string {
	log.Print(config)
	returningColumnNames := []string{}

	if len(config.SelectFields) == 0 {
		columns, _ := db.Migrator().ColumnTypes(_DataEntity)
		for _, column := range columns {
			if !IfElementExists(config.OmitFields, column.Name()) || config.PostSkipOmit {
				returningColumnNames = append(returningColumnNames, column.Name())
			}
		}

		return returningColumnNames
	}

	for _, columnName := range config.SelectFields {
		if !IfElementExists(config.OmitFields, strcase.ToSnake(columnName)) || config.PostSkipOmit {
			returningColumnNames = append(returningColumnNames, columnName)
		}
	}

	return returningColumnNames
}

func IfElementExists(arr []string, e string) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if strcase.ToSnake(v) == e {
			return true
		}
	}
	return false
}
