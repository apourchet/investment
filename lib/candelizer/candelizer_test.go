package candelizer_test

import "github.com/apourchet/investment/lib/candelizer"
import "testing"

func TestBasicCandelizer(t *testing.T) {
	c := candelizer.NewCandelizer(3)
	c.Step(1)
	c.Step(3)
	c.Step(2)
	if c.Open() != 1 || c.Close() != 2. || c.High() != 3. || c.Low() != 1 {
		t.Fail()
	}
}
