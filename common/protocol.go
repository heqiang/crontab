package common

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"expr"`
}

type Response struct {
	Error int         `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}
