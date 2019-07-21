package function

import (
	faasflow "github.com/s8sg/faas-flow"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
	"os"
)

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {

	dag := flow.Dag()
	dag.Node("n1").Modify(func(data []byte) ([]byte, error) {
		data = []byte(string(data) + "modifier1")
		return data, nil
	})
	dag.Node("n2").Apply("func1")
	dag.Node("n3").Apply("func2").Modify(func(data []byte) ([]byte, error) {
		data = []byte(string(data) + "modifier2")
		return data, nil
	})
	dag.Node("n4", faasflow.Aggregator(func(data map[string][]byte) ([]byte, error) {
		str1 := string(data["n2"])
		str2 := string(data["n3"])
		return []byte(str1 + str2), nil
	})).Callback("http://gateway:8080/function/fake-storage")

	dag.Edge("n1", "n2")
	dag.Edge("n1", "n3")
	dag.Edge("n2", "n4")
	dag.Edge("n3", "n4")
	return
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	consulss, err := consulStateStore.GetConsulStateStore(
		os.Getenv("consul_url"),
		os.Getenv("consul_dc"),
	)
	if err != nil {
		return nil, err
	}
	return consulss, nil
}

// ProvideDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	miniods, err := minioDataStore.InitFromEnv()
	if err != nil {
		return nil, err
	}
	return miniods, nil
}
