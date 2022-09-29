package config

type logging struct {
	Level string
}

type server struct {
	Debug            bool
	StaticPath       string
	ServerPort       int
	ServerHost       string
	RandomCodeLength int
	TokenTimeToLive  int
	SmsGatewayToken  string
}

type database struct {
	Uri          string
	DatabaseName string
	AsynqRedis   string
}

var (
	Logging = &logging{}

	Server = &server{}

	Database = &database{}
)
