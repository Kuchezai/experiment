package logger

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO: implement a normal logger, not the one built into gin
func LogsGinToJSON() gin.HandlerFunc {
	return gin.LoggerWithFormatter(

		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})
			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["response_time"] = params.Latency.String()

			if params.StatusCode >= 400 {
				errorMessage := params.ErrorMessage
				re := regexp.MustCompile(`Error #[0-9]+: `)
				errorMessage = strings.TrimRight(re.ReplaceAllString(errorMessage, ""), "\n")
				log["error"] = errorMessage
			}

			s, _ := json.Marshal(log)
			logLine := string(s) + "\n"

			return logLine
		},
	)
}
