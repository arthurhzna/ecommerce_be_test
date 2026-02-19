import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Config AppConfig

type AppConfig struct {
	Port                  int
	ApiKey                string
	Database              Database
	RateLimiterMaxRequest float64
	RateLimiterTimeSecond int
	JwtSecretKey          string
	JwtExpirationTime     int
}

type Database struct {
	Host                  string
	Port                  int
	Name                  string
	Username              string
	Password              string
	MaxOpenConnections    int
	MaxLifeTimeConnection int
	MaxIdleConnections    int
	MaxIdleTime           int
}

func Init() {
	config := AppConfig{

}
