package BeeSQLite

// Row - container for fields in one row of the table
type Row map[string]Field

// Result - container for all rows fetched from database
type Result []Row

func (r *Result) Count() int {
	return len(*r)
}

func (r *Result) IsNotEmpty() bool {
	if r.Count() > 0 {
		return true
	}
	return false
}

func (r *Result) First() Row {
	if r.IsNotEmpty() {
		return (*r)[0]
	}
	return nil
}

func (r *Result) Last() Row {
	n := r.Count()
	if n > 0 {
		return (*r)[n-1]
	}
	return nil
}
