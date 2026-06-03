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

type CreateEventInput struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
	Event     *TrackingEvent
}

type CreatePageViewInput struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
	URL       string
	Title     string
	Referrer  string
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
	currentSchemaVersion := "v1.0.0"
	eventBuilder.SetSchemaVersion(currentSchemaVersion)
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

	if ev.Target.Custom != nil {
		if err := saveCustomProps(ctx, tx, visitorEvent.ID, ev.Target.Custom); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if ev.Target.SearchFilters != nil {
		if err := saveSearchFilters(ctx, tx, visitorEvent.ID, ev.Target.SearchFilters); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(ev.Target.Items) > 0 {
		if err := saveCommerceItems(ctx, tx, visitorEvent.ID, ev.Target.Items); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return visitorEvent, nil
}

func saveTargetProps(ctx context.Context, tx *ent.Tx, eventID int64, tp TargetProps) error {
	tpBuilder := tx.EventTargetProp.Create().SetEventID(eventID)

	if tp.ElementType != nil {
		tpBuilder.SetElementType(string(*tp.ElementType))
	}
	if tp.ElementID != nil {
		tpBuilder.SetElementID(*tp.ElementID)
	}
	if tp.ElementName != nil {
		tpBuilder.SetElementName(*tp.ElementName)
	}
	if tp.ElementText != nil {
		tpBuilder.SetElementText(*tp.ElementText)
	}
	if tp.ElementClass != nil {
		tpBuilder.SetElementClass(*tp.ElementClass)
	}
	if tp.ElementSelector != nil {
		tpBuilder.SetElementSelector(*tp.ElementSelector)
	}
	if tp.Component != nil {
		tpBuilder.SetComponent(*tp.Component)
	}
	if tp.PageName != nil {
		tpBuilder.SetPageName(*tp.PageName)
	}
	if tp.ButtonName != nil {
		tpBuilder.SetButtonName(*tp.ButtonName)
	}
	if tp.ButtonID != nil {
		tpBuilder.SetButtonID(*tp.ButtonID)
	}
	if tp.Href != nil {
		tpBuilder.SetHref(*tp.Href)
	}
	if tp.FormName != nil {
		tpBuilder.SetFormName(*tp.FormName)
	}
	if tp.FormID != nil {
		tpBuilder.SetFormID(*tp.FormID)
	}
	if tp.FormField != nil {
		tpBuilder.SetFormField(*tp.FormField)
	}
	if tp.FormValueHash != nil {
		tpBuilder.SetFormValueHash(*tp.FormValueHash)
	}
	if tp.FormFieldCount != nil {
		tpBuilder.SetFormFieldCount(*tp.FormFieldCount)
	}
	if tp.FormCompletionPct != nil {
		tpBuilder.SetFormCompletionPct(*tp.FormCompletionPct)
	}
	if tp.MediaType != nil {
		tpBuilder.SetMediaType(string(*tp.MediaType))
	}
	if tp.MediaID != nil {
		tpBuilder.SetMediaID(*tp.MediaID)
	}
	if tp.MediaTitle != nil {
		tpBuilder.SetMediaTitle(*tp.MediaTitle)
	}
	if tp.MediaProvider != nil {
		tpBuilder.SetMediaProvider(*tp.MediaProvider)
	}
	if tp.MediaDurationMs != nil {
		tpBuilder.SetMediaDurationMs(*tp.MediaDurationMs)
	}
	if tp.MediaCurrentTimeMs != nil {
		tpBuilder.SetMediaCurrentTimeMs(*tp.MediaCurrentTimeMs)
	}
	if tp.MediaPercent != nil {
		tpBuilder.SetMediaPercent(*tp.MediaPercent)
	}
	if tp.MediaQuality != nil {
		tpBuilder.SetMediaQuality(*tp.MediaQuality)
	}
	if tp.MediaSoundOn != nil {
		tpBuilder.SetMediaSoundOn(*tp.MediaSoundOn)
	}
	if tp.MediaFullscreen != nil {
		tpBuilder.SetMediaFullscreen(*tp.MediaFullscreen)
	}
	if tp.ScrollPercent != nil {
		tpBuilder.SetScrollPercent(*tp.ScrollPercent)
	}
	if tp.ScrollDirection != nil {
		tpBuilder.SetScrollDirection(string(*tp.ScrollDirection))
	}
	if tp.ScrollDepthPx != nil {
		tpBuilder.SetScrollDepthPx(*tp.ScrollDepthPx)
	}
	if tp.MouseButton != nil {
		tpBuilder.SetMouseButton(string(*tp.MouseButton))
	}
	if tp.ClickX != nil {
		tpBuilder.SetClickX(*tp.ClickX)
	}
	if tp.ClickY != nil {
		tpBuilder.SetClickY(*tp.ClickY)
	}
	if tp.ElementX != nil {
		tpBuilder.SetElementX(*tp.ElementX)
	}
	if tp.ElementY != nil {
		tpBuilder.SetElementY(*tp.ElementY)
	}
	if tp.TransactionID != nil {
		tpBuilder.SetTransactionID(*tp.TransactionID)
	}
	if tp.Currency != nil {
		tpBuilder.SetCurrency(*tp.Currency)
	}
	if tp.Value != nil {
		tpBuilder.SetValue(*tp.Value)
	}
	if tp.Tax != nil {
		tpBuilder.SetTax(*tp.Tax)
	}
	if tp.Shipping != nil {
		tpBuilder.SetShipping(*tp.Shipping)
	}
	if tp.Coupon != nil {
		tpBuilder.SetCoupon(*tp.Coupon)
	}
	if tp.Affiliation != nil {
		tpBuilder.SetAffiliation(*tp.Affiliation)
	}
	if tp.PaymentType != nil {
		tpBuilder.SetPaymentType(*tp.PaymentType)
	}
	if tp.ShippingTier != nil {
		tpBuilder.SetShippingTier(*tp.ShippingTier)
	}
	if tp.SearchTerm != nil {
		tpBuilder.SetSearchTerm(*tp.SearchTerm)
	}
	if tp.SearchResultsCount != nil {
		tpBuilder.SetSearchResultsCount(*tp.SearchResultsCount)
	}
	if tp.AuthProvider != nil {
		tpBuilder.SetAuthProvider(string(*tp.AuthProvider))
	}
	if tp.AuthMethod != nil {
		tpBuilder.SetAuthMethod(*tp.AuthMethod)
	}
	if tp.IsNewUser != nil {
		tpBuilder.SetIsNewUser(*tp.IsNewUser)
	}
	if tp.ErrorCode != nil {
		tpBuilder.SetErrorCode(*tp.ErrorCode)
	}
	if tp.ErrorMessage != nil {
		tpBuilder.SetErrorMessage(*tp.ErrorMessage)
	}
	if tp.ErrorSeverity != nil {
		tpBuilder.SetErrorSeverity(string(*tp.ErrorSeverity))
	}
	if tp.ErrorStack != nil {
		tpBuilder.SetErrorStack(*tp.ErrorStack)
	}
	if tp.HttpStatus != nil {
		tpBuilder.SetHTTPStatus(*tp.HttpStatus)
	}
	if tp.ComponentName != nil {
		tpBuilder.SetComponentName(*tp.ComponentName)
	}
	if tp.SharePlatform != nil {
		tpBuilder.SetSharePlatform(*tp.SharePlatform)
	}
	if tp.ShareMethod != nil {
		tpBuilder.SetShareMethod(*tp.ShareMethod)
	}
	if tp.DurationMs != nil {
		tpBuilder.SetDurationMs(*tp.DurationMs)
	}
	if tp.Feature != nil {
		tpBuilder.SetFeature(*tp.Feature)
	}
	if tp.Priority != nil {
		tpBuilder.SetPriority(string(*tp.Priority))
	}

	_, err := tpBuilder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to save target props: %w", err)
	}
	return nil
}

func savePlatformClicks(ctx context.Context, tx *ent.Tx, eventID int64, pc PlatformClickIds) error {
	pcBuilder := tx.VisitorPlatformClick.Create().SetEventID(eventID)

	if pc.Gclid != nil {
		pcBuilder.SetGclid(*pc.Gclid)
	}
	if pc.Gbraid != nil {
		pcBuilder.SetGbraid(*pc.Gbraid)
	}
	if pc.Wbraid != nil {
		pcBuilder.SetWbraid(*pc.Wbraid)
	}
	if pc.Dclid != nil {
		pcBuilder.SetDclid(*pc.Dclid)
	}
	if pc.Fbclid != nil {
		pcBuilder.SetFbclid(*pc.Fbclid)
	}
	if pc.Fbp != nil {
		pcBuilder.SetFbp(*pc.Fbp)
	}
	if pc.Fbc != nil {
		pcBuilder.SetFbc(*pc.Fbc)
	}
	if pc.Ttclid != nil {
		pcBuilder.SetTtclid(*pc.Ttclid)
	}
	if pc.Ttp != nil {
		pcBuilder.SetTtp(*pc.Ttp)
	}
	if pc.Msclkid != nil {
		pcBuilder.SetMsclkid(*pc.Msclkid)
	}
	if pc.Twclid != nil {
		pcBuilder.SetTwclid(*pc.Twclid)
	}
	if pc.LiFatId != nil {
		pcBuilder.SetLiFatID(*pc.LiFatId)
	}
	if pc.ScCid != nil {
		pcBuilder.SetScCid(*pc.ScCid)
	}
	if pc.Epik != nil {
		pcBuilder.SetEpik(*pc.Epik)
	}
	if pc.LiverampAtsId != nil {
		pcBuilder.SetLiverampAtsID(*pc.LiverampAtsId)
	}

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

func saveCommerceItems(ctx context.Context, tx *ent.Tx, eventID int64, items []CommerceItem) error {
	for _, item := range items {
		itemBuilder := tx.EventCommerceItem.Create().
			SetEventID(eventID).
			SetItemID(item.ItemID)

		if item.ItemName != nil {
			itemBuilder.SetItemName(*item.ItemName)
		}
		if item.ItemBrand != nil {
			itemBuilder.SetItemBrand(*item.ItemBrand)
		}
		if item.ItemCategory != nil {
			itemBuilder.SetItemCategory(*item.ItemCategory)
		}
		if item.ItemCategory2 != nil {
			itemBuilder.SetItemCategory2(*item.ItemCategory2)
		}
		if item.ItemCategory3 != nil {
			itemBuilder.SetItemCategory3(*item.ItemCategory3)
		}
		if item.ItemVariant != nil {
			itemBuilder.SetItemVariant(*item.ItemVariant)
		}
		if item.Price != nil {
			itemBuilder.SetPrice(*item.Price)
		}
		if item.Currency != nil {
			itemBuilder.SetCurrency(*item.Currency)
		}
		if item.Quantity != nil {
			itemBuilder.SetQuantity(*item.Quantity)
		}
		if item.Coupon != nil {
			itemBuilder.SetCoupon(*item.Coupon)
		}
		if item.Discount != nil {
			itemBuilder.SetDiscount(*item.Discount)
		}
		if item.Index != nil {
			itemBuilder.SetIndex(*item.Index)
		}
		if item.ItemListID != nil {
			itemBuilder.SetItemListID(*item.ItemListID)
		}
		if item.ItemListName != nil {
			itemBuilder.SetItemListName(*item.ItemListName)
		}
		if item.Affiliation != nil {
			itemBuilder.SetAffiliation(*item.Affiliation)
		}

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
