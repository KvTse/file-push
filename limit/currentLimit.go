package limit

// GLimit to control the concurrent of process
type GLimit struct {
	n int
	c chan struct{}
}

// New initialization GLimit struct
func New(n int) *GLimit {
	return &GLimit{
		n: n,
		c: make(chan struct{}, n),
	}
}

// Run f in a new goroutine but with limit.
func (g *GLimit) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
