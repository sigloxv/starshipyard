package index

func NewOptions() *Options {
	return &Options{
		Limit: -1,
	}
}

type Options struct {
	Limit   int
	Skip    int
	Reverse bool
}
