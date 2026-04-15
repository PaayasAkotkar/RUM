package rum

import (
	rumrpc "rum/app/misc/rum"
	"context"
	"log"
	"strings"
	"sync"
)

// Rum implements the core design of server

type Rum[In, Out any] struct {
	rumrpc.UnimplementedOnRumServiceServer

	light *Light[In, IDispatchResult]
	store *RumStore[In, Out]

	post              chan ILinks[In, Out]
	deleteService     chan ILinks[In, Out]
	activateService   chan ILinks[In, Out]
	deactivateService chan ILinks[In, Out]
	deleteProfile     chan ILinks[In, Out]
	activateProfile   chan ILinks[In, Out]
	deactivateProfile chan ILinks[In, Out]

	ctx context.Context
	mu  sync.Mutex
	wg  sync.WaitGroup
	sb  strings.Builder
}

func New[In, Out any](ctx context.Context, store *RumStore[In, Out]) *Rum[In, Out] {

	return &Rum[In, Out]{
		store:             store,
		light:             NewLight[In, IDispatchResult](),
		post:              make(chan ILinks[In, Out]),
		deleteService:     make(chan ILinks[In, Out]),
		activateService:   make(chan ILinks[In, Out]),
		deactivateService: make(chan ILinks[In, Out]),
		deleteProfile:     make(chan ILinks[In, Out]),
		activateProfile:   make(chan ILinks[In, Out]),
		deactivateProfile: make(chan ILinks[In, Out]),
		ctx:               ctx,
	}
}

// Fetch returns the registered dispatches results
// func (r *Rum[In, Out]) Fetch(profile ISequence[In], service, event string) *IDispatchResult {
// 	a, err := r.Store.profile.GetProfile(profile)
// 	if err != nil {
// 		return nil
// 	}
// 	s, err := a.GetService(service)
// 	if err != nil {
// 		return nil
// 	}
// 	b := s.dispatch.GetResults(event)

// 	return b
// }

// Paper monitors the profile and returns the result of the serivce of created event
func (r *Rum[In, Out]) Paper(profile ISequence[In]) *IDispatchResult {
	log.Println("in tick fetch")
	return r.tickFetch(profile)
}

func (r *Rum[In, Out]) tickFetch(profile ISequence[In]) *IDispatchResult {
	// Use only Name+Rank for pub/sub and profile lookup — Input is a pointer
	// whose address won't match between subscriber and publisher.
	key := ISequence[In]{Name: profile.Name, Rank: profile.Rank}
	ch := r.light.Subscribe(key)
	defer r.light.Unsub(key, ch)

	for {
		select {
		case <-r.ctx.Done():
			return nil

		case result := <-ch:
			if result != nil && result.IsReady {
				return result
			}
		}
	}
}

// func (r *Rum[In, Out]) scanReady(profile ISequence[In]) *IDispatchResult {
// 	kit, err := r.GetProfile(profile)
// 	if err != nil {
// 		return nil
// 	}
// 	for _, svc := range kit.Services() {
// 		for n := range svc.GetDispatch().GetRegistry() {
// 			z := svc.GetDispatch().GetResults(n)
// 			if z != nil && z.IsReady {
// 				return z
// 			}
// 		}
// 	}
// 	return nil
// }
