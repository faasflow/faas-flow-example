# faas-flow examples
Super simple examples of faas-flow

## Dags
> Sync Chain
![sync-chain-dag](sync-chain-dag.png)
> Async Chain
![async-chain-dag](async-chain-dag.png)
> Parallel Branching
![parallel-branch-dag](parallel-branch-dag.png)
> Dynamic Branching
![dynamic-branch-dag](dynamic-branch-dag.png)
> Conditional Branching 
![conditional-branch-dag](conditional-branch-dag.png)

## Getting Started 
1. Deploy Openfaas
2. Deploy Consul as a statestore, follow : https://github.com/s8sg/faas-flow-consul-statestore or
https://learn.hashicorp.com/consul/datacenter-deploy/deployment-guide
3. Deploy Minio as a datastore, follow : https://docs.min.io/docs/minio-deployment-quickstart-guide.html
4. Deploy Jaguer Tracing, follow: https://www.jaegertracing.io/docs/1.8/getting-started/
5. Review your configuration at `flow.yml`
```yml
environment:
  gateway: "gateway:8080"
  enable_tracing: true
  trace_server: "jaegertracing:5775"  
  enable_hmac: false
  consul_url: "statestore_consul:8500"
  consul_dc: "dc1"
  s3_url: "minio:9000"
  s3_tls: false
```
6. Update the secrets for minio in `stack.yml`
```yml
    secrets:
      - s3-secret-key
      - s3-access-key
``` 
7. Deploy the flow-functions
```bash
faas template pull https://github.com/s8sg/faas-flow
faas build
faas deploy
```
8. Request the flows
```bash
curl -v http://127.0.0.1:8080/function/sync-chain
curl -v http://127.0.0.1:8080/function/async-chain
curl -v http://127.0.0.1:8080/function/parallel-branching
curl -v http://127.0.0.1:8080/function/dynamic-branching
curl -v http://127.0.0.1:8080/function/conditional-branching
``` 
9. Check the logs of storage function

## Tracing Information in faas-flow-tower
> Sync Chain 
![sync-chain-tracing](sync-chain-racing.png)
> Async Chain
![async-chain-tracing](async-chain-tracing.png)
> Parallel Branching
![parallel-branch-tracing](parallel-branch-tracing.png)
> Dynamic Branching
![dynamic-branch-tracing](dynamic-branch-tracing.png)
> Conditional Branching
![dynamic-branch-tracing](conditional-branch-tracing.png)
