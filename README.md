heaptimer
===

## Quick Start

````go
import "github.com/x-mod/heaptimer"

timer := heaptimer.New(time.Microsecond * 100)
go timer.Serve(context.TODO())
defer timer.Close()

timer.Push(v1, time.Now())
timer.Push(v2, time.Now().Add(time.Second * 2))
timer.PushWithDuration(v3, time.Second * 3)

v, ok := timer.Pop()
//OR
v, ok := <-timer.V
````

