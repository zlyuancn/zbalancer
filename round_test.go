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
	"strconv"
	"testing"
)

func Test_roundBalancer_Get(t *testing.T) {
	test := struct {
		ins   []Instance
		want  []interface{}
		count int
	}{
		[]Instance{
			NewInstance(1),
			NewInstance(2),
			NewInstance(3),
		},
		[]interface{}{1, 2, 3, 1, 2},
		5,
	}

	b, _ := NewBalancer(RoundBalancer)
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(); !reflect.DeepEqual(got.Instance(), test.want[i]) {
			t.Errorf("Get() = %v, want %v", got.Instance(), test.want[i])
		}
	}
}

func Test_roundBalancer_Upset(t *testing.T) {
	b, _ := NewBalancer(RoundBalancer)
	b.Update([]Instance{
		NewInstance(1),
		NewInstance(2),
		NewInstance(3),
		NewInstance(4),
		NewInstance(5),
	})

	_, _ = b.Get()
	_, _ = b.Get()
	_, _ = b.Get()

	test := struct {
		ins   []Instance
		want  []interface{}
		count int
	}{
		[]Instance{
			NewInstance(1),
			NewInstance(2),
			NewInstance(4),
			NewInstance(6),
		},
		[]interface{}{1, 2, 4, 6, 1},
		5,
	}
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(); !reflect.DeepEqual(got.Instance(), test.want[i]) {
			t.Errorf("Get() = %v, want %v", got.Instance(), test.want[i])
		}
	}
}

func Test_roundBalancer_Target(t *testing.T) {
	test := struct {
		ins    []Instance
		target []string
		want   []interface{}
		count  int
	}{
		[]Instance{
			NewInstance(1).SetName(strconv.Itoa(1)),
			NewInstance(2).SetName(strconv.Itoa(2)),
			NewInstance(3).SetName(strconv.Itoa(3)),
		},
		[]string{"3", "1", "2", "2", "1", "3"},
		[]interface{}{3, 1, 2, 2, 1, 3},
		3,
	}

	b, _ := NewBalancer(RoundBalancer)
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(WithTarget(test.target[i])); !reflect.DeepEqual(got.Instance(), test.want[i]) {
			t.Errorf("Get() = %v, want %v", got.Instance(), test.want[i])
		}
	}
}

func BenchmarkRoundBalancer_Get(b *testing.B) {
	balancer, _ := NewBalancer(RoundBalancer)
	balancer.Update([]Instance{
		NewInstance(1),
		NewInstance(2),
		NewInstance(3),
		NewInstance(4),
		NewInstance(5),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Get()
	}
}

func BenchmarkRoundBalancer_GetConcurrence(b *testing.B) {
	balancer, _ := NewBalancer(RoundBalancer)
	balancer.Update([]Instance{
		NewInstance(1),
		NewInstance(2),
		NewInstance(3),
		NewInstance(4),
		NewInstance(5),
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = balancer.Get()
		}
	})
}
