package disk

type Disk struct {
	MountPoint  string
	Size        int64
	Used        int64
	Free        int64
	UsedPercent float64
}
