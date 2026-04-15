// Package rum ....
// Flow ->
//
//				           {name : ["event1", "event2", "event3"]}
//			                                        V
//		                    Call the events as per the order & desc of the events.
//		                         Perform metric writing for each input.
//	                                Store completion event name & result.
package rum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Dispatcher controls registered agent functions and their results
type Dispatcher[in, out any] struct {
	registry   map[string]IRegister[in, out]
	rinput     map[string]in
	events     Stack[string]
	result     map[string]*IDispatchResult
	isComplete map[string]bool
	metric     map[string]map[int]IAgentResp // name -> count -> resp
}

func NewDispatcher[in, out any]() *Dispatcher[in, out] {
	return &Dispatcher[in, out]{
		registry:   make(map[string]IRegister[in, out]),
		rinput:     make(map[string]in),
		result:     make(map[string]*IDispatchResult),
		isComplete: make(map[string]bool),
		metric:     make(map[string]map[int]IAgentResp),
	}
}

// IAgentResp holds per-call metric data
type IAgentResp struct {
	Succeed *IMetricAgentSucceed `json:"succeed"`
	Fail    *IMetricAgentFail    `json:"fail"`
}

// get funcs

func (d *Dispatcher[in, out]) GetRegistry() map[string]IRegister[in, out] {
	return d.registry
}

func (d *Dispatcher[in, out]) GetEvents(limit int) []string {
	return d.events.Range(limit)
}

func (d *Dispatcher[in, out]) GetLatestRegistry() *string {
	return d.events.Latest()
}

func (d *Dispatcher[in, out]) GetResults(name string) *IDispatchResult {
	if _, ok := d.result[name]; !ok {
		for n := range d.registry {
			log.Println("names: ", n)
		}
		log.Println("dispatcher: not found name ", name)
		return nil
	}
	return d.result[name]
}

// GetMetric returns the latest metric entry for a named dispatch
func (d *Dispatcher[in, out]) GetMetric(name string) IAgentResp {
	return d.metric[name][d.metricCount(name)]
}

func (d *Dispatcher[in, out]) metricCount(name string) int {
	return len(d.metric[name])
}

// end

// set funcs

// Release deletes all the metrices
func (d *Dispatcher[in, out]) Release() {
	for r := range d.metric {
		delete(d.metric, r)
	}
}

func (d *Dispatcher[in, out]) Register(event string, fn IRegister[in, out]) {
	if _, ok := d.registry[event]; !ok {
		d.registry[event] = fn
		d.events.Push(event)
	}
}

func (d *Dispatcher[in, out]) Unregister(name string) {
	delete(d.registry, name)
	delete(d.rinput, name)
	delete(d.isComplete, name)
	delete(d.metric, name)
	delete(d.result, name)
	d.events.Erase(name)
}

// call invokes the named registered function and records its metric
func (d *Dispatcher[in, out]) call(ctx context.Context, name string, input in, policy *RetryPolicy) error {
	fn, ok := d.registry[name]
	if !ok {
		err := fmt.Errorf("service %s not found", name)
		d.writeMetric(name, IAgentResp{Fail: &IMetricAgentFail{At: time.Now(), Reason: err.Error()}})
		return err
	}

	max := 1
	interval := time.Duration(0)

	if policy != nil {
		max = policy.Max + 1 // +1 so Max=3 means 1 original + 3 retries
		interval = policy.Interval
	}
	log.Println("max: ", max)

	var lastErr error
	for attempt := range max {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
		}

		start := time.Now()
		res, err := fn.Fn(ctx, input)
		elapsed := time.Since(start)

		if err == nil {
			// success — write metric and result, return
			b, _ := json.Marshal(res)
			c, _ := json.Marshal(input)
			d.writeMetric(name, IAgentResp{Succeed: &IMetricAgentSucceed{
				TimeTaken:     elapsed,
				AgentReply:    string(b),
				ClientRequest: string(c),
				At:            time.Now(),
			}})
			r := NewDispatchResult()
			r.IsReady = true
			r.Result = b

			d.handleOutput(name, r)
			d.handleComplete(name, true)
			d.handleInput(name, input)
			return nil
		}

		// record each failure attempt
		d.writeMetric(name, IAgentResp{Fail: &IMetricAgentFail{
			At:     time.Now(),
			Reason: fmt.Sprintf("attempt %d: %s", attempt+1, err.Error()),
		}})
		lastErr = err
	}

	return lastErr
}

func (d *Dispatcher[in, out]) writeMetric(name string, resp IAgentResp) {
	if _, ok := d.metric[name]; !ok {
		d.metric[name] = make(map[int]IAgentResp)
	}
	d.metric[name][d.metricCount(name)+1] = resp
}

func (d *Dispatcher[in, out]) handleOutput(name string, res *IDispatchResult) {
	d.result[name] = res
}
func (d *Dispatcher[in, out]) handleInput(name string, input in) {
	d.rinput[name] = input
}
func (d *Dispatcher[in, out]) handleComplete(name string, complete bool) {
	d.isComplete[name] = complete
}

// end
