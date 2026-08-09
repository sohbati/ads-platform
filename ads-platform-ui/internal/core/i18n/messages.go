package i18n

// Messages holds all UI copy loaded from a language bundle file.
type Messages struct {
	Meta          MetaMessages                    `json:"meta"`
	Header        HeaderMessages                  `json:"header"`
	Nav           NavMessages                     `json:"nav"`
	Auth          AuthMessages                    `json:"auth"`
	Hero          HeroMessages                    `json:"hero"`
	Categories    CategoriesMessages              `json:"categories"`
	Cities        CitiesMessages                  `json:"cities"`
	CTA           CTAMessages                     `json:"cta"`
	Section       SectionMessages                 `json:"section"`
	Footer        FooterMessages                  `json:"footer"`
	Error         ErrorMessages                   `json:"error"`
	CategoryItems map[string]CategoryItemMessages `json:"category_items"`
	CityNames     map[string]string               `json:"city_names"`
}

type NavMessages struct {
	MenuAria     string `json:"menu_aria"`
	QueryAds     string `json:"query_ads"`
	MyInfo       string `json:"my_info"`
	NewAd        string `json:"new_ad"`
	Category     string `json:"category"`
	Location     string `json:"location"`
	UserDetails  string `json:"user_details"`
	UserAds      string `json:"user_ads"`
	MarkedAds    string `json:"marked_ads"`
	Setting      string `json:"setting"`
	Login        string `json:"login"`
	Logout       string `json:"logout"`
}

type AuthMessages struct {
	LoginTitle       string `json:"login_title"`
	MobileLabel      string `json:"mobile_label"`
	MobilePlaceholder string `json:"mobile_placeholder"`
	OtpLabel         string `json:"otp_label"`
	OtpPlaceholder   string `json:"otp_placeholder"`
	SendOtp          string `json:"send_otp"`
	VerifyOtp        string `json:"verify_otp"`
	LoginSuccess     string `json:"login_success"`
	LoginFailed      string `json:"login_failed"`
	OtpSent          string `json:"otp_sent"`
	WelcomeUser      string `json:"welcome_user"`
}

type SectionMessages struct {
	ComingSoon string `json:"coming_soon"`
}

type MetaMessages struct {
	SiteDescription string `json:"site_description"`
}

type HeaderMessages struct {
	HomeAria          string `json:"home_aria"`
	Categories        string `json:"categories"`
	SearchLabel       string `json:"search_label"`
	SearchPlaceholder string `json:"search_placeholder"`
	SearchSubmit      string `json:"search_submit"`
	CityAria          string `json:"city_aria"`
	PostAd            string `json:"post_ad"`
	AccountAria       string `json:"account_aria"`
}

type HeroMessages struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type CategoriesMessages struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ViewAll     string `json:"view_all"`
	Empty       string `json:"empty"`
}

type CitiesMessages struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AllCities   string `json:"all_cities"`
}

type CTAMessages struct {
	AriaLabel   string `json:"aria_label"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Button      string `json:"button"`
}

type FooterMessages struct {
	Tagline string `json:"tagline"`
	NavAria string `json:"nav_aria"`
	About   string `json:"about"`
	Help    string `json:"help"`
	Terms   string `json:"terms"`
	Privacy string `json:"privacy"`
}

type ErrorMessages struct {
	Title       string `json:"title"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	BackHome    string `json:"back_home"`
	Code        string `json:"code"`
}

type CategoryItemMessages struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
