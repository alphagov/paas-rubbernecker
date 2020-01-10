# Rubbernecker

An application that converts various tools the GOV.UK PaaS team is using into a
friendly kanban wall used for standup and quick reference of what's going on.

## Dependencies

There are various dependencies for Node as well as Golang that would need to be
installed before compilation of the code.

The make file has a target for them all:

```sh
make dependencies
```

## Building

The repo comes in with a make file which makes it easier to build the project
from various Golang, TypeScript and SASS source code.

The following command should build them all:

```sh
make build
```

## Testing

The tests can be fired with the use of make target which will run ginkgo suite
of tests:

```sh
make test
```

## Running

After the application has been compiled, you should be able to execute the
following:

```sh
./bin/rubbernecker
```

### Requirements

Following environment variables are required to be provided for the application
to work properly:

```sh
PIVOTAL_TRACKER_PROJECT_ID
PIVOTAL_TRACKER_API_TOKEN
PAGERDUTY_AUTHTOKEN
```

These can be provided in a form of flags. See the help section for more
details.

### Help

You can find some exciting functionality if you run:

```sh
./bin/rubbernecker --help
```


### Creative Commons usage

- [documentation icon](dist/img/documentation.svg) by Adrien Coquet
- [goat icon](dist/img/goat.svg) by Chanut is Industries
- [grinch icon](dist/img/grinch.svg) by Denis Shumaylov
- [brain icon](dist/img/brain.svg) by Sumit Saengthong
