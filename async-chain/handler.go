package function

import (
	faasflow "github.com/s8sg/faas-flow"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
)

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {
	dag := flow.Dag()
	dag.Node("n1").Apply("func1")
	dag.Node("n2").Apply("func2").
		Modify(func(data []byte) ([]byte, error) {
			data = []byte(string(data) + "modifier")
			return data, nil
		})
	dag.Node("n3").Callback("http://gateway:8080/function/fake-storage")
	dag.Edge("n1", "n2")
	dag.Edge("n2", "n3")
	flow.Finally(func(state string) {
		// cleanup code
	})

	return
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	return nil, nil
}

// DefineDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	// initialize minio DataStore
	miniods, err := minioDataStore.InitFromEnv()
	if err != nil {
		return nil, err
	}
	return miniods, nil
}
