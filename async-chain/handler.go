package function

import (
	faasflow "github.com/s8sg/faas-flow"
)

// Define provide definition of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {
	dag := flow.Dag()
	dag.Node("n1").Apply("func1")
	dag.Node("n2").Apply("func2").
		Modify(func(data []byte) ([]byte, error) {
			data = []byte(string(data) + "modifier")
			return data, nil
		})
	dag.Node("n3").Apply("fake-storage")
	dag.Edge("n1", "n2")
	dag.Edge("n2", "n3")
	flow.Finally(func(state string) {
		// cleanup code
	})

	return
}

// OverrideStateStore provides the override of the default StateStore
func OverrideStateStore() (faasflow.StateStore, error) {
	// NOTE: By default FaaS-Flow use consul as a state-store,
	//       This can be overridden with other synchronous KV store (e.g. ETCD)
	return nil, nil
}

// OverrideDataStore provides the override of the default DataStore
func OverrideDataStore() (faasflow.DataStore, error) {
	// NOTE: By default FaaS-Flow use minio as a data-store,
	//       This can be overridden with other synchronous KV store
	return nil, nil
}
