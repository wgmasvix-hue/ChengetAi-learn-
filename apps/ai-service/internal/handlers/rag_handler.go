package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/rag"
	"go.uber.org/zap"
)

// RAGHandler handles RAG (Retrieval Augmented Generation) endpoints
type RAGHandler struct {
	ragEngine *rag.Engine
	logger    *zap.SugaredLogger
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(ragEngine *rag.Engine, logger *zap.SugaredLogger) *RAGHandler {
	return &RAGHandler{
		ragEngine: ragEngine,
		logger:    logger,
	}
}

// RetrieveContext handles POST /rag/retrieve
func (h *RAGHandler) RetrieveContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req rag.RetrievalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Set defaults
	if req.TopK == 0 {
		req.TopK = 5
	}
	if req.MinSimilarity == 0 {
		req.MinSimilarity = 0.5
	}

	result, err := h.ragEngine.RetrieveByEmbedding(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("Error retrieving context", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// SemanticSearch handles POST /rag/search
func (h *RAGHandler) SemanticSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var searchReq struct {
		Query   string            `json:"query"`
		TopK    int               `json:"topK"`
		Filters map[string]string `json:"filters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if searchReq.TopK == 0 {
		searchReq.TopK = 5
	}

	result, err := h.ragEngine.RetrieveByKeyword(r.Context(), searchReq.Query, searchReq.TopK, searchReq.Filters)
	if err != nil {
		h.logger.Errorw("Error searching resources", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
