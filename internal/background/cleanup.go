package background

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenchrist/go_home/internal/value"
	"gorm.io/gorm"
)

func CleanupExpiredSensorValues(db *gorm.DB) {

	go func() {
		for {
			log.Debug().Msg("Cleaning up values")
			values := make([]value.SensorValue, 0)
			result := db.Where("expires_at < ?", time.Now()).Find(&values)
			if result.Error != nil {
				log.Error().Err(result.Error).Msg("failed to fetch expired values")
			}

			for _, value := range values {
				result := db.Delete(value)
				if result.Error != nil {
					log.Error().Err(result.Error).Uint("value_id", value.ID).Msg("Failed to delete valuev")
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()
}
