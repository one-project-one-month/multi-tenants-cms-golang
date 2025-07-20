package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

type (
	// StripeHandler handles Stripe payment operations
	StripeHandler struct {
		secretKey        string
		webhookSecret    string
		allowedCountries []string
		defaultCurrency  string
	}

	PaymentRequest struct {
		Amount      int64             `json:"amount"`
		Currency    string            `json:"currency"`
		Description string            `json:"description"`
		CustomerID  string            `json:"customer_id"`
		Metadata    map[string]string `json:"metadata"`
	}

	CheckoutSessionRequest struct {
		Amount          int64             `json:"amount"`
		Currency        string            `json:"currency"`
		ProductName     string            `json:"product_name"`
		Quantity        int64             `json:"quantity"`
		SuccessURL      string            `json:"success_url"`
		CancelURL       string            `json:"cancel_url"`
		CustomerEmail   string            `json:"customer_email"`
		PaymentMethod   string            `json:"payment_method"`
		BillingAddress  bool              `json:"billing_address"`
		ShippingAddress bool              `json:"shipping_address"`
		Metadata        map[string]string `json:"metadata"`
	}

	PaymentResponse struct {
		Success       bool              `json:"success"`
		PaymentIntent string            `json:"payment_intent_id,omitempty"`
		ClientSecret  string            `json:"client_secret,omitempty"`
		SessionID     string            `json:"session_id,omitempty"`
		SessionURL    string            `json:"session_url,omitempty"`
		Error         string            `json:"error,omitempty"`
		Metadata      map[string]string `json:"metadata,omitempty"`
	}

	PaymentStatusResponse struct {
		ID       string            `json:"id"`
		Status   string            `json:"status"`
		Amount   int64             `json:"amount"`
		Currency string            `json:"currency"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}
)

// NewStripeHandler creates a new StripeHandler with the provided configuration
func NewStripeHandler(secretKey, webhookSecret string, allowedCountries []string, defaultCurrency string) *StripeHandler {
	if defaultCurrency == "" {
		defaultCurrency = "usd"
	}

	return &StripeHandler{
		secretKey:        secretKey,
		webhookSecret:    webhookSecret,
		allowedCountries: allowedCountries,
		defaultCurrency:  defaultCurrency,
	}
}

func (h *StripeHandler) createPaymentIntentHandler(w http.ResponseWriter, r *http.Request) {
	stripe.Key = h.secretKey

	var payReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&payReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payReq.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if payReq.Currency == "" {
		payReq.Currency = h.defaultCurrency
	}

	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(payReq.Amount),
		Currency:    stripe.String(payReq.Currency),
		Description: stripe.String(payReq.Description),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	if payReq.CustomerID != "" {
		params.Customer = stripe.String(payReq.CustomerID)
	}

	if payReq.Metadata != nil {
		params.Metadata = payReq.Metadata
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Payment intent creation failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := PaymentResponse{
		Success:       true,
		PaymentIntent: pi.ID,
		ClientSecret:  pi.ClientSecret,
		Metadata:      pi.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *StripeHandler) createCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	stripe.Key = h.secretKey

	var sessionReq CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&sessionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if sessionReq.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if sessionReq.Currency == "" {
		sessionReq.Currency = h.defaultCurrency
	}
	if sessionReq.SuccessURL == "" || sessionReq.CancelURL == "" {
		http.Error(w, "SuccessURL and CancelURL are required", http.StatusBadRequest)
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(sessionReq.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(sessionReq.ProductName),
					},
					UnitAmount: stripe.Int64(sessionReq.Amount),
				},
				Quantity: stripe.Int64(sessionReq.Quantity),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(sessionReq.SuccessURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(sessionReq.CancelURL),
	}

	if sessionReq.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(sessionReq.CustomerEmail)
	}

	if sessionReq.BillingAddress {
		params.BillingAddressCollection = stripe.String("required")
	}

	if sessionReq.ShippingAddress {
		if len(h.allowedCountries) == 0 {
			h.allowedCountries = []string{"US", "GB", "CA"}
		}
		params.ShippingAddressCollection = &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice(h.allowedCountries),
		}
	}

	if sessionReq.Metadata != nil {
		params.Metadata = sessionReq.Metadata
	}

	s, err := session.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Checkout session creation failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := PaymentResponse{
		Success:    true,
		SessionID:  s.ID,
		SessionURL: s.URL,
		Metadata:   s.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *StripeHandler) webhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	if h.webhookSecret == "" {
		http.Error(w, "Webhook secret not configured", http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Webhook signature verification failed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			http.Error(w, "Failed to parse payment intent", http.StatusBadRequest)
			return
		}
		h.handlePaymentSuccess(&pi)

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			http.Error(w, "Failed to parse payment intent", http.StatusBadRequest)
			return
		}
		h.handlePaymentFailure(&pi)

	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			http.Error(w, "Failed to parse checkout session", http.StatusBadRequest)
			return
		}
		h.handleCheckoutSuccess(&session)

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *StripeHandler) handlePaymentSuccess(pi *stripe.PaymentIntent) {
	fmt.Printf("Payment succeeded: %s\n", pi.ID)
	fmt.Printf("Amount: %d %s\n", pi.Amount, pi.Currency)
	fmt.Printf("Customer: %s\n", pi.Customer.ID)
}

func (h *StripeHandler) handlePaymentFailure(pi *stripe.PaymentIntent) {
	fmt.Printf("Payment failed: %s\n", pi.ID)
	fmt.Printf("Failure reason: %s\n", pi.LastPaymentError.Msg)
}

func (h *StripeHandler) handleCheckoutSuccess(session *stripe.CheckoutSession) {
	fmt.Printf("Checkout session completed: %s\n", session.ID)
	fmt.Printf("Amount: %d\n", session.AmountTotal)
	fmt.Printf("Customer: %s\n", session.Customer.ID)
}

func (h *StripeHandler) getPaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	stripe.Key = h.secretKey

	paymentIntentID := r.URL.Query().Get("payment_intent_id")
	if paymentIntentID == "" {
		http.Error(w, "Payment intent ID is required", http.StatusBadRequest)
		return
	}

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve payment: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := PaymentStatusResponse{
		ID:       pi.ID,
		Status:   string(pi.Status),
		Amount:   pi.Amount,
		Currency: string(pi.Currency),
		Metadata: pi.Metadata,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *StripeHandler) RegisterHandlers(mux *runtime.ServeMux) {
	mux.HandlePath("POST", "/lms/v1/stripe/create-payment-intent", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h.createPaymentIntentHandler(w, r)
	})
	mux.HandlePath("GET", "/lms/v1/stripe/payment-status", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h.getPaymentStatusHandler(w, r)
	})
	mux.HandlePath("POST", "/lms/v1/stripe/create-checkout-session", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h.createCheckoutSessionHandler(w, r)
	})
	mux.HandlePath("POST", "/lms/v1/stripe/webhook", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h.webhookHandler(w, r)
	})
}
