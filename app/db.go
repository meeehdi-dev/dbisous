package app

type DbService struct{}

func (d *DbService) GetDb() string {
	return "Hello DB!"
}
