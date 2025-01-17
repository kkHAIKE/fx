// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPSignalE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package fx

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syscall"
	"testing"
)

func assertUnsentSignalError(
	t *testing.T,
	err error,
	expected *unsentSignalError,
) {
	t.Helper()

	actual := new(unsentSignalError)

	assert.ErrorContains(t, err, "channels are blocked")
	if assert.ErrorAs(t, err, &actual, "is unsentSignalError") {
		assert.Equal(t, expected, actual)
	}
}

func TestSignal(t *testing.T) {
	t.Parallel()
	recv := newSignalReceivers()
	a := recv.Done()
	_ = recv.Done() // we never listen on this

	expected := ShutdownSignal{
		Signal: syscall.SIGTERM,
	}

	require.NoError(t, recv.Broadcast(expected), "first broadcast should succeed")

	assertUnsentSignalError(t, recv.Broadcast(expected), &unsentSignalError{
		Signal: expected,
		Total:  2,
		Unsent: 2,
	})

	assert.Equal(t, expected.Signal, <-a)
	assert.Equal(t, expected.Signal, <-recv.Done(), "expect cached signal")
}
