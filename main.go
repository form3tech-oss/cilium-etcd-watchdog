// Copyright 2020 Form3 Financial Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	etcdclusterv1beta2 "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	etcdclustersclient "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned"
	log "github.com/sirupsen/logrus"
	etcdv3client "go.etcd.io/etcd/clientv3"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	podutils "k8s.io/kubernetes/pkg/api/v1/pod"

	"github.com/form3tech-oss/cilium-etcd-watchdog/internal/version"
)

type EtcdClusterQuorumStatus int

const (
	EtcdClusterQuorumStatusLost EtcdClusterQuorumStatus = iota
	EtcdClusterQuorumStatusOK
	EtcdClusterQuorumStatusUnknown
)

func birthCry(kubeClient kubernetes.Interface) {
	v, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		log.Fatalf("Failed to check Kubernetes version: %v", err)
	}
	log.Infof("cilium-etcd-watchdog %s (Kubernetes %v)", version.Version, v)
}

func createClients(pathToKubeconfig string) (kubernetes.Interface, etcdclustersclient.Interface, error) {
	c, err := clientcmd.BuildConfigFromFlags("", pathToKubeconfig)
	if err != nil {
		return nil, nil, err
	}
	k, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, nil, err
	}
	e, err := etcdclustersclient.NewForConfig(c)
	if err != nil {
		return nil, nil, err
	}
	return k, e, nil
}

