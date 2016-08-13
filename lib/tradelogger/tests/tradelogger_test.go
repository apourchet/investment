package test

import "time"
import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/apourchet/investment/lib/tradelogger"
)

func TestMain(t *testing.M) {
	logger, err := NewLocalLogger("/tmp/testlog")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
	logger.Log(&Item{time.Now(), "TestTag", "TestMessage"})
	_, err = ioutil.ReadFile("/tmp/testlog")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}
