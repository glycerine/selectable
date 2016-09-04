/*
selectable.Barrier is a select-friendly rendezvous
point.

copyright (c) 2016 Jason E. Aten
license: MIT
*/
package selectable

// Barrier is a select-friendly synchronizer.
//
// As a rendezvous point, Barrier is similar to
// a conditon variable or a waitGroup, but
// neither of those are channel based, so
// they don't play nicely with select.
//
type Barrier struct {

	// a select statement that wants to
	// to wait on a Barrier should
	// read a wait channel from Wait()
	// then wait on that channel. It will
	// be closed when true is sent to
	// ReleaseAndReset.
	//
	// sample use:
	/*
			   b := selectable.NewBarrier()
		       go func() {
			   for {
			       select {
			           case <-b.Wait():
			              // ReleaseAndReset was called
			           case <-b.Done
			              // *always* include a <-b.Done in
			              // your select to avoid deadlock
			              // on shutdown.
			              ...
			       }
			    }
	*/
	waitForRelease chan chan struct{}

	// Send on the ReleaseAndReset channel
	// to release all waiting go routines
	// and establish a new channel to
	// wait on.
	//
	// e.g. b.ReleaseAndReset <- struct{}{}
	//
	ReleaseAndReset chan struct{}

	// RequestStop is used to stop the Barrier
	// goroutine.
	//
	// Clients may close(b.RequestStop) in order to
	// shutdown the Barrier b. Sending false on
	// b.RequestStop will do the same thing. Waiting
	// goroutines will not be released by these
	// actions.
	//
	// If client instead sends true on b.RequestStop,
	// all waiting goroutines will be released before
	// shutting down.
	//
	RequestStop chan bool

	// Done will be closed when the
	// RequestStop has succeeded.
	Done chan struct{}

	// waitCh is handed out by waitForRelease
	// and renewed by ReleaseAndReset (after
	// being closed).
	waitCh chan struct{}
}

// NewBarrier creates a new select-friendly
// Barrier.
func NewBarrier() *Barrier {
	b := &Barrier{
		ReleaseAndReset: make(chan struct{}),

		waitForRelease: make(chan chan struct{}),

		RequestStop: make(chan bool),
		Done:        make(chan struct{}),
		waitCh:      make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-b.ReleaseAndReset:
				close(b.waitCh)
				b.waitCh = make(chan struct{})
			case b.waitForRelease <- b.waitCh:
				// only send b.waitCh, nothing else.
			case withRelease := <-b.RequestStop:
				//
				// we close b.waitForRelease so
				// that the '<- <- b.waitForReceive'
				// pattern doesn't deadlock on shutdown
				// on the first receive that lacks a
				// select.
				//
				// Once closed, b.waitForRelease
				// will return a nil channel,
				// which will never be chosen
				// in a select{}, thus allowing
				// the b.Done channel's close
				// to be detected.
				if withRelease {
					close(b.waitCh)
				}
				close(b.waitForRelease)
				close(b.Done)
				return
			}
		}
	}()
	return b
}

// Wait returns a channel to wait on. The
// channel will be closed when
//
//    b.ReleaseAndReset <- struct{}{}
//
// is invoked.
//
// If the Barrier is shutting down, the
// returned channel may be nil. Hence
// you should always invoke Wait() from
// within a select{} statement. Include
// case <-b.Done:
// in your select to handle shutdown
// gracefully.
//
func (b *Barrier) Wait() chan struct{} {
	return <-b.waitForRelease
}
