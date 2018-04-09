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
