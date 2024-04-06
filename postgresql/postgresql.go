package postgresql

type Postgresql struct{
	Port int
	Host string
	User string
	Password string
	Database string
}

func NewPostgresql(
	host string,
	port int,
	database string,
	user string,
	password string,
) *Postgresql{
	return &Postgresql{
		Host: host,
		Port: port,
		Database: database,
		User: user,
		Password: password,
	}
}