package helpers

import (
	"context"

	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	orderRepo "github.com/arthurhzna/ecommerce_be_test/repositories/order"
	orderItemRepo "github.com/arthurhzna/ecommerce_be_test/repositories/order_item"
	paymentRepo "github.com/arthurhzna/ecommerce_be_test/repositories/payment"
	productRepo "github.com/arthurhzna/ecommerce_be_test/repositories/product"
	userRepo "github.com/arthurhzna/ecommerce_be_test/repositories/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) FindByUUID(ctx context.Context, productUUID uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, productUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) FindByName(ctx context.Context, name string) (*models.Product, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id uint) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUUID(ctx context.Context, orderUUID uuid.UUID) (*models.Order, error) {
	args := m.Called(ctx, orderUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

type MockOrderItemRepository struct {
	mock.Mock
}

func (m *MockOrderItemRepository) CreateBulk(ctx context.Context, items []*models.OrderItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockOrderItemRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.OrderItem), args.Error(1)
}

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	args := m.Called(ctx, payment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

type MockRepositoryRegistry struct {
	mock.Mock
	UserRepo      *MockUserRepository
	ProductRepo   *MockProductRepository
	OrderRepo     *MockOrderRepository
	OrderItemRepo *MockOrderItemRepository
	PaymentRepo   *MockPaymentRepository
}

func NewMockRepositoryRegistry() *MockRepositoryRegistry {
	return &MockRepositoryRegistry{
		UserRepo:      new(MockUserRepository),
		ProductRepo:   new(MockProductRepository),
		OrderRepo:     new(MockOrderRepository),
		OrderItemRepo: new(MockOrderItemRepository),
		PaymentRepo:   new(MockPaymentRepository),
	}
}

func (m *MockRepositoryRegistry) GetUser() userRepo.IUserRepository {
	return m.UserRepo
}

func (m *MockRepositoryRegistry) GetProduct() productRepo.IProductRepository {
	return m.ProductRepo
}

func (m *MockRepositoryRegistry) GetOrder() orderRepo.IOrderRepository {
	return m.OrderRepo
}

func (m *MockRepositoryRegistry) GetOrderItem() orderItemRepo.IOrderItemRepository {
	return m.OrderItemRepo
}

func (m *MockRepositoryRegistry) GetPayment() paymentRepo.IPaymentRepository {
	return m.PaymentRepo
}
