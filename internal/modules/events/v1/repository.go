package v1

import (
	"context"
	"fmt"
	"time"

	"tracking/internal/ent"

	"github.com/google/uuid"
)

type Repository interface {
	CreateEvent(ctx context.Context, input CreateEventInput) (*ent.VisitorEvent, error)
	CreatePageView(ctx context.Context, input CreatePageViewInput) (*ent.PageView, error)
}

type entRepository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	if client == nil {
		panic("ent client cannot be nil")
	}
	return &entRepository{client: client}
}


func (r *entRepository) CreateEvent(ctx context.Context, input CreateEventInput) (*ent.VisitorEvent, error) {

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	ev := input.Event
	occurredAt, _ := time.Parse(time.RFC3339, ev.OccurredAt)

	eventBuilder := tx.VisitorEvent.Create().
		SetUUID(uuid.NewString()).
		SetSessionID(input.SessionID).
		SetVisitorID(input.VisitorID).
		SetDeviceID(input.DeviceID).
		SetPageID(ev.PageID).
		SetEventType(string(ev.Type)).
		SetEventCategory(string(ev.Category)).
		SetEventName(ev.Name).
		SetOccurredAt(occurredAt).
		SetSequence(ev.Sequence)

	now := time.Now()
	eventBuilder.SetServerReceivedAt(now)
	eventBuilder.SetSchemaVersion(Version)
	if ev.SdkName != nil {
		eventBuilder.SetSdkName(*ev.SdkName)
	}
	if ev.SdkVersion != nil {
		eventBuilder.SetSdkVersion(*ev.SdkVersion)
	}

	visitorEvent, err := eventBuilder.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create visitor event: %w", err)
	}

	if err := saveTargetProps(ctx, tx, visitorEvent.ID, ev.Target); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := savePlatformClicks(ctx, tx, visitorEvent.ID, ev.ClickIds); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := saveCustomProps(ctx, tx, visitorEvent.ID, ev.Target.Custom); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := saveSearchFilters(ctx, tx, visitorEvent.ID, ev.Target.SearchFilters); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := saveCommerceItems(ctx, tx, visitorEvent.ID, ev.Target.Items); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return visitorEvent, nil
}

func saveTargetProps(ctx context.Context, tx *ent.Tx, eventID int64, tp TargetProps) error {
	tpBuilder := tx.EventTargetProp.Create().SetEventID(eventID)

	setElementProps(tpBuilder, tp)
	setFormProps(tpBuilder, tp)
	setMediaProps(tpBuilder, tp)
	setScrollAndClickProps(tpBuilder, tp)
	setCommerceProps(tpBuilder, tp)
	setAuthProps(tpBuilder, tp)
	setErrorProps(tpBuilder, tp)
	setMiscComponentProps(tpBuilder, tp)
	setMiscActionProps(tpBuilder, tp)

	_, err := tpBuilder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to save target props: %w", err)
	}
	return nil
}

func setElementProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.ElementType != nil {
		b.SetElementType(string(*tp.ElementType))
	}
	if tp.ElementID != nil {
		b.SetElementID(*tp.ElementID)
	}
	if tp.ElementName != nil {
		b.SetElementName(*tp.ElementName)
	}
	if tp.ElementText != nil {
		b.SetElementText(*tp.ElementText)
	}
	if tp.ElementClass != nil {
		b.SetElementClass(*tp.ElementClass)
	}
	if tp.ElementSelector != nil {
		b.SetElementSelector(*tp.ElementSelector)
	}
	if tp.ElementX != nil {
		b.SetElementX(*tp.ElementX)
	}
	if tp.ElementY != nil {
		b.SetElementY(*tp.ElementY)
	}
}

func setFormProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.FormName != nil {
		b.SetFormName(*tp.FormName)
	}
	if tp.FormID != nil {
		b.SetFormID(*tp.FormID)
	}
	if tp.FormField != nil {
		b.SetFormField(*tp.FormField)
	}
	if tp.FormValueHash != nil {
		b.SetFormValueHash(*tp.FormValueHash)
	}
	if tp.FormFieldCount != nil {
		b.SetFormFieldCount(*tp.FormFieldCount)
	}
	if tp.FormCompletionPct != nil {
		b.SetFormCompletionPct(*tp.FormCompletionPct)
	}
}

func setMediaProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.MediaType != nil {
		b.SetMediaType(string(*tp.MediaType))
	}
	if tp.MediaID != nil {
		b.SetMediaID(*tp.MediaID)
	}
	if tp.MediaTitle != nil {
		b.SetMediaTitle(*tp.MediaTitle)
	}
	if tp.MediaProvider != nil {
		b.SetMediaProvider(*tp.MediaProvider)
	}
	if tp.MediaDurationMs != nil {
		b.SetMediaDurationMs(*tp.MediaDurationMs)
	}
	if tp.MediaCurrentTimeMs != nil {
		b.SetMediaCurrentTimeMs(*tp.MediaCurrentTimeMs)
	}
	if tp.MediaPercent != nil {
		b.SetMediaPercent(*tp.MediaPercent)
	}
	if tp.MediaQuality != nil {
		b.SetMediaQuality(*tp.MediaQuality)
	}
	if tp.MediaSoundOn != nil {
		b.SetMediaSoundOn(*tp.MediaSoundOn)
	}
	if tp.MediaFullscreen != nil {
		b.SetMediaFullscreen(*tp.MediaFullscreen)
	}
}

func setScrollAndClickProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.ScrollPercent != nil {
		b.SetScrollPercent(*tp.ScrollPercent)
	}
	if tp.ScrollDirection != nil {
		b.SetScrollDirection(string(*tp.ScrollDirection))
	}
	if tp.ScrollDepthPx != nil {
		b.SetScrollDepthPx(*tp.ScrollDepthPx)
	}
	if tp.MouseButton != nil {
		b.SetMouseButton(string(*tp.MouseButton))
	}
	if tp.ClickX != nil {
		b.SetClickX(*tp.ClickX)
	}
	if tp.ClickY != nil {
		b.SetClickY(*tp.ClickY)
	}
}

func setCommerceProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.TransactionID != nil {
		b.SetTransactionID(*tp.TransactionID)
	}
	if tp.Currency != nil {
		b.SetCurrency(*tp.Currency)
	}
	if tp.Value != nil {
		b.SetValue(*tp.Value)
	}
	if tp.Tax != nil {
		b.SetTax(*tp.Tax)
	}
	if tp.Shipping != nil {
		b.SetShipping(*tp.Shipping)
	}
	if tp.Coupon != nil {
		b.SetCoupon(*tp.Coupon)
	}
	if tp.Affiliation != nil {
		b.SetAffiliation(*tp.Affiliation)
	}
	if tp.PaymentType != nil {
		b.SetPaymentType(*tp.PaymentType)
	}
	if tp.ShippingTier != nil {
		b.SetShippingTier(*tp.ShippingTier)
	}
}

func setAuthProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.AuthProvider != nil {
		b.SetAuthProvider(string(*tp.AuthProvider))
	}
	if tp.AuthMethod != nil {
		b.SetAuthMethod(*tp.AuthMethod)
	}
	if tp.IsNewUser != nil {
		b.SetIsNewUser(*tp.IsNewUser)
	}
}

func setErrorProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.ErrorCode != nil {
		b.SetErrorCode(*tp.ErrorCode)
	}
	if tp.ErrorMessage != nil {
		b.SetErrorMessage(*tp.ErrorMessage)
	}
	if tp.ErrorSeverity != nil {
		b.SetErrorSeverity(string(*tp.ErrorSeverity))
	}
	if tp.ErrorStack != nil {
		b.SetErrorStack(*tp.ErrorStack)
	}
	if tp.HttpStatus != nil {
		b.SetHTTPStatus(*tp.HttpStatus)
	}
}

func setMiscComponentProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.Component != nil {
		b.SetComponent(*tp.Component)
	}
	if tp.PageName != nil {
		b.SetPageName(*tp.PageName)
	}
	if tp.ButtonName != nil {
		b.SetButtonName(*tp.ButtonName)
	}
	if tp.ButtonID != nil {
		b.SetButtonID(*tp.ButtonID)
	}
	if tp.Href != nil {
		b.SetHref(*tp.Href)
	}
	if tp.ComponentName != nil {
		b.SetComponentName(*tp.ComponentName)
	}
}

func setMiscActionProps(b *ent.EventTargetPropCreate, tp TargetProps) {
	if tp.SharePlatform != nil {
		b.SetSharePlatform(*tp.SharePlatform)
	}
	if tp.ShareMethod != nil {
		b.SetShareMethod(*tp.ShareMethod)
	}
	if tp.DurationMs != nil {
		b.SetDurationMs(*tp.DurationMs)
	}
	if tp.Feature != nil {
		b.SetFeature(*tp.Feature)
	}
	if tp.Priority != nil {
		b.SetPriority(string(*tp.Priority))
	}
	if tp.SearchTerm != nil {
		b.SetSearchTerm(*tp.SearchTerm)
	}
	if tp.SearchResultsCount != nil {
		b.SetSearchResultsCount(*tp.SearchResultsCount)
	}
}

func setPrimaryClickIDs(b *ent.VisitorPlatformClickCreate, pc PlatformClickIds) {
	if pc.Gclid != nil {
		b.SetGclid(*pc.Gclid)
	}
	if pc.Gbraid != nil {
		b.SetGbraid(*pc.Gbraid)
	}
	if pc.Wbraid != nil {
		b.SetWbraid(*pc.Wbraid)
	}
	if pc.Dclid != nil {
		b.SetDclid(*pc.Dclid)
	}
	if pc.Fbclid != nil {
		b.SetFbclid(*pc.Fbclid)
	}
	if pc.Fbp != nil {
		b.SetFbp(*pc.Fbp)
	}
	if pc.Fbc != nil {
		b.SetFbc(*pc.Fbc)
	}
}

