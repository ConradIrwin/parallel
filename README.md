Package parallel is an implementation of structured concurrency for go.
https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/

It is designed to help reason about parallel code by ensuring that
go-routines are started and stopped in a strictly nested pattern: a child
goroutine will never outlive its parent.

See https://pkg.go.dev/github.com/ConradIrwin/parallel for full documentation.

```go
parallel.Do(func(p *parallel.P) {
    p.Go(doSomethingSlow)
    p.Go(doSomethingElse)
})

parallel.Each([]int{1,2,3}, func(i int) {
    time.Sleep(i * time.Second)
    fmt.Println(i)
})
```
