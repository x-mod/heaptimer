heaptimer
===

## Quick Start

````go
import "github.com/x-mod/heaptimer"

timer := heaptimer.New(time.Microsecond * 100)
go timer.Serve(context.TODO())
defer timer.Close()

timer.Push(v1, time.Second * 1)
timer.Push(v2, time.Second * 2)
timer.Push(v3, time.Second * 3)

v, ok := timer.Pop()
//OR
v, ok := <-timer.C
````

