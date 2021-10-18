package golist

import "testing"

func TestTaskContext(t *testing.T) {
	m := "my message"
	a := 123

	c := &taskContext{
		setMessage: func(msg string) {
			if msg != m {
				t.Errorf("expected set-message to be %q, got %q", m, msg)
			}
		},
		println: func(a ...interface{}) error {
			return nil
		},
		printfln: func(f string, a ...interface{}) error {
			return nil
		},
	}
	func(c TaskContext) {
		c.SetMessage(m)
		if err := c.Println(a); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}(c)
}
