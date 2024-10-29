//go:build integration
// +build integration

package gogobosh_test

import (
	. "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"time"
)

func TestIntegrationBOSHAllProxy(t *testing.T) {
	g := NewGomegaWithT(t)

	boshAllProxy := os.Getenv("BOSH_ALL_PROXY")
	boshClientSecret := os.Getenv("BOSH_CLIENT_SECRET")
	if boshClientSecret == "" || boshAllProxy == "" {
		t.Skip("BOSH_ALL_PROXY test requires BOSH_ALL_PROXY and BOSH_CLIENT_SECRET env vars to run")
	}

	config := &Config{
		BOSHAddress:       "https://192.168.56.6:25555",
		ClientID:          "admin",
		ClientSecret:      boshClientSecret,
		SkipSslValidation: true,
	}
	client, err := NewClient(config)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(client).NotTo(BeNil())

	info, err := client.GetInfo()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(info.UUID).ShouldNot(BeEmpty())
}

func TestIntegration(t *testing.T) {
	g := NewGomegaWithT(t)

	boshClientSecret := os.Getenv("BOSH_CLIENT_SECRET")
	if boshClientSecret == "" {
		t.Skip("Integration test requires BOSH_CLIENT_SECRET env var to run")
	}

	config := &Config{
		BOSHAddress:       "https://192.168.56.6:25555",
		ClientID:          "admin",
		ClientSecret:      boshClientSecret,
		SkipSslValidation: true,
	}
	client, err := NewClient(config)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(client.UpdateCloudConfig(cloudConfig)).To(Succeed())
	cfg, err := client.GetCloudConfig(true)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cfg).To(HaveLen(1))
	g.Expect(cfg[0].Type).To(Equal("cloud"))

	task, err := client.UploadStemcell(
		"https://storage.googleapis.com/bosh-core-stemcells/1.92/bosh-stemcell-1.92-warden-boshlite-ubuntu-bionic-go_agent.tgz",
		"4d6823188f510215355643ad766300e076ec2e5a")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("queued"))
	task, err = client.WaitUntilDone(task, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("done"))

	task, err = client.UploadRelease(
		"https://bosh.io/d/github.com/cloudfoundry-community/nginx-release?v=1.21.6",
		"59dbc1e8dd5f4c85cca18dce1d5b70f11f9ddfcd")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("queued"))
	task, err = client.WaitUntilDone(task, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("done"))

	task, err = client.CreateDeployment(nginxManifest)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("queued"))
	task, err = client.WaitUntilDone(task, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("done"))

	vms, err := client.GetDeploymentVMs("nginx")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(vms).To(HaveLen(1))

	task, err = client.Stop("nginx", vms[0].JobName, vms[0].ID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("queued"))
	task, err = client.WaitUntilDone(task, time.Minute)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("done"))

	task, err = client.Start("nginx", "nginx", vms[0].ID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("queued"))
	task, err = client.WaitUntilDone(task, time.Minute)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(task.State).Should(Equal("done"))
}

const cloudConfig = `azs:
- name: z1
- name: z2
- name: z3

vm_types:
- name: default

disk_types:
- name: default
  disk_size: 3072

networks:
- name: default
  type: manual
  subnets:
  - azs: [z1, z2, z3]
    dns: [8.8.8.8]
    range: 10.244.0.0/24
    gateway: 10.244.0.1
    static: [10.244.0.34]
    reserved: []

compilation:
  workers: 5
  az: z1
  reuse_compilation_vms: true
  vm_type: default
  network: default`

const nginxManifest = `---
name: nginx

releases:
- name: nginx
  version: latest

stemcells:
# Centos 3421.11 _appears_ not to work under BOSH Lite
- alias: ubuntu
  os: ubuntu-bionic
  version: latest

instance_groups:
- name: nginx
  instances: 1
  azs: [ z1 ]
  vm_type: default
  persistent_disk_type: default
  stemcell: ubuntu
  networks:
  - name: default
    static_ips: [ 10.244.0.34 ]
  jobs:
  - name: nginx
    release: nginx
    properties:
      nginx_conf: |
        user nobody vcap; # group vcap can read most directories
        worker_processes  1;
        error_log /var/vcap/sys/log/nginx/error.log   info;
        events {
          worker_connections  1024;
        }
        http {
          include /var/vcap/packages/nginx/conf/mime.types;
          default_type  application/octet-stream;
          sendfile        on;
          ssi on;
          keepalive_timeout  65;
          server_names_hash_bucket_size 64;
          server {
            server_name _; # invalid value which will never trigger on a real hostname.
            listen *:80;
            # FIXME: replace all occurrences of 'example.com' with your server's FQDN
            access_log /var/vcap/sys/log/nginx/example.com-access.log;
            error_log /var/vcap/sys/log/nginx/example.com-error.log;
            root /var/vcap/data/nginx/document_root;
            index index.shtml;
          }
        }
      pre_start: |
        #!/bin/bash -ex
        NGINX_DIR=/var/vcap/data/nginx/document_root
        if [ ! -d $NGINX_DIR ]; then
          mkdir -p $NGINX_DIR
          cd $NGINX_DIR
          cat > index.shtml <<EOF
            <html><head><title>BOSH on IPv6</title>
            </head><body>
            <h2>Welcome to BOSH's nginx Release</h2>
            <h2>
            My hostname/IP: <b><!--# echo var="HTTP_HOST" --></b><br />
            Your IP: <b><!--# echo var="REMOTE_ADDR" --></b>
            </h2>
            </body></html>
        EOF
        fi
update:
  canaries: 1
  max_in_flight: 1
  serial: false
  canary_watch_time: 1000-60000
  update_watch_time: 1000-60000
`
