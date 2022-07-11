package enum

type Status int64

const (
	Reserved Status = iota + 1
	Booked
	Completed
	Cancelled
)

func (r Status) String() string {
	switch r {
	case Reserved:
		return "Reserved"
	case Booked:
		return "Booked"
	case Completed:
		return "Completed"
	case Cancelled:
		return "Cancelled"
	}
	return "Unknown"
}
