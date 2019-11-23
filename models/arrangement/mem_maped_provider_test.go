package arrangement

import (
	"fmt"
	"time"
)

func ExampleInitProvider() {

	provider := GetLoadedInMemoryMapedProvider()

	fmt.Printf("arrangement: %v\n", provider.GetAt(time.Now()))
	// Output:
	// .
}

// func TestInit(t *testing.T) {
// 	cases := struct {
// 		at time.Time
// 		want
// 	}
// 	for
// }
