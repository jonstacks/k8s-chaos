package pod

import (
	"math/rand"

	"k8s.io/client-go/pkg/api/v1"
)

// Shuffle shuffles the given array.
func Shuffle(pods []v1.Pod, source rand.Source) {
	random := rand.New(source)
	for i := len(pods) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		pods[i], pods[j] = pods[j], pods[i]
	}
}
