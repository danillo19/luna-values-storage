package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"luna-values-storage/internal/common"
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

	err := req.ParseMultipartForm(MaxMultipartMemory)
	if err != nil {
		r.logger.Error(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonMetaData := req.FormValue("meta-data")

	value := new(Value)
	if err := json.NewDecoder(bytes.NewReader([]byte(jsonMetaData))).Decode(value); err != nil {
		r.logger.Error(err)

		RespondWithError(ErrParsingFailed, w, req)
		return
	}

	fileHeaders, ok := req.MultipartForm.File["value"]
	if ok && len(fileHeaders) == 1 {
		fileHeader := fileHeaders[0]
		uploadFile := common.UploadFile{}
		file, err := fileHeader.Open()
		if err != nil {
			r.logger.Error(err)
			RespondWithError(fmt.Errorf("failed to retrive a file from multipart"), w, req)
			return
		}

		uploadFile.FromMultipart(file, fileHeader)

		if err := r.s3Client.UploadToFiles(ctx, uploadFile.Path, &uploadFile); err != nil {
			r.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		value.Val = r.s3Client.GetPublicURL("files", uploadFile.Path)
	}

	multipartValue := req.FormValue("value")
	if multipartValue != "" {
		value.Val = multipartValue
	}

	if multipartValue == "" && fileHeaders == nil {
		RespondWithError(fmt.Errorf("empty value"), w, req)
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
