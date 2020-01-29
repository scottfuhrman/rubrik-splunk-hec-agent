# Rubrik Splunk HEC Agent

Splunk's HTTP Event Collector, or HEC, interface allows systems to post event data directly to Splunk's REST API via HTTP, without the need to have intermediate agents or log aggregation services. Rubrik's Splunk HEC Agent is written in GoLang and runs as a binary, or container, and will pull data from Rubrik CDM, and feed it to Splunk's HEC interface. Providing a way of getting data into Splunk without installing an add-on.

## :hammer: Installation

Pull down the following dependencies:

```bash
go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm
go get github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk
go get github.com/rubrikinc/rubrik-splunk-hec-agent
```

Clone this repository to the machine configured with GoLang, browse to the root folder, and run the following command to build the package:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

This will build the package for the linux/amd64 architecture. For other architectures, replace the values of `GOOS` and `GOARCH` as described [here](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).

This results in an executable named `main` in the current folder. This can be run to start exposing metrics.

## :mag: Example

The following environment variables must be present:

```bash
export rubrik_cdm_node_ip=192.168.0.1
export rubrik_cdm_username='admin'
export rubrik_cdm_password='MyPassword123!'

export SPLUNK_HEC_TOKEN='1234abcd-2345-67ef-a12b-1234abcd5678'
export SPLUNK_URL='https://mysplunkserver:8088/services/collector/event'
export SPLUNK_INDEX=development
```

Once these are present, the HEC agent can be run from within the cloned repo using:

```bash
go run main.go
```

Run this, or the binary created above. Once the agent is running, the results should be visible in Splunk using the following query syntax:

```none
(index="development") (source="rubrikhec")
```

Replacing the index name with whatever you specified with your `SPLUNK_INDEX` environment variable.

## :blue_book: Documentation

Here are some resources to get you started! If you find any challenges from this project are not properly documented or are unclear, please raise an issueand let us know! This is a fun, safe environment - don't worry if you're a GitHub newbie! :heart:

* Quick Start Guide
* [Rubrik API Documentation](https://github.com/rubrikinc/api-documentation)

## :muscle: How You Can Help

We glady welcome contributions from the community. From updating the documentation to adding more functions for Python, all ideas are welcome. Thank you in advance for all of your issues, pull requests, and comments! :star:

* [Contributing Guide](CONTRIBUTING.md)
* [Code of Conduct](CODE_OF_CONDUCT.md)

## :pushpin: License

* [MIT License](LICENSE)

## :point_right: About Rubrik Build

We encourage all contributors to become members. We aim to grow an active, healthy community of contributors, reviewers, and code owners. Learn more in our [Welcome to the Rubrik Build Community](https://github.com/rubrikinc/welcome-to-rubrik-build) page.

We'd  love to hear from you! Email us: build@rubrik.com :love_letter:
