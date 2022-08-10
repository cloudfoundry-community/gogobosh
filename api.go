package gogobosh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const maxRetries = 360 // equates to 5m with 1s sleep

// GetStemcells from given BOSH
func (c *Client) GetStemcells() ([]Stemcell, error) {
	r := c.NewRequest("GET", "/stemcells")
	resp, err := c.DoRequest(r)
	if err != nil {
		return []Stemcell{}, fmt.Errorf("error requesting stemcells: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Stemcell{}, fmt.Errorf("error reading stemcells response: %w", err)
	}

	var stemcells []Stemcell
	err = json.Unmarshal(resBody, &stemcells)
	if err != nil {
		return []Stemcell{}, fmt.Errorf("error unmarshalling stemcells response: %w", err)
	}
	return stemcells, nil
}

func (c *Client) UploadStemcell(url, sha1 string) (Task, error) {
	r := c.NewRequest("POST", "/stemcells")
	in := struct {
		Location string `json:"location"`
		SHA1     string `json:"sha1"`
	}{
		Location: url,
		SHA1:     sha1,
	}

	b, err := json.Marshal(&in)
	if err != nil {
		return Task{}, fmt.Errorf("error marshalling upload request: %w", err)
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		return Task{}, fmt.Errorf("error requesting stemcell upload: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading stemcell upload response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshaling stemcell upload response: %w", err)
	}
	return task, nil
}

// GetReleases from given BOSH
func (c *Client) GetReleases() ([]Release, error) {
	r := c.NewRequest("GET", "/releases")
	resp, err := c.DoRequest(r)
	if err != nil {
		return []Release{}, fmt.Errorf("error requesting releases: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Release{}, fmt.Errorf("error reading releases response: %w", err)
	}

	var releases []Release
	err = json.Unmarshal(resBody, &releases)
	if err != nil {
		return []Release{}, fmt.Errorf("error unmarshalling releases response: %w", err)
	}
	return releases, nil
}

func (c *Client) UploadRelease(url, sha1 string) (Task, error) {
	r := c.NewRequest("POST", "/releases")
	in := struct {
		Location string `json:"location"`
		SHA1     string `json:"sha1"`
	}{
		Location: url,
		SHA1:     sha1,
	}

	b, err := json.Marshal(&in)
	if err != nil {
		return Task{}, fmt.Errorf("error marshalling upload release request: %w", err)
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		return Task{}, fmt.Errorf("error requesting release upload: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading upload release task response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshalling upload release task response: %w", err)
	}
	return task, nil
}

// GetDeployments from given BOSH
func (c *Client) GetDeployments() ([]Deployment, error) {
	r := c.NewRequest("GET", "/deployments")
	resp, err := c.DoRequest(r)
	if err != nil {
		return []Deployment{}, fmt.Errorf("error requesting deployments: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Deployment{}, fmt.Errorf("error reading deployments response: %w", err)
	}

	var deployments []Deployment
	err = json.Unmarshal(resBody, &deployments)
	if err != nil {
		return []Deployment{}, fmt.Errorf("error unmarshalling deployments response: %w", err)
	}
	return deployments, nil
}

// GetDeployment from given BOSH
func (c *Client) GetDeployment(name string) (Manifest, error) {
	r := c.NewRequest("GET", "/deployments/"+name)
	resp, err := c.DoRequest(r)

	if err != nil {
		return Manifest{}, fmt.Errorf("error requesting deployment manifest: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Manifest{}, fmt.Errorf("error reading deployment manifest response: %w", err)
	}

	var manifest Manifest
	err = json.Unmarshal(resBody, &manifest)
	if err != nil {
		return Manifest{}, fmt.Errorf("error unmarshalling deployment manifest response: %w", err)
	}
	return manifest, nil
}

// DeleteDeployment from given BOSH
func (c *Client) DeleteDeployment(name string) (Task, error) {
	resp, err := c.DoRequest(c.NewRequest("DELETE", "/deployments/"+name))
	if err != nil {
		return Task{}, fmt.Errorf("error deleting deployment: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return Task{}, fmt.Errorf("deployment %s not found", name)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading delete deployment response: %w", err)
	}

	var task Task
	err = json.Unmarshal(b, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshalling delete deployment response: %w", err)
	}
	return task, nil
}

// CreateDeployment from given BOSH
func (c *Client) CreateDeployment(manifest string) (Task, error) {
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.body = buffer
	r.header["Content-Type"] = "text/yaml"

	resp, err := c.DoRequest(r)

	if err != nil {
		return Task{}, fmt.Errorf("error creating deployment: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading create deployment response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshalling create deployment response: %w", err)
	}
	return task, nil
}

// GetDeploymentVMs from given BOSH
func (c *Client) GetDeploymentVMs(name string) ([]VM, error) {
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	resp, err := c.DoRequest(r)

	if err != nil {
		return []VM{}, fmt.Errorf("error requesting deployment VMs task: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []VM{}, fmt.Errorf("error reading deployment VMs task response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return []VM{}, fmt.Errorf("error unmarshalling deployment VMs task response: %w", err)
	}
	for i := 0; i <= maxRetries; i++ {
		if i == maxRetries {
			return []VM{}, fmt.Errorf("timed out getting deployment VMs task results after %d tries", maxRetries)
		}

		taskStatus, err := c.GetTask(task.ID)
		if err != nil {
			log.Printf("Error getting task %v, retrying...", err)
		}
		if taskStatus.State == "done" {
			break
		}
		time.Sleep(time.Second)
	}

	var vms []VM
	output, err := c.GetTaskResult(task.ID)
	if err != nil {
		return []VM{}, fmt.Errorf("error getting deployment VMs task result: %w", err)
	}
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				return []VM{}, fmt.Errorf("error unmarshalling deployment VMs response: %w", err)
			}
			vms = append(vms, vm)
		}
	}
	return vms, nil
}

// GetTasksByQuery from given BOSH
func (c *Client) GetTasksByQuery(query url.Values) ([]Task, error) {
	requestUrl := "/tasks?" + query.Encode()
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return []Task{}, fmt.Errorf("error requesting tasks by query: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Task{}, fmt.Errorf("error reading tasks by query response: %w", err)
	}

	var tasks []Task
	err = json.Unmarshal(resBody, &tasks)
	if err != nil {
		return []Task{}, fmt.Errorf("error unmarshalling tasks by query response: %w", err)
	}
	return tasks, nil
}

// GetTasks retrieves all BOSH tasks
func (c *Client) GetTasks() ([]Task, error) {
	return c.GetTasksByQuery(nil)
}

// GetTask retrieves the specified task from BOSH
func (c *Client) GetTask(id int) (Task, error) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Task{}, fmt.Errorf("error requesting task: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading task response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshalling task response: %w", err)
	}
	return task, nil
}

