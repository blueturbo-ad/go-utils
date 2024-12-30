package baseopenrtb

// openrtb 的基础数据结构，后续可以更具新版本协议作出修改，后面所有的adx的请求全部继承这个结构，然后作出相应的修改
type Publisher struct {
	Id     *string      `json:"id,omitempty"`
	Name   *string      `json:"name,omitempty"`
	Cat    *[]string    `json:"cat,omitempty"`
	Domain *string      `json:"domain,omitempty"`
	Ext    *interface{} `json:"ext,omitempty"`
}

type Producer struct {
	Id     *string      `json:"id,omitempty"`
	Name   *string      `json:"name,omitempty"`
	Cat    *[]string    `json:"cat,omitempty"`
	Domain *string      `json:"domain,omitempty"`
	Ext    *interface{} `json:"ext,omitempty"`
}

type Segment struct {
	Id    *string      `json:"id,omitempty"`
	Name  *string      `json:"name,omitempty"`
	Value *string      `json:"value,omitempty"`
	Ext   *interface{} `json:"ext,omitempty"`
}

type Data struct {
	Id      *string      `json:"id,omitempty"`
	Name    *string      `json:"name,omitempty"`
	Segment []Segment    `json:"segment,omitempty"`
	Ext     *interface{} `json:"ext,omitempty"`
}

type Content struct {
	Id                 *string      `json:"id,omitempty"`
	Episode            *int         `json:"episode,omitempty"`
	Title              *string      `json:"title,omitempty"`
	Series             *string      `json:"series,omitempty"`
	Season             *string      `json:"season,omitempty"`
	Artist             *string      `json:"artist,omitempty"`
	Genre              *string      `json:"genre,omitempty"`
	Album              *string      `json:"album,omitempty"`
	Isrc               *string      `json:"isrc,omitempty"`
	Producer           *Producer    `json:"producer,omitempty"`
	Url                *string      `json:"url,omitempty"`
	Cat                *[]string    `json:"cat,omitempty"`
	Prodq              *int         `json:"prodq,omitempty"`
	VideoQuality       *int         `json:"videoquality,omitempty"`
	Context            *int         `json:"context,omitempty"`
	ContentRating      *string      `json:"contentrating,omitempty"`
	UserRating         *string      `json:"userrating,omitempty"`
	QagMediaRating     *int         `json:"qagmediarating,omitempty"`
	Keywords           *string      `json:"keywords,omitempty"`
	Livestream         *int         `json:"livestream,omitempty"`
	SourceRelationship *int         `json:"sourcerelationship,omitempty"`
	Len                *int         `json:"len,omitempty"`
	Language           *string      `json:"language,omitempty"`
	Embeddable         *int         `json:"embeddable,omitempty"`
	Data               []Data       `json:"data,omitempty"`
	Ext                *interface{} `json:"ext,omitempty"`
}

type Site struct {
	Id            *string      `json:"id,omitempty"`
	Name          *string      `json:"name,omitempty"`
	Domain        *string      `json:"domain,omitempty"`
	Cat           *[]string    `json:"cat,omitempty"`
	SectionCat    *[]string    `json:"sectioncat,omitempty"`
	PageCat       *[]string    `json:"pagecat,omitempty"`
	Page          *string      `json:"page,omitempty"`
	Ref           *string      `json:"ref,omitempty"`
	Search        *string      `json:"search,omitempty"`
	Mobile        *int         `json:"mobile,omitempty"`
	PrivacyPolicy *int         `json:"privacypolicy,omitempty"`
	Publisher     *Publisher   `json:"publisher,omitempty"`
	Content       *Content     `json:"content,omitempty"`
	Keywords      *string      `json:"keywords,omitempty"`
	Ext           *interface{} `json:"ext,omitempty"`
}

