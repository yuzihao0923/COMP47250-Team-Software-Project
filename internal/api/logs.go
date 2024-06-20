package api

import (
    "COMP47250-Team-Software-Project/internal/log"
    "COMP47250-Team-Software-Project/pkg/serializer"
    "net/http"
)

func HandleLogs(w http.ResponseWriter, r *http.Request) {
    logEntries := log.GetLogEntries()
    w.Header().Set("Content-Type", "application/json")
    err := serializer.JSONSerializerInstance.SerializeToWriter(logEntries, w)
    if err != nil {
        log.WriteErrorResponse(w, http.StatusInternalServerError, err)
    }
}
