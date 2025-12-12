package entities

import "time"

// Budget aggregates income and expense for a given period.
type Budget struct {
	Income      int64      `db:"income" json:"income"`
	Expense     int64      `db:"expense" json:"expense"`
	LastUpdated *time.Time `db:"last_updated" json:"last_updated,omitempty"`
}
