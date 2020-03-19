package namedlocker

import (
	"runtime"
	"sync"
	"testing"
)

func TestStore(t *testing.T) {
	sto := New()
	sto.Lock("hello")
	sto.TryUnlock("hello")

	sto.RLock("hello")
	sto.TryRUnlock("hello")
}

func TestStorePanicRLock(t *testing.T) {
	sto := New()
	func() {
		defer recover()
		sto.TryRUnlock("hello")
	}()
}

func TestStorePanicLock(t *testing.T) {
	sto := New()
	func() {
		defer recover()
		sto.TryUnlock("hello")
	}()
}

func BenchmarkSyncStore(b *testing.B) {
	b.ReportAllocs()
	sto := New()
	k := ""
	for i := 0; i < b.N; i++ {
		sto.Lock(k)
		runtime.Gosched()
		sto.Unlock(k)
	}
}

func BenchmarkSyncNormal(b *testing.B) {
	b.ReportAllocs()
	lk := sync.RWMutex{}
	for i := 0; i < b.N; i++ {
		lk.Lock()
		runtime.Gosched()
		lk.Unlock()
	}
}

func BenchmarkAsyncStore(b *testing.B) {
	b.ReportAllocs()
	sto := New()
	k := ""
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sto.Lock(k)
			runtime.Gosched()
			sto.Unlock(k)
		}
	})
}

func BenchmarkAsyncNormal(b *testing.B) {
	b.ReportAllocs()
	lk := sync.RWMutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lk.Lock()
			runtime.Gosched()
			lk.Unlock()
		}
	})
}
