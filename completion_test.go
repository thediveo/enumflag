// Copyright 2023 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package enumflag

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "github.com/thediveo/success"
	"golang.org/x/exp/slices"
)

const dummyCommandName = "enumflag-testing"

// See:
// https://serverfault.com/questions/506612/standard-place-for-user-defined-bash-completion-d-scripts/1013395#1013395
const bashComplDirEnv = "BASH_COMPLETION_USER_DIR"

var _ = FDescribe("flag enum completions", Ordered, func() {

	var enumflagTestingPath string
	var completionsDir string

	BeforeAll(func() {
		By("building a CLI binary for testing")
		enumflagTestingPath = Successful(gexec.Build("./test/enumflag-testing"))
		DeferCleanup(func() {
			gexec.CleanupBuildArtifacts()
		})

		By("creating a temporary directory for storing completion scripts")
		completionsDir = Successful(os.MkdirTemp("", "bash-completions-*"))
		DeferCleanup(func() {
			os.RemoveAll(completionsDir)
		})

		By("telling the CLI binary to give us a completion script and storing it away")
		session := Successful(
			gexec.Start(exec.Command(enumflagTestingPath, "completion", "bash"),
				GinkgoWriter, GinkgoWriter))
		Eventually(session).Within(5 * time.Second).ProbeEvery(100 * time.Millisecond).
			Should(gexec.Exit(0))
		completionScript := session.Out.Contents()
		Expect(completionScript).To(MatchRegexp(`^# bash completion V2 for`))
		Expect(os.WriteFile(filepath.Join(completionsDir, "enumflag-testing"),
			completionScript, 0770)).
			To(Succeed())
	})

	Bash := func() (*gexec.Session, io.WriteCloser) {
		GinkgoHelper()
		By("creating a new test bash session")
		bashCmd := exec.Command("/bin/bash", "--rcfile", "/etc/profile", "-i")
		// Run the silly interactive subshell in its own session so we don't get
		// funny surprises such as the subshell getting suspended...
		bashCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		bashCmd.Env = append(slices.Clone(os.Environ()),
			bashComplDirEnv+"="+completionsDir,
			"PATH="+filepath.Dir(enumflagTestingPath)+":"+os.Getenv("PATH"),
		)
		stdin := Successful(bashCmd.StdinPipe())
		session := Successful(
			gexec.Start(bashCmd, GinkgoWriter, GinkgoWriter))
		DeferCleanup(func() {
			By("killing the test bash session")
			session.Kill().Wait(2 * time.Second)
		})
		return session, stdin
	}

	It("tests the test", func() {
		bash, bashIn := Bash()
		By("listing the canary in the first search PATH directory")
		Expect(bashIn.Write([]byte("ls -l ${PATH%%:*}\n"))).Error().NotTo(HaveOccurred())
		Eventually(bash.Out).Should(gbytes.Say(dummyCommandName))

		By("ensuring the canary is in the PATH and gets completed")
		Expect(bashIn.Write([]byte(dummyCommandName[:len(dummyCommandName)-4] + "\t"))).Error().NotTo(HaveOccurred())
		Eventually(bash.Err).Should(gbytes.Say(dummyCommandName))
	})

	It("completes test command", func() {
		bash, bashIn := Bash()

		Expect(bashIn.Write([]byte(dummyCommandName + " t\t"))).Error().NotTo(HaveOccurred())
		Eventually(bash.Err).Should(gbytes.Say(dummyCommandName + " test"))
	})

})
