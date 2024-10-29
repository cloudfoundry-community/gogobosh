package gogobosh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetStemcells from given BOSH
func (c *Client) GetStemcells() ([]Stemcell, error) {
	r := c.NewRequest("GET", "/stemcells")
	var stemcells []Stemcell
	err := c.DoRequestAndUnmarshal(r, &stemcells)
	if err != nil {
		return []Stemcell{}, fmt.Errorf("error getting stemcells: %w", err)
	}
	return stemcells, nil
}

// UploadStemcell to the given BOSH
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

	var task Task
	err = c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error uploading stemcell %s: %w", url, err)
	}
	return task, nil
}

// GetReleases from the given BOSH
func (c *Client) GetReleases() ([]Release, error) {
	r := c.NewRequest("GET", "/releases")
	var releases []Release
	err := c.DoRequestAndUnmarshal(r, &releases)
	if err != nil {
		return []Release{}, fmt.Errorf("error requesting releases: %w", err)
	}
	return releases, nil
}

// UploadRelease to the given BOSH
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

	var task Task
	err = c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error uploading release: %w", err)
	}
	return task, nil
}

// GetDeployments returns all deployments from the given BOSH
func (c *Client) GetDeployments() ([]Deployment, error) {
	r := c.NewRequest("GET", "/deployments")
	var deployments []Deployment
	err := c.DoRequestAndUnmarshal(r, &deployments)
	if err != nil {
		return []Deployment{}, fmt.Errorf("error requesting deployments: %w", err)
	}
	return deployments, nil
}

// GetDeployment returns a specific deployment by name from the given BOSH
func (c *Client) GetDeployment(name string) (Manifest, error) {
	r := c.NewRequest("GET", "/deployments/"+name)
	var manifest Manifest
	err := c.DoRequestAndUnmarshal(r, &manifest)
	if err != nil {
		return Manifest{}, fmt.Errorf("error requesting deployment manifest: %w", err)
	}
	return manifest, nil
}

// DeleteDeployment from given BOSH
func (c *Client) DeleteDeployment(name string) (Task, error) {
	r := c.NewRequest("DELETE", "/deployments/"+name)
	var task Task
	err := c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error deleting deployment %s: %w", name, err)
	}
	return task, nil
}

// CreateDeployment deploys the given deployment manifest
func (c *Client) CreateDeployment(manifest string) (Task, error) {
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.body = buffer
	r.header["Content-Type"] = "text/yaml"

	var task Task
	err := c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error creating deployment: %w", err)
	}
	return task, nil
}

// GetDeploymentVMs returns all the VMs that make up the specified deployment
func (c *Client) GetDeploymentVMs(name string) ([]VM, error) {
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	var task Task
	err := c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return []VM{}, fmt.Errorf("error requesting deployment %s VMs: %w", name, err)
	}

	task, err = c.WaitUntilDone(task, time.Minute*5)
	if err != nil {
		return []VM{}, fmt.Errorf("error waiting for deployment %s VM task to complete: %w", name, err)
	}

	var vms []VM
	output, err := c.GetTaskResult(task.ID)
	if err != nil {
		return []VM{}, fmt.Errorf("error getting deployment %s VMs task result: %w", name, err)
	}
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				return []VM{}, fmt.Errorf("error unmarshalling deployment %s VMs response: %w", name, err)
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
	var tasks []Task
	err := c.DoRequestAndUnmarshal(r, &tasks)
	if err != nil {
		return []Task{}, fmt.Errorf("error requesting tasks by query %s: %w", query.Encode(), err)
	}
	return tasks, nil
}

// GetTasks returns all BOSH tasks
func (c *Client) GetTasks() ([]Task, error) {
	return c.GetTasksByQuery(nil)
}

// GetTask returns the specified task from BOSH
func (c *Client) GetTask(id int) (Task, error) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	var task Task
	err := c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error getting task %s: %w", stringID, err)
	}
	return task, nil
}

// GetTaskOutput returns the completed tasks output
func (c *Client) GetTaskOutput(id int, typ string) ([]string, error) {
	r := c.NewRequest("GET", "/tasks/"+strconv.Itoa(id)+"/output?type="+typ)

	res, err := c.DoRequest(r)
	if err != nil {
		return []string{}, fmt.Errorf("error requesting task output: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, fmt.Errorf("error reading task output response: %w", err)
	}

	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"), nil
}

// GetTaskResult returns the tasks result
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
	var cfg []Cfg
	err := c.DoRequestAndUnmarshal(r, &cfg)
	if err != nil {
		return []Cfg{}, fmt.Errorf("error cloud config: %w", err)
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

	var task Task
	err = c.DoRequestAndUnmarshal(r, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error cleaning up BOSH: %w", err)
	}
	return task, nil
}

func (c *Client) Restart(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("restart", deployment, instanceGroup, instanceID, true)
}

func (c *Client) RestartNoConverge(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("restart", deployment, instanceGroup, instanceID, false)
}

func (c *Client) Stop(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("stopped", deployment, instanceGroup, instanceID, true)
}

func (c *Client) StopNoConverge(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("stopped", deployment, instanceGroup, instanceID, false)
}

func (c *Client) Start(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("started", deployment, instanceGroup, instanceID, true)
}

func (c *Client) StartNoConverge(deployment, instanceGroup, instanceID string) (Task, error) {
	return c.vmAction("started", deployment, instanceGroup, instanceID, false)
}

func (c *Client) vmAction(action, deployment, instanceGroup, instanceID string, converge bool) (Task, error) {
	var p string
	if converge {
		p = fmt.Sprintf("/deployments/%s/jobs/%s/%s?state=%s",
			deployment, instanceGroup, instanceID, action)
	} else {
		p = fmt.Sprintf("/deployments/%s/instance_groups/%s/%s/actions/%s",
			deployment, instanceGroup, instanceID, action)
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
				doneCh <- Result{
					Task:  Task{},
					Error: fmt.Errorf("error getting task %d status: %w", curTask.ID, err),
				}
				return
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
