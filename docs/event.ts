export enum EventType {
    NAVIGATION = "navigation",
    INTERACTION = "interaction",
    ENGAGEMENT = "engagement",
    AUTHENTICATION = "authentication",
    COMMUNICATION = "communication",
    COMMERCE = "commerce",
    SYSTEM = "system",
    CUSTOM = "custom",
}

export enum EventCategory {
    PAGE_VIEW = "navigation.page_view",
    PAGE_LEAVE = "navigation.page_leave",
    PAGE_REDIRECT = "navigation.page_redirect",
    HISTORY_CHANGE = "navigation.history_change",

    CLICK = "interaction.click",
    RIGHT_CLICK = "interaction.right_click",
    DOUBLE_CLICK = "interaction.double_click",
    HOVER = "interaction.hover",
    SCROLL = "interaction.scroll",
    KEYBOARD = "interaction.keyboard",
    TOUCH = "interaction.touch",
    DRAG = "interaction.drag",
    FOCUS = "interaction.focus",
    BLUR = "interaction.blur",
    COPY = "interaction.copy",
    PASTE = "interaction.paste",

    FORM_START = "interaction.form_start",
    FORM_SUBMIT = "interaction.form_submit",
    FORM_ABANDON = "interaction.form_abandon",
    FORM_ERROR = "interaction.form_error",
    INPUT_CHANGE = "interaction.input_change",

    READ = "engagement.read",
    CONTENT_VIEW = "engagement.content_view",
    CONTENT_SHARE = "engagement.content_share",
    CONTENT_SAVE = "engagement.content_save",
    DOWNLOAD = "engagement.download",
    SEARCH = "engagement.search",

    VIDEO_WATCH = "engagement.video_watch",
    VIDEO_PLAY = "engagement.video_play",
    VIDEO_PAUSE = "engagement.video_pause",
    VIDEO_SEEK = "engagement.video_seek",
    VIDEO_COMPLETE = "engagement.video_complete",
    VIDEO_BUFFER = "engagement.video_buffer",
    VIDEO_QUALITY = "engagement.video_quality",
    VIDEO_PROGRESS = "engagement.video_progress",

    // ── Authentication ─────────────────────────────────────
    LOGIN = "authentication.login",
    LOGOUT = "authentication.logout",
    SIGNUP = "authentication.signup",
    PASSWORD_RESET = "authentication.password_reset",
    TOKEN_REFRESH = "authentication.token_refresh",
    MFA_CHALLENGE = "authentication.mfa_challenge",

    // ── Communication ──────────────────────────────────────
    MESSAGE_SENT = "communication.message_sent",
    CONTACT_FORM = "communication.contact_form",
    CHAT_OPEN = "communication.chat_open",
    CHAT_CLOSE = "communication.chat_close",
    NOTIFICATION_CLICK = "communication.notification_click",
    NOTIFICATION_DISMISS = "communication.notification_dismiss",
    SUBSCRIBE = "communication.subscribe",
    UNSUBSCRIBE = "communication.unsubscribe",
    LEAD_SUBMIT = "communication.lead_submit",

    // ── Commerce ───────────────────────────────────────────
    VIEW_ITEM = "commerce.view_item",
    VIEW_ITEM_LIST = "commerce.view_item_list",
    SELECT_ITEM = "commerce.select_item",
    ADD_TO_CART = "commerce.add_to_cart",
    REMOVE_FROM_CART = "commerce.remove_from_cart",
    VIEW_CART = "commerce.view_cart",
    BEGIN_CHECKOUT = "commerce.begin_checkout",
    ADD_PAYMENT_INFO = "commerce.add_payment_info",
    ADD_SHIPPING_INFO = "commerce.add_shipping_info",
    PURCHASE = "commerce.purchase",
    REFUND = "commerce.refund",
    ADD_TO_WISHLIST = "commerce.add_to_wishlist",
    PROMO_VIEW = "commerce.promo_view",
    PROMO_CLICK = "commerce.promo_click",

