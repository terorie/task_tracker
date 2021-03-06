package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"src/task_tracker/storage"
	"time"
)

type CreateWorkerRequest struct {
}

type CreateWorkerResponse struct {
	Ok       bool      `json:"ok"`
	Message  string    `json:"message,omitempty"`
	WorkerId uuid.UUID `json:"id,omitempty"`
}

type GetWorkerResponse struct {
	Ok      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Worker  *storage.Worker `json:"worker,omitempty"`
}

func (api *WebAPI) WorkerCreate(r *Request) {

	workerReq := &CreateWorkerRequest{}
	if !r.GetJson(workerReq) {
		return
	}

	identity := getIdentity(r)

	if !canCreateWorker(r, workerReq, identity) {

		logrus.WithFields(logrus.Fields{
			"identity":            identity,
			"createWorkerRequest": workerReq,
		}).Warn("Failed CreateWorkerRequest")

		r.Json(CreateWorkerResponse{
			Ok:      false,
			Message: "You are now allowed to create a worker",
		}, 403)
		return
	}

	id, err := api.workerCreate(workerReq, getIdentity(r))
	if err != nil {
		handleErr(err, r)
	} else {
		r.OkJson(CreateWorkerResponse{
			Ok:       true,
			WorkerId: id,
		})
	}
}

func (api *WebAPI) WorkerGet(r *Request) {

	id, err := uuid.Parse(r.Ctx.UserValue("id").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid UUID")

		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		return
	}

	worker := api.Database.GetWorker(id)

	if worker != nil {
		r.OkJson(GetWorkerResponse{
			Ok:     true,
			Worker: worker,
		})
	} else {
		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: "Worker not found",
		}, 404)
	}
}

func (api *WebAPI) workerCreate(request *CreateWorkerRequest, identity *storage.Identity) (uuid.UUID, error) {

	worker := storage.Worker{
		Id:       uuid.New(),
		Created:  time.Now().Unix(),
		Identity: identity,
	}

	api.Database.SaveWorker(&worker)
	return worker.Id, nil
}

func canCreateWorker(r *Request, cwr *CreateWorkerRequest, identity *storage.Identity) bool {
	return true
}

func getIdentity(r *Request) *storage.Identity {

	identity := storage.Identity{
		RemoteAddr: r.Ctx.RemoteAddr().String(),
		UserAgent:  string(r.Ctx.UserAgent()),
	}

	return &identity
}
