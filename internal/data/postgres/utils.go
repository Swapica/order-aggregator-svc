package postgres

type sqlString string

func (s sqlString) ToSql() (string, []interface{}, error) {
	return string(s), nil, nil
}

func isNilInterface(v interface{}) bool {
	// Add new types here when you use this function
	switch v := v.(type) {
	case *string:
		return v == nil
	case *int64:
		return v == nil
	case *uint8:
		return v == nil
	}
	return v == nil
	// return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) // 7 times slower
	//	value := reflect.ValueOf(v)
	//	return v == nil || (value.Kind() == reflect.Ptr && value.IsNil()) // 5 times slower
}
