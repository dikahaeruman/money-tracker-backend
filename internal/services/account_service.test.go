package services

import (
	"context"
	"money-tracker-backend/internal/models"
	"testing"
)

type mockAccountRepository struct {
	accounts        map[string]*models.Account
	userAccounts    map[int][]*models.Account
	CreateAccountFn func(ctx context.Context, account *models.Account) (*models.Account, error)
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{
		accounts:     make(map[string]*models.Account),
		userAccounts: make(map[int][]*models.Account),
	}
}

func (m *mockAccountRepository) GetAccountByID(ctx context.Context, accountID string) (*models.Account, error) {
	return nil, nil
}

func (m *mockAccountRepository) GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error) {
	return nil, nil
}

func (m *mockAccountRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	return nil, nil
}

func (m *mockAccountRepository) DeleteAccount(ctx context.Context, accountID string) error {
	return nil
}

func (m *mockAccountRepository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	if m.CreateAccountFn != nil {
		return m.CreateAccountFn(ctx, account)
	}
	return nil, nil
}

func TestCreateAccount(t *testing.T) {
	// Create a mock repository
	mockRepo := &mockAccountRepository{}

	// Create an instance of the account service with the mock repository
	accountService := NewAccountService(mockRepo)

	// Create a test account
	testAccount := &models.Account{
		UserID:      1,
		AccountName: "Test Account",
		Balance:     1000.00,
		Currency:    "USD",
	}

	// Mock the CreateAccount method of the repository
	mockRepo.CreateAccountFn = func(ctx context.Context, account *models.Account) (*models.Account, error) {
		account.ID = "test-account-id"
		return account, nil
	}

	// Call the CreateAccount method
	createdAccount, err := accountService.CreateAccount(context.Background(), testAccount)

	// Assert that no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Assert that the created account is not nil
	if createdAccount == nil {
		t.Error("Expected created account to not be nil")
	}

	// Assert that the created account has the expected values
	if createdAccount.ID != "test-account-id" {
		t.Errorf("Expected account ID to be 'test-account-id', got %s", createdAccount.ID)
	}
	if createdAccount.UserID != testAccount.UserID {
		t.Errorf("Expected UserID to be %d, got %d", testAccount.UserID, createdAccount.UserID)
	}
	if createdAccount.AccountName != testAccount.AccountName {
		t.Errorf("Expected AccountName to be %s, got %s", testAccount.AccountName, createdAccount.AccountName)
	}
	if createdAccount.Balance != testAccount.Balance {
		t.Errorf("Expected Balance to be %.2f, got %.2f", testAccount.Balance, createdAccount.Balance)
	}
	if createdAccount.Currency != testAccount.Currency {
		t.Errorf("Expected Currency to be %s, got %s", testAccount.Currency, createdAccount.Currency)
	}
}
