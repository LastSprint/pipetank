# PeptideTank

## Project structure

- `cmd` - contains `main` packages for different executables commands.
  - `api` - executable for clients API (`gRPC`)
  - `ui` - executable for the tool's WebUI (with front-end API)
  - `tools` - directory that contains different tools/scripts (executables) for the project.
  - `raw_events_collector` - executable for consuming raw events from MongoDB ChangeStream and storing them in UI-friendly aggregate.
- `e2e_tests` - directory that contains end-to-end tests for the project.
- `internal` - directory that contains internal packages for the project.
  - `apps` - directory that contains different applications for the project. Contains implementations of `cmd` executables.
  - `repo` - contains repository layer for the project.
- `pkg` - contains requsable components for the project.
  - `client` - contains client implementation for this service clients (`gRPC`) 