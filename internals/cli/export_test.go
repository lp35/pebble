// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"

	"github.com/canonical/go-flags"

	"github.com/canonical/pebble/client"
)

func RunMain() error {
	return Run(RunOptionsForTest())
}

func RunOptionsForTest() *RunOptions {
	o := &RunOptions{
		ClientConfig: newClientConfig,
	}
	return withDefaultRunOptions(o)
}

var clientConfigBaseURL string

func FakeClientConfigBaseURL(baseURL string) (restore func()) {
	clientConfigBaseURL = baseURL
	return func() {
		clientConfigBaseURL = ""
	}
}

func newClientConfig() (*client.Config, error) {
	config := client.Config{BaseURL: clientConfigBaseURL}
	return &config, nil
}

func Client() *client.Client {
	cfg, err := newClientConfig()
	if err != nil {
		panic("cannot build client config: " + err.Error())
	}
	cli, err := client.New(cfg)
	if err != nil {
		panic("cannot create client: " + err.Error())
	}
	return cli
}

var (
	CanUnicode      = canUnicode
	ColorTable      = colorTable
	MonoColorTable  = mono
	ColorColorTable = color
	NoEscColorTable = noesc

	MaybePresentWarnings = maybePresentWarnings

	MaybeCopyPebbleDir = maybeCopyPebbleDir

	WithDefaultRunOptions = withDefaultRunOptions
)

func FakeIsStdoutTTY(t bool) (restore func()) {
	oldIsStdoutTTY := isStdoutTTY
	isStdoutTTY = t
	return func() {
		isStdoutTTY = oldIsStdoutTTY
	}
}

func FakeIsStdinTTY(t bool) (restore func()) {
	oldIsStdinTTY := isStdinTTY
	isStdinTTY = t
	return func() {
		isStdinTTY = oldIsStdinTTY
	}
}

func PebbleMain() (exitCode int) {
	oldOsExit := osExit
	osExit = func(code int) {
		panic(&exitStatus{code})
	}
	defer func() {
		osExit = oldOsExit
		if v := recover(); v != nil {
			if e, ok := v.(*exitStatus); ok {
				exitCode = e.code
			} else {
				panic(v)
			}
		}
	}()
	if err := RunMain(); err != nil {
		fmt.Fprintf(Stderr, "error: %v\n", err)
		osExit(1)
	}
	return
}

func ParserForTest() *flags.Parser {
	runOpts := RunOptionsForTest()

	return Parser(&ParserOptions{
		Client: withClient{getClient: func() (*client.Client, error) {
			return Client(), nil
		}},
		GetClientConfig: runOpts.ClientConfig,
		PebbleDir:       runOpts.PebbleDir,
	})
}

// ParserWithClientForTest creates a test parser that uses the given client
// directly, for tests that need to inspect the specific client instance.
func ParserWithClientForTest(c *client.Client) *flags.Parser {
	runOpts := RunOptionsForTest()
	return Parser(&ParserOptions{
		Client: withClient{getClient: func() (*client.Client, error) {
			return c, nil
		}},
		GetClientConfig: runOpts.ClientConfig,
		PebbleDir:       runOpts.PebbleDir,
	})
}
