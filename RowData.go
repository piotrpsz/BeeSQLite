package BeeSQLite

type Row struct {
	data map[string]Field
}

func (r *Row) IsEmpty() bool {
	if len(r.data) > 0 {
		return false
	}
	return true
}

func (r *Row) IsNotEmpty() bool {
	if len(r.data) > 0 {
		return true
	}
	return false
}
