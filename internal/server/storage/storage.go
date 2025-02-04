package storage

import (
	"context"
	"time"

	v1 "github.com/plainq/plainq/internal/server/schema/v1"
)

// QueueStorage encapsulates interaction with queue storage.
type QueueStorage interface {
	// CreateQueue creates new queue.
	CreateQueue(ctx context.Context, input *v1.CreateQueueRequest) (*v1.CreateQueueResponse, error)

	// DescribeQueue returns information about specified queue.
	DescribeQueue(ctx context.Context, input *v1.DescribeQueueRequest) (*v1.DescribeQueueResponse, error)

	// ListQueues returns a list of existing queues.
	ListQueues(ctx context.Context, input *v1.ListQueuesRequest) (*v1.ListQueuesResponse, error)

	// PurgeQueue purges all messages from the queue.
	PurgeQueue(ctx context.Context, input *v1.PurgeQueueRequest) (*v1.PurgeQueueResponse, error)

	// DeleteQueue deletes a queue if it's not empty. Also supports DeleteQueueInput.Force
	// to delete queue with messages.
	DeleteQueue(ctx context.Context, input *v1.DeleteQueueRequest) (*v1.DeleteQueueResponse, error)

	// Send sends message to the queue.
	Send(ctx context.Context, input *v1.SendRequest) (*v1.SendResponse, error)

	// Receive receives message form the queue.
	Receive(ctx context.Context, input *v1.ReceiveRequest) (*v1.ReceiveResponse, error)

	// Delete delete messages from the queue.
	Delete(ctx context.Context, input *v1.DeleteRequest) (*v1.DeleteResponse, error)
}

type AccountStorage interface {
	// CreateAccount creates record with account information in database.
	CreateAccount(ctx context.Context, account Account) error

	// GetAccountByID fetches account record from database by given id.
	GetAccountByID(ctx context.Context, id string) (*Account, error)

	// GetAccountByEmail fetches account record from database by given email.
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)

	// SetAccountVerified update 'verified' field of account record in database.
	SetAccountVerified(ctx context.Context, email string, verified bool) error

	// SetAccountPassword update account 'password' field of account record in database.
	SetAccountPassword(ctx context.Context, id, password string) error

	// DeleteAccount deletes account record from database by given id.
	DeleteAccount(ctx context.Context, id string) error

	// CreateRefreshToken creates refresh token record in database.
	CreateRefreshToken(ctx context.Context, token RefreshToken) error

	// DeleteRefreshToken deletes given token from database.
	DeleteRefreshToken(ctx context.Context, token string) error

	// DeleteRefreshTokenByTokenID deletes given token from database by its id.
	DeleteRefreshTokenByTokenID(ctx context.Context, tid string) error

	// PurgeRefreshTokens deletes all refresh token records related to given account.
	PurgeRefreshTokens(ctx context.Context, aid string) error

	// DenyAccessToken denies access token by given token string.
	DenyAccessToken(ctx context.Context, token string, ttl time.Duration) error
}

// Account represents user account with all its properties.
type Account struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session represents an auth session.
type Session struct {
	// AccessToken to be used for accessing resources.
	AccessToken string

	// RefreshToken to be used to generate a new pair of tokens.
	RefreshToken string

	// Time of token creation.
	CreatedAt time.Time

	// Time of token expiry.
	ExpiresAt time.Time
}

// RefreshToken represents refresh token.
type RefreshToken struct {
	ID        string
	AID       string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
