/*
Copyright 2018 Sysdig.

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

package cache

import (
	"time"
	"sync"
)

// Cache is a simple structure to handle caching with a duration timeout
type Cache struct {
	Timeout  time.Duration
	deadline time.Time
	loaded   bool
	data     interface{}
	mutex    sync.Mutex
}

func (c *Cache) SetData(data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.loaded = true
	c.data = data
	c.deadline = time.Now().Add(c.Timeout)
}

func (c *Cache) Data() (data interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.loaded || time.Now().After(c.deadline) {
		return nil, false
	}
	return c.data, true
}
