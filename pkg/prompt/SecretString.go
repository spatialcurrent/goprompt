// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package prompt

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func SecretString(question string, stdout bool) (string, error) {
	if len(question) > 0 {
		if stdout {
			_, _ = fmt.Fprintf(os.Stdout, "%s: ", question)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s: ", question)
		}
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return "", fmt.Errorf("error making terminal raw: %w", err)
	}

	value := ""
	var inputErr error

	wg := &sync.WaitGroup{}
	wg.Add(1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)
	go func() {
		s := <-signals
		signal.Stop(signals)
		inputErr = fmt.Errorf("received signal %q", s)
		// stop waiting
		wg.Done()
	}()

	go func() {
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			inputErr = fmt.Errorf("error reading secret string from terminal: %w", err)
		} else {
			value = string(bytes.TrimSpace(b))
		}
		// stop waiting
		wg.Done()
	}()

	wg.Wait()

	if oldState != nil {
		_ = terminal.Restore(0, oldState)
	}

	// add new line
	if stdout {
		_, _ = fmt.Fprintln(os.Stdout, "")
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "")
	}

	return value, inputErr
}
