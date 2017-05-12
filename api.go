package gogobosh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
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
	task := Task{}
	r := c.NewRequest("DELETE", "/deployments/"+name+"?force=true")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting deleting deployment %v", err)
		return task, err
	}
	defer resp.Body.Close()
	url, _ := resp.Location()
	re, _ := regexp.Compile(`(\d+)$`)
	stringID := re.FindStringSubmatch(url.Path)
	id, err := strconv.Atoi(stringID[0])
	if err != nil {
		return task, err
	}
	task, err = c.GetTask(id)
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
		log.Printf("Error reading task request %v", resBody)
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
func (c *Client) GetTasks() ([]Task, error) {
	r := c.NewRequest("GET", "/tasks")
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

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int) []string {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID+"/output?type=result")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting task %v", err)
		return nil
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return nil
	}

	return strings.Split(string(resBody), "\n")
}
