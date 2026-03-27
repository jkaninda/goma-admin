package dto

// ListSnapshotsRq is the request for listing config snapshots.
type ListSnapshotsRq struct {
	Page int `query:"page" default:"0" description:"Page number (0-based)"`
	Size int `query:"size" default:"20" description:"Page size (max 100)"`
}

// SnapshotByIDRq identifies a single config snapshot.
type SnapshotByIDRq struct {
	ID int `param:"id" required:"true" description:"Snapshot ID"`
}
