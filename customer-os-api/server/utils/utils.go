package utils

func ToIntPtr(v *int64) *int {
	if v == nil {
		return nil
	}
	var output = int(*v)
	return &output
}

func ToInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}
	var output = int64(*v)
	return &output
}
