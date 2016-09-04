# selectable
a select{}-friendly barrier. Also called a rendezvous point for goroutines.

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
