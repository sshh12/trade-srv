package indexers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sshh12/go-tdameritrade/tdameritrade"
	events "github.com/sshh12/trade-srv/events"
	"golang.org/x/oauth2"
)

const tdaSource string = "tda"

type tdHTTPHeaderStore struct{}

func (s *tdHTTPHeaderStore) StoreToken(token *oauth2.Token, w http.ResponseWriter, req *http.Request) error {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:    "refreshToken",
			Value:   token.RefreshToken,
			Expires: token.Expiry,
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:    "accessToken",
			Value:   token.AccessToken,
			Expires: token.Expiry,
		},
	)
	return nil
}

func (s tdHTTPHeaderStore) GetToken(req *http.Request) (*oauth2.Token, error) {
	refreshToken, err := req.Cookie("refreshToken")
	if err != nil {
		return nil, err
	}
	accessToken, err := req.Cookie("accessToken")
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
		Expiry:       refreshToken.Expires,
	}, nil
}

func (s tdHTTPHeaderStore) StoreState(state string, w http.ResponseWriter, req *http.Request) error {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "state",
			Value: state,
		},
	)
	return nil
}

func (s tdHTTPHeaderStore) GetState(req *http.Request) (string, error) {
	cookie, err := req.Cookie("state")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

type tdHandlers struct {
	authenticator *tdameritrade.Authenticator
	channel       chan *tdameritrade.Client
}

func (h *tdHandlers) Authenticate(w http.ResponseWriter, req *http.Request) {
	redirectURL, err := h.authenticator.StartOAuth2Flow(w, req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)
}

func (h *tdHandlers) Callback(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := h.authenticator.FinishOAuth2Flow(ctx, w, req)
	h.channel <- client
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/done", http.StatusTemporaryRedirect)
}

func (h *tdHandlers) Done(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Done"))
}

func startTDAIndexer(es *events.EventStream, opts *IndexerOptions) error {
	if opts.TDAConsumerKey == "" {
		log.Error("No TDA consumer key provided")
		return nil
	}
	ctx := context.Background()
	tdaAuth := tdameritrade.NewAuthenticator(
		&tdHTTPHeaderStore{},
		oauth2.Config{
			ClientID: opts.TDAConsumerKey,
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
				AuthURL:  "https://auth.tdameritrade.com/auth",
			},
			RedirectURL: "http://localhost:8283/callback",
		},
	)
	clientChan := make(chan *tdameritrade.Client)
	handlers := &tdHandlers{authenticator: tdaAuth, channel: clientChan}
	http.HandleFunc("/", handlers.Authenticate)
	http.HandleFunc("/done", handlers.Done)
	http.HandleFunc("/callback", handlers.Callback)
	go http.ListenAndServe(":8283", nil)
	client := <-clientChan
	extendedHours := true
	phOpts := &tdameritrade.PriceHistoryOptions{
		FrequencyType:         "minute",
		NeedExtendedHoursData: &extendedHours,
	}
	for {
		symbols, err := es.GetSymbols()
		if err != nil {
			log.Error(err)
			return err
		}
		rand.Shuffle(len(symbols), func(i, j int) {
			symbols[i], symbols[j] = symbols[j], symbols[i]
		})
		for _, sym := range symbols {
			ph, _, err := client.PriceHistory.PriceHistory(ctx, sym.Sym, phOpts)
			if err != nil {
				log.Error(err)
				continue
			}
			ticks := make([]events.TDAOHLCV, 0)
			for _, candle := range ph.Candles {
				ticks = append(ticks, events.TDAOHLCV{
					ID:     fmt.Sprintf("%s%d", sym.Sym, candle.Datetime),
					Symbol: sym.Sym,
					Date:   candle.Datetime,
					Open:   candle.Open,
					Close:  candle.Close,
					High:   candle.High,
					Low:    candle.Low,
					Volume: candle.Volume,
				})
			}
			if err := es.OnMinOHLCVs(ticks); err != nil {
				log.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}
