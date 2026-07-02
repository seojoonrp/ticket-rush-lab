package service

import (
	"context"
	"log"
	"seojoonrp/ticket-rush-lab/internal/repository"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	seatID primitive.ObjectID
	userID string
}

type WorkerPool struct {
	seatRepo    *repository.SeatRepo
	bookingRepo *repository.BookingRepo

	jobs chan Job
	wg   sync.WaitGroup
}

func NewWorkerPool(sr *repository.SeatRepo, br *repository.BookingRepo, workers, buffer int) *WorkerPool {
	wp := &WorkerPool{
		seatRepo:    sr,
		bookingRepo: br,

		jobs: make(chan Job, buffer),
	}

	for range workers {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for j := range wp.jobs {
		if err := wp.doWrite(j); err != nil {
			log.Printf("[WORKER POOL] error persist failed seat=%s user=%s: %v",
				j.seatID.Hex(), j.userID, err)
		}
	}
}

func (wp *WorkerPool) doWrite(j Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	seat, err := wp.seatRepo.FindByID(ctx, j.seatID)
	if err != nil {
		return err
	}

	if err := wp.seatRepo.UpdateOnBook(ctx, j.seatID, j.userID); err != nil {
		return err
	}

	_, err = wp.bookingRepo.Create(ctx, seat.ShowID, j.seatID, j.userID)
	if err != nil {
		return err
	}

	return nil
}

func (wp *WorkerPool) Submit(j Job) {
	wp.jobs <- j
}

// verify 전에 큐가 비었는지 확인 (drain용)
func (wp *WorkerPool) QueueDepth() int {
	return len(wp.jobs)
}

func (wp *WorkerPool) Shutdown() {
	close(wp.jobs)
	wp.wg.Wait()
}
