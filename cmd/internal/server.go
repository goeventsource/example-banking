package internal

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	banking "github.com/goeventsource/example-banking"

	"github.com/goeventsource/example-banking/internal"
)

type Server struct {
	svc    *internal.Service
	codec  *banking.RootEncodeDecoder
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(
	svc *internal.Service,
	codec *banking.RootEncodeDecoder,
) *Server {
	s := &Server{
		svc:    svc,
		codec:  codec,
		logger: slog.Default(),
		mux:    http.NewServeMux(),
	}
	s.mux.HandleFunc("/open", s.OpenAccount)
	s.mux.HandleFunc("/activate", s.Activate)
	s.mux.HandleFunc("/withdraw", s.Withdraw)
	s.mux.HandleFunc("/deposit", s.Deposit)
	s.mux.HandleFunc("/close", s.Close)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) OpenAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		Currency string
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currency := money.GetCurrency(reqBody.Currency)
	if currency == nil {
		s.logger.InfoContext(ctx, "invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.svc.OpenAccount(ctx, *currency)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not open account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := s.codec.Encode(acc)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not encode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(resBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not write response")
		return
	}
}

func (s *Server) Activate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		ID      uuid.UUID
		AgentID uuid.UUID
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqBody.ID == (uuid.UUID{}) || reqBody.AgentID == (uuid.UUID{}) {
		s.logger.InfoContext(ctx, "invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.svc.Activate(ctx, reqBody.ID, reqBody.AgentID)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not activate")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := s.codec.Encode(acc)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not encode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not write response")
		return
	}
}

func (s *Server) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		ID       uuid.UUID
		Amount   money.Amount
		Currency string
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currency := money.GetCurrency(reqBody.Currency)

	if reqBody.ID == (uuid.UUID{}) || reqBody.Amount == 0 || currency == nil {
		s.logger.InfoContext(ctx, "invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.svc.Withdraw(ctx, reqBody.ID, reqBody.Amount, *currency)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not withdraw")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := s.codec.Encode(acc)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not encode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not write response")
		return
	}
}

func (s *Server) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		ID       uuid.UUID
		Amount   money.Amount
		Currency string
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currency := money.GetCurrency(reqBody.Currency)

	if reqBody.ID == (uuid.UUID{}) || reqBody.Amount == 0 || currency == nil {
		s.logger.InfoContext(ctx, "invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.svc.Deposit(ctx, reqBody.ID, reqBody.Amount, *currency)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not deposit")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := s.codec.Encode(acc)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not encode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not write response")
		return
	}
}

func (s *Server) Close(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody struct {
		ID      uuid.UUID
		AgentID uuid.UUID
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqBody.ID == (uuid.UUID{}) || reqBody.AgentID == (uuid.UUID{}) {
		s.logger.InfoContext(ctx, "invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.svc.Close(ctx, reqBody.ID, reqBody.AgentID)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not close")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := s.codec.Encode(acc)
	if err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not encode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resBody); err != nil {
		s.logger.With(slog.Any("err", err)).ErrorContext(ctx, "could not write response")
		return
	}
}
