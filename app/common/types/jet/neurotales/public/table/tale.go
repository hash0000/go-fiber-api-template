//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Tale = newTaleTable("public", "tale", "")

type taleTable struct {
	postgres.Table

	// Columns
	ID                   postgres.ColumnInteger
	UserID               postgres.ColumnInteger
	Name                 postgres.ColumnString
	FileName             postgres.ColumnString
	IsPayed              postgres.ColumnBool
	TaleGenerationID     postgres.ColumnString
	CreatedAt            postgres.ColumnTimestampz
	ChildData            postgres.ColumnString
	BackgroundCharacters postgres.ColumnString
	Preferences          postgres.ColumnString
	Moral                postgres.ColumnString
	OpenAiAnswer         postgres.ColumnString
	FabulaImgToTextJSON  postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TaleTable struct {
	taleTable

	EXCLUDED taleTable
}

// AS creates new TaleTable with assigned alias
func (a TaleTable) AS(alias string) *TaleTable {
	return newTaleTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TaleTable with assigned schema name
func (a TaleTable) FromSchema(schemaName string) *TaleTable {
	return newTaleTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TaleTable with assigned table prefix
func (a TaleTable) WithPrefix(prefix string) *TaleTable {
	return newTaleTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TaleTable with assigned table suffix
func (a TaleTable) WithSuffix(suffix string) *TaleTable {
	return newTaleTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTaleTable(schemaName, tableName, alias string) *TaleTable {
	return &TaleTable{
		taleTable: newTaleTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newTaleTableImpl("", "excluded", ""),
	}
}

func newTaleTableImpl(schemaName, tableName, alias string) taleTable {
	var (
		IDColumn                   = postgres.IntegerColumn("id")
		UserIDColumn               = postgres.IntegerColumn("user_id")
		NameColumn                 = postgres.StringColumn("name")
		FileNameColumn             = postgres.StringColumn("file_name")
		IsPayedColumn              = postgres.BoolColumn("is_payed")
		TaleGenerationIDColumn     = postgres.StringColumn("tale_generation_id")
		CreatedAtColumn            = postgres.TimestampzColumn("created_at")
		ChildDataColumn            = postgres.StringColumn("child_data")
		BackgroundCharactersColumn = postgres.StringColumn("background_characters")
		PreferencesColumn          = postgres.StringColumn("preferences")
		MoralColumn                = postgres.StringColumn("moral")
		OpenAiAnswerColumn         = postgres.StringColumn("open_ai_answer")
		FabulaImgToTextJSONColumn  = postgres.StringColumn("fabula_img_to_text_json")
		allColumns                 = postgres.ColumnList{IDColumn, UserIDColumn, NameColumn, FileNameColumn, IsPayedColumn, TaleGenerationIDColumn, CreatedAtColumn, ChildDataColumn, BackgroundCharactersColumn, PreferencesColumn, MoralColumn, OpenAiAnswerColumn, FabulaImgToTextJSONColumn}
		mutableColumns             = postgres.ColumnList{UserIDColumn, NameColumn, FileNameColumn, IsPayedColumn, TaleGenerationIDColumn, CreatedAtColumn, ChildDataColumn, BackgroundCharactersColumn, PreferencesColumn, MoralColumn, OpenAiAnswerColumn, FabulaImgToTextJSONColumn}
	)

	return taleTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                   IDColumn,
		UserID:               UserIDColumn,
		Name:                 NameColumn,
		FileName:             FileNameColumn,
		IsPayed:              IsPayedColumn,
		TaleGenerationID:     TaleGenerationIDColumn,
		CreatedAt:            CreatedAtColumn,
		ChildData:            ChildDataColumn,
		BackgroundCharacters: BackgroundCharactersColumn,
		Preferences:          PreferencesColumn,
		Moral:                MoralColumn,
		OpenAiAnswer:         OpenAiAnswerColumn,
		FabulaImgToTextJSON:  FabulaImgToTextJSONColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}