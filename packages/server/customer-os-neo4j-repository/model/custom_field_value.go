package model

import "time"

type CustomFieldValue struct {
	Str     *string    `json:"string,omitempty"`
	Int     *int64     `json:"int,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
	Bool    *bool      `json:"bool,omitempty"`
	Decimal *float64   `json:"decimal,omitempty"`
}

func (c *CustomFieldValue) RealValue() any {
	if c.Int != nil {
		return *c.Int
	} else if c.Decimal != nil {
		return *c.Decimal
	} else if c.Time != nil {
		return *c.Time
	} else if c.Bool != nil {
		return *c.Bool
	} else if c.Str != nil {
		return *c.Str
	}
	return nil
}
