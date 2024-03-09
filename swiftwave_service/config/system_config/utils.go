package system_config

import (
	"fmt"
	"gorm.io/gorm"
)

var config *SystemConfig
var configVersion uint = 0

func Fetch(db *gorm.DB) (*SystemConfig, error) {
	if config != nil {
		// Fetch the latest version of the config
		var record SystemConfig
		tx := db.First(&record).Select("config_version")
		if tx.Error != nil {
			return nil, tx.Error
		}
		// if the version is the same, return the cached config
		if record.ConfigVersion == configVersion {
			return config, nil
		}
	}
	// fetch first record
	var record SystemConfig
	tx := db.First(&record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	config = &record
	configVersion = record.ConfigVersion
	return config, nil
}

func Update(db *gorm.DB, config *SystemConfig) error {
	tx := db.Save(config)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a AMQPConfig) URI() string {
	return fmt.Sprintf("%s://%s:%s@%s", a.Protocol, a.User, a.Password, a.Host)
}
