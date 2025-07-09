package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/OmerYesilkaya/fileuploader/internal/api"
	"github.com/OmerYesilkaya/fileuploader/internal/utils"
)

type FileHandler struct {
	Ctx *api.AppContext
}

func handleSigleFileUpload(w http.ResponseWriter, fileHeader *multipart.FileHeader) {
	file, err := fileHeader.Open()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to open file: "+err.Error())
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}

	// write to disk
	filePath := fmt.Sprintf("%s/%s", "./uploads", fileHeader.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create file on disk: "+err.Error())
		return
	}
	_, err = io.Copy(outFile, buf)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to write file to disk: "+err.Error())
		return
	}
	defer outFile.Close()
}

func (h *FileHandler) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB limit
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		return
	}
	files := r.MultipartForm.File["file"]
	response := make([]map[string]string, 0, len(files))
	for _, fileHeader := range files {
		handleSigleFileUpload(w, fileHeader)
		response = append(response, map[string]string{
			"filename": fileHeader.Filename,
			"size":     fmt.Sprintf("%d", fileHeader.Size),
			"status":   "uploaded",
		})
	}

	// id: string;
	// filename: string;
	// size: number;
	// status: string;
	// url: string;

	err = utils.ResponseSuccess(w, http.StatusCreated, "files_uploaded_successfully", response)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to encode response: "+err.Error())
		return
	}
}
