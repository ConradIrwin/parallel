package parallel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ConradIrwin/parallel"
)

func Test_Parallel(t *testing.T) {

	a, b, c := false, false, false

	parallel.Do(func(p *parallel.P) {
		p.Go(func() {
			time.Sleep(1 * time.Millisecond)
			a = true
		})
		p.Go(func() {
			time.Sleep(2 * time.Millisecond)
			b = true
		})
		c = true
	})

	if !a || !b || !c {
		t.Fatal("parallel.Do returned before callbacks ran")
	}
}

func Test_Parallel_Panic(t *testing.T) {
	a := false

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("parallel.Do did not panic")
		}
		if !a {
			t.Fatal("parallel.Do panicked before goroutines returned")
		}
	}()
	parallel.Do(func(p *parallel.P) {
		p.Go(func() {
			time.Sleep(10 * time.Millisecond)
			a = true
		})

		p.Go(func() {
			panic("hi")
		})

		panic("oops")
	})
	t.Fatal("parallel.Do did not panic")
}

func Test_Parallel_Misuse(t *testing.T) {
	var p any
	func() {
		defer func() {
			p = recover()
		}()
		var a *parallel.P

		parallel.Do(func(p *parallel.P) {
			a = p
		})
		a.Go(func() {
			panic("oops")
		})
	}()

	if fmt.Sprint(p) != "parallel: cannot call Go after Do has returned" {
		t.Fatal("parallel.Go panicked with", p)
	}
}

var doSomethingElse = func() {}
var doSomethingSlow = func() {}

func ExampleDo() {
	parallel.Do(func(p *parallel.P) {
		p.Go(doSomethingSlow)
		p.Go(doSomethingElse)
	})
}
func ExampleEach() {
	parallel.Each([]int{1, 2, 3}, func(i int) {
		time.Sleep(time.Duration(i) * time.Millisecond)
		fmt.Println(i)
	})
	// Output: 1
	// 2
	// 3
}
