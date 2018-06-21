package core



const (
	QUEUE_DATA_CENTERLOAD = "ane_its_ai_data_centerLoad_queue"

	QUEUE_BIZ_ORDER = "ane_its_ai_biz_order_queue"

	QUEUE_BIZ_EWB = "ane_its_ai_biz_ewb_queue"

	QUEUE_BASIC_ROUTE = "ane_its_ai_basic_route_queue"

	QUEUE_BIZ_EWBSLIST = "ane_its_ai_biz_ewbsList_queue"

	QUEUE_DATA_SITELOAD = "ane_its_ai_data_siteLoad_queue"

	QUEUE_DATA_CENTERUNLOAD = "ane_its_ai_data_centerUnload_queue"

	QUEUE_DATA_CENTERPALLET = "ane_its_ai_data_centerPallet_queue"

	QUEUE_DATA_CENTERSORT = "ane_its_ai_data_centerSort_queue"

	QUEUE_BASIC_AREA = "ane_its_ai_basic_area_queue"

	QUEUE_BASIC_ATTEND = "ane_its_ai_basic_attend_queue"

	QUEUE_TRIGGER_SITESEND = "ane_its_ai_trigger_siteSend_queue"

	QUEUE_TRIGGER_SITEUPLOAD = "ane_its_ai_trigger_siteUpload_queue"

	QUEUE_TRIGGER_INOROUT = "ane_its_ai_trigger_inOrOut_queue"

	QUEUE_TRIGGER_STAYORLEAVE = "ane_its_ai_trigger_stayOrLeave_queue"

	QUEUE_TRIGGER_CENTERTRANSPORT = "ane_its_ai_trigger_centerTransport_queue"

	QUEUE_BASIC_SITE = "ane_its_ai_basic_site_queue"

	QUEUE_BASIC_VEHICLELINE = "ane_its_ai_basic_vehicleLine_queue"

	QUEUE_BASIC_PLATFORM = "ane_its_ai_basic_platform_queue"
)

var BasicInfoTableNames = []string{
	QUEUE_DATA_CENTERLOAD,
	QUEUE_BIZ_ORDER,
	QUEUE_BIZ_EWB,
	QUEUE_BASIC_ROUTE,
	QUEUE_BIZ_EWBSLIST,
	QUEUE_DATA_SITELOAD,

	QUEUE_DATA_CENTERUNLOAD,

	QUEUE_DATA_CENTERPALLET,

	QUEUE_DATA_CENTERSORT,
	QUEUE_BASIC_AREA,

	QUEUE_BASIC_ATTEND,

	QUEUE_TRIGGER_SITESEND,

	QUEUE_TRIGGER_SITEUPLOAD,

	QUEUE_TRIGGER_INOROUT,

	QUEUE_TRIGGER_STAYORLEAVE,

	QUEUE_TRIGGER_CENTERTRANSPORT,

	QUEUE_BASIC_SITE,

	QUEUE_BASIC_VEHICLELINE,

	QUEUE_BASIC_PLATFORM,
}

var LinkTableNames = []string{
	"Link_Ewb",
	"Link_Site",
	"Link_Vehicle",
	"Link_Operator",
}

// 若为uid ，即无对应的key值，即用uid,用多版本记录
var LinkKey = map[string]string{
	QUEUE_DATA_CENTERLOAD : "scanTime"                                  ,
	QUEUE_BIZ_ORDER : "ewbDate"                                         ,
	QUEUE_BIZ_EWB : "uid"                                               ,
	QUEUE_BASIC_ROUTE : "startDate"                                     ,
	QUEUE_BIZ_EWBSLIST : "taskStartTime"                                ,
	QUEUE_DATA_SITELOAD : "scanTime"                                    ,
	QUEUE_DATA_CENTERUNLOAD : "scanTime"                                ,
	QUEUE_DATA_CENTERPALLET : "palletTime"                              ,
	QUEUE_DATA_CENTERSORT : "sortTime"                                  ,
	QUEUE_BASIC_AREA : "uid"                                            ,
	QUEUE_BASIC_ATTEND : "attendanceDay"                                ,
	QUEUE_TRIGGER_SITESEND : "triggerTime"                              ,
	QUEUE_TRIGGER_SITEUPLOAD : "triggerTime"                            ,
	QUEUE_TRIGGER_INOROUT : "triggerTime"                               ,
	QUEUE_TRIGGER_STAYORLEAVE : "triggerTime"                           ,
	QUEUE_TRIGGER_CENTERTRANSPORT : "triggerTime"                       ,
	QUEUE_BASIC_SITE : "createdTime"                                    ,
	QUEUE_BASIC_VEHICLELINE : "createdTime"                             ,
	QUEUE_BASIC_PLATFORM : "uid"                                        ,
}

const(
	MONGO_TAG = iota
	MYSQL_TAG
)

var MultiRelationSavingTableNames = []string{
	QUEUE_BIZ_EWB,
	QUEUE_BASIC_ROUTE,
	QUEUE_DATA_CENTERSORT,
	QUEUE_BIZ_EWBSLIST,
	QUEUE_DATA_SITELOAD,
	QUEUE_DATA_CENTERUNLOAD,
	QUEUE_DATA_CENTERPALLET,
	QUEUE_DATA_CENTERLOAD,
}