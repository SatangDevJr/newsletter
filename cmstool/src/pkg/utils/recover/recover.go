package recover

import (
	"fmt"
	"runtime"
	"subscribetool/src/pkg/utils/logger"
)

func Recover(log logger.Logger) bool {
	r := recover()
	if r != nil {
		err, ok := r.(error)
		if !ok {
			go log.Error("", "middleware_Recover", nil, r)
		} else {

			stack := make([]byte, 4<<10)
			length := runtime.Stack(stack, true)
			go log.Error("", "middleware_Recover", nil, map[string]interface{}{
				"error": err,
				"stack": fmt.Sprintf("%s\n", stack[:length]),
			})
		}
		return true
	}
	return false
}
