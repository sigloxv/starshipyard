package css

func DefaultCSS() string {
	return (Framework() + Overrides())
}
