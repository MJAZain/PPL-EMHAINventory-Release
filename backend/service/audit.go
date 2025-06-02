package service

import (
	"encoding/json"
	"fmt"
	"go-gin-auth/config"
	"go-gin-auth/model"
	"reflect"
	"time"
)

func LogAudit(tableName string, recordID string, action string, changedBy string, before interface{}, after interface{}, Description string) error {
	// Bandingkan data sebelum dan sesudah
	var beforeChanged, afterChanged map[string]interface{}

	if before != nil {
		// Bandingkan struct jika ada before-nya
		beforeChanged, afterChanged = compareStructs(before, after)

		// Jika tidak ada perubahan, tidak perlu buat audit log
		if len(beforeChanged) == 0 && len(afterChanged) == 0 {
			return nil
		}
	} else {
		// Jika tidak ada before (misalnya saat INSERT), gunakan data after saja
		afterChanged = structToMap(after)
	}

	// beforeChanged, afterChanged := compareStructs(before, after)

	// // Jika tidak ada perubahan, tidak perlu log
	// if len(beforeChanged) == 0 && len(afterChanged) == 0 {
	// 	return nil
	// }

	// beforeJSON, _ := json.Marshal(before)
	// afterJSON, _ := json.Marshal(after)
	beforeJSON, _ := json.Marshal(beforeChanged)
	afterJSON, _ := json.Marshal(afterChanged)

	audit := model.AuditLog{
		TableName:   tableName,
		Description: Description,
		RecordID:    recordID,
		Action:      action,
		ChangedBy:   changedBy,
		ChangedAt:   time.Now(),
		BeforeData:  string(beforeJSON),
		AfterData:   string(afterJSON),
	}

	return config.DB.Create(&audit).Error
}
func compareStructs(before interface{}, after interface{}) (map[string]interface{}, map[string]interface{}) {
	beforeChanged := make(map[string]interface{})
	afterChanged := make(map[string]interface{})

	// beforeJSON, _ := json.MarshalIndent(before, "", "  ")
	// afterJSON, _ := json.MarshalIndent(after, "", "  ")

	// fmt.Println("Before Changed JSON:\n", string(beforeJSON))
	// fmt.Println("After Changed JSON:\n", string(afterJSON))

	beforeVal := reflect.ValueOf(before)
	afterVal := reflect.ValueOf(after)

	fmt.Println("Before Value:", beforeVal)
	fmt.Println("After Value:", afterVal)

	// beforeVal := reflect.ValueOf(beforeJSON)
	// afterVal := reflect.ValueOf(afterJSON)

	// Pastikan before dan after adalah pointer ke struct
	if beforeVal.Kind() == reflect.Ptr {
		beforeVal = beforeVal.Elem()
	}
	if afterVal.Kind() == reflect.Ptr {
		afterVal = afterVal.Elem()
	}

	for i := 0; i < beforeVal.NumField(); i++ {
		field := beforeVal.Type().Field(i)

		// Lewati jika field tidak bisa diakses (private/unexported)
		if !beforeVal.Field(i).CanInterface() || !afterVal.Field(i).CanInterface() {
			continue
		}

		// Ambil nama field dari tag JSON jika ada
		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		beforeField := beforeVal.Field(i).Interface()
		afterField := afterVal.Field(i).Interface()

		if !reflect.DeepEqual(beforeField, afterField) {
			beforeChanged[fieldName] = beforeField
			afterChanged[fieldName] = afterField
		}
	}

	return beforeChanged, afterChanged
}
func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		result[field.Name] = value
	}

	return result
}

func compareStructsOld(before interface{}, after interface{}) (map[string]interface{}, map[string]interface{}) {
	beforeChanged := make(map[string]interface{})
	afterChanged := make(map[string]interface{})

	beforeVal := reflect.ValueOf(before)
	afterVal := reflect.ValueOf(after)

	// Pastikan sebelum dan sesudah adalah pointer ke struct
	if beforeVal.Kind() == reflect.Ptr {
		beforeVal = beforeVal.Elem()
	}
	if afterVal.Kind() == reflect.Ptr {
		afterVal = afterVal.Elem()
	}

	for i := 0; i < beforeVal.NumField(); i++ {
		field := beforeVal.Type().Field(i)
		fieldName := field.Name
		beforeField := beforeVal.Field(i).Interface()
		afterField := afterVal.Field(i).Interface()

		// Bandingkan nilai before dan after
		if !reflect.DeepEqual(beforeField, afterField) {
			beforeChanged[fieldName] = beforeField
			afterChanged[fieldName] = afterField
		}
	}

	return beforeChanged, afterChanged
}