func setSecondaryClickIDs(b *ent.VisitorPlatformClickCreate, pc PlatformClickIds) {
	if pc.Ttclid != nil {
		b.SetTtclid(*pc.Ttclid)
	}
	if pc.Ttp != nil {
		b.SetTtp(*pc.Ttp)
	}
	if pc.Msclkid != nil {
		b.SetMsclkid(*pc.Msclkid)
	}
	if pc.Twclid != nil {
		b.SetTwclid(*pc.Twclid)
	}
	if pc.LiFatId != nil {
		b.SetLiFatID(*pc.LiFatId)
	}
	if pc.ScCid != nil {
		b.SetScCid(*pc.ScCid)
	}
	if pc.Epik != nil {
		b.SetEpik(*pc.Epik)
	}
	if pc.LiverampAtsId != nil {
		b.SetLiverampAtsID(*pc.LiverampAtsId)
	}
}

func savePlatformClicks(ctx context.Context, tx *ent.Tx, eventID int64, pc PlatformClickIds) error {
	pcBuilder := tx.VisitorPlatformClick.Create().SetEventID(eventID)

	setPrimaryClickIDs(pcBuilder, pc)
	setSecondaryClickIDs(pcBuilder, pc)

	_, err := pcBuilder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to save platform clicks: %w", err)
	}
	return nil
}

func saveCustomProps(ctx context.Context, tx *ent.Tx, eventID int64, customProps map[string]interface{}) error {
	for k, v := range customProps {
		valStr := fmt.Sprintf("%v", v)
		_, err := tx.EventCustomProp.Create().
			SetEventID(eventID).
			SetKey(k).
			SetValue(valStr).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to save custom prop %s: %w", k, err)
		}
	}
	return nil
}

func saveSearchFilters(ctx context.Context, tx *ent.Tx, eventID int64, filters map[string]string) error {
	for k, v := range filters {
		_, err := tx.EventTargetSearchFilter.Create().
			SetEventID(eventID).
			SetKey(k).
			SetValue(v).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to save search filter %s: %w", k, err)
		}
	}
	return nil
}

func setCommerceItemDetails(b *ent.EventCommerceItemCreate, item CommerceItem) {
	if item.ItemName != nil {
		b.SetItemName(*item.ItemName)
	}
	if item.ItemBrand != nil {
		b.SetItemBrand(*item.ItemBrand)
	}
	if item.ItemCategory != nil {
		b.SetItemCategory(*item.ItemCategory)
	}
	if item.ItemCategory2 != nil {
		b.SetItemCategory2(*item.ItemCategory2)
	}
	if item.ItemCategory3 != nil {
		b.SetItemCategory3(*item.ItemCategory3)
	}
	if item.ItemVariant != nil {
		b.SetItemVariant(*item.ItemVariant)
	}
	if item.Index != nil {
		b.SetIndex(*item.Index)
	}
	if item.ItemListID != nil {
		b.SetItemListID(*item.ItemListID)
	}
	if item.ItemListName != nil {
		b.SetItemListName(*item.ItemListName)
	}
	if item.Affiliation != nil {
		b.SetAffiliation(*item.Affiliation)
	}
}

func setCommerceItemPricing(b *ent.EventCommerceItemCreate, item CommerceItem) {
	if item.Price != nil {
		b.SetPrice(*item.Price)
	}
	if item.Currency != nil {
		b.SetCurrency(*item.Currency)
	}
	if item.Quantity != nil {
		b.SetQuantity(*item.Quantity)
	}
	if item.Coupon != nil {
		b.SetCoupon(*item.Coupon)
	}
	if item.Discount != nil {
		b.SetDiscount(*item.Discount)
	}
}

func saveCommerceItems(ctx context.Context, tx *ent.Tx, eventID int64, items []CommerceItem) error {
	for _, item := range items {
		itemBuilder := tx.EventCommerceItem.Create().
			SetEventID(eventID).
			SetItemID(item.ItemID)

		setCommerceItemDetails(itemBuilder, item)
		setCommerceItemPricing(itemBuilder, item)

		_, err := itemBuilder.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to save commerce item %s: %w", item.ItemID, err)
		}
	}
	return nil
}

func (r *entRepository) CreatePageView(ctx context.Context, input CreatePageViewInput) (*ent.PageView, error) {
	pv, err := r.client.PageView.
		Create().
		SetUUID(uuid.New().String()).
		SetSessionID(input.SessionID).
		SetVisitorID(input.VisitorID).
		SetDeviceID(input.DeviceID).
		SetURL(input.URL).
		SetTitle(input.Title).
		SetReferrer(input.Referrer).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create page view: %w", err)
	}
	return pv, nil
}
