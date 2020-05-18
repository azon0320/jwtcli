package utils

import "fmt"

type Wrapper struct {
	Map map[string]interface{}
}

func (w Wrapper) GetBool(key string) bool {
	if v, ok := w.Map[key].(bool); ok {
		return v
	}
	return false
}

func (w Wrapper) GetString(key string) string {
	return fmt.Sprint(w.Map[key])
}

func (w Wrapper) GetFloat64(key string) float64 {
	if v, ok := w.Map[key].(float64); ok {
		return v
	}
	return 0
}

func (w Wrapper) GetInt(key string) int {
	if v, ok := w.Map[key].(int); ok {
		return v
	}
	return 0
}

func WrapMap(m map[string]interface{}) *Wrapper {
	return &Wrapper{Map: m}
}
