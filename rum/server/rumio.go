package rum

import (
	"log"
)

// read returns the current active service to publish
// it walks the serviceStack and returns the latest active one
func (k *Kit[In, Out]) read() *Service[In, Out] {
	k.mu.RLock()
	defer k.mu.RUnlock()
	latest := k.serviceStack.Latest()
	if latest == nil {
		return nil
	}
	svc, ok := k.activeService[*latest]
	if !ok {
		return nil
	}
	return svc
}

type IWrite[In, Out any] struct {
	Service *Service[In, Out]
	Profile ISequence[In]
	Report  *ProfileMetric
}

// write performs the dispatching of the profile events as per the desc and writes the metrics
func (r *Rum[In, Out]) write(profile ISequence[In], desc *Service[In, Out]) *IWrite[In, Out] {
	ctx := r.ctx
	log.Println("writing..")

	kits := r.store.profile.Sort(profile.Name)

	for _, ser := range kits {
		if b := ser.GetBudget(); b != nil && b.Exhausted() {
			log.Println("budget error")
			r.onDeactivateService(profile.Name, ser.Name)
			break
		}

		if b := ser.GetBudget(); b != nil {
			b.Spend()
		}

		// Dispatch funcs
		for n := range ser.GetDispatch().GetRegistry() {
			var policy *RetryPolicy
			if f := ser.GetFormat(); f != nil {
				policy = f.Retry
			}
			mx := profile.Input
			if mx == nil {
				continue
			}

			if err := ser.GetDispatch().call(ctx, n, *mx, policy); err == nil {
				rp := ser.GetDispatch().GetMetric(n)
				if kit, err := r.store.GetKit(profile.Name); err == nil {
					if rp.Fail != nil {
						kit.AddFailReport(*rp.Fail)
					} else if rp.Succeed != nil {
						kit.AddSucceedReport(*rp.Succeed)
					}
				}
			}
		}

		r.handleServiceFormat(profile, ser)
	}

	if kit, err := r.store.GetKit(profile.Name); err == nil {
		r.handleProfileFormat(profile, kit)
	}

	finalMetric := NewProfileMetric()
	if kit, err := r.store.GetKit(profile.Name); err == nil {
		finalMetric.Metric[profile.Name] = kit.GetMetrics()
	}

	finalResult := &IDispatchResult{
		IsReady: true,
		Metric:  finalMetric,
	}

	for _, ser := range kits {
		for n := range ser.GetDispatch().GetRegistry() {
			z := ser.GetDispatch().GetResults(n)
			if z != nil && z.IsReady {
				finalResult.Result = z.Result
			}
		}
	}

	// Normalize key to Name and Rank only for pub/sub matching.
	key := ISequence[In]{Name: profile.Name, Rank: profile.Rank}
	r.light.Publish(key, finalResult)

	return &IWrite[In, Out]{
		Report:  &finalMetric,
		Service: desc,
		Profile: profile,
	}
}
