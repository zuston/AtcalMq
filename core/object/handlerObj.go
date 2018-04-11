package object

type CenterLoadObj struct {
	PlatformCode string
	ScanType int
	SiteCode string
	ScanTime string
	SiteId int
	Volume float32
	Weight int
	DataType int
	EwbNo string
	EwbsListNo string
	NextSiteCode string
	NextSiteId int
	OperatorCode string
	ScanMan string
	HewbNo string
}


type CenterUnLoadObj struct {
	Weight int
	EwbNo string
	EwbsListNo string
	PrevSiteCode string
	ScanMan string
	ScanTime string
	ScanType int
	Volume int
	HewbNo string
	OperatorCode string
	PrevSiteId int
	SiteCode string
}

type CenterSortObj struct {
	SiteId int
	SortTime string
	SortType int
	EwbNo string
	HewbNo string
	OperatorCode string
	PalletNo string
	SiteCode string
}

type CenterPalletObj struct {
	SiteCode string
	AreaCode string
	PlatformCode string
	PrevSiteName string
	ServicesType string
	NextSiteName string
	PalletType string
	Weight float32
	Volume float32
	EwbNo string
	EwbsListNo string
	HewbNo string
	PalletNo string
	PalletUserCode string
}

// 2018-04-11

type CenterTransportObj struct {
	NextSiteName string
	CarNo string
	TaskNo string
	TriggerTime string
	ViaOrder int
}

type EwbsListObj struct {
	EwbsListNo string
	OperationType int
	TaskStartTime string
	TaskType int
}

// 未写
type VehicleLineObj struct {
	ViaCount int
	StartSiteName string
	TaskDetail []taskDetailObj
	CarType string
	CarNo string
	CreatedTime string
	ModifiedTime string
	IsDeleted int
	EndSiteName string
	TaskNo string
	LineName string
	DriveName string
	PhoneNo string
	LoadModel string
}

type taskDetailObj struct {
	TaskNo string
	ActualOutTime string
	ViaOrder int
	ViaSiteName string
	NextSiteName string
	ActualArriveTime string
	PlanArriveTime string
	PlanOutTime string
}






