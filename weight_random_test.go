/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func Test_weightRandomBalancer_Get(t *testing.T) {
	tests := []struct {
		name          string
		ins           []Instance
		possibilities []float64
	}{
		{
			"testA",
			[]Instance{
				NewInstance("A").SetWeight(150),
				NewInstance("B").SetWeight(250),
				NewInstance("C").SetWeight(200),
				NewInstance("D").SetWeight(300),
				NewInstance("E").SetWeight(100),
			},
			[]float64{0.15, 0.25, 0.2, 0.3, 0.1},
		},
		{
			"testB",
			[]Instance{
				NewInstance("A"),
				NewInstance("B"),
				NewInstance("C"),
				NewInstance("D"),
				NewInstance("E"),
			},
			[]float64{0.2, 0.2, 0.2, 0.2, 0.2},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBalancer(WeightRandomBalancer)
			b.Update(test.ins)

			const count = 10000000
			result := make(map[string]int)
			for i := 0; i < count; i++ {
				in, _ := b.Get()
				result[in.Instance().(string)]++
			}
			t.Log(result)

			for i := 0; i < len(test.ins); i++ {
				name := test.ins[i].Instance().(string)
				wantP := test.possibilities[i]
				realP := float64(result[test.ins[i].Instance().(string)]) / float64(count)
				errP := math.Abs((realP - wantP) / wantP)
				t.Logf("The probability of %v is %.5f, and the error is %.5f", name, realP, errP)
				if errP >= 0.01 {
					t.Errorf("%v has a margin of error of more than 0.05", name)
				}
			}
		})
	}
}

func Test_weightRandomBalancer_Target(t *testing.T) {
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

	b, _ := NewBalancer(WeightRandomBalancer)
	b.Update(test.ins)

	for i := 0; i < test.count; i++ {
		if got, _ := b.Get(WithTarget(test.target[i])); !reflect.DeepEqual(got.Instance(), test.want[i]) {
			t.Errorf("Get() = %v, want %v", got.Instance(), test.want[i])
		}
	}
}

func BenchmarkWeightRandomBalancer_Get(b *testing.B) {
	balancer, _ := NewBalancer(WeightRandomBalancer)
	balancer.Update([]Instance{
		NewInstance("A").SetWeight(300),
		NewInstance("B").SetWeight(500),
		NewInstance("C").SetWeight(400),
		NewInstance("D").SetWeight(600),
		NewInstance("E").SetWeight(200),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Get()
	}
}

func BenchmarkWeightRandomBalancer_GetConcurrence(b *testing.B) {
	balancer, _ := NewBalancer(WeightRandomBalancer)
	balancer.Update([]Instance{
		NewInstance("A").SetWeight(300),
		NewInstance("B").SetWeight(500),
		NewInstance("C").SetWeight(400),
		NewInstance("D").SetWeight(600),
		NewInstance("E").SetWeight(200),
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = balancer.Get()
		}
	})
}
