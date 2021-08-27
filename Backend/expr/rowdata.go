package expr

type RowData struct {
	ColumnIndices map[string]int
	Data          []interface{}
}

type Invalid struct {
	Invalid string `json:"invalid"`
}

func (r *RowData) Get(key string) interface{} {
	idx, ok := r.ColumnIndices[key]
	if ok {
		return r.Data[idx]
	} else {
		return nil
	}
}

func (r *RowData) SetInvalid(key string, reason string) {
	idx, ok := r.ColumnIndices[key]
	if ok {
		r.Data[idx] = Invalid{reason}
	}
}
