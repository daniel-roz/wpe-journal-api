package handlers

import (
	"api/models"
	"api/repositories"
	"encoding/json"
	"errors"
	"net/http"
	"uuid"
)

type JournalHandler struct {
	Repository *repositories.JournalRepository
}
type JournalPreviewHandler struct {
	JournalPreviews []models.JournalPreview
}

func NewJournalHandler(repository *repositories.JournalRepository) *JournalHandler {
	return &JournalHandler{
		Repository: repository,
	}
}

func (h *JournalHandler) GetJournalsByCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category(r.PathValue("category"))
	if !category.IsValid() {
		http.Error(w, "Invalid category.", http.StatusBadRequest)
		return
	}

	match, err := h.Repository.FindJournalsByCategory(r.Context(), category)
	if err != nil {
		http.Error(w, "Failed to retrieve journals.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

func (h *JournalHandler) GetJournalByID(w http.ResponseWriter, r *http.Request) {
	category := models.Category(r.PathValue("category"))
	if !category.IsValid() {
		http.Error(w, "Invalid category.", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	journalID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid journal ID.", http.StatusBadRequest)
		return
	}

	journal, err := h.Repository.FindJournalByID(r.Context(), category, journalID)
	if errors.Is(err, repositories.ErrJournalNotFound) {
		http.Error(w, "Journal not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to find journal.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(journal)
}

func (h *JournalHandler) CreateJournal(w http.ResponseWriter, r *http.Request) {
	var journal models.Journal
	if err := json.NewDecoder(r.Body).Decode(&journal); err != nil {
		http.Error(w, "Invalid JSON.", http.StatusBadRequest)
		return
	}

	if !journal.Category.IsValid() {
		http.Error(w, "Invalid category.", http.StatusBadRequest)
		return
	}

	created, err := h.Repository.CreateJournal(r.Context(), journal)
	if err != nil {
		http.Error(w, "Failed to create journal.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *JournalHandler) UpdateJournal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	journalID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid journal ID.", http.StatusBadRequest)
		return
	}

	var update models.JournalUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid JSON.", http.StatusBadRequest)
		return
	}
	if update.Category != nil && !update.Category.IsValid() {
		http.Error(w, "Invalid category.", http.StatusBadRequest)
		return
	}

	updated, err := h.Repository.UpdateJournal(r.Context(), journalID, update)
	if errors.Is(err, repositories.ErrJournalNotFound) {
		http.Error(w, "Journal not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update journal.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *JournalHandler) DeleteJournal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	journalID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid journal ID.", http.StatusBadRequest)
		return
	}

	err = h.Repository.DeleteJournal(r.Context(), journalID)
	if errors.Is(err, repositories.ErrJournalNotFound) {
		http.Error(w, "Journal not found.", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to delete journal.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
