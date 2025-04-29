package models

import (
	"fmt"
	"reflect"

	"go-schedule/libs/types"

	"github.com/go-mysql-org/go-mysql/schema"
	"github.com/spf13/cast"
)

var (
	tableCache = make(map[string]*schema.Table, 0)
)

type BinlogMySQL struct{}

func NewBinlogMySQL() *BinlogMySQL {
	return &BinlogMySQL{}
}

func (m *BinlogMySQL) GetTableSchema(where, schemaName, tableName string) ([]types.TableColumnRows, error) {
	db := tools.GetMySQLClient(fmt.Sprintf("%s.slave", where))

	result := make([]types.TableColumnRows, 0)

	query := "SELECT COLUMN_NAME, DATA_TYPE, COLLATION_NAME, EXTRA FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION"
	rows, err := db.Query(query, schemaName, tableName)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var (
			columnName    string
			dataType      string
			collationName string
			extra         string
		)

		rows.Scan(
			&columnName,
			&dataType,
			&collationName,
			&extra,
		)

		result = append(result, types.TableColumnRows{
			Name:      columnName,
			Type:      dataType,
			Collation: collationName,
			Extra:     extra,
		})
	}

	defer rows.Close()

	return result, nil
}

func (m *BinlogMySQL) GetTableSchemaFromCache(where, schemaName, tableName string) (*schema.Table, error) {
	fullName := fmt.Sprintf("%s.%s", schemaName, tableName)

	if result, ok := tableCache[fullName]; ok {
		return result, nil
	}

	result := &schema.Table{
		Schema: schemaName,
		Name:   tableName,
	}

	data, _ := m.GetTableSchema(where, schemaName, tableName)

	for _, v := range data {
		result.AddColumn(v.Name, v.Type, v.Collation, v.Extra)
	}
	// fmt.Println(string(units.MarshalJson(result)))

	tableCache[fullName] = result

	return result, nil
}

// Cannot be used directly, the definition of each table is different
func (m *BinlogMySQL) GetBinlogFieldAndValeue(tableSchema *schema.Table, row []any) map[string]string {
	result := make(map[string]string, 0)

	for i, col := range tableSchema.Columns {
		value := row[i]
		t := reflect.TypeOf(value)
		// fmt.Printf("Type name: %v\n", t)

		switch col.Type {
		case schema.TYPE_NUMBER, schema.TYPE_MEDIUM_INT: // tinyint, smallint, int, bigint, year
			result[col.Name] = fmt.Sprintf("%v", value)
		case schema.TYPE_FLOAT: // float, double
			if v, ok := value.(float64); ok {
				result[col.Name] = cast.ToString(v)
			} else {
				result[col.Name] = fmt.Sprintf("%v", value)
			}
		case schema.TYPE_ENUM: // ENUM
			enumValues := []string{"", "0", "1"}

			if index, ok := value.(int64); ok && int(index) < len(enumValues) {
				result[col.Name] = enumValues[index]
			}
		// case schema.TYPE_SET: // SET - unused
		// 	if v, ok := value.(uint64); ok {
		// 		result[col.Name] = m.ParseSetValue(v, []string{"0", "1"})
		// 	}
		case schema.TYPE_STRING: // char, varchar, text
			if t == nil {
				result[col.Name] = ""
			} else if v, ok := value.(string); ok {
				result[col.Name] = v
			} else {
				result[col.Name] = fmt.Sprintf("%s", value)
			}
		case schema.TYPE_DATETIME, schema.TYPE_TIMESTAMP, schema.TYPE_DATE, schema.TYPE_TIME:
			if t, ok := value.(string); ok {
				result[col.Name] = t
			} else {
				result[col.Name] = fmt.Sprintf("%v", value)
			}
		case schema.TYPE_BIT: // byte
			if b, ok := value.([]byte); ok {
				result[col.Name] = string(b)
			}
		case schema.TYPE_JSON: // json
			if jsonBytes, ok := value.([]byte); ok {
				result[col.Name] = string(jsonBytes)
			} else {
				result[col.Name] = fmt.Sprintf("%s", value)
			}
		case schema.TYPE_DECIMAL: // DECIMAL
			if v, ok := value.(string); ok {
				result[col.Name] = v
			} else {
				result[col.Name] = fmt.Sprintf("%s", value)
			}
		case schema.TYPE_BINARY: // BINARY, VARBINARY
			if binaryData, ok := value.([]byte); ok {
				result[col.Name] = fmt.Sprintf("0x%x", binaryData)
			}
		case schema.TYPE_POINT:
			if pointData, ok := value.([]byte); ok {
				result[col.Name] = m.ParsePointValue(pointData)
			}
		default:
			fmt.Println("unknown type", col.Type, value)
		}
	}

	return result
}

func (m *BinlogMySQL) ParseSetValue(value uint64, setItems []string) []string {
	var result []string
	for i := uint(0); i < uint(len(setItems)); i++ {
		if (value & (1 << i)) != 0 {
			result = append(result, setItems[i])
		}
	}
	return result
}

func (m *BinlogMySQL) ParsePointValue(data []byte) string {
	if len(data) < 9 {
		return "Invalid POINT"
	}
	x := float64(int32(data[1]) | int32(data[2])<<8 | int32(data[3])<<16 | int32(data[4])<<24)
	y := float64(int32(data[5]) | int32(data[6])<<8 | int32(data[7])<<16 | int32(data[8])<<24)
	return fmt.Sprintf("POINT(%f, %f)", x, y)
}
