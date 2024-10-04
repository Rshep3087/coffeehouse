// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type BrewMethod string

const (
	BrewMethodChemex      BrewMethod = "chemex"
	BrewMethodV60         BrewMethod = "v60"
	BrewMethodFrenchpress BrewMethod = "french press"
	BrewMethodAeropress   BrewMethod = "aeropress"
)

func (e *BrewMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BrewMethod(s)
	case string:
		*e = BrewMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for BrewMethod: %T", src)
	}
	return nil
}

type NullBrewMethod struct {
	BrewMethod BrewMethod `json:"brew_method"`
	Valid      bool       `json:"valid"` // Valid is true if BrewMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBrewMethod) Scan(value interface{}) error {
	if value == nil {
		ns.BrewMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BrewMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBrewMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BrewMethod), nil
}

type TempUnit string

const (
	TempUnitF TempUnit = "F"
	TempUnitC TempUnit = "C"
)

func (e *TempUnit) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TempUnit(s)
	case string:
		*e = TempUnit(s)
	default:
		return fmt.Errorf("unsupported scan type for TempUnit: %T", src)
	}
	return nil
}

type NullTempUnit struct {
	TempUnit TempUnit `json:"temp_unit"`
	Valid    bool     `json:"valid"` // Valid is true if TempUnit is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTempUnit) Scan(value interface{}) error {
	if value == nil {
		ns.TempUnit, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TempUnit.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTempUnit) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TempUnit), nil
}

type WeightUnit string

const (
	WeightUnitG  WeightUnit = "g"
	WeightUnitOz WeightUnit = "oz"
)

func (e *WeightUnit) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = WeightUnit(s)
	case string:
		*e = WeightUnit(s)
	default:
		return fmt.Errorf("unsupported scan type for WeightUnit: %T", src)
	}
	return nil
}

type NullWeightUnit struct {
	WeightUnit WeightUnit `json:"weight_unit"`
	Valid      bool       `json:"valid"` // Valid is true if WeightUnit is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullWeightUnit) Scan(value interface{}) error {
	if value == nil {
		ns.WeightUnit, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.WeightUnit.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullWeightUnit) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.WeightUnit), nil
}

// holds coffee recipes users add
type Recipe struct {
	ID            int64           `json:"id"`
	RecipeName    string          `json:"recipe_name"`
	BrewMethod    BrewMethod      `json:"brew_method"`
	CoffeeWeight  float64         `json:"coffee_weight"`
	WeightUnit    WeightUnit      `json:"weight_unit"`
	GrindSize     int32           `json:"grind_size"`
	WaterWeight   float64         `json:"water_weight"`
	WaterUnit     string          `json:"water_unit"`
	WaterTemp     sql.NullFloat64 `json:"water_temp"`
	WaterTempUnit sql.NullString  `json:"water_temp_unit"`
}

type SavedRecipe struct {
	UserID    int32     `json:"user_id"`
	RecipeID  int32     `json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           int32     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Active       bool      `json:"active"`
	Version      int32     `json:"version"`
}
