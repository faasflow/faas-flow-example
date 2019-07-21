package function

import (
	faasflow "github.com/s8sg/faas-flow"
)

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {
	flow.SyncNode().Apply("func1").Apply("func2").
		Modify(func(data []byte) ([]byte, error) {
			data = []byte(string(data) + "modifier")
			return data, nil
		})
	return
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	return nil, nil
}

// ProvideDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	return nil, nil
}
