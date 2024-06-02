package service

import "database/sql"

type TableColumn struct {
	Name                   string        `json:"name"`
	DataType               string        `json:"dataType"`
	CharacterMaximumLength sql.NullInt64 `json:"characterMaximumLength"`
	IsNullable             string        `json:"isNullable"`
}

type TableData struct {
	Columns []string
	Rows    [][]interface{}
}

type TableKey struct {
	ConstraintName     string `json:"constraintName"`
	ConstraintType     string `json:"constraintType"`
	ColumnName         string `json:"columnName"`
	ForeignTableSchema string `json:"foreignTableSchema"`
	ForeignTableName   string `json:"foreignTableName"`
	ForeignColumnName  string `json:"foreignColumnName"`
}

type ForeignKey struct {
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

type InsertRows struct {
	TableName string          `json:"TableName"`
	Columns   []string        `json:"Columns"`
	Rows      [][]interface{} `json:"Rows"`
}

type UpdateRow struct {
	Сolumns []string      `json:"columns"`
	OldRow  []interface{} `json:"oldRow"`
	NewRow  []interface{} `json:"newRow"`
}

type Row struct {
	Row map[string]interface{} `json:"row"`
}

type CreatedTableData struct {
	TableName string          `json:"tableName"`
	Columns   []CreatedColumn `json:"columns"`
}

type CreatedColumn struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Primary    bool              `json:"primary"`
	NotNull    bool              `json:"notnull"`
	ForeignKey CreatedForeignKey `json:"foreignKey"`
}

type CreatedForeignKey struct {
	ForeignTable  string `json:"foreignTable"`
	ForeignColumn string `json:"foreignColumn"`
}

type User struct {
	Username string `json:"username"`
}

type Role struct {
	Rolename string `json:"rolename"`
}

type UserRoles struct {
	Roles []Role
	Users []User
}

type Filters struct {
	ColumnName  string
	FilterValue string
	Limit       int
	Offset      int
}
