package processor

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type processor struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	client       *pubsub.Client
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

// New takes a processor Config instance and attempts to create a new processor from it.
// If the Config.Validate() function fails, the processor will not be created and the resulting error will be returned.
func New(cfg *Config) (*processor, error) {
	err := cfg.Validate()
	if err != nil {
		zap.S().Errorw("configuration error for new client", "error", err)
		return nil, err
	}
	proc := &processor{}
	proc.ctx, proc.cancelFunc = context.WithCancel(context.Background())
	proc.client, err = pubsub.NewClient(proc.ctx, cfg.ProjectID)
	if err != nil {
		zap.S().Errorw("client instantiation error", "error", err)
		return nil, err
	}
	proc.topic = proc.client.Topic(cfg.TopicID)
	proc.subscription = proc.client.Subscription(cfg.SubscriptionID)

	if exists, err := proc.subscription.Exists(proc.ctx); !exists || err != nil {
		zap.S().Errorw("could not validate subscription", "error", err)
		return nil, errors.New("subscription does not exist")
	}
	return proc, nil
}
