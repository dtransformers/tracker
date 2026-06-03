package v1

// EventType corresponds to EventType enum
type EventType string

const (
	EventTypeNavigation     EventType = "navigation"
	EventTypeInteraction    EventType = "interaction"
	EventTypeEngagement     EventType = "engagement"
	EventTypeAuthentication EventType = "authentication"
	EventTypeCommunication  EventType = "communication"
	EventTypeCommerce       EventType = "commerce"
	EventTypeSystem         EventType = "system"
	EventTypeCustom         EventType = "custom"
)

type EventCategory string

const (
	EventCategoryPageView      EventCategory = "navigation.page_view"
	EventCategoryPageLeave     EventCategory = "navigation.page_leave"
	EventCategoryPageRedirect  EventCategory = "navigation.page_redirect"
	EventCategoryHistoryChange EventCategory = "navigation.history_change"

	EventCategoryClick       EventCategory = "interaction.click"
	EventCategoryRightClick  EventCategory = "interaction.right_click"
	EventCategoryDoubleClick EventCategory = "interaction.double_click"
	EventCategoryHover       EventCategory = "interaction.hover"
	EventCategoryScroll      EventCategory = "interaction.scroll"
	EventCategoryKeyboard    EventCategory = "interaction.keyboard"
	EventCategoryTouch       EventCategory = "interaction.touch"
	EventCategoryDrag        EventCategory = "interaction.drag"
	EventCategoryFocus       EventCategory = "interaction.focus"
	EventCategoryBlur        EventCategory = "interaction.blur"
	EventCategoryCopy        EventCategory = "interaction.copy"
	EventCategoryPaste       EventCategory = "interaction.paste"

	EventCategoryFormStart   EventCategory = "interaction.form_start"
	EventCategoryFormSubmit  EventCategory = "interaction.form_submit"
	EventCategoryFormAbandon EventCategory = "interaction.form_abandon"
	EventCategoryFormError   EventCategory = "interaction.form_error"
	EventCategoryInputChange EventCategory = "interaction.input_change"

	EventCategoryRead         EventCategory = "engagement.read"
	EventCategoryContentView  EventCategory = "engagement.content_view"
	EventCategoryContentShare EventCategory = "engagement.content_share"
	EventCategoryContentSave  EventCategory = "engagement.content_save"
	EventCategoryDownload     EventCategory = "engagement.download"
	EventCategorySearch       EventCategory = "engagement.search"

	EventCategoryVideoWatch    EventCategory = "engagement.video_watch"
	EventCategoryVideoPlay     EventCategory = "engagement.video_play"
	EventCategoryVideoPause    EventCategory = "engagement.video_pause"
	EventCategoryVideoSeek     EventCategory = "engagement.video_seek"
	EventCategoryVideoComplete EventCategory = "engagement.video_complete"
	EventCategoryVideoBuffer   EventCategory = "engagement.video_buffer"
	EventCategoryVideoQuality  EventCategory = "engagement.video_quality"
	EventCategoryVideoProgress EventCategory = "engagement.video_progress"

	EventCategoryLogin         EventCategory = "authentication.login"
	EventCategoryLogout        EventCategory = "authentication.logout"
	EventCategorySignup        EventCategory = "authentication.signup"
	EventCategoryPasswordReset EventCategory = "authentication.password_reset"
	EventCategoryTokenRefresh  EventCategory = "authentication.token_refresh"
	EventCategoryMFAChallenge  EventCategory = "authentication.mfa_challenge"

	EventCategoryMessageSent         EventCategory = "communication.message_sent"
	EventCategoryContactForm         EventCategory = "communication.contact_form"
	EventCategoryChatOpen            EventCategory = "communication.chat_open"
	EventCategoryChatClose           EventCategory = "communication.chat_close"
	EventCategoryNotificationClick   EventCategory = "communication.notification_click"
	EventCategoryNotificationDismiss EventCategory = "communication.notification_dismiss"
	EventCategorySubscribe           EventCategory = "communication.subscribe"
	EventCategoryUnsubscribe         EventCategory = "communication.unsubscribe"
	EventCategoryLeadSubmit          EventCategory = "communication.lead_submit"

	EventCategoryViewItem        EventCategory = "commerce.view_item"
	EventCategoryViewItemList    EventCategory = "commerce.view_item_list"
	EventCategorySelectItem      EventCategory = "commerce.select_item"
	EventCategoryAddToCart       EventCategory = "commerce.add_to_cart"
	EventCategoryRemoveFromCart  EventCategory = "commerce.remove_from_cart"
	EventCategoryViewCart        EventCategory = "commerce.view_cart"
	EventCategoryBeginCheckout   EventCategory = "commerce.begin_checkout"
	EventCategoryAddPaymentInfo  EventCategory = "commerce.add_payment_info"
	EventCategoryAddShippingInfo EventCategory = "commerce.add_shipping_info"
	EventCategoryPurchase        EventCategory = "commerce.purchase"
	EventCategoryRefund          EventCategory = "commerce.refund"
	EventCategoryAddToWishlist   EventCategory = "commerce.add_to_wishlist"
	EventCategoryPromoView       EventCategory = "commerce.promo_view"
	EventCategoryPromoClick      EventCategory = "commerce.promo_click"

	EventCategoryError          EventCategory = "system.error"
	EventCategoryWarning        EventCategory = "system.warning"
	EventCategoryPerformance    EventCategory = "system.performance"
	EventCategoryConsentUpdate  EventCategory = "system.consent_update"
	EventCategoryAppCrash       EventCategory = "system.app_crash"
	EventCategoryExperimentSeen EventCategory = "system.experiment_seen"
)

