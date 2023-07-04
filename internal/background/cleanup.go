package background

import (
	"log"
	"time"

	"github.com/soerenchrist/go_home/internal/value"
	"gorm.io/gorm"
)

func CleanupExpiredSensorValues(db *gorm.DB) {

	go func() {
		for {
			log.Println("Cleaning up values")
			values := make([]value.SensorValue, 0)
			result := db.Where("expires_at < ?", time.Now()).Find(&values)
			if result.Error != nil {
				log.Fatalf("failed to fetch expired values: %v", result.Error)
			}

			for _, value := range values {
				result := db.Delete(value)
				if result.Error != nil {
					log.Fatalf("Failed to delete value %d: %v", value.ID, result.Error)
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()
}