type Metric struct {
	Type   *string      `json:"type,omitempty"`
	Value  *float32     `json:"value,omitempty"`
	Vendor *string      `json:"vendor,omitempty"`
	Ext    *interface{} `json:"ext,omitempty"`
}
type Format struct {
	W      *int `json:"w,omitempty"`
	H      *int `json:"h,omitempty"`
	WRatio *int `json:"wratio,omitempty"`
	HRatio *int `json:"hratio,omitempty"`
}

type Banner struct {
	Format   *[]Format `json:"format,omitempty"`
	W        *int      `json:"w,omitempty"`
	H        *int      `json:"h,omitempty"`
	BType    *[]int    `json:"btype,omitempty"`
	BAttr    *[]int    `json:"battr,omitempty"`
	Pos      *int      `json:"pos,omitempty"`
	Mimes    *[]string `json:"mimes,omitempty"`
	TopFrame *int      `json:"topframe,omitempty"`
	ExpDir   *[]int    `json:"expdir,omitempty"`
	Api      *[]int    `json:"api,omitempty"`
	Id       *string   `json:"id,omitempty"`
	Vcm      *int      `json:"vcm,omitempty"`
	WMax     *int      `json:"wmax,omitempty"`
	WMin     *int      `json:"wmin,omitempty"`
	HMax     *int      `json:"hmax,omitempty"`
	HMin     *int      `json:"hmin,omitempty"`
	Required *int      `json:"required,omitempty"`
}

type Video struct {
	Mimes          *[]string    `json:"mimes,omitempty"`
	MinDuration    *int32       `json:"minduration,omitempty"`
	MaxDuration    *int32       `json:"maxduration,omitempty"`
	Protocols      *[]int       `json:"protocols,omitempty"`
	W              *int         `json:"w,omitempty"`
	H              *int         `json:"h,omitempty"`
	StartDelay     *int         `json:"startdelay,omitempty"`
	Placement      *int         `json:"placement,omitempty"`
	Linearity      *int         `json:"linearity,omitempty"`
	Skip           *int         `json:"skip,omitempty"`
	SkipMin        *int         `json:"skipmin,omitempty"`
	SkipAfter      *int         `json:"skipafter,omitempty"`
	Sequence       *int         `json:"sequence,omitempty"`
	BAttr          *[]int       `json:"battr,omitempty"`
	MaxExtended    *int         `json:"maxextended,omitempty"`
	MinBitrate     *int         `json:"minbitrate,omitempty"`
	MaxBitrate     *int         `json:"maxbitrate,omitempty"`
	BoxingAllowed  *int         `json:"boxingallowed,omitempty"`
	PlaybackMethod *[]int       `json:"playbackmethod,omitempty"`
	PlaybackEnd    *int         `json:"playbackend,omitempty"`
	Delivery       *[]int       `json:"delivery,omitempty"`
	Pos            *int         `json:"pos,omitempty"`
	CompanionAd    *[]Banner    `json:"companionad,omitempty"`
	Api            *[]int       `json:"api,omitempty"`
	CompanionType  *[]int       `json:"companiontype,omitempty"`
	Ext            *interface{} `json:"ext,omitempty"`
}

type Audio struct {
	Mimes         *string `json:"mimes,omitempty"`
	MinDuration   *int    `json:"minduration,omitempty"`
	MaxDuration   *int    `json:"maxduration,omitempty"`
	Protocols     *[]int  `json:"protocols,omitempty"`
	StartDelay    *int    `json:"startdelay,omitempty"`
	Sequence      *int    `json:"sequence,omitempty"`
	BAttr         *[]int  `json:"battr,omitempty"`
	MaxExtended   *int    `json:"maxextended,omitempty"`
	MinBitrate    *int    `json:"minbitrate,omitempty"`
	MaxBitrate    *int    `json:"maxbitrate,omitempty"`
	Delivery      *[]int  `json:"delivery,omitempty"`
	Api           *[]int  `json:"api,omitempty"`
	CompanionType *[]int  `json:"companiontype,omitempty"`
	MaxSeq        *int    `json:"maxseq,omitempty"`
	Feed          *int    `json:"feed,omitempty"`
	Stitched      *int    `json:"stitched,omitempty"`
	Nvol          *int    `json:"nvol,omitempty"`
}

