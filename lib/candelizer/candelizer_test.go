package candelizer_test

import "github.com/apourchet/investment/lib/candelizer"
import "testing"
import "time"

func TestBasicCandelizer(t *testing.T) {
	c := candelizer.NewCandelizer(3)
	c.Step(1, time.Now())
	c.Step(3, time.Now())
	c.Step(2, time.Now())
	if c.Open() != 1 || c.Close() != 2. || c.High() != 3. || c.Low() != 1 {
		t.Fail()
	}
}
