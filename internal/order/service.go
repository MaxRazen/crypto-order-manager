package order

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/MaxRazen/crypto-order-manager/internal/deadline"
	"github.com/MaxRazen/crypto-order-manager/internal/logger"
	"github.com/MaxRazen/crypto-order-manager/internal/ordergrpc"
	"github.com/MaxRazen/crypto-order-manager/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Validate(req *ordergrpc.CreateOrderRequest) (*NewOrder, error) {
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
	if req.GetQuantity().GetType() == ordergrpc.QuantityType_UNKNOWN_QUANTITY_TYPE {
		return nil, status.Error(codes.InvalidArgument, "quantity.type must be specified")
	}
	if match, err := regexp.MatchString("^(\\d+)(\\.\\d+)?$", price); err != nil || !match {
		return nil, status.Error(codes.InvalidArgument, "quantity.value does not match the expected pattern")
	}

	// validate Deadline properties
	reqDeadlines := req.GetDeadlines()
	deadlines := make([]deadline.Deadline, 0, len(reqDeadlines))
	for idx, d := range reqDeadlines {
		// validate Deadline.type property
		dType := d.Type.String()
		if d.Type == ordergrpc.DeadlineType_UNKNOWN_DEADLINE_TYPE {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].type must be specified", idx))
		}

		// validate Deadline.action property
		dAction := d.Action.String()
		if d.Action == ordergrpc.DeadlineAction_UNKNOWN_DEADLINE_ACTION {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].action must be specified", idx))
		}

		// validate Deadline.value property
		dValue, err := deadline.ParseDeadlineValue(dAction, d.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("deadlines[%v].value has invalid value", idx))
		}

		deadlines = append(deadlines, deadline.Deadline{
			Type:   dType,
			Action: dAction,
			Value:  dValue,
		})
	}

	data := NewOrder{
		Pair:     pair,
		Market:   market,
		Action:   action,
		Behavior: behavior,
		Price:    price,
		Quantity: Quantity{
			Type:  quantityType,
			Value: quantityValue,
		},
		Deadlines: deadlines,
	}

	return &data, nil
}

func Create(ctx context.Context, log *logger.Logger, ds *datastore.Client, data *NewOrder) (*Order, error) {
	dsk := storage.NewIDKey(OrdersKind, storage.GenerateID())

	o := Order{
		Status:    StatusNew,
		RuleId:    data.RuleId,
		Pair:      data.Pair,
		Market:    data.Market,
		Action:    data.Action,
		Behavior:  data.Behavior,
		Price:     data.Price,
		Quantity:  data.Quantity,
		Deadlines: data.Deadlines,
		CreatedAt: time.Now(),
	}

	dsk, err := ds.Put(ctx, dsk, o)
	if err != nil {
		return nil, err
	}
	o.DSKey = dsk

	return &o, nil
}
