/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/job"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/gorilla/mux"
)

// queueTaskInput - Inputed task info
type queueTaskInput struct {
	TaskID       string `json:"task_id"`
	WorkDuration string `json:"work_duration"`
}

// handler - Worker handler
type handler struct {
	worker job.IWorker
}

// queueTask - Task queuing to worker
func (h *handler) queueTask(w http.ResponseWriter, req *http.Request) {
	var input queueTaskInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		logger.WithError(err).Info("failed to read POST body")
		renderResponse(w, http.StatusBadRequest, `{"error": "failed to read POST body"}`)
		return
	}
	defer req.Body.Close()

	// parse the work duration from the request body.
	workDuration, errParse := time.ParseDuration(input.WorkDuration)
	if errParse != nil {
		logger.WithError(errParse).Info("faile to parse work duration in request")
		renderResponse(w, http.StatusBadRequest, `{"error": "failed to parse work duration in request"}`)
		return
	}

	// queue the task in background task manager
	if err := h.worker.QueueTask(input.TaskID, workDuration, job.TaskInfo{}); err != nil {
		logger.WithError(err).Info("failed to queue task")
		if err == job.ErrWorkerBusy {
			w.Header().Set("Retry-After", "60")
			renderResponse(w, http.StatusServiceUnavailable, `{"error": "workers are busy, try again later"}`)
			return
		}
		renderResponse(w, http.StatusInternalServerError, `{"error": "failed to queue task"}`)
		return
	}

	renderResponse(w, http.StatusAccepted, `{"status": "task queued successfully"}`)
}

// renderResponse - Redering response
func renderResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// main - entry point
func main() {
	logger.New()

	ctx := context.Background()
	gracePeriod := 5 * time.Second
	workerCount := 10
	buffer := 100
	httpAddr := ":4444"

	logger.Info("start workers")
	w := job.NewWorker(workerCount, buffer)
	w.Start(ctx)

	h := handler{worker: w}

	router := mux.NewRouter()
	router.HandleFunc("/queue-task", h.queueTask).Methods("POST")

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("starting http server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatalf("listen failed")
		}
	}()

	<-done
	logger.Info("http server stopped")

	ctxTimeout, cancel := context.WithTimeout(ctx, gracePeriod)
	defer func() {
		w.Stop()
		cancel()
	}()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		logger.WithError(err).Fatalf("http server shutdown failed")
	}
}
