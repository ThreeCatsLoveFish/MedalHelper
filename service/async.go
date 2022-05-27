package service

import (
	"MedalHelper/dto"
	"MedalHelper/util"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sethvargo/go-retry"
)

// SyncAction implement IConcurrency, support synchronous actions
type SyncAction struct{}

func (a *SyncAction) Exec(user User, job *sync.WaitGroup, child IExec) []dto.MedalInfo {
	fail := make([]dto.MedalInfo, 0, len(user.medalsLow))
	for _, medal := range user.remainMedals {
		retryTime := util.GlobalConfig.CD.Retry
		if retryTime == 0 {
			if ok := child.Do(user, medal); !ok {
				fail = append(fail, medal)
			}
		} else {
			backup := retry.NewFibonacci(time.Duration(retryTime) * time.Second)
			backup = retry.WithMaxRetries(uint64(util.GlobalConfig.CD.MaxTry), backup)
			ctx := context.Background()
			err := retry.Do(ctx, backup, func(ctx context.Context) error {
				if ok := child.Do(user, medal); !ok {
					return retry.RetryableError(errors.New("action fail"))
				}
				return nil
			})
			if err != nil {
				fail = append(fail, medal)
			}
		}
	}
	child.Finish(user, fail)
	job.Done()
	return fail
}

// AsyncAction implement IConcurrency, support asynchronous actions
type AsyncAction struct{}

func (a *AsyncAction) Exec(user User, job *sync.WaitGroup, child IExec) []dto.MedalInfo {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	fail := make([]dto.MedalInfo, 0, len(user.medalsLow))
	for _, medal := range user.remainMedals {
		wg.Add(1)
		retryTime := util.GlobalConfig.CD.Retry
		if retryTime == 0 {
			go func(medal dto.MedalInfo) {
				if ok := child.Do(user, medal); !ok {
					mu.Lock()
					fail = append(fail, medal)
					mu.Unlock()
				}
				wg.Done()
			}(medal)
		} else {
			backup := retry.NewFibonacci(time.Duration(retryTime) * time.Second)
			backup = retry.WithMaxRetries(uint64(util.GlobalConfig.CD.MaxTry), backup)
			go func(medal dto.MedalInfo) {
				ctx := context.Background()
				err := retry.Do(ctx, backup, func(ctx context.Context) error {
					if ok := child.Do(user, medal); !ok {
						return retry.RetryableError(errors.New("action fail"))
					}
					return nil
				})
				if err != nil {
					mu.Lock()
					fail = append(fail, medal)
					mu.Unlock()
				}
				wg.Done()
			}(medal)
		}
	}
	wg.Wait()
	child.Finish(user, fail)
	job.Done()
	return fail
}
