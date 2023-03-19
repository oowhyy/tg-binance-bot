package calc

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/oowhyy/tg-binance-bot/client/binance"
)

type Calc struct {
	client     *binance.Client
	cache      *binance.ExchangeInfo
	lastUpdate time.Time
	staleTime  time.Duration
}

// type coinSymbol string

type Triangle struct {
	Coins []string
	Coef  float64
}
type Graph map[string]*Vertex

type Vertex struct {
	name      string
	connected map[string]float64
}

const DefaultTrianglesNum = 7

func New(cl *binance.Client) *Calc {
	return &Calc{
		client:    cl,
		staleTime: 10 * time.Second,
	}
}

func (calc *Calc) BestTriangles(limit int) ([]Triangle, error) {
	if calc.cache == nil || calc.lastUpdate.Add(calc.staleTime).Before(time.Now()) {
		log.Print("NEW CAHCE")
		res, err := calc.client.ExchangeInfo()
		if err != nil {
			return nil, err
		}
		calc.cache = res
		calc.lastUpdate = time.Now()
	}
	bookTicker, err := calc.client.Ticker()
	if err != nil {
		return nil, err
	}
	gr, err := makeGraph(calc.cache, bookTicker)
	if err != nil {
		return nil, fmt.Errorf("wrong graph data format: %w", err)
	}
	res := max3Cycles(*gr)
	// sort descending
	sort.Slice(res, func(i, j int) bool { return res[i].Coef > res[j].Coef })
	return res[:limit], nil
}

func makeGraph(info *binance.ExchangeInfo, data []binance.BookTicker) (*Graph, error) {
	dict := map[string][2]string{}
	res := Graph{}
	edgeCount := 0
	for _, s := range info.Symbols {
		if s.Status == binance.SymbolStatusTypeTrading {
			dict[s.Symbol] = [2]string{s.BaseAsset, s.QuoteAsset}
			if _, ok := res[s.BaseAsset]; !ok {
				ver := &Vertex{s.BaseAsset, map[string]float64{}}
				res[s.BaseAsset] = ver
			}
			if _, ok := res[s.QuoteAsset]; !ok {
				ver := &Vertex{s.QuoteAsset, map[string]float64{}}
				res[s.QuoteAsset] = ver
			}
		}
	}
	var base, quote string
	for _, ticker := range data {
		if _, ok := dict[ticker.Symbol]; !ok {
			continue
		}
		base = dict[ticker.Symbol][0]
		quote = dict[ticker.Symbol][1]
		bid, err := strconv.ParseFloat(ticker.BidPrice, 64)
		if err != nil {
			return nil, err
		}
		ask, err := strconv.ParseFloat(ticker.AskPrice, 64)
		if err != nil {
			return nil, err
		}

		res[base].connected[quote] = bid
		res[quote].connected[base] = 1 / ask
		edgeCount += 2
	}
	log.Printf("ok graph: %d vertices, %d edges", len(res), edgeCount)
	return &res, nil
}

func max3Cycles(graph Graph) []Triangle {
	res := []Triangle{}
	//iterate first vertex
	for first, p1 := range graph {
		// iterate second vertex
		for second, wei1 := range p1.connected {
			// iterate third vertex
			for third, wei2 := range graph[second].connected {
				// is a triangle
				if wei3, ok := graph[third].connected[first]; ok {
					tr := Triangle{[]string{first, second, third}, wei1*wei2*wei3 - 1}
					res = append(res, tr)
				}
			}
		}
	}
	return res
}
