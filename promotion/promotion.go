package promotion

import (
	"time"

	"github.com/jinzhu/gorm"
)

type DiscountRule struct {
	gorm.Model
	DiscountID string
	Type       string
	Value      string
}

type DiscountBenefit struct {
	gorm.Model
	DiscountID string
	Type       string
	Value      string
}

type DiscountCode struct {
	gorm.Model
	DiscountID     string
	Code           string
	AvailableTimes uint
}

type Discount struct {
	gorm.Model
	Name     string
	Benefits []DiscountBenefit
	Rules    []DiscountRule
	Codes    []DiscountCode
	Begins   *time.Time
	Expires  *time.Time
}
