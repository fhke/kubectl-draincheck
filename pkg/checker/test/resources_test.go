package test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createResources(t *testing.T) []expected {
	// prime return slice
	var ret []expected

	// seed random
	rand.Seed(time.Now().UnixNano())

	// Create resources 5 times
	for i := 0; i < 5; i++ {
		// create pods that have multiple PDBs acting on them
		ret = append(
			ret,
			createResourcesFor(
				t,
				"default",
				rand.Int31n(9)+1, // between 1 & 10 pods
				2,                // two PDBs
				1,                // allow disruptions
				"multiple pod disruption budgets are acting on the same pod",
			)...,
		)

		// create pods that don't allow disruptions
		ret = append(
			ret,
			createResourcesFor(
				t,
				"default",
				rand.Int31n(9)+1, // between 1 & 10 pods
				1,                // one PDB
				0,                // don't allow disruptions
				"pod disruption budget allows no disruptions",
			)...,
		)
	}

	// create unmanaged pod
	_, err := clientset.CoreV1().Pods("default").Create(
		context.TODO(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "unmanaged",
				Namespace: "default",
				Labels: map[string]string{
					"test-id": "unmanaged",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "pause",
						Image: "google/pause:asm",
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod should not return error")

	// create PDB affecting unmanaged pod
	_, err = clientset.PolicyV1().PodDisruptionBudgets("default").Create(
		context.TODO(),
		&policyv1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "unmanaged",
				Namespace: "default",
				Labels: map[string]string{
					"test-id": "unmanaged",
				},
			},
			Spec: policyv1.PodDisruptionBudgetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"test-id": "unmanaged",
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod disruption budget should not return error")

	ret = append(
		ret,
		expected{
			podName:   "unmanaged",
			namespace: "default",
			pdbNames:  []string{"unmanaged"},
			reason:    "pod has no owner references",
		},
	)

	// create unmanaged pod without any pdbs
	_, err = clientset.CoreV1().Pods("default").Create(
		context.TODO(),
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "unmanaged-2",
				Namespace: "default",
				Labels: map[string]string{
					"test-id": "unmanaged-2",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "pause",
						Image: "google/pause:asm",
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod should not return error")

	ret = append(
		ret,
		expected{
			podName:   "unmanaged-2",
			namespace: "default",
			pdbNames:  []string{},
			reason:    "pod has no owner references",
		},
	)

	return ret
}

func createResourcesFor(t *testing.T, namespace string, podCount int32, pdbCount, maxUnavailable int, reason string) []expected {
	var (
		// unique ID to use in names & labels
		id = strconv.FormatInt(time.Now().UnixNano(), 10)
		// labels applied to all pods
		podLabels = map[string]string{
			"test-id": id,
			"app":     "nginx",
		}
		// Name of created statefulset
		stsName = fmt.Sprintf("tst-%s", id)
		// names of created pod disruption budgets
		pdbNames = make([]string, pdbCount)
	)

	// create a statefulset
	_, err := clientset.AppsV1().StatefulSets(namespace).Create(
		context.TODO(),
		&appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      stsName,
				Namespace: namespace,
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: podLabels,
				},
				PodManagementPolicy: appsv1.ParallelPodManagement,
				Replicas:            &podCount,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: podLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "pause",
								Image: "google/pause:asm",
							},
						},
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	require.NoErrorf(t, err, "Creating deployment %s should not return error", stsName)

	// create pod disruption budgets
	for i := 0; i < pdbCount; i++ {
		pdb, err := clientset.PolicyV1().PodDisruptionBudgets(namespace).Create(
			context.TODO(),
			&policyv1.PodDisruptionBudget{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("tst-%s-%d", id, i),
					Namespace: namespace,
				},
				Spec: policyv1.PodDisruptionBudgetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: podLabels,
					},
					MaxUnavailable: intStr(maxUnavailable),
				},
			},
			metav1.CreateOptions{},
		)
		require.NoError(t, err, "Creating pod disruption budget should not error")
		pdbNames[i] = pdb.Name
	}

	// create return structs
	ret := make([]expected, podCount)
	for i := 0; i < int(podCount); i++ {
		ret[i] = expected{
			podName:   fmt.Sprintf("%s-%d", stsName, i),
			namespace: namespace,
			pdbNames:  pdbNames,
			reason:    reason,
		}
	}
	return ret
}

// Wait for pods to start
func waitPods(t *testing.T, es []expected, interval time.Duration, maxTries int) {
	var tries int
	for range time.NewTicker(interval).C {
		var notRunning int
		for _, e := range es {
			// get pod
			pod, err := clientset.CoreV1().Pods(e.namespace).Get(
				context.TODO(),
				e.podName,
				metav1.GetOptions{},
			)
			require.NoErrorf(t, err, "Getting pod %s/%s should not error", e.namespace, e.podName)
			// Check that pod is running
			if pod.Status.Phase != corev1.PodRunning {
				t.Logf("Pod %s/%s is not running (phase %s)", e.namespace, e.podName, pod.Status.Phase)
				notRunning++
			} else {
				t.Logf("Pod %s/%s is running", e.namespace, e.podName)
			}
		}

		if notRunning > 0 {
			t.Logf("Waiting for %d pods to start", notRunning)
		} else {
			t.Log("All pods running")
			return
		}
		tries++
		if tries > maxTries {
			t.Fatal("Exceeded max retries while waiting for pods to start")
		}
	}
}

func intStr(i int) *intstr.IntOrString {
	out := intstr.FromInt(i)
	return &out
}