type Native struct {
	Request *string      `json:"request,omitempty"`
	Ver     *string      `json:"ver,omitempty"`
	Api     *[]int       `json:"api,omitempty"`
	BAttr   *[]int       `json:"battr,omitempty"`
	Ext     *interface{} `json:"ext,omitempty"`
}

type Deal struct {
	Id          *string   `json:"id,omitempty"`
	BidFloor    *float32  `json:"bidfloor,omitempty"`
	BidFloorCur *string   `json:"bidfloorcur,omitempty"`
	At          *int      `json:"at,omitempty"`
	WSeat       *[]string `json:"wseat,omitempty"`
	WAdomain    *[]string `json:"wadomain,omitempty"`
}

type Pmp struct {
	PrivateAuction *int   `json:"private_auction,omitempty"`
	Deals          []Deal `json:"deals,omitempty"`
}

type Imp struct {
	Id                *string      `json:"id,omitempty"`
	Metric            []Metric     `json:"metric,omitempty"`
	Banner            *Banner      `json:"banner,omitempty"`
	Video             *Video       `json:"video,omitempty"`
	Audio             *Audio       `json:"audio,omitempty"`
	Native            *Native      `json:"native,omitempty"`
	Pmp               *Pmp         `json:"pmp,omitempty"`
	DisplayManager    *string      `json:"displaymanager,omitempty"`
	DisplayManagerver *string      `json:"displaymanagerver,omitempty"`
	Instl             *int         `json:"instl,omitempty"`
	TagId             *string      `json:"tagid,omitempty"`
	BidFloor          *float64     `json:"bidfloor,omitempty"`
	BidFloorCur       *string      `json:"bidfloorcur,omitempty"`
	ClickBrowser      *int         `json:"clickbrowser,omitempty"`
	Secure            *int         `json:"secure,omitempty"`
	IframeBuster      *[]string    `json:"iframebuster,omitempty"`
	Ext               *interface{} `json:"ext,omitempty"`
}

type App struct {
	Id            *string      `json:"id,omitempty"`
	Name          *string      `json:"name,omitempty"`
	Bundle        *string      `json:"bundle,omitempty"`
	Domain        *string      `json:"domain,omitempty"`
	StoreUrl      *string      `json:"storeurl,omitempty"`
	Cat           *[]string    `json:"cat,omitempty"`
	SectionCat    *[]string    `json:"sectioncat,omitempty"`
	PageCat       *[]string    `json:"pagecat,omitempty"`
	Ver           *string      `json:"ver,omitempty"`
	PrivacyPolicy *int         `json:"privacypolicy,omitempty"`
	Paid          *int         `json:"paid,omitempty"`
	Publisher     *Publisher   `json:"publisher,omitempty"`
	Content       *Content     `json:"content,omitempty"`
	Keywords      *string      `json:"keywords,omitempty"`
	Ext           *interface{} `json:"ext,omitempty"`
}

type Geo struct {
	Lat           *float32     `json:"lat,omitempty"`
	Lon           *float32     `json:"lon,omitempty"`
	Type          *int         `json:"type,omitempty"`
	Accuracy      *int         `json:"accuracy,omitempty"`
	Iastfix       *int         `json:"iastfix,omitempty"`
	IpService     *int         `json:"ipservice,omitempty"`
	Country       *string      `json:"country,omitempty"`
	Region        *string      `json:"region,omitempty"`
	RegionFips104 *string      `json:"regionfips104,omitempty"`
	Metro         *string      `json:"metro,omitempty"`
	City          *string      `json:"city,omitempty"`
	Zip           *string      `json:"zip,omitempty"`
	UtcOffset     *int         `json:"utcoffset,omitempty"`
	Ext           *interface{} `json:"ext,omitempty"`
}

