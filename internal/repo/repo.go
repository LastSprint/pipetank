package repo

import (
	"context"

	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/LastSprint/pipetank/pkg/utils"
)

type Repo struct {
	client *mdb.Client
	clock  utils.Clock
}

func NewRepo(
	ctx context.Context,
	client *mdb.Client,
	clock utils.Clock,
) (*Repo, error) {
	r := &Repo{
		client: client,
		clock:  clock,
	}

	err := r.registerIndexes(ctx)

	return r, err
}

func (r *Repo) registerIndexes(ctx context.Context) error {
	err := r.createRawEventIndexes(ctx)
	if err != nil {
		return err
	}

	err = r.createSingleStageIndexes(ctx)
	if err != nil {
		return err
	}

	return nil
}
