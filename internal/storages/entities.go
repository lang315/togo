package storages

type User struct {
	ID       string `sql:",pk,varchar(26)" json:"id"`
	Password string `sql:",varchar(255),notnull" json:"-"`
	MaxToDo  int    `sql:"notnull" json:"max_to_do"`
}

type Task struct {
	ID        string `sql:",pk,varchar(26)" json:"-"`
	Content   string `sql:",notnull" json:"content"`
	UserID    string `sql:",varchar(26),notnull" json:"-"`
	CreatedAt int64  `sql:",notnull" json:"created_at" json:"created_at"`
}
