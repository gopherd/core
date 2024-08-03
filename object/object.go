package object

type Object map[string]any

func (o Object) Get(key string) any {
	return o[key]
}

func (o Object) Set(key string, value any) {
	o[key] = value
}

func (o Object) Delete(key string) {
	delete(o, key)
}

func (o Object) Keys() []string {
	keys := make([]string, 0, len(o))
	for key := range o {
		keys = append(keys, key)
	}
	return keys
}

func (o Object) Values() []any {
	values := make([]any, 0, len(o))
	for _, value := range o {
		values = append(values, value)
	}
	return values
}

func (o Object) Len() int {
	return len(o)
}

func (o Object) Copy() Object {
	clone := make(Object, len(o))
	for key, value := range o {
		clone[key] = value
	}
	return clone
}

func (o Object) Merge(other Object) {
	for key, value := range other {
		o[key] = value
	}
}

func (o Object) Contains(key string) bool {
	_, ok := o[key]
	return ok
}

func (o Object) Clear() {
	for key := range o {
		delete(o, key)
	}
}

func (o Object) Int(key string, def int) int {
	v, ok := o[key]
	if !ok {
		return def
	}
	return v.(int)
}

func (o Object) Int64(key string, def int64) int64 {
	v, ok := o[key]
	if !ok {
		return def
	}
	return v.(int64)
}

func (o Object) Float64(key string, def float64) float64 {
	v, ok := o[key]
	if !ok {
		return def
	}
	if f, ok := v.(float32); ok {
		return float64(f)
	}
	if i, ok := v.(int); ok {
		return float64(i)
	}
	if i, ok := v.(int64); ok {
		return float64(i)
	}
	return v.(float64)
}

func (o Object) Bool(key string, def bool) bool {
	v, ok := o[key]
	if !ok {
		return def
	}
	return v.(bool)
}

func (o Object) String(key, def string) string {
	v, ok := o[key]
	if !ok {
		return def
	}
	if bytes, ok := v.([]byte); ok {
		return string(bytes)
	}
	return v.(string)
}

func (o Object) Bytes(key string, def []byte) []byte {
	v, ok := o[key]
	if !ok {
		return def
	}
	if s, ok := v.(string); ok {
		return []byte(s)
	}
	return v.([]byte)
}

func (o Object) Object(key string, def Object) Object {
	v, ok := o[key]
	if !ok {
		return def
	}
	return v.(Object)
}

func (o Object) Array(key string, def []any) []any {
	v, ok := o[key]
	if !ok {
		return def
	}
	return v.([]any)
}

func (o Object) RequiredInt(key string) int {
	return o[key].(int)
}

func (o Object) RequiredInt64(key string) int64 {
	return o[key].(int64)
}

func (o Object) RequiredFloat64(key string) float64 {
	v := o[key]
	if f, ok := v.(float32); ok {
		return float64(f)
	}
	if i, ok := v.(int); ok {
		return float64(i)
	}
	if i, ok := v.(int64); ok {
		return float64(i)
	}
	return v.(float64)
}

func (o Object) RequiredBool(key string) bool {
	return o[key].(bool)
}

func (o Object) RequiredString(key string) string {
	v := o[key]
	if bytes, ok := v.([]byte); ok {
		return string(bytes)
	}
	return v.(string)
}

func (o Object) RequiredBytes(key string) []byte {
	v := o[key]
	if s, ok := v.(string); ok {
		return []byte(s)
	}
	return v.([]byte)
}

func (o Object) RequiredObject(key string) Object {
	return o[key].(Object)
}

func (o Object) RequiredArray(key string) []any {
	return o[key].([]any)
}
