package scheduler

var (
	stateChanBuf int
	statesCapacity int
)

func init() {
	stateChanBuf = 10
	statesCapacity = 5
}
