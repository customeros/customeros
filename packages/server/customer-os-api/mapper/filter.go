package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
)

// MapFilterToCommonModel converts a GraphQL Filter to a common model Filter
func MapFilterToCommonModel(where *model.Filter) *commonmodel.Filter {
	if where == nil {
		return nil
	}

	commonWhere := &commonmodel.Filter{}
	if where.Filter != nil {
		commonWhere.Filter = &commonmodel.FilterItem{
			Property:  where.Filter.Property,
			Operation: where.Filter.Operation,
			Value: commonmodel.AnyTypeValue{
				Str:       where.Filter.Value.Str,
				Int:       where.Filter.Value.Int,
				Time:      where.Filter.Value.Time,
				Bool:      where.Filter.Value.Bool,
				Float:     where.Filter.Value.Float,
				ArrayStr:  where.Filter.Value.ArrayStr,
				ArrayInt:  where.Filter.Value.ArrayInt,
				ArrayBool: where.Filter.Value.ArrayBool,
				ArrayTime: where.Filter.Value.ArrayTime,
			},
		}
	}
	if where.And != nil {
		commonWhere.And = make([]*commonmodel.Filter, len(where.And))
		for i, f := range where.And {
			commonWhere.And[i] = &commonmodel.Filter{
				Filter: &commonmodel.FilterItem{
					Property:  f.Filter.Property,
					Operation: commonmodel.ComparisonOperator(f.Filter.Operation),
					Value: commonmodel.AnyTypeValue{
						Str:       f.Filter.Value.Str,
						Int:       f.Filter.Value.Int,
						Time:      f.Filter.Value.Time,
						Bool:      f.Filter.Value.Bool,
						Float:     f.Filter.Value.Float,
						ArrayStr:  f.Filter.Value.ArrayStr,
						ArrayInt:  f.Filter.Value.ArrayInt,
						ArrayBool: f.Filter.Value.ArrayBool,
						ArrayTime: f.Filter.Value.ArrayTime,
					},
				},
			}
		}
	}
	if where.Or != nil {
		commonWhere.Or = make([]*commonmodel.Filter, len(where.Or))
		for i, f := range where.Or {
			commonWhere.Or[i] = &commonmodel.Filter{
				Filter: &commonmodel.FilterItem{
					Property:  f.Filter.Property,
					Operation: commonmodel.ComparisonOperator(f.Filter.Operation),
					Value: commonmodel.AnyTypeValue{
						Str:       f.Filter.Value.Str,
						Int:       f.Filter.Value.Int,
						Time:      f.Filter.Value.Time,
						Bool:      f.Filter.Value.Bool,
						Float:     f.Filter.Value.Float,
						ArrayStr:  f.Filter.Value.ArrayStr,
						ArrayInt:  f.Filter.Value.ArrayInt,
						ArrayBool: f.Filter.Value.ArrayBool,
						ArrayTime: f.Filter.Value.ArrayTime,
					},
				},
			}
		}
	}
	if where.Not != nil {
		commonWhere.Not = &commonmodel.Filter{
			Filter: &commonmodel.FilterItem{
				Property:  where.Not.Filter.Property,
				Operation: commonmodel.ComparisonOperator(where.Not.Filter.Operation),
				Value: commonmodel.AnyTypeValue{
					Str:       where.Not.Filter.Value.Str,
					Int:       where.Not.Filter.Value.Int,
					Time:      where.Not.Filter.Value.Time,
					Bool:      where.Not.Filter.Value.Bool,
					Float:     where.Not.Filter.Value.Float,
					ArrayStr:  where.Not.Filter.Value.ArrayStr,
					ArrayInt:  where.Not.Filter.Value.ArrayInt,
					ArrayBool: where.Not.Filter.Value.ArrayBool,
					ArrayTime: where.Not.Filter.Value.ArrayTime,
				},
			},
		}
	}

	return commonWhere
}
