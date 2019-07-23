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
	conditionalDags := dag.ConditionalBranch("C",
		[]string{"c1", "c2"}, // possible conditions
		func(response []byte) []string {
			// for each returned condition the corresponding branch will execute
			// this function executes in the runtime of condition C
			return []string{"c1", "c2"}
		},
		faasflow.Aggregator(func(data map[string][]byte) ([]byte, error) {
			str1 := string(data["c1"])
			str2 := string(data["c2"])
			return []byte(str1 + str2), nil
		}),
	)

	conditionalDags["c2"].Node("n1").Apply("func1").Modify(func(data []byte) ([]byte, error) {
		data = []byte(string(data) + "modifier2")
		return data, nil
	})
	foreachDag := conditionalDags["c1"].ForEachBranch("F",
		func(data []byte) map[string][]byte {
			// for each returned key in the hashmap a new branch will be executed
			// this function executes in the runtime of foreach F
			return map[string][]byte{"f1": data, "f2": data}
		},
		faasflow.Aggregator(func(data map[string][]byte) ([]byte, error) {
			str1 := string(data["f1"])
			str2 := string(data["f2"])
			return []byte(str1 + str2), nil
		}),
	)
	foreachDag.Node("n1").Modify(func(data []byte) ([]byte, error) {
		data = []byte(string(data) + "modifier3")
		return data, nil
	})
	dag.Node("n2").Callback("http://gateway:8080/function/fake-storage")
	dag.Edge("n1", "C")
	dag.Edge("C", "n2")
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

// DefineDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	miniods, err := minioDataStore.InitFromEnv()
	if err != nil {
		return nil, err
	}
	return miniods, nil
}
