package app

type Task struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}
