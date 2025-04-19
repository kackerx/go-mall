package do

import "time"

type CommodityCategory struct {
	ID       int64
	Level    int
	ParentID int64
	Name     string
	Icon     string
	Rank     int

	CreatedBy string
	UpdatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}
