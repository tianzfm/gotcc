package dao

type TaskGroupFlowPO struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Defination  string `db:"defination"`
	Version     int    `db:"version"`
	IsActive    int    `db:"is_active"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	CreateUser  string `db:"create_user"`
	UpdatedUser string `db:"updated_user"`
}
