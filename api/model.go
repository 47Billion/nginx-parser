package nginx_parser

type Entry struct {
	SubID                      string `csv:"SubID"`
	A                          string `json:"a" csv:"a"`
	Ab                         string `csv:"ab"`
	Adwords                    string `csv:"adwords"`
	Aid                        string `csv:"aid"`
	Aif1                       string `csv:"aif1"`
	Aif5                       string `csv:"aif5"`
	Aifa                       string `csv:"aifa"`
	An                         string `csv:"An"`
	And1                       string `csv:"and1"`
	Andi                       string `csv:"andi"`
	Api_key                    string `csv:"api_key"`
	Append_app_conv_trk_params string `csv:"append_app_conv_trk_params"`
	Av                         string `csv:"av"`
	Br                         string `csv:"br"`
	C                          string `csv:"c"`
	Ca                         string `csv:"ca"`
	Campaignid                 string `csv:"campaignid"`
	Cid                        string `csv:"cid"`
	City                       string `csv:"city"`
	Citycode                   string `csv:"citycode"`
	Cl                         string `csv:"cl"`
	Click_t                    string `csv:"click_t"`
	Click_ts                   string `csv:"click_ts"`
	Country                    string `csv:"country"`
	Cr                         string `csv:"cr"`
	Creative                   string `csv:"creative"`
	D                          string `csv:"d"`
	Ddl_enabled                string `csv:"ddl_enabled"`
	Ddl_to                     string `csv:"ddl_to"`
	De                         string `csv:"de"`
	Did                        string `csv:"did"`
	Dk                         string `csv:"dk"`
	Dn                         string `csv:"dn"`
	Dnt                        string `csv:"dnt"`
	E                          string `csv:"e"`
	Event_type                 string `csv:"event_type"`
	Gclid                      string `csv:"gclid"`
	Goal_AC                    string `csv:"goal_AC"`
	Goal_FP                    string `csv:"goal_FP"`
	//Goal_ac                    string `csv:"goal_ac"`
	//Goal_fp                    string `csv:"goal_fp"`
	Goal_id                    string `csv:"goal_id"`
	H                          string `csv:"h"`
	Hop_cnt                    string `csv:"hop_cnt"`
	Hops                       string `csv:"hops"`
	I                          string `csv:"i"`
	Id                         string `csv:"id"`
	Idfa                       string `csv:"idfa"`
	Idfv                       string `csv:"idfv"`
	Ifa1                       string `csv:"ifa1"`
	Ifa5                       string `csv:"ifa5"`
	Install_ts                 string `csv:"install_ts"`
	Ip                         string `csv:"ip"`
	Is_view_thru               string `csv:"is_view_thru"`
	K                          string `csv:"k"`
	Keyword                    string `csv:"keyword"`
	Lag                        string `csv:"lag"`
	Loc_physical_ms            string `csv:"loc_physical_ms"`
	Lpurl                      string `csv:"lpurl"`
	Ma                         string `csv:"ma"`
	Mac1                       string `csv:"mac1"`
	Matchtype                  string `csv:"matchtype"`
	Mm_uuid                    string `csv:"mm_uuid"`
	Mo                         string `csv:"mo"`
	N                          string `csv:"n"`
	Network                    string `csv:"network"`
	O                          string `csv:"o"`
	Odin                       string `csv:"odin"`
	Op                         string `csv:"op"`
	P                          string `csv:"p"`
	Pid                        string `csv:"pid"`
	Pl                         string `csv:"pl"`
	Pm                         string `csv:"pm"`
	Pn                         string `csv:"pn"`
	Pr                         string `csv:"pr"`
	Pref                       string `csv:"pref"`
	Product_id                 string `csv:"product_id"`
	R                          string `csv:"r"`
	Rand                       string `csv:"rand"`
	Re                         string `csv:"re"`
	Referrer                   string `csv:"referrer"`
	Region                     string `csv:"region"`
	Rt                         string `csv:"rt"`
	S                          string `csv:"s"`
	Sc                         string `csv:"sc"`
	Scs                        string `csv:"scs"`
	Sdk                        string `csv:"sdk"`
	Seq                        string `csv:"seq"`
	Siteid                     string `csv:"siteid"`
	Src                        string `csv:"src"`
	St                         string `csv:"st"`
	Strategy                   string `csv:"strategy"`
	Su                         string `csv:"su"`
	//Subid                      string `csv:"subid"`
	T                          string `csv:"t"`
	Tt_adv_id                  string `csv:"tt_adv_id"`
	U                          string `csv:"u"`
	Udi1                       string `csv:"udi1"`
	Udi5                       string `csv:"udi5"`
	Udid                       string `csv:"udid"`
	Ut                         string `csv:"ut"`
	Utm_content                string `csv:"utm_content"`
	V                          string `csv:"v"`
	Vca                        string `csv:"vca"`
	Vcr                        string `csv:"vcr"`
	Ve                         string `csv:"ve"`
	VoucherKey                 string `csv:"voucherKey"`
	Voucherkey                 string `csv:"voucherkey"`
	RemoteAddress              string    `csv:"remoteAddress"`
	Path                       string    `csv:"path"`
	TimeLocal                  string    `csv:"timeLocal"`
}

type Record struct {
	RemoteAddress string
	TimeLocal     string
	Path          string
	Args          map[string][]string
}