func getCiliumEtcdEtcdClusterResource(etcdClient etcdclustersclient.Interface, clusterName, clusterNamespace string) (*etcdclusterv1beta2.EtcdCluster, error) {
	c, err := etcdClient.EtcdV1beta2().EtcdClusters(clusterNamespace).Get(clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getEtcdClusterQuorumStatus(kubeClient kubernetes.Interface, etcdCluster *etcdclusterv1beta2.EtcdCluster, dialTimeout, opTimeout time.Duration) (EtcdClusterQuorumStatus, error) {
	// Compute the quorum size for Cilium's etcd cluster.
	q := (etcdCluster.Spec.Size / 2) + (etcdCluster.Spec.Size % 2)

	// List pods belonging to Cilium's etcd cluster.
	p, err := listEtcdClusterPods(kubeClient, etcdCluster)
	if err != nil {
		return EtcdClusterQuorumStatusUnknown, fmt.Errorf("failed to list pods for etcdcluster: %v", err)
	}

	// Check how many pods are actually running and ready.
	e := make([]string, 0)
	r := 0
	for _, pp := range p {
		if isPodRunningAndReady(pp) {
			r++
		}
		// Compute and store the etcd endpoint corresponding to this pod for eventual later usage.
		e = append(e, fmt.Sprintf("http://%s:2379", pp.Status.PodIP))
	}
	// If there are less pods running and ready than the quorum size, report quorum as having been lost.
	if r < q {
		return EtcdClusterQuorumStatusLost, nil
	}

	// Create an etcd client and see if we can actually write something to the cluster.
	etcdConfig := etcdv3client.Config{
		DialTimeout: dialTimeout,
		Endpoints:   e,
	}
	etcdClient, err := etcdv3client.New(etcdConfig)
	if err != nil {
		return EtcdClusterQuorumStatusUnknown, fmt.Errorf("failed to create etcd client: %v", err)
	}
	defer etcdClient.Close()
	ctx, fn := context.WithTimeout(context.Background(), opTimeout)
	defer fn()
	if _, err = etcdClient.Put(ctx, "cilium-etcd-watchdog", time.Now().String()); err != nil {
		return EtcdClusterQuorumStatusUnknown, fmt.Errorf("failed to wrtie to Cilium's etcd cluster: %v", err)
	}
	return EtcdClusterQuorumStatusOK, nil
}

func main() {
	// Parse command-line flags.
	clusterBootstrapGracePeriod := flag.Duration("cluster-bootstrap-grace-period", 2*time.Minute, "The period, measured since the creation of the EtcdCluster resource, during which we wait for the cluster to be marked as available")
	clusterName := flag.String("cluster-name", "cilium-etcd", "The name of the etcd cluster used by Cilium")
	clusterNamespace := flag.String("cluster-namespace", "kube-system", "The namespace where the etcd cluster used by Cilium is deployed")
	etcdClientDialTimeout := flag.Duration("etcd-client-dial-timeout", 2*time.Second, "The timeout to use when attempting to connect to Cilium's etcd cluster")
	etcdClientOpTimeout := flag.Duration("etcd-client-op-timeout", 5*time.Second, "The timeout to use when attempting to read or write from Cilium's etcd cluster")
	logLevel := flag.String("log-level", log.InfoLevel.String(), "The log level to use")
	maxQuorumStatusCheckFailures := flag.Int("max-quorum-status-check-failures", 3, "The maximum number of failed attempts to check if Cilium's etcd cluster has quorum before forcibly declaring quorum to be lost")
	pathToKubeconfig := flag.String("path-to-kubeconfig", "", "The path to the kubeconfig file to use")
	pollingInterval := flag.Duration("polling-interval", 10*time.Second, "The interval at which to poll the health of the etcd cluster used by Cilium")
	flag.Parse()

	// Configure logging.
	if v, err := log.ParseLevel(*logLevel); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	} else {
		log.SetLevel(v)
	}
	klog.SetOutput(ioutil.Discard)

	// Create a Kubernetes clienta and an 'etcd.database.coreos.com/v1beta2' API client.
	kubeClient, etcdClient, err := createClients(*pathToKubeconfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Birth cry.
	birthCry(kubeClient)

	f := 0
	t := time.NewTicker(*pollingInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			log.Debug("Inspecting the health of the etcd cluster used by Cilium...")

			// Grab the target 'EtcdCluster' resource.
			c, err := getCiliumEtcdEtcdClusterResource(etcdClient, *clusterName, *clusterNamespace)
			if err != nil {
				if kubeerrors.IsNotFound(err) {
					log.Info("Waiting for Cilium's etcd 'EtcdCluster' resource to be created")
				} else {
					log.Errorf("Failed to get Cilium's etcd 'EtcdCluster' resource: %v", err)
				}
				continue
			}

			// If the cluster has not been marked as available and the grace period didn't elapse yet, do nothing.
			t := c.CreationTimestamp.Time.Add(*clusterBootstrapGracePeriod)
			if !etcdClusterHasBootstrapped(c) && time.Now().Sub(t) < 0 {
				log.Infof("Waiting until after %s for Cilium's etcd cluster to be marked as available", t)
				continue
			}

			// After the grace period has elapsed, it is assumed that the cluster must have quorum.
			// If it doesn't, we eventually delete the target 'EtcdCluster' and wait for it to be recreated by whatever created it in the first place.
			q, err := getEtcdClusterQuorumStatus(kubeClient, c, *etcdClientDialTimeout, *etcdClientOpTimeout)
			switch q {
			case EtcdClusterQuorumStatusOK:
				f = 0
				log.Debug("Cilium's etcd cluster has quorum")
				continue
			case EtcdClusterQuorumStatusUnknown:
				f++
				log.Warnf("Failed to check if Cilium's etcd cluster has quorum: %v", err)
				if f < *maxQuorumStatusCheckFailures {
					continue
				}
				fallthrough
			case EtcdClusterQuorumStatusLost:
				f = 0
				log.Warn("Cilium's etcd has lost quorum, deleting the 'EtcdCluster' resource")
				if err := etcdClient.EtcdV1beta2().EtcdClusters(c.Namespace).Delete(c.Name, metav1.NewDeleteOptions(0)); err != nil {
					log.Warnf("Failed to delete the 'EtcdCluster' resource: %v", err)
				}
			}
		}
	}
}

func etcdClusterHasBootstrapped(etcdCluster *etcdclusterv1beta2.EtcdCluster) bool {
	for _, c := range etcdCluster.Status.Conditions {
		if c.Type == etcdclusterv1beta2.ClusterConditionAvailable {
			// We only check for the presence of the 'Available' condition, and not for its value, as if it is present then it has been 'true' at least once.
			return true
		}
	}
	return false
}

func isPodRunningAndReady(pod corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodRunning && podutils.IsPodReady(&pod)
}

func listEtcdClusterPods(kubeClient kubernetes.Interface, etcdCluster *etcdclusterv1beta2.EtcdCluster) ([]corev1.Pod, error) {
	p, err := kubeClient.CoreV1().Pods(etcdCluster.Namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=etcd,etcd_cluster=%s", etcdCluster.Name),
	})
	if err != nil {
		return nil, err
	}
	return p.Items, nil
}
