package serializer

type ServiceListItem struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	ServiceDesc string `json:"service_desc"`
	LoadType    int    `json:"load_type" `
	ServiceAddr string `json:"service_addr"`
	Qps         int64  `json:"qps" `
	Qpd         int64  `json:"qpd" `
	TotalNode   int    `json:"total_node"`
}

type ServiceList struct {
	Total int                `json:"total"`
	List  []*ServiceListItem `json:"list"`
}
