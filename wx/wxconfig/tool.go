package wxconfig

type Data map[string]interface{}

// Json所有数字都是float64
func (d Data) Float64(k string) float64 {
	v, ok := d[k]
	if !ok {
		return 0.00
	}
	switch u := v.(type) {
	case float64:
		return u
	default:
		return 0.00
	}
}

func (d Data) Int64(k string) int64 {
	return int64(d.Float64(k))
}

func (d Data) Int(k string) int {
	return int(d.Float64(k))
}
