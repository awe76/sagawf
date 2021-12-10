# Saga Workflow POC in go-micro 
Saga coordinator demonstrates distributed transaction processing

## Getting Started
### Clone the Saga Workflow
Clone the coordinator locally.
```shell
git clone git@github.com:awe76/sagawf.git
```

### Install go-micro
https://github.com/asim/go-micro/tree/master/cmd/micro

### Lanch service
```shell
cd sagawf
micro run
```

## Execute test call
```shell
micro call sagawf Sagawf.RunWorkflow '{"name":"default workflow","start":"s1","end":"s2","payload":"1", "operations":[{"name":"op1","from":"s1","to":"s2"},{"name":"op2","from":"s1","to":"s3"},{"name":"op3","from":"s3","to":"s2"}]}'
```

## Execution result sample

### Successful result:
```shell
{"state":{"s1":{"state":{"input":"\"1\""}},"s2":{"state":{"op1":"0.87289363","op3":"0.38564107"}},"s3":{"state":{"op2":"0.79769564"}}},"workflow_ref":{"id":1,"name":"default workflow"}}
```

log:
```shell
op2 operation is started
op1 operation is started
op1 operation is completed
op2 operation is completed
op3 operation is started
op3 operation is completed
default workflow 1 workflow is completed
workflow state: map[s1:map[input:1] s2:map[op1:0.87289363 op3:0.38564107] s3:map[op2:0.79769564]]
```

### Rollbacked result:
```shell
{"state":{"s1":{"state":{"input":"\"1\""}},"s2":{"state":{"op1":"0.032219958"}}},"workflow_ref":{"id":2,"is_rollback":true,"name":"default workflow"}}
```

log:
```shell
op1 operation is started
op2 operation is started
op1 operation is completed
op2 operation is failed
op1 operation rollback is started
op1 operation rollback is completed
default workflow 2 workflow is rollbacked
workflow state: map[s1:map[input:1] s2:map[op1:0.032219958]]
```
