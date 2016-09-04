# selectable
selectable.Barrier is a select{}-friendly barrier. A rendezvous point for goroutines.

[docs at https://godoc.org/github.com/glycerine/selectable](https://godoc.org/github.com/glycerine/selectable])

the cental logic inside the Barrier looks like this:

```
      for {
         select {
         case <-b.ReleaseAndReset:
            close(b.waitCh) // release all waiting goroutines
            b.waitCh = make(chan struct{}) // and make a new wait channel
            
         case b.waitForRelease <- b.waitCh:
            // only send b.waitCh, nothing else.
            
         case withRelease := <-b.RequestStop:
            if withRelease {
               close(b.waitCh) // release all waiting goroutines
            }
            close(b.waitForRelease) // return nil wait channel from now all
            close(b.Done) // signal that shutdown is complete
            return
         }
      }

...
// Wait returns a channel to wait on. The
// channel will be closed when
// `b.ReleaseAndReset <- struct{}{}`
// is invoked.
func (b *Barrier) Wait() chan struct{} {
	return <-b.waitForRelease
}
```

and the Barrier is used like this:

```
      b := selectable.NewBarrier()
      go func() {
         for {
            select {
               case <-b.Wait(): // wait here for release.
                  // ReleaseAndReset <- struct{}{} was invoked
                  
               case <-b.Done:
                  // Since b.Wait() could return a nil channel
                  // if the Barrier is shutting down,
                  // *always* include a <-b.Done in
                  // your select to avoid deadlock
                  // on shutdown.
                  return
               ...
            }
          }
       }()

       ...
       b.ReleaseAndReset <- struct{}{} // release the kraken!
       ...
```

[docs at https://godoc.org/github.com/glycerine/selectable](https://godoc.org/github.com/glycerine/selectable])

Author: Jason E. Aten, Ph.D.

license: MIT
