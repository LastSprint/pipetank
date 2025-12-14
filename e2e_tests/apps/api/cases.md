# Test cases

Dependencies:
1. MongoDB
2. gRPC App (cmd/grpc_api)

Purpose:

Black-box testing of the whole system.

## Scenarios

### Close to real-life clients

Given:
- 3 clients
- 2 processes
- 3 stages


Client 1 (process 1):
- Execution 1
  - stage 1 - success
  - stage 2 - success
  - stage 3 - success
- Execution 2
  - stage 1 - started - failed (bcs of stage 3)
  - stage 2 - started - failed (bcs of stage 3)
  - stage 3 - failed
- Execution 3 (partial)
  - stage 1 - started - succeeded
  - stage 2 - metadata update
  - stage 3 - (no start) - failed (with failure and output)

Client 2 (process 1):
- Execution 1
  - stage 1 - success
  - stage 2 - success
    - update
    - update
  - stage 3 - success
- Execution 2
  - stage 1 - started - failed (bcs of stage 2)
  - stage 2 - failed
  - stage 3 - not started

Client 3 (process 2 at the same time with client 1 and client 2):
- Execution 1
  - stage 1 - success
  - stage 2 - success
  - stage 3 - failed (does not break previous stages)
- Execution 2
    - stage 1 - started
    - Here the client fails and does not send any events
  - Execution 3 (here client restores, but starts with new Execution ID)
    - stage 1 - success 