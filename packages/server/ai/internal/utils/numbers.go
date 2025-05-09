package utils

// Helper function to safely get float32 pointer
func Float32Ptr(f float64) *float32 {
	v := float32(f)
	return &v
}

// Helper function to safely get int32 pointer
func Int32Ptr(i int) *int32 {
	v := int32(i)
	return &v
}
