package util

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	defaultName                   = ""
	defaultPort            uint16 = 8091
	defaultRemoteName             = ""
	defaultRemotePort      uint16 = 8090
	defaultTimeout                = time.Second * 5
	defaultRemoteActorName        = "tree-service"
)

// Flags the CLI is able to understand.
type Flags struct {
	// Token of the tree service session
	Token string

	// Id of the tree to work with
	Id int

	// Name to bind the CLI to
	Name string

	// Port to bind the CLI to
	Port uint16

	// Name of the remote service to connect to
	RemoteName string

	// Port of the remove service to connect to
	RemotePort uint16

	// Name of the remote service actor
	RemoteActorName string

	// Timeout of the command execution
	Timeout time.Duration
}

// Get all program flags received via command line.
func GetProgramFlags() *Flags {
	token := flag.String(
		"token",
		"",
		"Token of the tree service session",
	)
	id := flag.Int(
		"id",
		-1,
		"Id of the tree to work with",
	)

	bind := flag.String(
		"bind",
		fmt.Sprintf("%s:%d", defaultName, defaultPort),
		"the name and port to bind the CLI to",
	)
	remote := flag.String(
		"remote",
		fmt.Sprintf("%s:%d", defaultRemoteName, defaultRemotePort),
		"the name and port of the service to connect to",
	)
	remoteActorName := flag.String(
		"remote-name",
		defaultRemoteActorName,
		"the name of the tree service remote actor",
	)

	timeout := flag.Duration(
		"timeout",
		defaultTimeout,
		"After what time the command execution will be cancelled",
	)

	flag.Parse()

	name, port, e := parseNamePort(bind)
	if e != nil {
		log.Printf("Could not understand argument --bind=\"%s\". Using default instead.", *bind)
		name, port = defaultName, defaultPort
	}

	remoteName, remotePort, e := parseNamePort(remote)
	if e != nil {
		log.Printf("Could not understand argument --remote=\"%s\". Using default instead.", *remote)
		remoteName, remotePort = defaultRemoteName, defaultRemotePort
	}

	return &Flags{
		Token:           *token,
		Id:              *id,
		Name:            name,
		Port:            port,
		RemoteName:      remoteName,
		RemotePort:      remotePort,
		RemoteActorName: *remoteActorName,
		Timeout:         *timeout,
	}
}

// Parse name and port of the passed string.
// The string needs the form "{name}:{port}".
func parseNamePort(srcPtr *string) (string, uint16, error) {
	src := *srcPtr

	parts := strings.Split(src, ":")
	if len(parts) != 2 {
		return "", 0, errors.New("expected string in format \"{name}:{port}\"")
	}

	port, e := strconv.ParseUint(parts[1], 10, 16)
	if e != nil {
		return "", 0, errors.New(fmt.Sprintf("could not parse port %s", parts[1]))
	}

	return parts[0], uint16(port), nil
}

func (f *Flags) String() string {
	return fmt.Sprintf("{ token: %s, id: %v }", f.Token, f.Id)
}
