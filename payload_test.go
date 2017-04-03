package gogobosh_test

const stemcells = `[
  {
    "name": "bosh-warden-boshlite-ubuntu-trusty-go_agent",
    "operating_system": "ubuntu-trusty",
    "version":  "3126",
    "cid":  "c3705a0d-0dd3-4b67-52b5-50533a432244"
  }
]`

const releases = `[
  {
    "name": "bosh-warden-cpi",
    "release_versions": [
      {
        "version": "28",
        "commit_hash": "4c36884a",
        "uncommitted_changes": false,
        "currently_deployed": true
      }
    ]
  }
]`

const deployments = `[
  {
    "name": "cf-warden",
    "cloud_config": "none",
    "releases": [
      {
        "name": "cf",
        "version": "223"
      }
    ],
    "stemcells": [
      {
        "name": "bosh-warden-boshlite-ubuntu-trusty-go_agent",
        "version": "3126"
      }
    ]
  }
]`

const deploymentTask = `{
  "id": 2,
  "state": "processing",
  "description": "run errand acceptance_tests from deployment cf-warden"
}`

const tasks = `[
  {
    "id": 1180,
    "state": "processing",
    "description": "run errand acceptance_tests from deployment cf-warden"
  }
]`

const manifest = `{
  "manifest": "---\nfoo: bar\n"
}`

const task = `{
  "id": 2,
  "state": "done",
  "description": "run errand acceptance_tests from deployment cf-warden"
}`

const vms = `{"vm_cid":"ec974048-3352-4ba4-669d-beab87b16bcb","disk_cid":null,"ips":["10.244.0.142"],"dns":[],"agent_id":"c5e7c705-459e-41c0-b640-db32d8dc6e71","job_name":"doppler_z1","index":0,"job_state":"running","state":"started","resource_pool":"medium_z1","vm_type":"default","vitals":{"cpu":{"sys":"9.1","user":"2.1","wait":"1.7"},"disk":{"ephemeral":{"inode_percent":"11","percent":"36"},"persistent":{"inode_percent":"11","percent":"36"},"system":{"inode_percent":"11","percent":"36"}},"load":["0.61","0.74","1.10"],"mem":{"kb":"2520960","percent":"41"},"swap":{"kb":"102200","percent":"10"}},"processes":[{"name":"doppler","state":"running","uptime":{"secs":11794845},"mem":{"kb":2252,"percent":16.5},"cpu":{"total":0.9}},{"name":"syslog_drain_binder","state":"running","uptime":{"secs":11794845},"mem":{"kb":2252,"percent":16.5},"cpu":{"total":0.9}},{"name":"metron_agent","state":"running","uptime":{"secs":11794845},"mem":{"kb":2252,"percent":16.5},"cpu":{"total":0.9}}],"resurrection_paused":false,"az":"z1","id":"4a9278c8-e93a-4d6a-b22c-13560208da9e","bootstrap":true,"ignore":false}`
