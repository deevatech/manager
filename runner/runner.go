package runner

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	. "github.com/deevatech/manager/types"
	"github.com/fsouza/go-dockerclient"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Context struct {
	Client            *docker.Client
	ImageName         string
	ImageTag          string
	Args              []string
	RunnerPort        docker.Port
	ContainerID       string
	ContainerHostPort string
}

var context *Context

func init() {
	// Make sure we download the image before we get started
	context = NewContext()

	if err := context.PullImage(); err != nil {
		log.Fatal(err)
	}
}

func NewClient() (*docker.Client, error) {
	log.Println("Checking docker environment...")

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Println("Connecting to Docker failed:", err)
		return nil, err
	}

	return client, nil
}

func NewContext() *Context {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	return &Context{
		Client:     client,
		ImageName:  "deeva/runner-service",
		ImageTag:   "latest",
		Args:       []string{},
		RunnerPort: "8080",
	}
}

func Run(p RunParams) (*RunResults, error) {
	return context.Start(p)
}

func (ctx Context) PullImage() error {
	skipPull := os.Getenv("DEEVA_MANAGER_SKIP_PULL")
	if len(skipPull) != 0 {
		log.Println("Skipping pull image...")
		return nil
	}

	opts := docker.PullImageOptions{
		Repository:   ctx.ImageName,
		Registry:     "hub.docker.com",
		Tag:          ctx.ImageTag,
		OutputStream: os.Stdout,
	}

	log.Printf("Downloading image %s:%s from %s...", opts.Repository, opts.Tag, opts.Registry)

	if err := ctx.Client.PullImage(opts, docker.AuthConfiguration{}); err != nil {
		log.Printf("Failed pulling image %s : %s", ctx.ImageName, err)
		return err
	}

	return nil
}

func (ctx *Context) Start(p RunParams) (*RunResults, error) {
	log.Printf("Starting container with context %v", ctx)

	// create new container
	container, err := ctx.Client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: ctx.ImageName,
			Cmd:   ctx.Args,
		},
	})
	if err != nil {
		log.Printf("Failed creating container: %s", err)
		return nil, err
	} else {
		ctx.ContainerID = container.ID

		defer func() {
			if err := ctx.Client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}); err != nil {
				log.Println(err)
				return
			}
			log.Printf("Removed container %s", container.ID)
		}()
	}
	log.Printf("Created container: %s", container.ID)

	if err = ctx.Client.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		log.Printf("Starting container %+v ... failed: %v", container.ID, err)
		return nil, err
	} else {
		defer func() {
			if err := ctx.Client.StopContainer(container.ID, 0); err != nil {
				log.Println(err)
				return
			}
			log.Printf("Stopped container %s", container.ID)
		}()
	}
	log.Printf("Started container: %+v", spew.Sdump(container))

	// wait for container to wake up
	if err := waitStarted(ctx.Client, container.ID, 1*time.Second); err != nil {
		log.Printf("Couldn't reach runner container %s, aborting!", container.ID)
		return nil, err
	}
	log.Println("Container started!")

	if container, err = ctx.Client.InspectContainer(container.ID); err != nil {
		log.Printf("Couldn't inspect runner container %s, aborting!", container.ID)
		return nil, err
	}
	log.Println("Container inspected!")

	// determine IP address
	containerIP := strings.TrimSpace(container.NetworkSettings.IPAddress)

	// wait for runner-service to wake up
	hostport := fmt.Sprintf("%s:%s", containerIP, ctx.RunnerPort)
	if err := waitReachable(hostport, 1*time.Second); err != nil {
		log.Printf("Couldn't reach runner application in container %s via %s, aborting!", container.ID, hostport)
		return nil, err
	} else {
		ctx.ContainerHostPort = hostport
	}
	log.Println("Container reached!")

	return context.Request(p)
}

// waitReachable waits for hostport to became reachable for the maxWait time.
func waitReachable(hostport string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)

	for time.Now().Before(done) {
		if c, err := net.Dial("tcp", hostport); err == nil {
			c.Close()
			return nil
		}

		time.Sleep(20 * time.Millisecond)
	}

	return fmt.Errorf("cannot connect %v for %v", hostport, maxWait)
}

// waitStarted waits for a container to start for the maxWait time.
func waitStarted(client *docker.Client, id string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)

	for time.Now().Before(done) {
		c, err := client.InspectContainer(id)

		if err != nil {
			break
		}

		if c.State.Running {
			return nil
		}

		time.Sleep(20 * time.Millisecond)
	}

	return fmt.Errorf("cannot start container %s for %v", id, maxWait)
}
