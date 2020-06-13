package op_error

import (
	"log"
)

func ErrorDeal(err error, msg string) {
	log.Printf("%s: %s", msg, err)

}

func PanicDeal(err error, msg string) {
	if err != nil {
		//log.info("%s: %s", msg, err)
	}
}
