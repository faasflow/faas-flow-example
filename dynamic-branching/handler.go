package function

import (
	faasflow "github.com/s8sg/faas-flow"
)

// Define provide definition of the workflow
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
	dag.Node("n2").Apply("fake-storage")
	dag.Edge("n1", "C")
	dag.Edge("C", "n2")
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
