package throwerror

import (
	"log"
	"os"
)

func ThrowError(err error, msg string) {
	log.Fatalf("%s: %v\n", msg, err)
	os.Exit(1)
}
