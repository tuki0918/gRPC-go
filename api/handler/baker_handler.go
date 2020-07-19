package handler

import (
	"math/rand"
	"sync"
	"time"
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pancake/maker/gen/api"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BakerHandler struct {
	report *report
}

type report struct {
	sync.Mutex
	data map[api.Pancake_Menu]int
}

func NewBakerHandler() * BakerHandler {
	return &BakerHandler {
		report: &report {
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

func (h *BakerHandler) Bake(
	ctx context.Context,
	req *api.BakerRequest,
	) (*api.BakerResponse, error) {

		if req.Menu == api.Pancake_UNKNOWN
			|| req.Menu > api.Pancake_SPICY_CURRY {
				return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んで下さい！")
			}

	}

	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	return &api.BakerResponse {
		Pancake: &api.Pancake {
			Menu: req.Menu,
			ChefName: "gami"
			TechnicalScore: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos: int32(now.Nanoseconds()),
			},
		},
	}, nil
)

func (h *BakerHandler) Report(
	ctx context.Context,
	req *api.ReportRequest,
	) (*api.ReportResponse, error) {

		counts := make([]*api.Report_BakeCount, 0)

		h.report.Lock()
		for k, v := range h.report.data {
			counts = append(counts, &api.Report_BakeCount{
				Menu: k,
				Count: int32(v),
			})
		}
		h.report.Unlock()

		return &api.ReportResponse{
			Report: &api.Report{
				BakeCounts: counts,
			},
		}, nil

	}
)
