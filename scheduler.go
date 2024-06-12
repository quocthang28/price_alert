package main

import (
	"fmt"
	"log"
	"time"
	_ "time/tzdata"

	"github.com/go-co-op/gocron"
)

const (
	PriceAlertJob = "PriceAlertJob"
)

type AppScheduler struct {
	s *gocron.Scheduler
	//jobs []*gocron.Job
}

func newAppScheduler() *AppScheduler {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal(err)
	}

	s := gocron.NewScheduler(loc)
	s.StartAsync()

	appScheduler := &AppScheduler{
		s: s,
	}

	return appScheduler
}

func (as *AppScheduler) schedulePriceAlertJob(priceAlertFunc func(symbols []string), symbols []string, interval string, reschedule bool) {
	if reschedule {
		err := as.s.RemoveByTag(PriceAlertJob)
		if err != nil {
			log.Println(err)
		}
	}

	now := time.Now()
	nextSchedule := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	//_, err := as.s.Every(10).Seconds().StartAt(nextSchedule).Do(priceAlertFunc, symbols)

	_, err := as.s.Tag(PriceAlertJob).Every(interval).StartAt(nextSchedule).Do(priceAlertFunc, symbols)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
	}

	//as.jobs = append(as.jobs, job)
}
