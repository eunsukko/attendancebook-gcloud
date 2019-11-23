package eventlist

import "time"

//
type Provider interface {
	GetAt(t time.Time) Eventlist
	//
	// when from already exist, then update the saved elist
	SetFrom(from time.Time, elist Eventlist)
}
