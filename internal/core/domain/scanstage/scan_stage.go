package scanstage

const (
	Enqueued = "Enqueued"
	Aborted  = "Aborted"
	Done     = "Done"
	Failed   = "Failed"
)

type ScanStage struct {
	Value    int32
	Stage    string
	SubStage string
}

func (stage ScanStage) IsCompleted() bool {
	return stage.Stage == Done || stage.Stage == Failed || stage.Stage == Aborted
}

func (stage ScanStage) IsFailed() bool {
	return stage.Stage == Failed || stage.Stage == Aborted
}

func (stage ScanStage) IsQueued() bool {
	return stage.Stage == Enqueued
}
