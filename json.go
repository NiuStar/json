package json

import (
	"encoding/json"
	"github.com/NiuStar/reflect"
	reflect2 "reflect"
	"github.com/NiuStar/xsql3/Type"
	"fmt"
)

func Unmarshal(body []byte,v Type.IHandler) (error) {

	type_ := reflect.GetReflectType(v)

	fmt.Println("type_.Kind():",type_.Kind(),type_)

	if type_.Kind() == reflect2.Interface {
		type_ = type_.Elem()
	}
	fmt.Println("reflect.GetReflectType(v).Kind():",type_.Kind())
	if reflect.GetReflectType(v).Kind() == reflect2.Interface {
		return json.Unmarshal(body,v)
	}

	var data interface{}
	err := json.Unmarshal(body,&data)
	if err != nil {
		return err
	}

	valueType := reflect.GetReflectValue(v)

	scanMap(valueType,data,reflect2.ValueOf(v.(Type.DBOperation).TableName()))

	body1,err1 := json.Marshal(v)
	if err1 != nil {
		return err1
	}

	fmt.Println("v _result:",string(body1))
	return nil
	//

	/*body1,err1 := json.Marshal(scanMap(v,data))
	if err1 != nil {
		return err1
	}
	err = json.Unmarshal(body1,v)
	if err != nil {
		return err
	}*/


	//rv := reflect2.ValueOf(v)
	//rv.Elem().Set(reflect2.ValueOf(scanMap(data)))
	return nil
}
//客户端传参加上String等服务器数据库类型,v是转换后的结构体对象，m是客户端传来转换为MAP的对象
func scanMap(valueType reflect2.Value,m interface{},tableName reflect2.Value) bool {
	listName := make(map[string]int)

	//fmt.Println("valueType.Kind():",valueType.Kind())
	switch valueType.Kind() {
	case reflect2.Struct:
		type_ := valueType.Type()
		for i := 0 ;i < type_.NumField();i ++ {
			if len(type_.Field(i).Tag.Get("json")) > 0 {
				listName[type_.Field(i).Tag.Get("json")] = i
			}
		}

		if Type.IsTabelType(valueType.Type()) {
			if m != nil {
				valueType.Addr().MethodByName("SetValue").Call([]reflect2.Value{reflect2.ValueOf(m)})
				//fmt.Println("SetValue Over:",valueType.Addr().MethodByName("Value").Call([]reflect2.Value{})[0].Interface())
			}
			valueType.Addr().MethodByName("SetTableName").Call([]reflect2.Value{tableName})
			valueType.Addr().MethodByName("SetParent").Call([]reflect2.Value{valueType.Addr()})
			return true
		//	return valueType.MethodByName("SetValue").Call([]reflect2.Value{reflect2.ValueOf(m)})
		} else {
			for i:=0;i<valueType.NumField();i++ {

				jsonName := valueType.Type().Field(i).Name
				if []byte(jsonName)[0] < 'A' || []byte(jsonName)[0] > 'Z' {
					continue
				}
				if valueType.Type().Field(i).Tag.Get("json") != "" {
					jsonName = valueType.Type().Field(i).Tag.Get("json")
				}

				fmt.Println("valueType jsonName:",jsonName)

				if m != nil {
					jsonName := valueType.Type().Field(i).Name
					if len(valueType.Type().Field(i).Tag.Get("json")) > 0 {
						jsonName = valueType.Type().Field(i).Tag.Get("json")
					}
					if scanMap(valueType.Field(i),m.(map[string]interface{})[jsonName],tableName) {
						valueType.Field(i).FieldByName("Names").SetString(jsonName)
					}
				} else {
					if scanMap(valueType.Field(i),nil,tableName) {
						valueType.Field(i).FieldByName("Names").SetString(jsonName)
					}
				}
			}
			//fmt.Println("SetValue Over:",valueType)
		}

	case reflect2.Ptr:
		for ;reflect2.Ptr == valueType.Kind(); {
			valueType = valueType.Elem()
		}
		return scanMap(valueType,m,tableName)
	case reflect2.Interface:
		for ;reflect2.Interface == valueType.Kind(); {
			valueType = valueType.Elem()
		}
		return scanMap(valueType,m,tableName)
	default:
		return false
	}
	return false

/*
	switch type_.Kind() {


	case reflect2.Map:
		list := v.(map[string]interface{})
		for key,value := range v.(map[string]interface{}) {
			valueType.Field(listName[key]).Set(scanMap(value))
		}
		return list
	case reflect2.Array:
		var list []interface{}
		for _,value := range v.([]interface{}) {
			list = append(list,scanMap(value))
		}
		return list
	case reflect2.String:
		return map[string]interface{}{"value":v}
	case reflect2.Int:
	case reflect2.Int8:
	case reflect2.Int16:
	case reflect2.Int32:
	case reflect2.Int64:
	case reflect2.Uint:
	case reflect2.Uint8:
	case reflect2.Uint16:
	case reflect2.Uint32:
	case reflect2.Uint64:
	case reflect2.Uintptr:
		return map[string]interface{}{"value":v}
	case reflect2.Float32:
		//return xsql3.Float{Values:float64(v.(float32))}
	case reflect2.Float64:
		return map[string]interface{}{"value":v}
		//return xsql3.Float{Values:v.(float64)}
	default:
		return v
	}
	return nil*/
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {

	fmt.Println("MarshalIndent 0:",v)
	list := scanStruct(reflect2.ValueOf(v))
	fmt.Println("MarshalIndent 1:",list)
	body , err := json.MarshalIndent(list,prefix,indent)
	fmt.Println("MarshalIndent 3:",string(body))
	return body,err
}

//返回给客户端传参去掉String等服务器数据库类型
func scanStruct(v reflect2.Value) (result interface{}) {

	//fmt.Println("v:",v,v.Kind())
	switch v.Kind() {
	case reflect2.Struct:
		if Type.IsTabelType(v.Type()) {

			//不确定代码，如果没有int或string没有值，是否返回为空还是默认值？？？暂时返回为空
			if v.Addr().MethodByName("IsNil").Call([]reflect2.Value{})[0].Bool() {
				return nil
			}
			return v.Addr().MethodByName("Value").Call([]reflect2.Value{})[0].Interface()
		}
		var list = make(map[string]interface{})
		for i:=0;i<v.NumField();i++ {
			jsonName := v.Type().Field(i).Name

			if []byte(jsonName)[0] < 'A' || []byte(jsonName)[0] > 'Z' || jsonName == v.Type().Field(i).Type.Name() {
				continue
			}
			if v.Type().Field(i).Tag.Get("json") != "" {
				jsonName = v.Type().Field(i).Tag.Get("json")
			}
			list[jsonName] = scanStruct(v.Field(i))
		}
		return list
	case reflect2.Ptr:
		for ;reflect2.Ptr == v.Kind(); {
			v = v.Elem()
		}
		return scanStruct(v)
	case reflect2.Interface:
		for ;reflect2.Interface == v.Kind(); {
			v = v.Elem()
		}
		return scanStruct(v)
	default:
		return v.Interface()
	}
	fmt.Println("我也不知道错在哪了")
	return nil
}


