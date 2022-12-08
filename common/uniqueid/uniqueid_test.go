package uniqueid

import (
	"fmt"
	"sync"
	"testing"
)

func TestGenId(t *testing.T) {
	wg := &sync.WaitGroup{}
	m := &sync.Map{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got := GenId()
			//got := GenUuid()
			//got := GenSn("TEST")
			_, ok := m.Load(got)
			if ok {
				fmt.Println("repeat id")
				return
			}
			m.Store(got, &struct{}{})
			println(fmt.Sprintf("%X", got))
			//println(fmt.Sprintf("%x", got))
		}()
	}
	wg.Wait()
}

func TestGenId2(t *testing.T) {
	d := GenId()
	fmt.Println(d)
	fmt.Println(fmt.Sprintf("%X", d))
}
