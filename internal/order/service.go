package order

import (
	"context"
	"fmt"
	"regexp"

	"github.com/MaxRazen/crypto-order-manager/internal/app"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Validate(req *ordergrpc.CreateOrderRequest) (*CreationData, error) {
	// validate Pair property
	pair := req.GetPair()
	if match, err := regexp.MatchString("^([A-Z]+)-([A-Z]+)$", pair); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "pair does not match the expected pattern")
	}

	// validate Market property
	market := req.GetMarket()
	if market == "" {
		return nil, status.Error(codes.InvalidArgument, "market must be specified")
	}

	// validate Action propery
	action := req.GetAction().String()
	if req.GetAction() == ordergrpc.ActionType_UNKNOWN_ACTION {
		return nil, status.Error(codes.InvalidArgument, "action must be specified")
	}

	// validate Behavior propery
	behavior := req.GetBehavior().String()
	if req.GetBehavior() == ordergrpc.Behavior_UNKNOWN_BEHAVIOR {
		return nil, status.Error(codes.InvalidArgument, "behavior must be specified")
	}

	// validate Price property
	price := req.GetPrice()
	if match, err := regexp.MatchString("^(\\d+)(\\.\\d+)?$", price); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "price does not match the expected pattern")
	}

	// validate Quantity properties
	quantity := req.GetQuantity()
	quantityType := quantity.GetType().String()
	quantityValue := quantity.GetValue()
	if match, err := regexp.MatchString("^(\\d+)(\\.\\d+)?$", price); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "quantity.value does not match the expected pattern")
	}

	// validate Deadline properties
	reqDeadlines := req.GetDeadlines()
	deadlines := make([]Deadline, 0, len(reqDeadlines))
	for idx, d := range reqDeadlines {
		// validate Deadline.type property
		dType := d.Type.String()
		if d.Type == ordergrpc.DeadlineType_UNKNOWN_DEADLINE_TYPE {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].type must be specified", idx))
		}

		// validate Deadline.value property
		dValue := d.Value
		if dValue == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].value must be specified", idx))
		}

		// validate Deadline.action property
		dAction := d.Action.String()
		if d.Action == ordergrpc.DeadlineAction_UNKNOWN_DEADLINE_ACTION {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].action must be specified", idx))
		}

		deadlines = append(deadlines, Deadline{
			Type:   dType,
			Value:  dValue,
			Action: dAction,
		})
	}

	data := CreationData{
		Pair:     pair,
		Market:   market,
		Action:   action,
		Behavior: behavior,
		Quantity: Quantity{
			Type:  quantityType,
			Value: quantityValue,
		},
		Deadlines: deadlines,
	}

	return &data, nil
}

func SaveOrder(ctx context.Context, log *logger.Logger, app *app.App, data *CreationData) error {
	dk := storage.NewIDKey(IncomingOrdersKind, storage.GenerateID())

	if _, err := app.Storage.Put(ctx, dk, data); err != nil {
		return err
	}

	return nil
}

func NewTracker() {

}
