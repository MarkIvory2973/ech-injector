package logs

import (
	"fmt"
	"os"
)

func Info(scope string, format string, arguments ...any) {
	message := fmt.Sprintf(format, arguments...)
	fmt.Fprintf(os.Stdout, "info(%s): %s\n", scope, message)
}

func Warning(scope string, format string, err error, arguments ...any) {
	if err != nil {
		format = fmt.Sprintf("%s: %v", format, err)
	}

	message := fmt.Sprintf(format, arguments...)
	fmt.Fprintf(os.Stderr, "warning(%s): %s\n", scope, message)
}

func Fatal(scope string, format string, err error, arguments ...any) {
	if err != nil {
		format = fmt.Sprintf("%s: %v", format, err)
	}

	message := fmt.Sprintf(format, arguments...)
	fmt.Fprintf(os.Stderr, "fatal(%s): %s\n", scope, message)
}
