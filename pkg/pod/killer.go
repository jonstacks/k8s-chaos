package pod

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// Killer randomly kills pods that match the given name to simulate some
// chaos. You can set the period that the pod Killer runs to change how
// aggressive it is.
type Killer struct {
	client    *kubernetes.Clientset
	namespace string
	period    time.Duration
	regex     *regexp.Regexp
	maxKill   int

	cancel    chan bool
	isRunning bool
}

// NewKiller creates a new pod.Killer from the options supplied.
func NewKiller(client *kubernetes.Clientset, opts *KillerOpts) *Killer {
	if opts == nil {
		opts = DefaultKillerOpts()
	}
	return &Killer{
		client:    client,
		namespace: opts.Namespace,
		regex:     opts.Regex,
		maxKill:   opts.MaxKill,
		period:    opts.Period,
		cancel:    make(chan bool, 1),
		isRunning: false,
	}
}

// Start starts the killer randomly killing pods
func (k *Killer) Start() {
	k.isRunning = true

	// Do first run, then start ticker for subsequent calls.
	if err := k.terminate(); err != nil {
		fmt.Println(err)
	}

	ticker := time.Tick(k.period)
	go func() {
		for {
			select {
			case <-k.cancel:
				k.isRunning = false
				return
			case <-ticker:
				if err := k.terminate(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

// Stop waits for the killer to stop randomly killing pods and then returns.
func (k *Killer) Stop() {
	k.cancel <- true
	for k.isRunning {
	}
}

// IsRunning returns whether or not the pod.Killer is running.
func (k *Killer) IsRunning() bool {
	return k.isRunning
}

// Namespace returns the namespace the killer is targeting.
func (k *Killer) Namespace() string {
	return k.namespace
}

// Internal terminate method called which does the actual pod retreival,
// filtering, and termination.
func (k *Killer) terminate() error {
	resp, err := k.client.Pods(k.namespace).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error getting pods: %s", err)
	}
	pods := resp.Items
	fmt.Printf("Found %d pods.\n", len(pods))

	pods = k.filter(pods)
	fmt.Printf("-> Filtered candidates down to %d pods.\n", len(pods))

	k.destroy(pods)
	return nil
}

func (k *Killer) filter(pods []v1.Pod) []v1.Pod {
	if k.regex == nil {
		return pods
	}

	matching := make([]v1.Pod, 0)
	for _, p := range pods {
		if k.regex.MatchString(p.Name) {
			matching = append(matching, p)
		}
	}
	return matching
}

func (k *Killer) destroy(pods []v1.Pod) {
	podsDeleted := 0
	Shuffle(pods, rand.NewSource(time.Now().UnixNano()))

	for _, p := range pods {
		// Short-circuit if we've deleted the max number of pods.
		if podsDeleted >= k.maxKill {
			return
		}

		err := k.client.Pods(k.namespace).Delete(p.Name, &metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("-> Error deleting pod %s: %s\n", p.Name, err)
		} else {
			podsDeleted++
			fmt.Printf("-> Deleted pod: %s\n", p.Name)
		}
	}
}
