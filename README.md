# selectable
a select{}-friendly barrier. Also called a rendezvous point for goroutines.

[docs at https://godoc.org/github.com/glycerine/selectable](https://godoc.org/github.com/glycerine/selectable])

the cental logic inside the Barrier looks like this:

```
      for {
         select {
         case <-b.ReleaseAndReset:
            close(b.waitCh)
            b.waitCh = make(chan struct{})
         case b.waitForRelease <- b.waitCh:
            // only send b.waitCh, nothing else.
         case withRelease := <-b.RequestStop:
            if withRelease {
               close(b.waitCh)
            }
            close(b.waitForRelease)
            close(b.Done)
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
               case <-b.Wait():
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
       b.ReleaseAndReset <- struct{}{}
       ...
```

[docs at https://godoc.org/github.com/glycerine/selectable](https://godoc.org/github.com/glycerine/selectable])

Author: Jason E. Aten, Ph.D.

license: MIT
