package utils

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGenRandomStr(t *testing.T) {
	mu := &sync.Mutex{}
	res := make([]any, 0)

	var wg sync.WaitGroup
	for _ = range make([]int, 100) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			v := GenRandomStr(5)

			mu.Lock()
			res = append(res, v)
			mu.Unlock()
		}()
	}
	wg.Wait()

	assert.True(t, IsUnique(res))
}
