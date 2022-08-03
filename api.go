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

// GetStemcells from given BOSH
func (c *Client) GetStemcells() ([]Stemcell, error) {
	r := c.NewRequest("GET", "/stemcells")
	resp, err := c.DoRequest(r)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error requesting stemcells  %v", err)
		return nil, err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading stemcells request %v", resBody)
		return nil, err
	}

	stemcells := []Stemcell{}
	err = json.Unmarshal(resBody, &stemcells)
	if err != nil {
		log.Printf("Error unmarshaling stemcells %v", err)
	}
	return stemcells, err
}

func (c *Client) UploadStemcell(url, sha1 string) (Task, error) {
	task := Task{}
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
		return task, err
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		log.Printf("Error requesting stemcell upload %v", err)
		return task, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return task, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}

// GetReleases from given BOSH
func (c *Client) GetReleases() ([]Release, error) {
	r := c.NewRequest("GET", "/releases")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting releases  %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading releases request %v", resBody)
		return nil, err
	}

	releases := []Release{}
	err = json.Unmarshal(resBody, &releases)
	if err != nil {
		log.Printf("Error unmarshaling releases %v", err)
	}
	return releases, err
}

func (c *Client) UploadRelease(url, sha1 string) (Task, error) {
	task := Task{}
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
		return task, err
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		log.Printf("Error requesting release upload %v", err)
		return task, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return task, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}

// GetDeployments from given BOSH
func (c *Client) GetDeployments() ([]Deployment, error) {
	r := c.NewRequest("GET", "/deployments")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployments  %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployments request %v", resBody)
		return nil, err
	}
	deployments := []Deployment{}
	err = json.Unmarshal(resBody, &deployments)
	if err != nil {
		log.Printf("Error unmarshaling deployments %v", err)
	}
	return deployments, err
}

// GetDeployment from given BOSH
func (c *Client) GetDeployment(name string) (Manifest, error) {
	manifest := Manifest{}
	r := c.NewRequest("GET", "/deployments/"+name)
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployment manifest %v", err)
		return manifest, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment manifest request %v", resBody)
		return manifest, err
	}
	err = json.Unmarshal(resBody, &manifest)
	if err != nil {
		log.Printf("Error unmarshaling deployment manifest %v", err)
	}
	return manifest, err
}

// DeleteDeployment from given BOSH
func (c *Client) DeleteDeployment(name string) (Task, error) {
	var task Task
	resp, err := c.DoRequest(c.NewRequest("DELETE", "/deployments/"+name))
	if err != nil {
		log.Printf("Error requesting deleting deployment %v", err)
		return task, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return task, fmt.Errorf("deployment %s not found", name)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task response %v", err)
		return task, err
	}
	err = json.Unmarshal(b, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}

// CreateDeployment from given BOSH
func (c *Client) CreateDeployment(manifest string) (Task, error) {
	task := Task{}
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.body = buffer
	r.header["Content-Type"] = "text/yaml"

	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting create deployment %v", err)
		return task, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task response %v", resBody)
		return task, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}

// GetDeploymentVMs from given BOSH
func (c *Client) GetDeploymentVMs(name string) ([]VM, error) {
	var task Task
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deployment vms %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment vms request %v", resBody)
		return nil, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
		return nil, err
	}
	for {
		taskStatus, err := c.GetTask(task.ID)
		if err != nil {
			log.Printf("Error getting task %v", err)
		}
		if taskStatus.State == "done" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	vms := []VM{}
	output := c.GetTaskResult(task.ID)
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				log.Printf("Error unmarshaling vms %v %v", value, err)
				return nil, err
			}
			vms = append(vms, vm)
		}
	}
	return vms, err
}

// GetTasks from given BOSH
func (c *Client) GetTasksByQuery(query url.Values) ([]Task, error) {
	requestUrl := "/tasks?" + query.Encode()
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting tasks  %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading tasks request %v", resBody)
		return nil, err
	}
	tasks := []Task{}
	err = json.Unmarshal(resBody, &tasks)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
	}
	return tasks, err
}

func (c *Client) GetTasks() ([]Task, error) {
	return c.GetTasksByQuery(nil)
}

// GetTask from given BOSH
func (c *Client) GetTask(id int) (Task, error) {
	task := Task{}
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting task %v", err)
		return task, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return task, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}

// GetTaskOutput ...
func (c *Client) GetTaskOutput(id int, typ string) ([]string, error) {
	r := c.NewRequest("GET", "/tasks/"+strconv.Itoa(id)+"/output?type="+typ)

	res, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"), nil
}

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int) []string {
	l, _ := c.GetTaskOutput(id, "result")
	return l
}

func (c *Client) GetTaskEvents(id int) ([]TaskEvent, error) {
	raw, err := c.GetTaskOutput(id, "event")
	if err != nil {
		return nil, err
	}

	events := make([]TaskEvent, len(raw))
	for i := range raw {
		err = json.Unmarshal([]byte(raw[i]), &events[i])
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

// GetCloudConfig from given BOSH
func (c *Client) GetCloudConfig(latest bool) ([]Cfg, error) {
	cfg := []Cfg{}

	qs := "?latest=true"
	if !latest {
		qs = "?latest=false"
	}
	r := c.NewRequest("GET", "/configs"+qs)
	resp, err := c.DoRequest(r)
	if err != nil {
		return cfg, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cfg, err
	}
	return cfg, json.Unmarshal(resBody, &cfg)
}

// UpdateCloudConfig
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
		return err
	}

	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"

	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//Cleanup will post to the cleanup endpoint of bosh, passing along the removeall flag passed in as a bool
func (c *Client) Cleanup(removeall bool) (Task, error) {
	task := Task{}
	r := c.NewRequest("POST", "/cleanup")
	var requestBody struct {
		Config struct {
			RemoveAll bool `json:"remove_all"`
		} `json:"config"`
	}
	requestBody.Config.RemoveAll = removeall
	b, err := json.Marshal(&requestBody)
	if err != nil {
		return task, err
	}
	r.body = bytes.NewBuffer(b)
	r.header["Content-Type"] = "application/json"
	resp, err := c.DoRequest(r)
	if err != nil {
		return task, err
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return task, err
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
	}
	return task, err
}
