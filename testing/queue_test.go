// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"github.com/globocom/tsuru/queue"
	. "launchpad.net/gocheck"
)

func (s *S) TestFakeQPutAndGet(c *C) {
	q := FakeQ{}
	msg := queue.Message{Action: "do-something"}
	err := q.Put(&msg, 0)
	c.Assert(err, IsNil)
	m, err := q.Get(1e6)
	c.Assert(err, IsNil)
	c.Assert(m.Action, Equals, msg.Action)
}

func (s *S) TestFakeQGetTimeout(c *C) {
	q := FakeQ{}
	m, err := q.Get(1e6)
	c.Assert(m, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Timed out.")
}

func (s *S) TestFakeQPutWithTimeout(c *C) {
	q := FakeQ{}
	msg := queue.Message{Action: "do-something"}
	err := q.Put(&msg, 1e6)
	c.Assert(err, IsNil)
	_, err = q.Get(1e3)
	c.Assert(err, NotNil)
	_, err = q.Get(1e9)
	c.Assert(err, IsNil)
}

func (s *S) TestFakeQDelete(c *C) {
	q := FakeQ{}
	err := q.Delete(nil)
	c.Assert(err, IsNil)
}

func (s *S) TestFakeQRelease(c *C) {
	q := FakeQ{}
	msg := queue.Message{Action: "do-something"}
	err := q.Put(&msg, 0)
	c.Assert(err, IsNil)
	m, err := q.Get(1e6)
	c.Assert(err, IsNil)
	err = q.Release(m, 0)
	c.Assert(err, IsNil)
	_, err = q.Get(1e6)
	c.Assert(err, IsNil)
}

func (s *S) TestFakeHandlerStart(c *C) {
	h := fakeHandler{}
	c.Assert(h.running, Equals, false)
	h.Start()
	c.Assert(h.running, Equals, true)
}

func (s *S) TestFakeHandlerStop(c *C) {
	h := fakeHandler{}
	h.Start()
	h.Stop()
	c.Assert(h.running, Equals, false)
}
