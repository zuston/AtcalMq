package object

type PushMapperObj struct {
	TABLENAME string
	QUEUENAME string
	REALTIONS []PushRelation
}

type PushRelation struct {
	// table db column name simplify name
	CN string
	// queue name simplify name
	QN string
}

