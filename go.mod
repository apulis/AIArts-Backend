module github.com/apulis/AIArtsBackend

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/crewjam/saml v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jinzhu/gorm v1.9.14
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/afero v1.3.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3
	golang.org/x/tools/gopls v0.5.2 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.1 // indirect
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.0
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
)

replace github.com/crewjam/saml => github.com/apulis/saml v0.4.3-0.20201026095950-84556b93e766
