package models

//go:generate reform

// reform:News
type News struct {
	ID      uint64 `reform:"Id,pk"`
	Title   string `reform:"Title"`
	Content string `reform:"Content"`
}

// reform:NewsCategories
type NewsCategory struct {
	NewsID     uint64 `reform:"NewsId"`
	CategoryID uint64 `reform:"CategoryId"`
}
