package container

type todo struct {
	db database
}

func newTodo(db database) *todo {
	return &todo{db: db}
}

type database interface {
	client() string
}

type driver struct {
	dbanme string
}

func (d driver) client() string {
	return d.dbanme
}

var _ database = &driver{}

func newDriver(dbanme string) *driver {
	return &driver{dbanme: dbanme}
}

type namer interface {
	getName() string
}

type ager interface {
	getAge() int
}

type driverer interface {
	getDb() database
}

type userer interface {
	namer
	ager
	driverer
}

type user struct {
	name string
	age  int
	db   database
}

func newUser(name string, age int) *user {
	return &user{name: name, age: age}
}

func newUserWithDriver(db database) *user {
	return &user{db: db}
}

func (u *user) getName() string {
	return u.name
}

func (u *user) getAge() int {
	return u.age
}

func (u *user) getDb() database {
	return u.db
}
