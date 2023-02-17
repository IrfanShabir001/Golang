package shared

var Configs MapPropertySource

func GetConfigs() *MapPropertySource {
	return &Configs
}

func InitConfigs() {
	Configs = MapPropertySource{
		Data: map[string]interface{}{
			"database-url" : "database",
			"database-username" : "database",
			"database-password" : "database",
			"database-name" : "database",
			"listening-port": "5000",
			"sql-username":   "root",
			"sql-password":   "",
			"sql-url":        "tcp(127.0.0.1:3306)",
			"sql-name":       "test",
			"students-collection": "students",
			"disable-auth":   false,
		},
	}
}
