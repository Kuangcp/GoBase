package category

type CategoryTree struct {
	ID       uint           `json:"value"`
	Name     string         `json:"label"`
	Children []*CategoryTree `json:"children,omitempty"`
}