type ChannelType string

const (
	ChannelTypeOrganicSearch    ChannelType = "organic_search"
	ChannelTypePaidSearch       ChannelType = "paid_search"
	ChannelTypePaidSocial       ChannelType = "paid_social"
	ChannelTypeOrganicSocial    ChannelType = "organic_social"
	ChannelTypeEmail            ChannelType = "email"
	ChannelTypeSMS              ChannelType = "sms"
	ChannelTypePushNotification ChannelType = "push_notification"
	ChannelTypeReferral         ChannelType = "referral"
	ChannelTypeDisplay          ChannelType = "display"
	ChannelTypeAffiliate        ChannelType = "affiliate"
	ChannelTypeDirect           ChannelType = "direct"
	ChannelTypeUnknown          ChannelType = "unknown"
)

type AdPlatform string

const (
	AdPlatformGoogle    AdPlatform = "google"
	AdPlatformMeta      AdPlatform = "meta"
	AdPlatformTiktok    AdPlatform = "tiktok"
	AdPlatformTwitter   AdPlatform = "twitter"
	AdPlatformLinkedin  AdPlatform = "linkedin"
	AdPlatformSnapchat  AdPlatform = "snapchat"
	AdPlatformPinterest AdPlatform = "pinterest"
	AdPlatformMicrosoft AdPlatform = "microsoft"
	AdPlatformAmazon    AdPlatform = "amazon"
	AdPlatformCriteo    AdPlatform = "criteo"
)

type ContentType string

const (
	ContentTypeHome     ContentType = "home"
	ContentTypeProduct  ContentType = "product"
	ContentTypeCategory ContentType = "category"
	ContentTypeArticle  ContentType = "article"
	ContentTypeBlog     ContentType = "blog"
	ContentTypeLanding  ContentType = "landing"
	ContentTypeCheckout ContentType = "checkout"
	ContentTypeCart     ContentType = "cart"
	ContentTypeAccount  ContentType = "account"
	ContentTypeSearch   ContentType = "search"
	ContentTypeThankYou ContentType = "thank_you"
	ContentTypeError    ContentType = "error"
	ContentTypeCustom   ContentType = "custom"
)

type TargetElementType string

const (
	TargetElementTypeButton       TargetElementType = "button"
	TargetElementTypeLink         TargetElementType = "link"
	TargetElementTypeImage        TargetElementType = "image"
	TargetElementTypeVideo        TargetElementType = "video"
	TargetElementTypeForm         TargetElementType = "form"
	TargetElementTypeInput        TargetElementType = "input"
	TargetElementTypeSelect       TargetElementType = "select"
	TargetElementTypeCheckbox     TargetElementType = "checkbox"
	TargetElementTypeRadio        TargetElementType = "radio"
	TargetElementTypeTextarea     TargetElementType = "textarea"
	TargetElementTypeBanner       TargetElementType = "banner"
	TargetElementTypeModal        TargetElementType = "modal"
	TargetElementTypeNavItem      TargetElementType = "nav_item"
	TargetElementTypeProductCard  TargetElementType = "product_card"
	TargetElementTypeCarousel     TargetElementType = "carousel"
	TargetElementTypeTab          TargetElementType = "tab"
	TargetElementTypeAccordion    TargetElementType = "accordion"
	TargetElementTypeSearchBar    TargetElementType = "search_bar"
	TargetElementTypeNotification TargetElementType = "notification"
	TargetElementTypeToast        TargetElementType = "toast"
	TargetElementTypeDropdown     TargetElementType = "dropdown"
	TargetElementTypeUnknown      TargetElementType = "unknown"
)

