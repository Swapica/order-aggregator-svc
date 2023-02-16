package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AddOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddOrder(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	order := request.DBModel()
	q := OrdersQ(r).FilterByOrderID(order.OrderID).FilterByChain(&order.SrcChain)
	log := Log(r).WithFields(logan.F{"order_id": order.OrderID, "src_chain": order.SrcChain, "dest_chain": order.DestChain})

	conflict, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Debug("order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	srcChain := ChainsQ(r).FilterByChainID(order.SrcChain).Get()
	if srcChain == nil {
		log.Debug("src_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	destChain := ChainsQ(r).FilterByChainID(order.DestChain).Get()
	if destChain == nil {
		log.Debug("dest_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	newOrder, err := q.Insert(order)
	if err != nil {
		log.WithError(err).Error("failed to add order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, responses.NewOrder(newOrder, srcChain.Key, destChain.Key))
}
