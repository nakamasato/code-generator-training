package main

import (
	examplecomv1alpha1 "code-generator-training/pkg/api/example.com/v1alpha1"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	beforeNum := 3
	afterNum := 5
	replica := int32(beforeNum)
	foo := examplecomv1alpha1.Foo{
		Spec: examplecomv1alpha1.FooSpec{
			DeploymentName: "test",
			Replicas:       &replica,
		},
	}

	copiedFoo := foo
	deepCopiedFoo := foo.DeepCopy()

	checkReplica(t, int32(beforeNum), *copiedFoo.Spec.Replicas)
	checkReplica(t, int32(beforeNum), *deepCopiedFoo.Spec.Replicas)

	replica = int32(afterNum)

	checkReplica(t, int32(afterNum), *copiedFoo.Spec.Replicas)
	checkReplica(t, int32(beforeNum), *deepCopiedFoo.Spec.Replicas)
}

func checkReplica(t *testing.T, want, got int32) {
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
