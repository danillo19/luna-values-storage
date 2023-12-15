package http

import (
	"encoding/json"
	"net/http"
)

func (r Resolver) GetValue(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	value, err := r.valueService.GetValue(ctx, req.URL.Query().Get("id"))
	if err != nil {
		r.logger.Error(err)

		RespondWithError(err, w, req)
		return
	}

	RespondOK(valueDomainToHttp(value), w, req)
}

func (r Resolver) SetValue(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	value := new(Value)
	if err := json.NewDecoder(req.Body).Decode(value); err != nil {
		r.logger.Error(err)

		RespondWithError(ErrParsingFailed, w, req)
		return
	}

	newValue, err := r.valueService.SetValue(ctx, value.IntoDomainType())
	if err != nil {
		r.logger.Error(err)

		RespondWithError(err, w, req)
		return
	}

	RespondOK(valueDomainToHttp(newValue), w, req)
}

func (r Resolver) DeleteValue(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")
	if err := r.valueService.DeleteValue(ctx, id); err != nil {
		r.logger.Error(err)

		RespondWithError(err, w, req)
		return
	}

	RespondOK(id, w, req)
}
