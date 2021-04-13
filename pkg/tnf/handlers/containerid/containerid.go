// Copyright (C) 2021 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package containerid

import (
	"regexp"
	"time"

	"github.com/test-network-function/test-network-function/pkg/tnf"
	"github.com/test-network-function/test-network-function/pkg/tnf/dependencies"
	"github.com/test-network-function/test-network-function/pkg/tnf/identifier"
	"github.com/test-network-function/test-network-function/pkg/tnf/reel"
)

// ContainerID provides a way to find an id of a container from inside of it.
type ContainerID struct {
	result  int
	timeout time.Duration
	args    []string
	id      string
}

const (
	// SuccessfulOutputRegex matches a cgroup name that should be generated by crio and includes the container id
	// inside of it in a known location
	SuccessfulOutputRegex = `crio-(\w+)\.scope`
)

// Args returns the command line args for the test.
func (id *ContainerID) Args() []string {
	return id.args
}

// GetIdentifier returns the tnf.Test specific identifier.
func (id *ContainerID) GetIdentifier() identifier.Identifier {
	return identifier.ContainerIDIdentifier
}

// Timeout returns the timeout in seconds for the test.
func (id *ContainerID) Timeout() time.Duration {
	return id.timeout
}

// Result returns the test result.
func (id *ContainerID) Result() int {
	return id.result
}

// ReelFirst returns a step which expects the container id within the test timeout.
func (id *ContainerID) ReelFirst() *reel.Step {
	return &reel.Step{
		Expect:  []string{SuccessfulOutputRegex},
		Timeout: id.timeout,
	}
}

// ReelMatch parses the the result of "/proc/self/cgroup" looking for a cgroup generated by crio
// and resolve the container id from it
func (id *ContainerID) ReelMatch(_, _, match string) *reel.Step {
	re := regexp.MustCompile(SuccessfulOutputRegex)
	matched := re.FindStringSubmatch(match)
	if matched != nil {
		id.result = tnf.SUCCESS
		id.id = matched[1]
	} else {
		id.result = tnf.FAILURE
	}
	return nil
}

// ReelTimeout returns a step which kills the container id test by sending it ^C.
func (id *ContainerID) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;  container id requires no intervention on eof.
func (id *ContainerID) ReelEOF() {
}

// Command returns command line args for getting the cgroups of the host machine
func Command() []string {
	return []string{"cat", dependencies.CgroupProcfsPath}
}

// NewContainerID creates a new container id test which lists all cgroups of the host from inside a pod
// and resolve the id of the container itself from it
func NewContainerID(timeout time.Duration) *ContainerID {
	return &ContainerID{
		result:  tnf.ERROR,
		timeout: timeout,
		args:    Command(),
	}
}

// GetID returns the container id
func (id *ContainerID) GetID() string {
	return id.id
}