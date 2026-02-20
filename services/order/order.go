package order

import (
	"context"

	errProduct "github.com/arthurhzna/ecommerce_be_test/constants/error/product"
	statusConstants "github.com/arthurhzna/ecommerce_be_test/constants/status"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
)

type OrderService struct {
	repository repositories.IRepositoryRegistry
}

type IOrderService interface {
	CreateOrder(ctx context.Context, userID uint, req *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetMyOrders(ctx context.Context, userID uint) (*dto.OrderListResponse, error)
}

func NewOrderService(repository repositories.IRepositoryRegistry) IOrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	var totalAmount float64
	var orderItems []*models.OrderItem

	for _, item := range req.Items {
		product, err := s.repository.GetProduct().FindByUUID(ctx, item.ProductUUID)
		if err != nil {
			return nil, errProduct.ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, errProduct.ErrInsufficientStock
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItem := &models.OrderItem{
			OrderID:   0,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	order := &models.Order{
		UserID:    userID,
		PaymentID: nil,
		Amount:    totalAmount,
		Status:    statusConstants.OrderStatusPending,
		PaidAt:    nil,
	}

	createdOrder, err := s.repository.GetOrder().Create(ctx, order)
	if err != nil {
		return nil, err
	}

	for _, item := range orderItems {
		item.OrderID = createdOrder.ID
	}

	err = s.repository.GetOrderItem().CreateBulk(ctx, orderItems)
	if err != nil {
		return nil, err
	}

	for _, orderItem := range orderItems {
		err = s.repository.GetProduct().UpdateStock(ctx, orderItem.ProductID, orderItem.Quantity)
		if err != nil {
			return nil, err
		}
	}

	orderWithItems, err := s.repository.GetOrder().FindByID(ctx, createdOrder.ID)
	if err != nil {
		return nil, err
	}

	itemsResponse := make([]dto.OrderItemResponse, len(orderWithItems.OrderItems))
	for i, item := range orderWithItems.OrderItems {
		itemsResponse[i] = dto.OrderItemResponse{
			UUID:        item.UUID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Price * float64(item.Quantity),
		}
	}

	return &dto.OrderResponse{
		UUID:      orderWithItems.UUID,
		UserID:    orderWithItems.UserID,
		Amount:    orderWithItems.Amount,
		Status:    orderWithItems.Status.GetStatusString(),
		Items:     itemsResponse,
		CreatedAt: orderWithItems.CreatedAt,
		PaidAt:    orderWithItems.PaidAt,
	}, nil
}

func (s *OrderService) GetMyOrders(ctx context.Context, userID uint) (*dto.OrderListResponse, error) {
	orders, err := s.repository.GetOrder().FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	ordersResponse := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		itemsResponse := make([]dto.OrderItemResponse, len(order.OrderItems))
		for j, item := range order.OrderItems {
			itemsResponse[j] = dto.OrderItemResponse{
				UUID:        item.UUID,
				ProductID:   item.ProductID,
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				Price:       item.Price,
				Subtotal:    item.Price * float64(item.Quantity),
			}
		}

		ordersResponse[i] = dto.OrderResponse{
			UUID:      order.UUID,
			UserID:    order.UserID,
			Amount:    order.Amount,
			Status:    order.Status.GetStatusString(),
			Items:     itemsResponse,
			CreatedAt: order.CreatedAt,
			PaidAt:    order.PaidAt,
		}
	}

	return &dto.OrderListResponse{
		Orders: ordersResponse,
	}, nil
}
