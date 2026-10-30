package models

type LevelEdit struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}
type UnitEdit struct {
	LevelCode   string `json:"level_code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Published   bool   `json:"published"`
}
type Resource struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Kind        string `json:"kind"`
	URL         string `json:"url"`
	Attribution string `json:"attribution"`
	Active      bool   `json:"active"`
}