    // ── System ─────────────────────────────────────────────
    ERROR = "system.error",
    WARNING = "system.warning",
    PERFORMANCE = "system.performance",
    CONSENT_UPDATE = "system.consent_update",
    APP_CRASH = "system.app_crash",
    EXPERIMENT_SEEN = "system.experiment_seen",
}




export enum ChannelType {
    ORGANIC_SEARCH = "organic_search",
    PAID_SEARCH = "paid_search",
    PAID_SOCIAL = "paid_social",
    ORGANIC_SOCIAL = "organic_social",
    EMAIL = "email",
    SMS = "sms",
    PUSH_NOTIFICATION = "push_notification",
    REFERRAL = "referral",
    DISPLAY = "display",
    AFFILIATE = "affiliate",
    DIRECT = "direct",
    UNKNOWN = "unknown",
}

export enum AdPlatform {
    GOOGLE = "google",
    META = "meta",
    TIKTOK = "tiktok",
    TWITTER = "twitter",
    LINKEDIN = "linkedin",
    SNAPCHAT = "snapchat",
    PINTEREST = "pinterest",
    MICROSOFT = "microsoft",
    AMAZON = "amazon",
    CRITEO = "criteo",
}


export enum ContentType {
    HOME = "home",
    PRODUCT = "product",
    CATEGORY = "category",
    ARTICLE = "article",
    BLOG = "blog",
    LANDING = "landing",
    CHECKOUT = "checkout",
    CART = "cart",
    ACCOUNT = "account",
    SEARCH = "search",
    THANK_YOU = "thank_you",
    ERROR = "error",
    CUSTOM = "custom",
}

export enum TargetElementType {
    BUTTON = "button",
    LINK = "link",
    IMAGE = "image",
    VIDEO = "video",
    FORM = "form",
    INPUT = "input",
    SELECT = "select",
    CHECKBOX = "checkbox",
    RADIO = "radio",
    TEXTAREA = "textarea",
    BANNER = "banner",
    MODAL = "modal",
    NAV_ITEM = "nav_item",
    PRODUCT_CARD = "product_card",
    CAROUSEL = "carousel",
    TAB = "tab",
    ACCORDION = "accordion",
    SEARCH_BAR = "search_bar",
    NOTIFICATION = "notification",
    TOAST = "toast",
    DROPDOWN = "dropdown",
    UNKNOWN = "unknown",
}

export enum MouseButton {
    LEFT = "left",
    RIGHT = "right",
    MIDDLE = "middle",
}

export enum ScrollDirection {
    UP = "up",
    DOWN = "down",
}

export enum MediaType {
    VIDEO = "video",
    AUDIO = "audio",
    LIVE = "live",
}


export enum CurrencyCode {
    USD = "USD",
    EUR = "EUR",
    GBP = "GBP",
    MAD = "MAD",
    CAD = "CAD",
    AUD = "AUD",
    JPY = "JPY",
    AED = "AED",
    SAR = "SAR",
}

export enum AuthProvider {
    EMAIL = "email",
    GOOGLE = "google",
    FACEBOOK = "facebook",
    APPLE = "apple",
    TWITTER = "twitter",
    GITHUB = "github",
    SSO = "sso",
}

export enum ErrorSeverity {
    INFO = "info",
    WARNING = "warning",
    ERROR = "error",
    FATAL = "fatal",
}

export enum Priority {
    LOW = "low",
    MEDIUM = "medium",
    HIGH = "high",
    CRITICAL = "critical",
}

export enum ConsentCategory {
    NECESSARY = "necessary",
    ANALYTICS = "analytics",
    MARKETING = "marketing",
    PERSONALIZATION = "personalization",
    FUNCTIONAL = "functional",
}


export interface PlatformClickIds {

    // ── Google ────────────────────────────────────────────────
    /**
     * Google Ads Click ID
     * Source: URL param `?gclid=`  /  cookie `_gcl_aw`
     * Used in: GA4 Measurement Protocol, Google Ads Conversion Import
     * Lifetime: 90 days
     */
    gclid?: string;

