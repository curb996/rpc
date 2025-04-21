package collector

type Pool struct {
	ch chan func()
}

func NewPool(size int) *Pool {
	p := &Pool{ch: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go func() {
			for f := range p.ch {
				f()
			}
		}()
	}
	return p
}
func (p *Pool) Submit(f func()) { p.ch <- f }
func (p *Pool) Close()          { close(p.ch) }
