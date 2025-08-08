package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

// Получение переменных окружения
func GetConfigsPath(configs []any) error {
	viper.AutomaticEnv()
	for _, cfg := range configs {
		val := reflect.ValueOf(cfg).Elem()
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			// Получаем тег `mapstructure`
			envKey := fieldType.Tag.Get("mapstructure")
			if envKey == "" {
				continue
			}

			envValue := viper.Get(envKey)
			if envValue == nil {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(viper.GetString(envKey))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInt(viper.GetInt64(envKey))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				field.SetUint(viper.GetUint64(envKey))
			case reflect.Float32, reflect.Float64:
				field.SetFloat(viper.GetFloat64(envKey))
			case reflect.Bool:
				field.SetBool(viper.GetBool(envKey))
			default:
				return fmt.Errorf("type `%v` not support env `%s`", field.Kind().String(), envKey)
			}
		}
	}

	return nil
}