    /**
     * Google App Campaign Click ID (privacy-safe, iOS 14+, app-to-web)
     * Source: URL param `?gbraid=`
     */
    gbraid?: string;

    /**
     * Google Web-to-App Click ID (privacy-safe)
     * Source: URL param `?wbraid=`
     */
    wbraid?: string;

    /**
     * DoubleClick / Display & Video 360 Click ID
     * Source: URL param `?dclid=`
     */
    dclid?: string;

    // ── Meta (Facebook / Instagram) ──────────────────────────
    /**
     * Facebook Ads Click ID — appended to destination URLs by Meta Ads Manager.
     * Source: URL param `?fbclid=`
     * Used to populate the `_fbc` cookie (see below).
     */
    fbclid?: string;

    /**
     * Facebook Browser ID — generated by the Meta Pixel on first page load.
     * Source: cookie `_fbp`  (format: `fb.1.<timestamp>.<random>`)
     * Required field for Meta CAPI user_data.
     * Lifetime: 90 days
     */
    fbp?: string;

    /**
     * Facebook Click Cookie — derived from fbclid by the Meta Pixel.
     * Source: cookie `_fbc`  (format: `fb.1.<timestamp>.<fbclid>`)
     * Required field for Meta CAPI user_data when fbclid is present.
     */
    fbc?: string;

    // ── TikTok ────────────────────────────────────────────────
    /**
     * TikTok Ads Click ID — appended to destination URLs by TikTok Ads.
     * Source: URL param `?ttclid=`
     * Used in TikTok Events API (server-side pixel).
     * Lifetime: 7 days
     */
    ttclid?: string;

    /**
     * TikTok Pixel Browser ID — generated by the TikTok Pixel on first load.
     * Source: cookie `_ttp`
     * Required field for TikTok Events API user object.
     */
    ttp?: string;

    // ── Microsoft / Bing Ads ─────────────────────────────────
    /**
     * Microsoft Advertising (Bing) Click ID
     * Source: URL param `?msclkid=`  /  auto-tagging (UET tag)
     * Used in Microsoft Advertising Conversion Tracking.
     * Lifetime: 90 days
     */
    msclkid?: string;

    // ── Twitter / X ───────────────────────────────────────────
    /**
     * Twitter (X) Ads Click ID
     * Source: URL param `?twclid=`
     * Used in Twitter Conversions API.
     */
    twclid?: string;

    // ── LinkedIn ──────────────────────────────────────────────
    /**
     * LinkedIn First-Party Ad Tracking ID
     * Source: URL param `?li_fat_id=`  /  cookie `li_fat_id`
     * Used in LinkedIn Conversions API.
     * Lifetime: 30 days
     */
    li_fat_id?: string;

    // ── Snapchat ──────────────────────────────────────────────
    /**
     * Snapchat Ads Click ID
     * Source: URL param `?ScCid=`
     * Used in Snapchat Conversions API.
     */
    ScCid?: string;

    // ── Pinterest ─────────────────────────────────────────────
    /**
     * Pinterest Click ID
     * Source: URL param `?epik=`  /  cookie `_epik`
     * Used in Pinterest Conversions API.
     */
    epik?: string;

    // ── LiveRamp ──────────────────────────────────────────────
    /**
     * LiveRamp Authenticated Traffic Solution (ATS) ID.
     * Enables identity resolution and addressability across publishers
     * without third-party cookies.
     */
    liveramp_ats_id?: string;
}



export interface CommerceItem {
    item_id: string;
    item_name?: string;
    item_brand?: string;
    item_category?: string;     // "Clothing"
    item_category2?: string;    // "T-Shirts"
    item_category3?: string;
    item_variant?: string;      // "Blue / XL"
    price?: number;
    currency?: CurrencyCode | string;
    quantity?: number;
    coupon?: string;
    discount?: number;
    index?: number;             // Position in list (for list-click attribution)
    item_list_id?: string;      // "homepage_recommendations"
    item_list_name?: string;    // "Recommended for You"
    affiliation?: string;       // "Online Store"
}


