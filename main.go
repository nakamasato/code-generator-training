package main

import (
	"code-generator-training/pkg/client/clientset/versioned/typed/example.com/v1alpha1"
	"flag"
	"fmt"
	"path"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()
	config, e := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if e != nil {
		panic(e.Error())
	}
	// v1alpha1 package
	client, e := v1alpha1.NewForConfig(config)
	if e != nil {
		panic(e.Error())
	}
	fooList, e := client.Foos("default").List(context.TODO(), metav1.ListOptions{})
	if len(fooList.Items) == 0 {
		fmt.Println("Foo not found.")
	}
	for i, foo := range fooList.Items {
		fmt.Printf("%d\t%s\t%d\n", i, foo.Name, *foo.Spec.Replicas)
	}
}

func init() {
	home := homedir.HomeDir()
	flag.StringVar(&kubeconfig, "kubeconfig", path.Join(home, ".kube", "config"), "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
