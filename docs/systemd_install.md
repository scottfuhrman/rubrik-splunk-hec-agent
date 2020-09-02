# Installing Rubrik HEC Agent as a service on Linux

## Overview

This document describes the creation of a `systemd` service to run the Rubrik HEC Agent.

## Steps

Move the executable to /usr/bin.

Create the below file in /etc/systemd/system/rubrik_hec_agent.service:

```none
[Unit]
Description=Rubrik Splunk HEC Agent

[Service]
Environment=rubrik_cdm_node_ip=rubrik.demo.com
Environment=rubrik_cdm_username=svc_prometheus
Environment=rubrik_cdm_password=Mypassword123!
Environment=SPLUNK_HEC_TOKEN=3b67fb99-9935-44ef-a35c-69d466c9328a
Environment=SPLUNK_URL=https://172.21.11.23:8088/services/collector/event
Environment=SPLUNK_INDEX=development
ExecStart=/usr/bin/rubrik_hec_agent

[Install]
WantedBy=multi-user.target
```

Run the following commands to reload the daemon service, and then start the created service:

systemctl daemon-reload
systemctl start rubrik_hec_agent.service

We can now check the status of the service using the systemctl status command:

```none
# systemctl status rubrik_hec_agent.service -l
● rubrik_hec_agent.service - Rubrik Splunk HEC Agent
   Loaded: loaded (/etc/systemd/system/rubrik_hec_agent.service; disabled; vendor preset: disabled)
   Active: active (running) since Wed 2020-09-02 08:39:32 BST; 3s ago
 Main PID: 22664 (rubrik_hec_agen)
   CGroup: /system.slice/rubrik_hec_agent.service
           └─22664 /usr/bin/rubrik_hec_agent

Sep 02 08:39:32 th-prometheus.rangers.lab systemd[1]: Started Rubrik Splunk HEC Agent.
Sep 02 08:39:32 th-prometheus.rangers.lab rubrik_hec_agent[22664]: 2020/09/02 08:39:32 Cluster name: DEVOPS-1
Sep 02 08:39:33 th-prometheus.rangers.lab rubrik_hec_agent[22664]: 2020/09/02 08:39:33 Posted rubrik:storagesummary event.
#
```

The service will now start with the system. Note that the above example has been simplified through using environment variables directly in the service file definition, it is more usual to use an EnvironmentFile, and secure this file, as the credentials are considered sensitive. More details on this can be found [here](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#EnvironmentFile=).