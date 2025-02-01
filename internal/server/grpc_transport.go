package server

import (
	"context"

	v1 "github.com/plainq/plainq/internal/server/schema/v1"
	"github.com/plainq/servekit/respond"
)

func (s *PlainQ) ListQueues(ctx context.Context, r *v1.ListQueuesRequest) (*v1.ListQueuesResponse, error) {
	output, listErr := s.storage.ListQueues(ctx, r)
	if listErr != nil {
		return respond.ErrorGRPC[*v1.ListQueuesResponse](ctx, listErr)
	}

	return output, nil
}

func (s *PlainQ) DescribeQueue(ctx context.Context, r *v1.DescribeQueueRequest) (*v1.DescribeQueueResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.DescribeQueueResponse](ctx, err)
	}

	output, createErr := s.storage.DescribeQueue(ctx, r)
	if createErr != nil {
		return respond.ErrorGRPC[*v1.DescribeQueueResponse](ctx, createErr)
	}

	return output, nil
}

func (s *PlainQ) CreateQueue(ctx context.Context, r *v1.CreateQueueRequest) (*v1.CreateQueueResponse, error) {
	output, createErr := s.storage.CreateQueue(ctx, r)
	if createErr != nil {
		return respond.ErrorGRPC[*v1.CreateQueueResponse](ctx, createErr)
	}

	return output, nil
}

func (s *PlainQ) DeleteQueue(ctx context.Context, r *v1.DeleteQueueRequest) (*v1.DeleteQueueResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.DeleteQueueResponse](ctx, err)
	}

	if _, err := s.storage.DeleteQueue(ctx, r); err != nil {
		return respond.ErrorGRPC[*v1.DeleteQueueResponse](ctx, err)
	}

	return &v1.DeleteQueueResponse{}, nil
}

func (s *PlainQ) PurgeQueue(ctx context.Context, r *v1.PurgeQueueRequest) (*v1.PurgeQueueResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.PurgeQueueResponse](ctx, err)
	}

	output, purgeErr := s.storage.PurgeQueue(ctx, r)
	if purgeErr != nil {
		return respond.ErrorGRPC[*v1.PurgeQueueResponse](ctx, purgeErr)
	}

	return output, nil
}

func (s *PlainQ) Send(ctx context.Context, r *v1.SendRequest) (*v1.SendResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.SendResponse](ctx, err)
	}

	output, sendErr := s.storage.Send(ctx, r)
	if sendErr != nil {
		return respond.ErrorGRPC[*v1.SendResponse](ctx, sendErr)
	}

	return output, nil
}

func (s *PlainQ) Receive(ctx context.Context, r *v1.ReceiveRequest) (*v1.ReceiveResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.ReceiveResponse](ctx, err)
	}

	output, receiveErr := s.storage.Receive(ctx, r)
	if receiveErr != nil {
		return respond.ErrorGRPC[*v1.ReceiveResponse](ctx, receiveErr)
	}

	return output, nil
}

func (s *PlainQ) Delete(ctx context.Context, r *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := validateQueueIDFromRequest(r); err != nil {
		return respond.ErrorGRPC[*v1.DeleteResponse](ctx, err)
	}

	output, deleteErr := s.storage.Delete(ctx, r)
	if deleteErr != nil {
		return respond.ErrorGRPC[*v1.DeleteResponse](ctx, deleteErr)
	}

	return output, nil
}
