package main

import (
	examplecomv1alpha1 "code-generator-training/pkg/api/example.com/v1alpha1"
	lister "code-generator-training/pkg/client/listers/example.com/v1alpha1"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

func TestDeepCopy(t *testing.T) {
	beforeNum := 3
	afterNum := 5
	replica := int32(beforeNum)
	foo := getFoo("foo-sample", "default", &replica, map[string]string{})

	copiedFoo := foo
	deepCopiedFoo := foo.DeepCopy()

	checkReplica(t, int32(beforeNum), *copiedFoo.Spec.Replicas)
	checkReplica(t, int32(beforeNum), *deepCopiedFoo.Spec.Replicas)

	replica = int32(afterNum)

	checkReplica(t, int32(afterNum), *copiedFoo.Spec.Replicas)
	checkReplica(t, int32(beforeNum), *deepCopiedFoo.Spec.Replicas)
}

func TestListerList(t *testing.T) {
	replica := int32(3)
	// Prepare Indexer
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	// Add deployment with label
	err := indexer.Add(getFoo("foo-with-label", "default", &replica, map[string]string{"watch": "true"}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}
	// Add deployment without label
	err = indexer.Add(getFoo("foo-without-label", "default", &replica, map[string]string{}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}

	// Prepare Lister
	fooLister := lister.NewFooLister(indexer)
	selector := labels.SelectorFromSet(labels.Set{"watch": "true"})

	// List
	ret, _ := fooLister.List(selector)

	// Expect the result length to be 1 (only hit Foo with label)
	if len(ret) != 1 {
		t.Errorf("want %d, got %d", 1, len(ret))
	}
}

func TestNamespaceLister(t *testing.T) {
	replica := int32(3)
	// Prepare Indexer
	indexer := cache.NewIndexer(
		cache.MetaNamespaceKeyFunc,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	// Add deployment with label in default namespace
	err := indexer.Add(getFoo("foo-with-label", "default", &replica, map[string]string{"watch": "true"}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}
	// Add deployment with label in target namespace
	err = indexer.Add(getFoo("foo-with-label", "target", &replica, map[string]string{"watch": "true"}))
	if err != nil {
		t.Errorf("indexer.Add failed %v\n", err)
	}

	// Prepare Lister
	fooLister := lister.NewFooLister(indexer)
	selector := labels.SelectorFromSet(labels.Set{"watch": "true"})

	// Foos() get namespaceLister
	fooNamespaceLister := fooLister.Foos("target")
	ret, _ := fooNamespaceLister.List(selector)

	// Expect the result length to be 1 (only hit Foo in target namespace)
	if len(ret) != 1 {
		t.Errorf("want %d, got %d", 1, len(ret))
	}
}

func checkReplica(t *testing.T, want, got int32) {
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func getFoo(name, namespace string, replicas *int32, labels map[string]string) *examplecomv1alpha1.Foo {
	return &examplecomv1alpha1.Foo{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: examplecomv1alpha1.FooSpec{
			DeploymentName: name,
			Replicas:       replicas,
		},
	}
}