export interface TargetProps {

    // ── DOM Element Identification ────────────────────────────
    /** Semantic element type */
    element_type?: TargetElementType;
    /** Element `id` attribute */
    element_id?: string;
    /** Element `name` attribute */
    element_name?: string;
    /** Visible text / accessible label */
    element_text?: string;
    /** CSS class string */
    element_class?: string;
    /** Full CSS selector path (for session replay tools) */
    element_selector?: string;
    /** Component/module name (React, Vue, Angular component) */
    component?: string;

    // ── Human-Readable Semantic Labels ────────────────────────
    page_name?: string;
    button_name?: string;
    button_id?: string;
    href?: string;

    // ── Form ──────────────────────────────────────────────────
    form_name?: string;
    form_id?: string;
    /** Which field triggered the event */
    form_field?: string;
    /**
     * SHA-256 hashed field value.
     * NEVER store plain-text form values — hash client-side.
     */
    form_value_hash?: string;
    form_field_count?: number;
    /** 0–100 % of fields completed at abandon/submit time */
    form_completion_pct?: number;

    // ── Media / Video ─────────────────────────────────────────
    media_type?: MediaType;
    media_id?: string;
    media_title?: string;
    /** "youtube" | "vimeo" | "native" | "mux" */
    media_provider?: string;
    media_duration_ms?: number;
    media_current_time_ms?: number;
    /** 0–100 % playback position */
    media_percent?: number;
    /** "1080p" | "720p" | "auto" */
    media_quality?: string;
    media_sound_on?: boolean;
    media_fullscreen?: boolean;

    // ── Scroll ────────────────────────────────────────────────
    scroll_percent?: number;
    scroll_direction?: ScrollDirection;
    scroll_depth_px?: number;

    // ── Mouse / Pointer ───────────────────────────────────────
    mouse_button?: MouseButton;
    /** Click X relative to viewport */
    click_x?: number;
    click_y?: number;
    /** Click X relative to element top-left */
    element_x?: number;
    element_y?: number;

    // ── Commerce ──────────────────────────────────────────────
    /** Order / transaction ID (for Purchase events) */
    transaction_id?: string;
    /** ISO 4217 currency code */
    currency?: CurrencyCode | string;
    /** Revenue value (use for Purchase, AddToCart, etc.) */
    value?: number;
    tax?: number;
    shipping?: number;
    coupon?: string;
    /** Merchant or store name */
    affiliation?: string;
    payment_type?: string;       // "credit_card" | "paypal" | "apple_pay"
    shipping_tier?: string;      // "standard" | "express" | "overnight"
    items?: CommerceItem[];

    // ── Search ────────────────────────────────────────────────
    search_term?: string;
    search_results_count?: number;
    search_filters?: Record<string, string>;

    // ── Authentication ────────────────────────────────────────
    auth_provider?: AuthProvider;
    auth_method?: string;        // "password" | "magic_link" | "sso"
    is_new_user?: boolean;

    // ── Error / System ────────────────────────────────────────
    error_code?: string;
    error_message?: string;
    error_severity?: ErrorSeverity;
    error_stack?: string;
    http_status?: number;
    component_name?: string;     // Which component crashed

    // ── Engagement / Sharing ──────────────────────────────────
    /** "twitter" | "facebook" | "whatsapp" | "copy_link" */
    share_platform?: string;
    /** "native_share" | "button" */
    share_method?: string;

    // ── Timing ────────────────────────────────────────────────
    /** Generic elapsed duration (hover time, read time, etc.) in ms */
    duration_ms?: number;

    // ── Feature Flags / Custom ────────────────────────────────
    feature?: string;
    priority?: Priority;
    /**
     * Escape hatch for product-specific properties.
     * Strongly prefer explicit fields above over this bag.
     */
    custom?: Record<string, unknown>;
}