// GetTaskOutput ...
func (c *Client) GetTaskOutput(id int, typ string) ([]string, error) {
	r := c.NewRequest("GET", "/tasks/"+strconv.Itoa(id)+"/output?type="+typ)

	res, err := c.DoRequest(r)
	if err != nil {
		return []string{}, fmt.Errorf("error requesting task output: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, fmt.Errorf("error reading task output response: %w", err)
	}

	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"), nil
}

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int) ([]string, error) {
	return c.GetTaskOutput(id, "result")
}

// GetTaskEvents retrieves the events for the specified task
func (c *Client) GetTaskEvents(id int) ([]TaskEvent, error) {
	raw, err := c.GetTaskOutput(id, "event")
	if err != nil {
		return []TaskEvent{}, fmt.Errorf("error getting the task events: %w", err)
	}

	events := make([]TaskEvent, len(raw))
	for i := range raw {
		err = json.Unmarshal([]byte(raw[i]), &events[i])
		if err != nil {
			return []TaskEvent{}, fmt.Errorf("error unmarshalling the task events: %w", err)
		}
	}

	return events, nil
}

// GetCloudConfig from given BOSH
func (c *Client) GetCloudConfig(latest bool) ([]Cfg, error) {
	qs := "?latest=true"
	if !latest {
		qs = "?latest=false"
	}
	r := c.NewRequest("GET", "/configs"+qs)
	resp, err := c.DoRequest(r)
	if err != nil {
		return []Cfg{}, fmt.Errorf("error requesting cloud config: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Cfg{}, fmt.Errorf("error unmarshalling cloud config response: %w", err)
	}

	var cfg []Cfg
	err = json.Unmarshal(resBody, &cfg)
	if err != nil {
		return []Cfg{}, fmt.Errorf("error unmarshalling the cloud config: %w", err)
	}
	return cfg, nil
}

// UpdateCloudConfig updates the cloud config with the specified config
func (c *Client) UpdateCloudConfig(config string) error {
	r := c.NewRequest("POST", "/configs")
	in := struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}{
		Name:    "default",
		Type:    "cloud",
		Content: config,
	}
	b, err := json.Marshal(&in)
	if err != nil {
		return fmt.Errorf("error marshalling the cloud config update: %w", err)
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("error updating the cloud config: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// Cleanup will post to the cleanup endpoint of bosh, passing along the removeAll flag passed in as a bool
func (c *Client) Cleanup(removeAll bool) (Task, error) {
	r := c.NewRequest("POST", "/cleanup")
	var requestBody struct {
		Config struct {
			RemoveAll bool `json:"remove_all"`
		} `json:"config"`
	}
	requestBody.Config.RemoveAll = removeAll
	b, err := json.Marshal(&requestBody)
	if err != nil {
		return Task{}, fmt.Errorf("error marshalling the cleanup request: %w", err)
	}
	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"
	resp, err := c.DoRequest(r)
	if err != nil {
		return Task{}, fmt.Errorf("error making the cleanup request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading the cleanup response: %w", err)
	}

	var task Task
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error unmarshalling the cleanup response: %w", err)
	}
	return task, err
}

func (c *Client) Restart(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("restart", deployment, jobName, instanceID, true)
}

func (c *Client) RestartNoConverge(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("restart", deployment, jobName, instanceID, false)
}

func (c *Client) Stop(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("stopped", deployment, jobName, instanceID, true)
}

func (c *Client) StopNoConverge(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("stopped", deployment, jobName, instanceID, false)
}

func (c *Client) Start(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("started", deployment, jobName, instanceID, true)
}

func (c *Client) StartNoConverge(deployment, jobName, instanceID string) (Task, error) {
	return c.vmAction("started", deployment, jobName, instanceID, false)
}

func (c *Client) vmAction(action, deployment, jobName, instanceID string, converge bool) (Task, error) {
	var p string
	if converge {
		p = fmt.Sprintf("/deployments/%s/jobs/%s/%s?state=%s",
			deployment, jobName, instanceID, action)
	} else {
		p = fmt.Sprintf("/deployments/%s/instance_groups/%s/%s/actions/%s",
			deployment, jobName, instanceID, action)
	}
	return c.executeVMAction(action, p)
}

func (c *Client) executeVMAction(action, actionPath string) (Task, error) {
	var task Task
	r := c.NewRequest("PUT", actionPath)
	r.header["Content-Type"] = "text/yaml"
	err := c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error creating VM %s task: %w", action, err)
	}
	return task, nil
}

func (c *Client) WaitUntilDone(task Task, timeout time.Duration) (Task, error) {
	type Result struct {
		Task  Task
		Error error
	}
	doneCh := make(chan Result)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func(taskID int) {
		for range ticker.C {
			curTask, err := c.GetTask(taskID)
			if err != nil {
				log.Printf("Failed getting task %d status, retrying: %s", taskID, err)
			}
			switch curTask.State {
			case "done":
				doneCh <- Result{
					Task: curTask,
				}
				return
			case "error":
				doneCh <- Result{
					Task:  curTask,
					Error: fmt.Errorf("task %d failed: %s", curTask.ID, curTask.Result),
				}
				return
			}
		}
	}(task.ID)

	for {
		select {
		case result := <-doneCh:
			return result.Task, result.Error
		case <-time.After(timeout):
			return task, fmt.Errorf("timed out waiting for task %d to complete", task.ID)
		}
	}
}
