module github.com/zlt-com/go-db

go 1.16

require (
	github.com/garyburd/redigo v1.6.2
	github.com/go-basic/uuid v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/zlt-com/go-common v0.0.0-20210525065252-11c6db91defb
	github.com/zlt-com/go-config v0.0.0-20210514005831-7dbcf4e20ee9
	github.com/zlt-com/go-logger v0.0.0-20210514013649-71002beb2252
)

replace (
	github.com/zlt-com/go-common => ../go-common
	github.com/zlt-com/go-config => ../go-config
	github.com/zlt-com/go-logger => ../go-logger
)
