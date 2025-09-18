package confighelper

import (
	"api-gateway-sql/config"
	"fmt"
)

// GetTargetAndDatabase used to get target and his database
func GetTargetAndDatabase(config *config.Config, targetName string) (*config.Target, *config.Database, error) {
	target, exist := config.GetTargetByName(targetName)
	if !exist {
		return nil, nil, fmt.Errorf("the specified target name %s does not exist", targetName)
	}

	database, exist := config.GetDatabaseByDataSourceName(target.DataSourceName)
	if !exist {
		return nil, nil, fmt.Errorf("the configured datasource name %s does not exist", target.DataSourceName)
	}

	return &target, &database, nil
}
