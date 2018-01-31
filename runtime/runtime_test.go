package runtime_test

import (
	"fmt"

	"github.com/ghettovoice/gosip/runtime"
)

func ExampleGetFrame() {
	frame, ok := runtime.GetFrame()

	fmt.Println(ok)
	fmt.Println(frame)
	// Output:
	// true
	// runtime_test.ExampleGetFrame (github.com/ghettovoice/gosip/runtime/runtime_test.go:10)
}

func ExampleGetFrameOffset() {
	fn := func() (*runtime.Frame, bool) {
		// offset by one frame to get function and line of caller
		return runtime.GetFrameOffset(1)
	}

	frame, ok := fn()

	fmt.Println(ok)
	fmt.Println(frame)
	// Output:
	// true
	// runtime_test.ExampleGetFrameOffset (github.com/ghettovoice/gosip/runtime/runtime_test.go:25)
}

func ExampleGetFrameOffset_veryBigOffset() {
	frame, ok := runtime.GetFrameOffset(50)

	fmt.Println(ok)
	fmt.Println(frame)
	// Output:
	// false
	// ???
}
