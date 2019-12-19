package mysql

import (
	"database/sql"
	"reflect"
	"fmt"
	"strings"
	"time"
)

func scanQueryOne(dest interface{}, rows *sql.Rows) error {
	defer rows.Close()
	if dest == nil {
		return nil
	}
	if !rows.Next() {
		return NO_DATA_TO_BIND
	}
	destType := reflect.TypeOf(dest)
	destTypeKind := destType.Kind()
	if destTypeKind != reflect.Ptr {
		panic("ptr")
	}
	destTypeElemKind := destType.Elem().Kind()
	if destTypeElemKind != reflect.Struct {
		panic("struct")
	}
	destValueElem := reflect.ValueOf(dest).Elem()
	if !destValueElem.CanSet() {
		panic("can set")
	}

	// 遍历查询结果
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columnTypes))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 新建一个元素实例
	err = rows.Scan(scanArgs...)
	if err != nil {
		return err
	}
	bindData(destValueElem, values, columnTypes)
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func token() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func scanQuery(dest interface{}, rows *sql.Rows) error {
	defer rows.Close()
	if dest == nil {
		return nil
	}
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("ptr")
	}
	listType := destType.Elem()
	if listType.Kind() != reflect.Slice {
		panic("slice")
	}
	var isPointer bool
	elemType := listType.Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		isPointer = true
	}
	if elemType.Kind() != reflect.Struct {
		panic("struct")
	}
	destValue := reflect.ValueOf(dest).Elem()
	if !destValue.CanSet() {
		return fmt.Errorf("kelp.db.mysql: dest can not set")
	}
	// 遍历查询结果
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columnTypes))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// 新建一个元素实例
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		elem := reflect.New(elemType).Elem()
		bindData(elem, values, columnTypes)
		if isPointer {
			destValue.Set(reflect.Append(destValue, elem.Addr()))
		} else {
			destValue.Set(reflect.Append(destValue, elem))
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func bindData(elem reflect.Value, values []interface{}, columnTypes []*sql.ColumnType) {
	for i, col := range values {
		key := columnTypes[i].Name()
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Type().Field(j)
			fieldName, ok := field.Tag.Lookup("column")
			if !ok {
				fieldName, ok = field.Tag.Lookup("json")
				if ok {
					fieldName = strings.Split(fieldName, ",")[0]
				}
			}
			if !ok {
				fieldName = strings.ToLower(field.Name)
			}
			if key == fieldName {
				eleField := elem.FieldByName(field.Name)
				if eleField.CanSet() {
					switch field.Type.Kind() {
					case reflect.Int:
						eleField.Set(reflect.ValueOf(ToInt(col)))
					case reflect.Int64:
						eleField.Set(reflect.ValueOf(ToInt64(col)))
					case reflect.Float64:
						eleField.Set(reflect.ValueOf(ToFloat(col)))
					case reflect.String:
						eleField.Set(reflect.ValueOf(ToString(col)))
					case reflect.Bool:
						eleField.Set(reflect.ValueOf(ToBool(col)))
					case reflect.Struct:
						switch {
						case field.Type.Name() == "Time":
							eleField.Set(reflect.ValueOf(ToTime(col)))
						}
					default:
						eleField.Set(reflect.ValueOf(col))
					}
				}
			}
		}
	}
}

func queryOne(conn connector, destObject interface{}, sql string, params []interface{}) error {
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	if err := scanQueryOne(destObject, rows); err != nil {
		return err
	}
	return nil
}

func query(conn connector, destObject interface{}, sql string, params []interface{}) error {
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	if err := scanQuery(destObject, rows); err != nil {
		return err
	}
	return nil
}

func insert(conn connector, sql string, params []interface{}) (int64, error) {
	ret, err := conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func execute(conn connector, sql string, params []interface{}) (int64, error) {
	ret, err := conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	eff, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return eff, nil
}
