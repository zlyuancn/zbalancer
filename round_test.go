/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"reflect"
	"testing"
)

func Test_roundBalancer_Get(t *testing.T) {
	test := struct {
		ins   []interface{}
		want  []interface{}
		count int
	}{
		[]interface{}{1, 2, 3},
		[]interface{}{1, 2, 3, 1, 2},
		5,
	}

	b, _ := NewBalancer(RoundBalancer)
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(); !reflect.DeepEqual(got, test.want[i]) {
			t.Errorf("Get() = %v, want %v", got, test.want[i])
		}
	}
}

func Test_roundBalancer_Upset(t *testing.T) {
	b, _ := NewBalancer(RoundBalancer)
	b.Update([]interface{}{1, 2, 3, 4, 5})

	b.Get()
	b.Get()
	b.Get()

	test := struct {
		ins   []interface{}
		want  []interface{}
		count int
	}{
		[]interface{}{1, 2, 4, 6},
		[]interface{}{6, 1, 2, 4, 6},
		5,
	}
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(); !reflect.DeepEqual(got, test.want[i]) {
			t.Errorf("Get() = %v, want %v", got, test.want[i])
		}
	}
}

func BenchmarkRoundBalancer_Get(b *testing.B) {
	balancer, _ := NewBalancer(RoundBalancer)
	balancer.Update([]interface{}{1, 2, 3, 4, 5})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Get()
	}
}

func BenchmarkRoundBalancer_GetConcurrence(b *testing.B) {
	balancer, _ := NewBalancer(RoundBalancer)
	balancer.Update([]interface{}{1, 2, 3, 4, 5})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = balancer.Get()
		}
	})
}
