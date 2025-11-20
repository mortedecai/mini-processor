package processor

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/censys/scan-takehome/internal/database/dal"
	"github.com/censys/scan-takehome/internal/database/models"
	"github.com/censys/scan-takehome/pkg/scanning"
)

type processor struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	client       *pubsub.Client
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	wg           sync.WaitGroup
	sigChannel   chan os.Signal
	scanEntryDB  dal.Scan
}

func (p *processor) receiveLoop() {
	defer p.wg.Done()
	err := p.subscription.Receive(p.ctx, p.HandleMessage)

	if err != nil && err != context.Canceled {
		zap.S().Errorw("receive message error", "error", err)
	}
}

func (p *processor) signalHandler() {
	sig := <-p.sigChannel
	zap.S().Infow("received sigint or sigterm", "signal", sig)
	p.cancelFunc()
}

func (p *processor) Start() {
	p.sigChannel = make(chan os.Signal, 1)
	signal.Notify(p.sigChannel, syscall.SIGINT, syscall.SIGTERM)
	go p.signalHandler()

	p.wg.Add(1)
	go p.receiveLoop()

	<-p.ctx.Done()
	p.wg.Done()
}

func (p *processor) Stop() {
	p.cancelFunc()
}

func (p *processor) HandleMessage(ctx context.Context, msg *pubsub.Message) {
	zap.S().Errorw("received message", "message", msg.Data)
	var tempScan scanning.Scan
	err := json.Unmarshal(msg.Data, &tempScan)
	if err != nil {
		zap.S().Errorw("failed to unmarshal scanning.Scan message", "error", err)
		msg.Nack()
		return
	}
	var scan scanning.Scan
	switch tempScan.DataVersion {
	case scanning.V1:
		scan.Data = &scanning.V1Data{}
	case scanning.V2:
		scan.Data = &scanning.V2Data{}
	default:
		zap.S().Errorw("unknown data version", "data_version", tempScan.DataVersion)
		msg.Nack()
		return
	}
	err = json.Unmarshal(msg.Data, &scan)
	entry, err := models.NewScanEntry(scan)
	if err != nil {
		zap.S().Errorw("failed to unmarshal full scan entry", "error", err)
		msg.Nack()
		return
	}
	if err = p.scanEntryDB.Upsert(entry); err != nil {
		zap.S().Errorw("failed to upsert full scan entry", "error", err)
		msg.Nack()
		return
	}
	msg.Ack()
}

// New takes a processor Config instance and attempts to create a new processor from it.
// If the Config.Validate() function fails, the processor will not be created and the resulting error will be returned.
func New(cfg *Config, seDB dal.Scan) (*processor, error) {
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
	proc.scanEntryDB = seDB
	return proc, nil
}
