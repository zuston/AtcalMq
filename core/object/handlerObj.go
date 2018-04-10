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


