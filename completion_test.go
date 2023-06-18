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

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"golang.org/x/exp/slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

const dummyCommandName = "enumflag-testing"

// See:
// https://serverfault.com/questions/506612/standard-place-for-user-defined-bash-completion-d-scripts/1013395#1013395
const bashComplDirEnv = "BASH_COMPLETION_USER_DIR"

type writer struct {
	io.WriteCloser
}

func (w *writer) WriteString(s string) {
	GinkgoHelper()
	Expect(w.WriteCloser.Write([]byte(s))).Error().NotTo(HaveOccurred())
}

var _ = FDescribe("flag enum completions end-to-end", Ordered, func() {

	var enumflagTestingPath string
	var completionsUserDir string

	BeforeAll(func() {
		By("building a CLI binary for testing")
		enumflagTestingPath = Successful(gexec.Build("./test/enumflag-testing"))
		DeferCleanup(func() {
			gexec.CleanupBuildArtifacts()
		})

		By("creating a temporary directory for storing completion scripts")
		completionsUserDir = Successful(os.MkdirTemp("", "bash-completions-*"))
		DeferCleanup(func() {
			os.RemoveAll(completionsUserDir)
		})
		// Notice how the bash-completion FAQ
		// https://github.com/scop/bash-completion/blob/master/README.md#faq
		// says that the completions must be inside a "completions" sub
		// directory of $BASH_COMPLETION_USER_DIR, and not inside
		// $BASH_COMPLETION_USER_DIR itself ... yeah, ¯\_(ツ)_/¯
		Expect(os.Mkdir(filepath.Join(completionsUserDir, "completions"), 0770)).To(Succeed())

		By("telling the CLI binary to give us a completion script that we then store away")
		session := Successful(
			gexec.Start(exec.Command(enumflagTestingPath, "completion", "bash"),
				GinkgoWriter, GinkgoWriter))
		Eventually(session).Within(5 * time.Second).ProbeEvery(100 * time.Millisecond).
			Should(gexec.Exit(0))
		completionScript := session.Out.Contents()
		Expect(completionScript).To(MatchRegexp(`^# bash completion V2 for`))
		Expect(os.WriteFile(filepath.Join(completionsUserDir, "completions", "enumflag-testing"),
			completionScript, 0770)).
			To(Succeed())
	})

	Bash := func() (*gexec.Session, *writer) {
		GinkgoHelper()
		By("creating a new test bash session")
		bashCmd := exec.Command("/bin/bash", "--rcfile", "/etc/profile", "-i")
		// Run the silly interactive subshell in its own session so we don't get
		// funny surprises such as the subshell getting suspended...
		bashCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		bashCmd.Env = append(slices.Clone(os.Environ()),
			bashComplDirEnv+"="+completionsUserDir,
			"PATH="+filepath.Dir(enumflagTestingPath)+":"+os.Getenv("PATH"),
		)
		stdin := &writer{Successful(bashCmd.StdinPipe())}
		session := Successful(
			gexec.Start(bashCmd, GinkgoWriter, GinkgoWriter))
		DeferCleanup(func() {
			By("killing the test bash session")
			stdin.Close()
			session.Kill().Wait(2 * time.Second)
		})
		return session, stdin
	}

	var bash *gexec.Session
	var bashin *writer

	BeforeEach(func() {
		bash, bashin = Bash()
	})

	It("tab-completes the canary's name in $PATH", func() {
		By("checking BASH_COMPLETION_USER_DIR")
		bashin.WriteString("echo $" + bashComplDirEnv + "\n")
		Eventually(bash.Out).Should(gbytes.Say(completionsUserDir))

		By("listing the canary in the first search PATH directory")
		bashin.WriteString("ls -l ${PATH%%:*}\n")
		Eventually(bash.Out).Should(gbytes.Say(dummyCommandName))

		By("ensuring the canary is in the PATH and gets completed")
		bashin.WriteString(dummyCommandName[:len(dummyCommandName)-4] + "\t")
		Eventually(bash.Err).Should(gbytes.Say(dummyCommandName))
	})

	It("completes canary's test subcommand", func() {
		bashin.WriteString(dummyCommandName + " t\t")
		Eventually(bash.Err).Should(gbytes.Say(dummyCommandName + " test"))
	})

	It("completes canary's \"mode\" enum flag name", func() {
		bashin.WriteString(dummyCommandName + " test --\t\t")
		Eventually(bash.Err).Should(gbytes.Say(
			`--help\s+\(help for test\)\s+--mode\s+\(sets foo mode\)`))
	})

	It("lists enum flag's values", func() {
		bashin.WriteString(dummyCommandName + " test --mode \t\t")
		Eventually(bash.Err).Should(gbytes.Say(
			`bar\s+\(bars the output\)\s+baz\s+\(bazs the output\)\s+foo\s+\(foos the output\)`))
		bashin.WriteString("\b=\t\t")
		Eventually(bash.Err).Should(gbytes.Say(
			`bar\s+\(bars the output\)\s+baz\s+\(bazs the output\)\s+foo\s+\(foos the output\)`))
	})

	It("completes enum flag's values", func() {
		bashin.WriteString(dummyCommandName + " test --mode ba\t\t")
		Eventually(bash.Err).Should(gbytes.Say(
			`bar\s+\(bars the output\)\s+baz\s+\(bazs the output\)`))
		bashin.WriteString("\b\bf\t")
		Eventually(bash.Err).Should(gbytes.Say(
			`oo `))
	})

})