type Device struct {
	Ua             *string      `json:"ua,omitempty"`
	Geo            *Geo         `json:"geo,omitempty"`
	Dnt            *int         `json:"dnt,omitempty"`
	Lmt            *int         `json:"lmt,omitempty"`
	Ip             *string      `json:"ip,omitempty"`
	Ipv6           *string      `json:"ipv6,omitempty"`
	DeviceType     *int32       `json:"devicetype,omitempty"`
	Make           *string      `json:"make,omitempty"`
	Model          *string      `json:"model,omitempty"`
	Os             *string      `json:"os,omitempty"`
	Osv            *string      `json:"osv,omitempty"`
	Hwv            *string      `json:"hwv,omitempty"`
	H              *int         `json:"h,omitempty"`
	W              *int         `json:"w,omitempty"`
	Ppi            *int         `json:"ppi,omitempty"`
	PxRatio        *float32     `json:"pxratio,omitempty"`
	Js             *int         `json:"js,omitempty"`
	GeoFetch       *int         `json:"geofetch,omitempty"`
	FlashVer       *string      `json:"flashver,omitempty"`
	Language       *string      `json:"language,omitempty"`
	Carrier        *string      `json:"carrier,omitempty"`
	MccMnc         *string      `json:"mccmnc,omitempty"`
	ConnectionType *int         `json:"connectiontype,omitempty"`
	Ifa            *string      `json:"ifa,omitempty"`
	DidSha1        *string      `json:"didsha1,omitempty"`
	DidMd5         *string      `json:"didmd5,omitempty"`
	DpidSha1       *string      `json:"dpidsha1,omitempty"`
	DpidMd5        *string      `json:"dpidmd5,omitempty"`
	MacSha1        *string      `json:"macsha1,omitempty"`
	MacMd5         *string      `json:"macmd5,omitempty"`
	Ext            *interface{} `json:"ext,omitempty"`
}

type User struct {
	Id         *string      `json:"id,omitempty"`
	BuyerUid   *string      `json:"buyeruid,omitempty"`
	Yob        *int         `json:"yob,omitempty"`
	Gender     *string      `json:"gender,omitempty"`
	Keywords   *string      `json:"keywords,omitempty"`
	CustomData *string      `json:"customdata,omitempty"`
	Geo        *Geo         `json:"geo,omitempty"`
	Data       []Data       `json:"data,omitempty"`
	Ext        *interface{} `json:"ext,omitempty"`
}

type Source struct {
	Fd     *int    `json:"fd,omitempty"`
	Tid    *string `json:"tid,omitempty"`
	Pchain *string `json:"pchain,omitempty"`
}

type Regs struct {
	Coppa *int         `json:"coppa,omitempty"`
	Ext   *interface{} `json:"ext,omitempty"`
}

type BidRequest struct {
	Id      *string      `json:"id"`
	Imp     []Imp        `json:"imp,omitempty"`
	Site    *Site        `json:"site,omitempty"`
	App     *App         `json:"app,omitempty"`
	Device  *Device      `json:"device,omitempty"`
	User    *User        `json:"user,omitempty"`
	Test    *int         `json:"test,omitempty"`
	At      *int         `json:"at,omitempty"`
	TMax    *int         `json:"tmax,omitempty"`
	WSeat   *[]string    `json:"wseat,omitempty"`
	BSeat   *[]string    `json:"bseat,omitempty"`
	AllImps *int         `json:"allimps,omitempty"`
	Cur     *[]string    `json:"cur,omitempty"`
	WLang   *[]string    `json:"wlang,omitempty"`
	BCat    *[]string    `json:"bcat,omitempty"`
	BAdv    *[]string    `json:"badv,omitempty"`
	BApp    *[]string    `json:"bapp,omitempty"`
	Source  *Source      `json:"source,omitempty"`
	Regs    *Regs        `json:"regs,omitempty"`
	Ext     *interface{} `json:"ext,omitempty"`
}
