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
	"testing"
)

func Test_weightRandomBalancer_Get(t *testing.T) {
	tests := []struct {
		name          string
		ins           []interface{}
		weights       []uint8
		possibilities []float64
	}{
		{
			"testA",
			[]interface{}{"A", "B", "C", "D", "E"},
			[]uint8{3, 5, 4, 6, 2},
			[]float64{0.15, 0.25, 0.2, 0.3, 0.1},
		},
		{
			"testB",
			[]interface{}{"A", "B", "C", "D", "E"},
			nil,
			[]float64{0.2, 0.2, 0.2, 0.2, 0.2},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, _ := NewBalancer(WeightRandomBalancer)
			b.Update(test.ins, WithUpdateWeights(test.weights))

			const count = 10000000
			result := make(map[string]int)
			for i := 0; i < count; i++ {
				in, _ := b.Get()
				result[in.(string)]++
			}
			t.Log(result)

			for i := 0; i < len(test.ins); i++ {
				name := test.ins[i].(string)
				wantP := test.possibilities[i]
				realP := float64(result[test.ins[i].(string)]) / float64(count)
				errP := math.Abs((realP - wantP) / wantP)
				t.Logf("The probability of %v is %.5f, and the error is %.5f", name, realP, errP)
				if errP >= 0.01 {
					t.Errorf("%v has a margin of error of more than 0.01", name)
				}
			}
		})
	}
}

func BenchmarkWeightRandomBalancer_Get(b *testing.B) {
	balancer, _ := NewBalancer(WeightRandomBalancer)
	balancer.Update([]interface{}{1, 2, 3, 4, 5}, WithUpdateWeights([]uint8{3, 5, 4, 6, 2}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Get()
	}
}

func BenchmarkWeightRandomBalancer_GetConcurrence(b *testing.B) {
	balancer, _ := NewBalancer(WeightRandomBalancer)
	balancer.Update([]interface{}{1, 2, 3, 4, 5}, WithUpdateWeights([]uint8{3, 5, 4, 6, 2}))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = balancer.Get()
		}
	})
}