type MouseButton string

const (
	MouseButtonLeft   MouseButton = "left"
	MouseButtonRight  MouseButton = "right"
	MouseButtonMiddle MouseButton = "middle"
)

type ScrollDirection string

const (
	ScrollDirectionUp   ScrollDirection = "up"
	ScrollDirectionDown ScrollDirection = "down"
)

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeLive  MediaType = "live"
)

type CurrencyCode string

const (
	CurrencyCodeUSD CurrencyCode = "USD"
	CurrencyCodeEUR CurrencyCode = "EUR"
	CurrencyCodeGBP CurrencyCode = "GBP"
	CurrencyCodeMAD CurrencyCode = "MAD"
	CurrencyCodeCAD CurrencyCode = "CAD"
	CurrencyCodeAUD CurrencyCode = "AUD"
	CurrencyCodeJPY CurrencyCode = "JPY"
	CurrencyCodeAED CurrencyCode = "AED"
	CurrencyCodeSAR CurrencyCode = "SAR"
)

type AuthProvider string

const (
	AuthProviderEmail    AuthProvider = "email"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderFacebook AuthProvider = "facebook"
	AuthProviderApple    AuthProvider = "apple"
	AuthProviderTwitter  AuthProvider = "twitter"
	AuthProviderGithub   AuthProvider = "github"
	AuthProviderSSO      AuthProvider = "sso"
)

type ErrorSeverity string

