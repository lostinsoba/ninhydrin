package schema

import "fmt"

const (
	setNamespace = "namespace"
)

func NamespaceKey() string {
	return setNamespace
}

const (
	sortedSetNamespaceTask = "namespace-task"
)

func NamespaceTaskKey(namespaceID string) string {
	return fmt.Sprintf("%s:%s", sortedSetNamespaceTask, namespaceID)
}
