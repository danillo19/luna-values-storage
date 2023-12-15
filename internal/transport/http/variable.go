package http

import (
	"encoding/json"
	"net/http"
)

func (r Resolver) GetVariable(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	variable, err := r.variableService.GetVariable(ctx, req.URL.Query().Get("id"))
	if err != nil {
		r.logger.Error(err)

		RespondWithError(err, w, req)
		return
	}

	RespondOK(variableDomainToHttp(variable), w, req)
}

func (r Resolver) SetVariable(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	variable := new(Variable)
	if err := json.NewDecoder(req.Body).Decode(variable); err != nil {
		r.logger.Error(err)
		RespondWithError(ErrParsingFailed, w, req)
		return
	}

	newVariable, err := r.variableService.SetVariable(ctx, variable.IntoDomainType())
	if err != nil {
		r.logger.Error(err)
		RespondWithError(err, w, req)
		return
	}

	RespondOK(variableDomainToHttp(newVariable), w, req)
}

func (r Resolver) DeleteVariable(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")
	if err := r.variableService.DeleteVariable(ctx, id); err != nil {
		r.logger.Error(err)
		RespondWithError(err, w, req)
		return
	}

	RespondOK(id, w, req)
}
