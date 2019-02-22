package res

import "time"

type BranchTree struct {
	Id        int64         `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Name      string        `json:"name"`
	Childrens []*BranchTree `json:"childrens"`
}
