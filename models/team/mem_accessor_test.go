package team

import (
	"fmt"
)

// func TestAccessor(t *testing.T) {
// 	accessor := NewInMemoryAccessor()

// }

func ExampleInitTeamAccessor() {
	accessor := GetLoadedInMemoryAccessor()

	ids := accessor.GetExistIDs()

	for _, id := range ids {
		team, _ := accessor.Get(id)

		fmt.Printf("team: %v\n", team)
	}

	// Output:
	// .
}
