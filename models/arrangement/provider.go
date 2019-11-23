package arrangement

import "time"

//
type Provider interface {
	GetAt(t time.Time) Arrangement
	//
	// when from already exist, then update the saved arrangement
	SetFrom(from time.Time, arrangement Arrangement)
}
