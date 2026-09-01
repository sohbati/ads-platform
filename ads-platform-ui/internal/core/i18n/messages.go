package i18n

// Messages holds all UI copy loaded from a language bundle file.
type Messages struct {
	Meta          MetaMessages                    `json:"meta"`
	Header        HeaderMessages                  `json:"header"`
	Nav           NavMessages                     `json:"nav"`
	Auth          AuthMessages                    `json:"auth"`
	Hero          HeroMessages                    `json:"hero"`
	Search        SearchMessages                  `json:"search"`
	Categories    CategoriesMessages              `json:"categories"`
	Cities        CitiesMessages                  `json:"cities"`
	Location      LocationMessages                `json:"location"`
	CTA           CTAMessages                     `json:"cta"`
	Section       SectionMessages                 `json:"section"`
	Footer        FooterMessages                  `json:"footer"`
	Error         ErrorMessages                   `json:"error"`
	NewAd         NewAdMessages                   `json:"new_ad"`
	MyAds         MyAdsMessages                   `json:"my_ads"`
	AdDetail      AdDetailMessages                `json:"ad_detail"`
	Appearance    AppearanceMessages              `json:"appearance"`
	ApiErrors     map[string]string               `json:"api_errors"`
	CategoryItems map[string]CategoryItemMessages `json:"category_items"`
	CityNames     map[string]string               `json:"city_names"`
}

type NavMessages struct {
	MenuAria    string `json:"menu_aria"`
	QueryAds    string `json:"query_ads"`
	MyInfo      string `json:"my_info"`
	NewAd       string `json:"new_ad"`
	Category    string `json:"category"`
	UserDetails string `json:"user_details"`
	UserAds     string `json:"user_ads"`
	MarkedAds   string `json:"marked_ads"`
	Setting     string `json:"setting"`
	Login       string `json:"login"`
	Logout      string `json:"logout"`
}

type AuthMessages struct {
	ModalTitle        string `json:"modal_title"`
	LoginTitle        string `json:"login_title"`
	MobileHeading     string `json:"mobile_heading"`
	MobileHint        string `json:"mobile_hint"`
	MobileLabel       string `json:"mobile_label"`
	MobilePlaceholder string `json:"mobile_placeholder"`
	OtpHeading        string `json:"otp_heading"`
	OtpHint           string `json:"otp_hint"`
	OtpLabel          string `json:"otp_label"`
	OtpPlaceholder    string `json:"otp_placeholder"`
	Next              string `json:"next"`
	SendOtp           string `json:"send_otp"`
	VerifyOtp         string `json:"verify_otp"`
	ChangeMobile      string `json:"change_mobile"`
	CloseAria         string `json:"close_aria"`
	TermsAccept       string `json:"terms_accept"`
	TermsLink         string `json:"terms_link"`
	PrivacyLink       string `json:"privacy_link"`
	LoginSuccess      string `json:"login_success"`
	LoginFailed       string `json:"login_failed"`
	OtpSent           string `json:"otp_sent"`
	ResendOtp         string `json:"resend_otp"`
	WelcomeUser       string `json:"welcome_user"`
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

type SearchMessages struct {
	ResultsTitle  string `json:"results_title"`
	ResultsFor    string `json:"results_for"`
	ResultsCount  string `json:"results_count"`
	Empty         string `json:"empty"`
	EmptyBrowse   string `json:"empty_browse"`
	Unavailable   string `json:"unavailable"`
	Negotiable    string `json:"negotiable"`
	Currency      string `json:"currency"`
	Prev          string `json:"prev"`
	Next          string `json:"next"`
	PageOf        string `json:"page_of"`
	LoadingMore   string `json:"loading_more"`
	LoadMoreError string `json:"load_more_error"`
	Retry         string `json:"retry"`
}

type CategoriesMessages struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ViewAll     string `json:"view_all"`
	Empty       string `json:"empty"`
	CloseAria   string `json:"close_aria"`
	LoadError   string `json:"load_error"`
}

type CitiesMessages struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AllCities   string `json:"all_cities"`
}

type LocationMessages struct {
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	SubtitleSingle    string `json:"subtitle_single"`
	SearchPlaceholder string `json:"search_placeholder"`
	Popular           string `json:"popular"`
	YourSelections    string `json:"your_selections"`
	SelectedCount     string `json:"selected_count"`
	Apply             string `json:"apply"`
	ClearAll          string `json:"clear_all"`
	CitySuffix        string `json:"city_suffix"`
	HeaderCount       string `json:"header_count"`
	LoadError         string `json:"load_error"`
	EmptySearch       string `json:"empty_search"`
	CloseAria         string `json:"close_aria"`
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

type NewAdMessages struct {
	Intro                   string `json:"intro"`
	LoadError               string `json:"load_error"`
	Category                string `json:"category"`
	CategoryPlaceholder     string `json:"category_placeholder"`
	City                    string `json:"city"`
	ChangeCity              string `json:"change_city"`
	CityRequired            string `json:"city_required"`
	Title                   string `json:"title"`
	TitlePlaceholder        string `json:"title_placeholder"`
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"description_placeholder"`
	Price                   string `json:"price"`
	PriceType               string `json:"price_type"`
	PriceFixed              string `json:"price_fixed"`
	PriceNegotiable         string `json:"price_negotiable"`
	PriceFree               string `json:"price_free"`
	PriceSalary             string `json:"price_salary"`
	Neighborhood            string `json:"neighborhood"`
	NeighborhoodPlaceholder string `json:"neighborhood_placeholder"`
	Details                 string `json:"details"`
	Pictures                string `json:"pictures"`
	PicturesHint            string `json:"pictures_hint"`
	PicturesAdd             string `json:"pictures_add"`
	PictureRemove           string `json:"picture_remove"`
	PictureView             string `json:"picture_view"`
	PictureClose            string `json:"picture_close"`
	Submit                  string `json:"submit"`
	Submitting              string `json:"submitting"`
	Success                 string `json:"success"`
	NeedCategory            string `json:"need_category"`
	NeedTitle               string `json:"need_title"`
	NeedDescription         string `json:"need_description"`
	TooManyPictures         string `json:"too_many_pictures"`
	EditHeading             string `json:"edit_heading"`
	EditIntro               string `json:"edit_intro"`
	SubmitEdit              string `json:"submit_edit"`
	SubmittingEdit          string `json:"submitting_edit"`
	PicturesReplaceHint     string `json:"pictures_replace_hint"`
	NotFound                string `json:"not_found"`
}

type MyAdsMessages struct {
	Empty       string `json:"empty"`
	Unavailable string `json:"unavailable"`
	PostCta     string `json:"post_cta"`
	Stats       string `json:"stats"`
}

type AdDetailMessages struct {
	NotFound     string `json:"not_found"`
	Unavailable  string `json:"unavailable"`
	PrevPhoto    string `json:"prev_photo"`
	NextPhoto    string `json:"next_photo"`
	PhotoCounter string `json:"photo_counter"`
	NoPhotos     string `json:"no_photos"`
	ContactLabel string `json:"contact_label"`
	ShowPhone    string `json:"show_phone"`
	CallPhone    string `json:"call_phone"`
}

type AppearanceMessages struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GroupAria   string `json:"group_aria"`
	Light       string `json:"light"`
	Dark        string `json:"dark"`
	Tide        string `json:"tide"`
}

type CategoryItemMessages struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
