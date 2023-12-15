package http

import (
	"luna-values-storage/internal/common"
	"net/http"
)

func (r Resolver) UploadFile(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	err := req.ParseMultipartForm(10 << 2) // Max file size is set to 10MB
	if err != nil {
		r.logger.Error(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadFile := common.UploadFile{}
	uploadFile.FromMultipart(file, header)

	if err := r.s3Client.UploadToFiles(ctx, uploadFile.Path, &uploadFile); err != nil {
		r.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondOK(r.s3Client.GetPublicURL("files", uploadFile.Path), w, req)
}
