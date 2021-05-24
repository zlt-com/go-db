module github.com/zlt-com/go-db

go 1.16

require (
	github.com/garyburd/redigo v1.6.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/zlt-com/go-common v0.0.0-20210513085452-3c4f1661ab09
	github.com/zlt-com/go-config v0.0.0-20210513094338-cc0e72f6cfb4
	github.com/zlt-com/go-logger v0.0.0-20210513095531-9e90dff15f9d
)

replace(
github.com/zlt-com/go-common => ../go-common
github.com/zlt-com/go-config => ../go-config
github.com/zlt-com/go-logger => ../go-logger
)