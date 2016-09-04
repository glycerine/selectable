package selectable

import (
	cv "github.com/glycerine/goconvey/convey"
	"sync"
	"testing"
)

func Test001BarrierReleasesAllGoroutines(t *testing.T) {
	cv.Convey("selectable.Barrier should release all waiting goroutines when ReleaseAndReset is sent a struct{}{}", t, func() {

		b := NewBarrier()
		m := 5 // number of re-uses of the barrier to conduct.

		for j := 0; j < m; j++ {
			// use the WaitGroup sync primitive to verify our Barrier
			n := 3
			wgUp := &sync.WaitGroup{}
			wgUp.Add(n)

			wgDown := &sync.WaitGroup{}
			wgDown.Add(n)

			for i := 0; i < n; i++ {
				go func() {
					ch := b.Wait() // normally we could receive immediately here;
					wgUp.Done()    // but we want to coordinate with wgUp too here.
					select {
					case <-ch: // normally <- b.Wait() is fine.
						wgDown.Done()
					case <-b.Done:
						// good form to always have this in our selects.
					}
				}()
			}
			wgUp.Wait()
			b.ReleaseAndReset <- struct{}{}
			wgDown.Wait()
			cv.So(true, cv.ShouldBeTrue) // we should get to here.
		}
	})
}

func Test002BarrierShutdownDoesntDeadlock(t *testing.T) {
	cv.Convey("selectable.Barrier on shutdown should not deadlock", t, func() {

		// use the WaitGroup sync primitive to verify our Barrier
		n := 3
		wgUp := &sync.WaitGroup{}
		wgUp.Add(n)

		wgDown := &sync.WaitGroup{}
		wgDown.Add(n)

		b := NewBarrier()
		for i := 0; i < n; i++ {
			go func() {
				ch := b.Wait() // normally we could receive immediately here;
				wgUp.Done()    // but we want to coordinate with wgUp too here.
				select {
				case <-ch: // normally <- b.Wait() is fine.
				case <-b.Done:
					wgDown.Done()
				}
			}()
		}
		wgUp.Wait()
		close(b.RequestStop)
		wgDown.Wait()
		cv.So(true, cv.ShouldBeTrue) // we should get to here.
	})
}
