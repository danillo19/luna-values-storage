package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"luna-values-storage/internal/common"
	"net/http"
)

const (
	MaxMultipartMemory = 500 * humanize.MiByte
)

func (r Resolver) UploadFile(w http.ResponseWriter, req *http.Request) {
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
		}

		uploadFile.FromMultipart(file, fileHeader)

		if err := r.s3Client.UploadToFiles(ctx, uploadFile.Path, &uploadFile); err != nil {
			r.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		value.Val = r.s3Client.GetPublicURL("files", uploadFile.Path)
	}

	multipartValue, ok := req.MultipartForm.Value["value"]
	if ok {
		value.Val = multipartValue
	}

	value.Val = multipartValue

}