const (
	ErrorSeverityInfo    ErrorSeverity = "info"
	ErrorSeverityWarning ErrorSeverity = "warning"
	ErrorSeverityError   ErrorSeverity = "error"
	ErrorSeverityFatal   ErrorSeverity = "fatal"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

type ConsentCategory string

const (
	ConsentCategoryNecessary       ConsentCategory = "necessary"
	ConsentCategoryAnalytics       ConsentCategory = "analytics"
	ConsentCategoryMarketing       ConsentCategory = "marketing"
	ConsentCategoryPersonalization ConsentCategory = "personalization"
	ConsentCategoryFunctional      ConsentCategory = "functional"
)

type PlatformClickIds struct {
	Gclid         *string `json:"gclid,omitempty"`
	Gbraid        *string `json:"gbraid,omitempty"`
	Wbraid        *string `json:"wbraid,omitempty"`
	Dclid         *string `json:"dclid,omitempty"`
	Fbclid        *string `json:"fbclid,omitempty"`
	Fbp           *string `json:"fbp,omitempty"`
	Fbc           *string `json:"fbc,omitempty"`
	Ttclid        *string `json:"ttclid,omitempty"`
	Ttp           *string `json:"ttp,omitempty"`
	Msclkid       *string `json:"msclkid,omitempty"`
	Twclid        *string `json:"twclid,omitempty"`
	LiFatId       *string `json:"li_fat_id,omitempty"`
	ScCid         *string `json:"ScCid,omitempty"`
	Epik          *string `json:"epik,omitempty"`
	LiverampAtsId *string `json:"liveramp_ats_id,omitempty"`
}

type CommerceItem struct {
	ItemID        string   `json:"item_id"`
	ItemName      *string  `json:"item_name,omitempty"`
	ItemBrand     *string  `json:"item_brand,omitempty"`
	ItemCategory  *string  `json:"item_category,omitempty"`
	ItemCategory2 *string  `json:"item_category2,omitempty"`
	ItemCategory3 *string  `json:"item_category3,omitempty"`
	ItemVariant   *string  `json:"item_variant,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	Currency      *string  `json:"currency,omitempty"`
	Quantity      *int     `json:"quantity,omitempty"`
	Coupon        *string  `json:"coupon,omitempty"`
	Discount      *float64 `json:"discount,omitempty"`
	Index         *int     `json:"index,omitempty"`
	ItemListID    *string  `json:"item_list_id,omitempty"`
	ItemListName  *string  `json:"item_list_name,omitempty"`
	Affiliation   *string  `json:"affiliation,omitempty"`
}

type TargetProps struct {
	ElementType        *TargetElementType     `json:"element_type,omitempty"`
	ElementID          *string                `json:"element_id,omitempty"`
	ElementName        *string                `json:"element_name,omitempty"`
	ElementText        *string                `json:"element_text,omitempty"`
	ElementClass       *string                `json:"element_class,omitempty"`
	ElementSelector    *string                `json:"element_selector,omitempty"`
	Component          *string                `json:"component,omitempty"`
	PageName           *string                `json:"page_name,omitempty"`
	ButtonName         *string                `json:"button_name,omitempty"`
	ButtonID           *string                `json:"button_id,omitempty"`
	Href               *string                `json:"href,omitempty"`
	FormName           *string                `json:"form_name,omitempty"`
	FormID             *string                `json:"form_id,omitempty"`
	FormField          *string                `json:"form_field,omitempty"`
	FormValueHash      *string                `json:"form_value_hash,omitempty"`
	FormFieldCount     *int                   `json:"form_field_count,omitempty"`
	FormCompletionPct  *float64               `json:"form_completion_pct,omitempty"`
	MediaType          *MediaType             `json:"media_type,omitempty"`
	MediaID            *string                `json:"media_id,omitempty"`
	MediaTitle         *string                `json:"media_title,omitempty"`
	MediaProvider      *string                `json:"media_provider,omitempty"`
	MediaDurationMs    *int64                 `json:"media_duration_ms,omitempty"`
	MediaCurrentTimeMs *int64                 `json:"media_current_time_ms,omitempty"`
	MediaPercent       *float64               `json:"media_percent,omitempty"`
	MediaQuality       *string                `json:"media_quality,omitempty"`
	MediaSoundOn       *bool                  `json:"media_sound_on,omitempty"`
	MediaFullscreen    *bool                  `json:"media_fullscreen,omitempty"`
	ScrollPercent      *float64               `json:"scroll_percent,omitempty"`
	ScrollDirection    *ScrollDirection       `json:"scroll_direction,omitempty"`
	ScrollDepthPx      *int                   `json:"scroll_depth_px,omitempty"`
	MouseButton        *MouseButton           `json:"mouse_button,omitempty"`
	ClickX             *int                   `json:"click_x,omitempty"`
	ClickY             *int                   `json:"click_y,omitempty"`
	ElementX           *int                   `json:"element_x,omitempty"`
	ElementY           *int                   `json:"element_y,omitempty"`
	TransactionID      *string                `json:"transaction_id,omitempty"`
	Currency           *string                `json:"currency,omitempty"`
	Value              *float64               `json:"value,omitempty"`
	Tax                *float64               `json:"tax,omitempty"`
	Shipping           *float64               `json:"shipping,omitempty"`
	Coupon             *string                `json:"coupon,omitempty"`
	Affiliation        *string                `json:"affiliation,omitempty"`
	PaymentType        *string                `json:"payment_type,omitempty"`
	ShippingTier       *string                `json:"shipping_tier,omitempty"`
	Items              []CommerceItem         `json:"items,omitempty"`
	SearchTerm         *string                `json:"search_term,omitempty"`
	SearchResultsCount *int                   `json:"search_results_count,omitempty"`
	SearchFilters      map[string]string      `json:"search_filters,omitempty"`
	AuthProvider       *AuthProvider          `json:"auth_provider,omitempty"`
	AuthMethod         *string                `json:"auth_method,omitempty"`
	IsNewUser          *bool                  `json:"is_new_user,omitempty"`
	ErrorCode          *string                `json:"error_code,omitempty"`
	ErrorMessage       *string                `json:"error_message,omitempty"`
	ErrorSeverity      *ErrorSeverity         `json:"error_severity,omitempty"`
	ErrorStack         *string                `json:"error_stack,omitempty"`
	HttpStatus         *int                   `json:"http_status,omitempty"`
	ComponentName      *string                `json:"component_name,omitempty"`
	SharePlatform      *string                `json:"share_platform,omitempty"`
	ShareMethod        *string                `json:"share_method,omitempty"`
	DurationMs         *int64                 `json:"duration_ms,omitempty"`
	Feature            *string                `json:"feature,omitempty"`
	Priority           *Priority              `json:"priority,omitempty"`
	Custom             map[string]interface{} `json:"custom,omitempty"`
}

type TrackingEvent struct {
	EventID          string           `json:"event_id"`
	PageID           int64            `json:"page_id,omitempty"`
	Type             EventType        `json:"type" binding:"required"`
	Category         EventCategory    `json:"category" binding:"required"`
	Name             string           `json:"name" binding:"required"`
	OccurredAt       string           `json:"occurred_at" binding:"required"`
	ServerReceivedAt string           `json:"server_received_at,omitempty" swagger:"ignore"`
	Sequence         int              `json:"sequence"`
	ClickIds         PlatformClickIds `json:"click_ids"`
	Target           TargetProps      `json:"target"`
	SchemaVersion    string           `json:"schema_version"`
	SdkName          *string          `json:"sdk_name,omitempty"`
	SdkVersion       *string          `json:"sdk_version,omitempty"`
}
