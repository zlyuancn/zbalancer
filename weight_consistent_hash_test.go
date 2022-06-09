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

func Test_weightConsistentHashBalancer_Get(t *testing.T) {
	tests := []struct {
		name string
		ins  []Instance
	}{
		{
			"testA",
			[]Instance{
				NewInstance("A").SetName("A").SetWeight(300),
				NewInstance("B").SetName("B").SetWeight(500),
				NewInstance("C").SetName("C").SetWeight(400),
				NewInstance("D").SetName("D").SetWeight(600),
				NewInstance("E").SetName("E").SetWeight(200),
			},
		},
		{
			"testB",
			[]Instance{
				NewInstance("A").SetName("A").SetWeight(300),
				NewInstance("B").SetName("B").SetWeight(300),
				NewInstance("C").SetName("C").SetWeight(300),
				NewInstance("D").SetName("D").SetWeight(300),
				NewInstance("E").SetName("E").SetWeight(300),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBalancer(WeightConsistentHashBalancer)
			b.Update(test.ins)

			const count = 10000000
			result := make(map[string]int)
			for i := 0; i < count; i++ {
				in, _ := b.Get(WithHashKey(strconv.Itoa(i)))
				result[in.Instance().(string)]++
			}
			t.Log(result)

			for i := 0; i < len(test.ins); i++ {
				name := test.ins[i].Instance().(string)
				p := float64(result[test.ins[i].Instance().(string)]) / float64(count)
				t.Logf("The probability of %v is %.5f", name, p)
			}
		})
	}
}

func Test_weightConsistentHashBalancer_Target(t *testing.T) {
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

	b, _ := NewBalancer(WeightConsistentHashBalancer)
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(WithTarget(test.target[i])); !reflect.DeepEqual(got.Instance(), test.want[i]) {
			t.Errorf("Get() = %v, want %v", got.Instance(), test.want[i])
		}
	}
}

func Test_weightConsistentHashBalancer_demotion(t *testing.T) {
	tests := []struct {
		name string
		ins  []Instance
	}{
		{
			"testA",
			[]Instance{
				NewInstance("A").SetName("A").SetWeight(300),
				NewInstance("B").SetName("B").SetWeight(500),
				NewInstance("C").SetName("C").SetWeight(400),
				NewInstance("D").SetName("D").SetWeight(600),
				NewInstance("E").SetName("E").SetWeight(200),
			},
		},
		{
			"testB",
			[]Instance{
				NewInstance("A").SetName("A").SetWeight(300),
				NewInstance("B").SetName("B").SetWeight(300),
				NewInstance("C").SetName("C").SetWeight(300),
				NewInstance("D").SetName("D").SetWeight(300),
				NewInstance("E").SetName("E").SetWeight(300),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBalancer(WeightConsistentHashBalancer)
			b.Update(test.ins)

			const count = 10000000
			result := make(map[string]int)
			for i := 0; i < count; i++ {
				in, _ := b.Get() // 不设置key迫使降级
				result[in.Instance().(string)]++
			}
			t.Log(result)

			for i := 0; i < len(test.ins); i++ {
				name := test.ins[i].Instance().(string)
				p := float64(result[test.ins[i].Instance().(string)]) / float64(count)
				t.Logf("The probability of %v is %.5f", name, p)
			}
		})
	}
}

func BenchmarkWeightConsistentHashBalancer_Get(b *testing.B) {
	balancer, _ := NewBalancer(WeightConsistentHashBalancer)
	balancer.Update([]Instance{
		NewInstance("A").SetName("A").SetWeight(300),
		NewInstance("B").SetName("B").SetWeight(500),
		NewInstance("C").SetName("C").SetWeight(400),
		NewInstance("D").SetName("D").SetWeight(600),
		NewInstance("E").SetName("E").SetWeight(200),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Get()
	}
}

func BenchmarkWeightConsistentHashBalancer_GetConcurrence(b *testing.B) {
	balancer, _ := NewBalancer(WeightConsistentHashBalancer)
	balancer.Update([]Instance{
		NewInstance("A").SetName("A").SetWeight(300),
		NewInstance("B").SetName("B").SetWeight(500),
		NewInstance("C").SetName("C").SetWeight(400),
		NewInstance("D").SetName("D").SetWeight(600),
		NewInstance("E").SetName("E").SetWeight(200),
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = balancer.Get()
		}
	})
}
