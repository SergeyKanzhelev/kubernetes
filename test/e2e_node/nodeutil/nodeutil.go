/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package nodeutil contains shared test utilities for node e2e tests.
// It exists so that test subpackages (e.g. standalone/) can reuse helpers
// without importing the main e2enode package (which would create a cycle).
package nodeutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/onsi/gomega"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	internalapi "k8s.io/cri-api/pkg/apis"
	remote "k8s.io/cri-client/pkg"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/cluster/ports"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/test/e2e/framework"
)

// KubeletHealthCheckURL is the URL to check kubelet health.
var KubeletHealthCheckURL = fmt.Sprintf("http://127.0.0.1:%d/healthz", ports.KubeletHealthzPort)

// KubeletCfg holds the kubelet configuration the test is running against.
// It is populated by the test suite setup (SynchronizedBeforeSuite) and
// should be treated as read-only by tests.
var KubeletCfg *kubeletconfig.KubeletConfiguration

// StaticPodPath returns the file path for a static pod manifest.
func StaticPodPath(dir, name, namespace string) string {
	return filepath.Join(dir, namespace+"-"+name+".yaml")
}

// DeleteStaticPod removes a static pod manifest file.
func DeleteStaticPod(dir, name, namespace string) error {
	file := StaticPodPath(dir, name, namespace)
	return os.Remove(file)
}

// FindKubeletServiceName searches the unit name among the services known to systemd.
// If the running parameter is true, restricts the search among currently running services;
// otherwise, also stopped, failed, exited (non-running in general) services are also considered.
func FindKubeletServiceName(running bool) string {
	cmdLine := []string{
		"systemctl", "list-units", "*kubelet*",
	}
	if running {
		cmdLine = append(cmdLine, "--state=running")
	}
	stdout, err := exec.Command("sudo", cmdLine...).CombinedOutput()
	framework.ExpectNoError(err)
	regex := regexp.MustCompile(`(kubelet-\w+)`)
	matches := regex.FindStringSubmatch(string(stdout))
	gomega.Expect(matches).ToNot(gomega.BeEmpty(), "Found more than one kubelet service running: %q", stdout)
	kubeletServiceName := matches[0]
	framework.Logf("Get running kubelet with systemctl: %v, %v", string(stdout), kubeletServiceName)
	return kubeletServiceName
}

// RestartKubelet restarts the kubelet service via systemctl.
func RestartKubelet(ctx context.Context, running bool) {
	kubeletServiceName := FindKubeletServiceName(running)
	// reset the kubelet service start-limit-hit
	stdout, err := exec.CommandContext(ctx, "sudo", "systemctl", "reset-failed", kubeletServiceName).CombinedOutput()
	framework.ExpectNoError(err, "Failed to reset kubelet start-limit-hit with systemctl: %v, %s", err, string(stdout))

	stdout, err = exec.CommandContext(ctx, "sudo", "systemctl", "restart", kubeletServiceName).CombinedOutput()
	framework.ExpectNoError(err, "Failed to restart kubelet with systemctl: %v, %s", err, string(stdout))
}

// GetCRIClient connects CRI and returns CRI runtime service clients and image service client.
func GetCRIClient(ctx context.Context) (internalapi.RuntimeService, internalapi.ImageManagerService, error) {
	// connection timeout for CRI service connection
	const connectionTimeout = 2 * time.Minute
	runtimeEndpoint := framework.TestContext.ContainerRuntimeEndpoint
	useStreaming := utilfeature.DefaultFeatureGate.Enabled(features.CRIListStreaming)
	r, err := remote.NewRemoteRuntimeServiceBuilder().
		WithEndpoint(runtimeEndpoint).
		WithConnectionTimeout(connectionTimeout).
		WithUseStreaming(useStreaming).
		Build(ctx)
	if err != nil {
		return nil, nil, err
	}
	imageManagerEndpoint := runtimeEndpoint
	if framework.TestContext.ImageServiceEndpoint != "" {
		//ImageServiceEndpoint is the same as ContainerRuntimeEndpoint if not
		//explicitly specified
		imageManagerEndpoint = framework.TestContext.ImageServiceEndpoint
	}
	i, err := remote.NewRemoteImageServiceBuilder().
		WithEndpoint(imageManagerEndpoint).
		WithConnectionTimeout(connectionTimeout).
		WithUseStreaming(useStreaming).
		Build(ctx)
	if err != nil {
		return nil, nil, err
	}
	return r, i, nil
}

// RemoveInitContainer stops and removes a container by its CRI container ID.
func RemoveInitContainer(ctx context.Context, ctrID string) {
	cricli, _, err := GetCRIClient(ctx)
	framework.ExpectNoError(err)
	splitID := strings.Split(ctrID, "://")
	gomega.Expect(splitID).To(gomega.HaveLen(2))
	ctrID = splitID[1]
	// Make sure the container is stopped before removing it. This may fail.
	_ = cricli.StopContainer(ctx, ctrID, 0)
	err = cricli.RemoveContainer(ctx, ctrID)
	framework.ExpectNoError(err)
}
