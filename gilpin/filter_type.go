package gilpin

type FilterType int

const (
	None = iota
	Sub
	Up
	Average
	Paeth
	Unknown
)

func (ft FilterType) String() string {
	return [...]string{"None", "Sub", "Up", "Average", "Paeth", "Unknown"}[ft]
}
