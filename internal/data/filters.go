//Filename: interla/data/filter.go

package data

import "alysianorales.net/TODO/internal/validator"

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	SortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	//Check page and page_size parameters
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be maximum of 1000")
	v.Check(f.Page > 0, "page_size", "must be greater than zero")
	v.Check(f.Page <= 100, "page_size", "must be maximum of 100")
	//Check that the sort parameter matches a value in the sort list
	v.Check(validator.In(f.Sort, f.SortList...), "sort", "invalid sort value")
}
