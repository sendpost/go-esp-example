package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	sendpost "github.com/sendpost/sendpost-go-sdk"
)

// ESPExample demonstrates a complete workflow that an ESP would typically follow
type ESPExample struct {
	client               *sendpost.APIClient
	accountAPIKey        string
	subAccountAPIKey     string
	createdSubAccountID  *int32
	createdSubAccountKey string
	createdWebhookID     *int64
	createdDomainID      string
	createdIPPoolID      *int64
	createdIPPoolName    string
	sentMessageID        string
}

// Configuration constants - Update these with your values
const (
	basePath       = "https://api.sendpost.io/api/v1"
	testFromEmail  = "sender@yourdomain.com"
	testToEmail    = "recipient@example.com"
	testDomainName = "yourdomain.com"
	webhookURL     = "https://your-webhook-endpoint.com/webhook"
)

// NewESPExample creates a new ESPExample instance
func NewESPExample() *ESPExample {
	// Get API keys from environment variables or use defaults
	accountAPIKey := os.Getenv("SENDPOST_ACCOUNT_API_KEY")
	if accountAPIKey == "" {
		accountAPIKey = "YOUR_ACCOUNT_API_KEY_HERE"
	}

	subAccountAPIKey := os.Getenv("SENDPOST_SUB_ACCOUNT_API_KEY")
	if subAccountAPIKey == "" {
		subAccountAPIKey = "YOUR_SUB_ACCOUNT_API_KEY_HERE"
	}

	// Create configuration
	cfg := sendpost.NewConfiguration()
	cfg.Servers = sendpost.ServerConfigurations{
		sendpost.ServerConfiguration{
			URL: basePath,
		},
	}

	// Create API client
	client := sendpost.NewAPIClient(cfg)

	return &ESPExample{
		client:           client,
		accountAPIKey:    accountAPIKey,
		subAccountAPIKey: subAccountAPIKey,
	}
}

// createAccountAuthContext creates a context with account API key authentication
func (e *ESPExample) createAccountAuthContext() context.Context {
	return context.WithValue(
		context.Background(),
		sendpost.ContextAPIKeys,
		map[string]sendpost.APIKey{
			"accountAuth": {Key: e.accountAPIKey},
		},
	)
}

// createSubAccountAuthContext creates a context with sub-account API key authentication
func (e *ESPExample) createSubAccountAuthContext() context.Context {
	return context.WithValue(
		context.Background(),
		sendpost.ContextAPIKeys,
		map[string]sendpost.APIKey{
			"subAccountAuth": {Key: e.subAccountAPIKey},
		},
	)
}

// ListSubAccounts lists all sub-accounts
func (e *ESPExample) ListSubAccounts() {
	fmt.Println("\n=== Step 1: Listing All Sub-Accounts ===")

	ctx := e.createAccountAuthContext()
	subAccountAPI := e.client.SubAccountAPI

	fmt.Println("Retrieving all sub-accounts...")
	subAccounts, resp, err := subAccountAPI.GetAllSubAccounts(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to list sub-accounts:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Retrieved %d sub-account(s)\n", len(subAccounts))
	for _, subAccount := range subAccounts {
		fmt.Printf("  - ID: %d\n", *subAccount.Id)
		if subAccount.Name != nil {
			fmt.Printf("    Name: %s\n", *subAccount.Name)
		}
		if subAccount.ApiKey != nil {
			fmt.Printf("    API Key: %s\n", *subAccount.ApiKey)
		}
		if subAccount.Type != nil {
			accountType := "Regular"
			if *subAccount.Type == 1 {
				accountType = "Plus"
			}
			fmt.Printf("    Type: %s\n", accountType)
		}
		if subAccount.Blocked != nil {
			blocked := "No"
			if *subAccount.Blocked {
				blocked = "Yes"
			}
			fmt.Printf("    Blocked: %s\n", blocked)
		}
		if subAccount.Created != nil {
			fmt.Printf("    Created: %d\n", *subAccount.Created)
		}
		fmt.Println()

		// Use first sub-account if none selected
		if e.createdSubAccountID == nil && subAccount.Id != nil {
			e.createdSubAccountID = subAccount.Id
			if subAccount.ApiKey != nil {
				e.createdSubAccountKey = *subAccount.ApiKey
			}
		}
	}
}

// CreateSubAccount creates a new sub-account
func (e *ESPExample) CreateSubAccount() {
	fmt.Println("\n=== Step 2: Creating Sub-Account ===")

	ctx := e.createAccountAuthContext()
	subAccountAPI := e.client.SubAccountAPI

	// Create new sub-account request
	name := fmt.Sprintf("ESP Client - %d", time.Now().Unix())
	createSubAccountRequest := sendpost.NewCreateSubAccountRequest()
	createSubAccountRequest.SetName(name)

	fmt.Printf("Creating sub-account: %s\n", name)

	subAccount, resp, err := subAccountAPI.CreateSubAccount(ctx).CreateSubAccountRequest(*createSubAccountRequest).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to create sub-account:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if subAccount.Id != nil {
		e.createdSubAccountID = subAccount.Id
	}
	if subAccount.ApiKey != nil {
		e.createdSubAccountKey = *subAccount.ApiKey
	}

	fmt.Println("✓ Sub-account created successfully!")
	if subAccount.Id != nil {
		fmt.Printf("  ID: %d\n", *subAccount.Id)
	}
	if subAccount.Name != nil {
		fmt.Printf("  Name: %s\n", *subAccount.Name)
	}
	if subAccount.ApiKey != nil {
		fmt.Printf("  API Key: %s\n", *subAccount.ApiKey)
	}
	if subAccount.Type != nil {
		accountType := "Regular"
		if *subAccount.Type == 1 {
			accountType = "Plus"
		}
		fmt.Printf("  Type: %s\n", accountType)
	}
}

// CreateWebhook creates a webhook for event notifications
func (e *ESPExample) CreateWebhook() {
	fmt.Println("\n=== Step 3: Creating Webhook ===")

	ctx := e.createAccountAuthContext()
	webhookAPI := e.client.WebhookAPI

	// Create new webhook
	enabled := true
	createWebhookRequest := sendpost.NewCreateWebhookRequest()
	createWebhookRequest.SetUrl(webhookURL)
	createWebhookRequest.SetEnabled(enabled)
	createWebhookRequest.SetProcessed(enabled)
	createWebhookRequest.SetDelivered(enabled)
	createWebhookRequest.SetDropped(enabled)
	createWebhookRequest.SetSoftBounced(enabled)
	createWebhookRequest.SetHardBounced(enabled)
	createWebhookRequest.SetOpened(enabled)
	createWebhookRequest.SetClicked(enabled)
	createWebhookRequest.SetUnsubscribed(enabled)
	createWebhookRequest.SetSpam(enabled)

	fmt.Println("Creating webhook...")
	fmt.Printf("  URL: %s\n", webhookURL)

	webhook, resp, err := webhookAPI.CreateWebhook(ctx).CreateWebhookRequest(*createWebhookRequest).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to create webhook:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if webhook.Id != nil {
		id := int64(*webhook.Id)
		e.createdWebhookID = &id
	}

	fmt.Println("✓ Webhook created successfully!")
	if webhook.Id != nil {
		fmt.Printf("  ID: %d\n", *webhook.Id)
	}
	if webhook.Url != nil {
		fmt.Printf("  URL: %s\n", *webhook.Url)
	}
	if webhook.Enabled != nil {
		fmt.Printf("  Enabled: %v\n", *webhook.Enabled)
	}
}

// ListWebhooks lists all webhooks
func (e *ESPExample) ListWebhooks() {
	fmt.Println("\n=== Step 4: Listing All Webhooks ===")

	ctx := e.createAccountAuthContext()
	webhookAPI := e.client.WebhookAPI

	fmt.Println("Retrieving all webhooks...")
	webhooks, resp, err := webhookAPI.GetAllWebhooks(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to list webhooks:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Retrieved %d webhook(s)\n", len(webhooks))
	for _, webhook := range webhooks {
		if webhook.Id != nil {
			fmt.Printf("  - ID: %d\n", *webhook.Id)
		}
		if webhook.Url != nil {
			fmt.Printf("    URL: %s\n", *webhook.Url)
		}
		if webhook.Enabled != nil {
			fmt.Printf("    Enabled: %v\n", *webhook.Enabled)
		}
		fmt.Println()
	}
}

// AddDomain adds a sending domain
func (e *ESPExample) AddDomain() {
	fmt.Println("\n=== Step 5: Adding Domain ===")

	ctx := e.createSubAccountAuthContext()
	domainAPI := e.client.DomainAPI

	// Create domain request
	domainRequest := sendpost.NewCreateDomainRequest()
	domainRequest.SetName(testDomainName)

	fmt.Printf("Adding domain: %s\n", testDomainName)

	domain, resp, err := domainAPI.SubaccountDomainPost(ctx).CreateDomainRequest(*domainRequest).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to add domain:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if domain.Id != nil {
		e.createdDomainID = strconv.FormatInt(int64(*domain.Id), 10)
	}

	fmt.Println("✓ Domain added successfully!")
	if domain.Id != nil {
		fmt.Printf("  ID: %d\n", *domain.Id)
	}
	if domain.Name != nil {
		fmt.Printf("  Domain: %s\n", *domain.Name)
	}
	if domain.Verified != nil {
		verified := "No"
		if *domain.Verified {
			verified = "Yes"
		}
		fmt.Printf("  Verified: %s\n", verified)
	}

	if domain.Dkim != nil && domain.Dkim.TextValue != nil {
		fmt.Printf("  DKIM Record: %s\n", *domain.Dkim.TextValue)
	}

	fmt.Println("\n⚠️  IMPORTANT: Add the DNS records shown above to your domain's DNS settings to verify the domain.")
}

// ListDomains lists all domains
func (e *ESPExample) ListDomains() {
	fmt.Println("\n=== Step 6: Listing All Domains ===")

	ctx := e.createSubAccountAuthContext()
	domainAPI := e.client.DomainAPI

	fmt.Println("Retrieving all domains...")
	domains, resp, err := domainAPI.GetAllDomains(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to list domains:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Retrieved %d domain(s)\n", len(domains))
	for _, domain := range domains {
		if domain.Id != nil {
			fmt.Printf("  - ID: %d\n", *domain.Id)
		}
		if domain.Name != nil {
			fmt.Printf("    Domain: %s\n", *domain.Name)
		}
		if domain.Verified != nil {
			verified := "No"
			if *domain.Verified {
				verified = "Yes"
			}
			fmt.Printf("    Verified: %s\n", verified)
		}
		fmt.Println()
	}
}

// SendTransactionalEmail sends a transactional email
func (e *ESPExample) SendTransactionalEmail() {
	fmt.Println("\n=== Step 7: Sending Transactional Email ===")

	ctx := e.createSubAccountAuthContext()
	emailAPI := e.client.EmailAPI

	// Create email message
	from := sendpost.NewEmailAddress()
	from.SetEmail(testFromEmail)
	from.SetName("Your Company")

	to := sendpost.NewRecipient()
	to.SetEmail(testToEmail)
	to.SetName("Customer")

	// Add custom fields
	customFields := map[string]interface{}{
		"customer_id": "67890",
		"order_value": "99.99",
	}
	to.SetCustomFields(customFields)

	toList := []sendpost.Recipient{*to}

	emailMessage := sendpost.NewEmailMessageObject()
	emailMessage.SetFrom(*from)
	emailMessage.SetTo(toList)
	emailMessage.SetSubject("Order Confirmation - Transactional Email")
	emailMessage.SetHtmlBody("<h1>Thank you for your order!</h1><p>Your order has been confirmed and will be processed shortly.</p>")
	emailMessage.SetTextBody("Thank you for your order! Your order has been confirmed and will be processed shortly.")

	// Enable tracking
	trackOpens := true
	trackClicks := true
	emailMessage.SetTrackOpens(trackOpens)
	emailMessage.SetTrackClicks(trackClicks)

	// Add custom headers
	headers := map[string]string{
		"X-Order-ID":   "12345",
		"X-Email-Type": "transactional",
	}
	emailMessage.SetHeaders(headers)

	// Use IP pool if available
	if e.createdIPPoolName != "" {
		emailMessage.SetIppool(e.createdIPPoolName)
		fmt.Printf("  Using IP Pool: %s\n", e.createdIPPoolName)
	}

	fmt.Println("Sending transactional email...")
	fmt.Printf("  From: %s\n", testFromEmail)
	fmt.Printf("  To: %s\n", testToEmail)
	fmt.Printf("  Subject: %s\n", emailMessage.GetSubject())

	responses, resp, err := emailAPI.SendEmail(ctx).EmailMessageObject(*emailMessage).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to send email:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if len(responses) > 0 {
		response := responses[0]
		if response.MessageId != nil {
			e.sentMessageID = *response.MessageId
		}

		fmt.Println("✓ Transactional email sent successfully!")
		if response.MessageId != nil {
			fmt.Printf("  Message ID: %s\n", *response.MessageId)
		}
		if response.To != nil {
			fmt.Printf("  To: %s\n", *response.To)
		}
	}
}

// SendMarketingEmail sends a marketing email
func (e *ESPExample) SendMarketingEmail() {
	fmt.Println("\n=== Step 8: Sending Marketing Email ===")

	ctx := e.createSubAccountAuthContext()
	emailAPI := e.client.EmailAPI

	// Create email message
	from := sendpost.NewEmailAddress()
	from.SetEmail(testFromEmail)
	from.SetName("Marketing Team")

	to := sendpost.NewRecipient()
	to.SetEmail(testToEmail)
	to.SetName("Customer 1")

	toList := []sendpost.Recipient{*to}

	emailMessage := sendpost.NewEmailMessageObject()
	emailMessage.SetFrom(*from)
	emailMessage.SetTo(toList)
	emailMessage.SetSubject("Special Offer - 20% Off Everything!")
	emailMessage.SetHtmlBody(
		"<html><body>" +
			"<h1>Special Offer!</h1>" +
			"<p>Get 20% off on all products. Use code: <strong>SAVE20</strong></p>" +
			"<p><a href=\"https://example.com/shop\">Shop Now</a></p>" +
			"</body></html>",
	)
	emailMessage.SetTextBody("Special Offer! Get 20% off on all products. Use code: SAVE20. Visit: https://example.com/shop")

	// Enable tracking
	trackOpens := true
	trackClicks := true
	emailMessage.SetTrackOpens(trackOpens)
	emailMessage.SetTrackClicks(trackClicks)

	// Add groups for analytics
	groups := []string{"marketing", "promotional"}
	emailMessage.SetGroups(groups)

	// Add custom headers
	headers := map[string]string{
		"X-Email-Type":  "marketing",
		"X-Campaign-ID": "campaign-001",
	}
	emailMessage.SetHeaders(headers)

	// Use IP pool if available
	if e.createdIPPoolName != "" {
		emailMessage.SetIppool(e.createdIPPoolName)
		fmt.Printf("  Using IP Pool: %s\n", e.createdIPPoolName)
	}

	fmt.Println("Sending marketing email...")
	fmt.Printf("  From: %s\n", testFromEmail)
	fmt.Printf("  To: %s\n", testToEmail)
	fmt.Printf("  Subject: %s\n", emailMessage.GetSubject())

	responses, resp, err := emailAPI.SendEmail(ctx).EmailMessageObject(*emailMessage).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to send email:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if len(responses) > 0 {
		response := responses[0]
		if e.sentMessageID == "" && response.MessageId != nil {
			e.sentMessageID = *response.MessageId
		}

		fmt.Println("✓ Marketing email sent successfully!")
		if response.MessageId != nil {
			fmt.Printf("  Message ID: %s\n", *response.MessageId)
		}
		if response.To != nil {
			fmt.Printf("  To: %s\n", *response.To)
		}
	}
}

// GetMessageDetails retrieves message details by message ID
func (e *ESPExample) GetMessageDetails() {
	fmt.Println("\n=== Step 9: Retrieving Message Details ===")

	if e.sentMessageID == "" {
		fmt.Println("✗ No message ID available. Please send an email first.")
		return
	}

	ctx := e.createAccountAuthContext()
	messageAPI := e.client.MessageAPI

	fmt.Printf("Retrieving message with ID: %s\n", e.sentMessageID)

	message, resp, err := messageAPI.GetMessageById(ctx, e.sentMessageID).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to get message:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Println("✓ Message retrieved successfully!")
	if message.MessageID != nil {
		fmt.Printf("  Message ID: %s\n", *message.MessageID)
	}
	if message.AccountID != nil {
		fmt.Printf("  Account ID: %d\n", *message.AccountID)
	}
	if message.SubAccountID != nil {
		fmt.Printf("  Sub-Account ID: %d\n", *message.SubAccountID)
	}
	if message.IpID != nil {
		fmt.Printf("  IP ID: %d\n", *message.IpID)
	}
	if message.PublicIP != nil {
		fmt.Printf("  Public IP: %s\n", *message.PublicIP)
	}
	if message.LocalIP != nil {
		fmt.Printf("  Local IP: %s\n", *message.LocalIP)
	}
	if message.EmailType != nil {
		fmt.Printf("  Email Type: %s\n", *message.EmailType)
	}
	if message.SubmittedAt != nil {
		fmt.Printf("  Submitted At: %d\n", *message.SubmittedAt)
	}
	if message.From.Email != nil {
		fmt.Printf("  From: %s\n", *message.From.Email)
	}
	if message.To != nil && message.To.Email != nil {
		fmt.Printf("  To: %s\n", *message.To.Email)
		if message.To.Name != nil {
			fmt.Printf("    Name: %s\n", *message.To.Name)
		}
	}
	if message.Subject != nil {
		fmt.Printf("  Subject: %s\n", *message.Subject)
	}
	if message.IpPool != nil && *message.IpPool != "" {
		fmt.Printf("  IP Pool: %s\n", *message.IpPool)
	}
	if message.Attempt != nil {
		fmt.Printf("  Delivery Attempts: %d\n", *message.Attempt)
	}
}

// GetSubAccountStats retrieves sub-account statistics
func (e *ESPExample) GetSubAccountStats() {
	fmt.Println("\n=== Step 10: Getting Sub-Account Statistics ===")

	if e.createdSubAccountID == nil {
		fmt.Println("✗ No sub-account ID available. Please create or list sub-accounts first.")
		return
	}

	ctx := e.createAccountAuthContext()
	statsAPI := e.client.StatsAPI

	// Get stats for the last 7 days
	toDate := time.Now()
	fromDate := toDate.AddDate(0, 0, -7)

	fmt.Printf("Retrieving stats for sub-account ID: %d\n", *e.createdSubAccountID)
	fmt.Printf("  From: %s\n", fromDate.Format("2006-01-02"))
	fmt.Printf("  To: %s\n", toDate.Format("2006-01-02"))

	stats, resp, err := statsAPI.AccountSubaccountStatSubaccountIdGet(ctx, int64(*e.createdSubAccountID)).
		From(fromDate.Format("2006-01-02")).
		To(toDate.Format("2006-01-02")).
		Execute()

	if err != nil {
		fmt.Printf("✗ Failed to get stats:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Println("✓ Stats retrieved successfully!")
	fmt.Printf("  Retrieved %d stat record(s)\n", len(stats))

	var totalProcessed, totalDelivered int64
	for _, stat := range stats {
		if stat.Date != nil {
			fmt.Printf("\n  Date: %s\n", *stat.Date)
		}
		if stat.Stat != nil {
			statData := stat.Stat
			if statData.Processed != nil {
				fmt.Printf("    Processed: %d\n", *statData.Processed)
				totalProcessed += int64(*statData.Processed)
			}
			if statData.Delivered != nil {
				fmt.Printf("    Delivered: %d\n", *statData.Delivered)
				totalDelivered += int64(*statData.Delivered)
			}
			if statData.Dropped != nil {
				fmt.Printf("    Dropped: %d\n", *statData.Dropped)
			}
			if statData.HardBounced != nil {
				fmt.Printf("    Hard Bounced: %d\n", *statData.HardBounced)
			}
			if statData.SoftBounced != nil {
				fmt.Printf("    Soft Bounced: %d\n", *statData.SoftBounced)
			}
			if statData.Unsubscribed != nil {
				fmt.Printf("    Unsubscribed: %d\n", *statData.Unsubscribed)
			}
			if statData.Spam != nil {
				fmt.Printf("    Spam: %d\n", *statData.Spam)
			}
		}
	}

	fmt.Println("\n  Summary (Last 7 days):")
	fmt.Printf("    Total Processed: %d\n", totalProcessed)
	fmt.Printf("    Total Delivered: %d\n", totalDelivered)
}

// GetAggregateStats retrieves aggregate statistics
func (e *ESPExample) GetAggregateStats() {
	fmt.Println("\n=== Step 11: Getting Aggregate Statistics ===")

	if e.createdSubAccountID == nil {
		fmt.Println("✗ No sub-account ID available. Please create or list sub-accounts first.")
		return
	}

	ctx := e.createAccountAuthContext()
	statsAPI := e.client.StatsAPI

	// Get aggregate stats for the last 7 days
	toDate := time.Now()
	fromDate := toDate.AddDate(0, 0, -7)

	fmt.Printf("Retrieving aggregate stats for sub-account ID: %d\n", *e.createdSubAccountID)
	fmt.Printf("  From: %s\n", fromDate.Format("2006-01-02"))
	fmt.Printf("  To: %s\n", toDate.Format("2006-01-02"))

	aggregateStat, resp, err := statsAPI.AccountSubaccountStatSubaccountIdAggregateGet(ctx, int64(*e.createdSubAccountID)).
		From(fromDate.Format("2006-01-02")).
		To(toDate.Format("2006-01-02")).
		Execute()

	if err != nil {
		fmt.Printf("✗ Failed to get aggregate stats:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Println("✓ Aggregate stats retrieved successfully!")
	if aggregateStat.Processed != nil {
		fmt.Printf("  Processed: %d\n", *aggregateStat.Processed)
	}
	if aggregateStat.Delivered != nil {
		fmt.Printf("  Delivered: %d\n", *aggregateStat.Delivered)
	}
	if aggregateStat.Dropped != nil {
		fmt.Printf("  Dropped: %d\n", *aggregateStat.Dropped)
	}
	if aggregateStat.HardBounced != nil {
		fmt.Printf("  Hard Bounced: %d\n", *aggregateStat.HardBounced)
	}
	if aggregateStat.SoftBounced != nil {
		fmt.Printf("  Soft Bounced: %d\n", *aggregateStat.SoftBounced)
	}
	if aggregateStat.Unsubscribed != nil {
		fmt.Printf("  Unsubscribed: %d\n", *aggregateStat.Unsubscribed)
	}
	if aggregateStat.Spam != nil {
		fmt.Printf("  Spam: %d\n", *aggregateStat.Spam)
	}
}

// ListIPs lists all IPs
func (e *ESPExample) ListIPs() {
	fmt.Println("\n=== Step 12: Listing All IPs ===")

	ctx := e.createAccountAuthContext()
	ipAPI := e.client.IPAPI

	fmt.Println("Retrieving all IPs...")
	ips, resp, err := ipAPI.GetAllIps(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to list IPs:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Retrieved %d IP(s)\n", len(ips))
	for _, ip := range ips {
		fmt.Printf("  - ID: %d\n", ip.Id)
		if ip.PublicIP != "" {
			fmt.Printf("    IP Address: %s\n", ip.PublicIP)
		}
		if ip.ReverseDNSHostname != nil {
			fmt.Printf("    Reverse DNS: %s\n", *ip.ReverseDNSHostname)
		}
		if ip.Created != 0 {
			fmt.Printf("    Created: %d\n", ip.Created)
		}
		fmt.Println()
	}
}

// CreateIPPool creates an IP pool
func (e *ESPExample) CreateIPPool() {
	fmt.Println("\n=== Step 13: Creating IP Pool ===")

	ctx := e.createAccountAuthContext()
	ipPoolsAPI := e.client.IPPoolsAPI
	ipAPI := e.client.IPAPI

	// First, get available IPs
	fmt.Println("Retrieving available IPs...")
	ips, resp, err := ipAPI.GetAllIps(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to get IPs:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if len(ips) == 0 {
		fmt.Println("⚠️  No IPs available. Please allocate IPs first.")
		return
	}

	// Create IP pool request
	poolName := fmt.Sprintf("Marketing Pool %d", time.Now().Unix())
	routingStrategy := int32(0) // 0 = RoundRobin, 1 = EmailProviderStrategy

	poolIPs := []sendpost.EIP{}
	// Add first available IP (you can add more)
	if len(ips) > 0 && ips[0].PublicIP != "" {
		eip := sendpost.NewEIP(ips[0].PublicIP)
		poolIPs = append(poolIPs, *eip)
	}

	poolRequest := sendpost.NewIPPoolCreateRequest()
	poolRequest.SetName(poolName)
	poolRequest.SetRoutingStrategy(routingStrategy)
	poolRequest.SetIps(poolIPs)

	// Set warmup interval (required, must be > 0)
	warmupInterval := int32(24) // 24 hours
	poolRequest.SetWarmupInterval(warmupInterval)

	// Set overflow strategy (0 = None, 1 = Use overflow pool)
	overflowStrategy := int32(0)
	poolRequest.SetOverflowStrategy(overflowStrategy)

	fmt.Printf("Creating IP pool: %s\n", poolName)
	fmt.Println("  Routing Strategy: Round Robin")
	fmt.Printf("  IPs: %d\n", len(poolIPs))
	fmt.Printf("  Warmup Interval: %d hours\n", warmupInterval)

	ipPool, resp, err := ipPoolsAPI.CreateIPPool(ctx).IPPoolCreateRequest(*poolRequest).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to create IP pool:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if ipPool.Id != nil {
		id := int64(*ipPool.Id)
		e.createdIPPoolID = &id
	}
	if ipPool.Name != nil {
		e.createdIPPoolName = *ipPool.Name
	}

	fmt.Println("✓ IP pool created successfully!")
	if ipPool.Id != nil {
		fmt.Printf("  ID: %d\n", *ipPool.Id)
	}
	if ipPool.Name != nil {
		fmt.Printf("  Name: %s\n", *ipPool.Name)
	}
	if ipPool.RoutingStrategy != nil {
		fmt.Printf("  Routing Strategy: %d\n", *ipPool.RoutingStrategy)
	}
	if ipPool.Ips != nil {
		fmt.Printf("  IPs in pool: %d\n", len(ipPool.Ips))
	}
}

// ListIPPools lists all IP pools
func (e *ESPExample) ListIPPools() {
	fmt.Println("\n=== Step 14: Listing All IP Pools ===")

	ctx := e.createAccountAuthContext()
	ipPoolsAPI := e.client.IPPoolsAPI

	fmt.Println("Retrieving all IP pools...")
	ipPools, resp, err := ipPoolsAPI.GetAllIPPools(ctx).Execute()

	if err != nil {
		fmt.Printf("✗ Failed to list IP pools:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Retrieved %d IP pool(s)\n", len(ipPools))
	for _, ipPool := range ipPools {
		if ipPool.Id != nil {
			fmt.Printf("  - ID: %d\n", *ipPool.Id)
		}
		if ipPool.Name != nil {
			fmt.Printf("    Name: %s\n", *ipPool.Name)
		}
		if ipPool.RoutingStrategy != nil {
			fmt.Printf("    Routing Strategy: %d\n", *ipPool.RoutingStrategy)
		}
		if ipPool.Ips != nil {
			fmt.Printf("    IPs in pool: %d\n", len(ipPool.Ips))
			for _, ip := range ipPool.Ips {
				if ip.PublicIP != "" {
					fmt.Printf("      - %s\n", ip.PublicIP)
				}
			}
		}
		fmt.Println()
	}
}

// GetAccountStats retrieves account-level statistics
func (e *ESPExample) GetAccountStats() {
	fmt.Println("\n=== Step 15: Getting Account-Level Statistics ===")

	ctx := e.createAccountAuthContext()
	statsAAPI := e.client.StatsAAPI

	// Get stats for the last 7 days
	toDate := time.Now()
	fromDate := toDate.AddDate(0, 0, -7)

	fmt.Println("Retrieving account-level stats...")
	fmt.Printf("  From: %s\n", fromDate.Format("2006-01-02"))
	fmt.Printf("  To: %s\n", toDate.Format("2006-01-02"))

	accountStats, resp, err := statsAAPI.GetAllAccountStats(ctx).
		From(fromDate.Format("2006-01-02")).
		To(toDate.Format("2006-01-02")).
		Execute()

	if err != nil {
		fmt.Printf("✗ Failed to get account stats:\n")
		fmt.Printf("  Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Println("✓ Account stats retrieved successfully!")
	fmt.Printf("  Retrieved %d stat record(s)\n", len(accountStats))

	for _, stat := range accountStats {
		if stat.Date != nil {
			fmt.Printf("\n  Date: %s\n", *stat.Date)
		}
		if stat.Stat != nil {
			statData := stat.Stat
			if statData.Processed != nil {
				fmt.Printf("    Processed: %d\n", *statData.Processed)
			}
			if statData.Delivered != nil {
				fmt.Printf("    Delivered: %d\n", *statData.Delivered)
			}
			if statData.Dropped != nil {
				fmt.Printf("    Dropped: %d\n", *statData.Dropped)
			}
			if statData.HardBounced != nil {
				fmt.Printf("    Hard Bounced: %d\n", *statData.HardBounced)
			}
			if statData.SoftBounced != nil {
				fmt.Printf("    Soft Bounced: %d\n", *statData.SoftBounced)
			}
			if statData.Opened != nil {
				fmt.Printf("    Opens: %d\n", *statData.Opened)
			}
			if statData.Clicked != nil {
				fmt.Printf("    Clicks: %d\n", *statData.Clicked)
			}
			if statData.Unsubscribed != nil {
				fmt.Printf("    Unsubscribed: %d\n", *statData.Unsubscribed)
			}
			if statData.Spams != nil {
				fmt.Printf("    Spams: %d\n", *statData.Spams)
			}
		}
	}
}

// RunCompleteWorkflow runs the complete ESP workflow
func (e *ESPExample) RunCompleteWorkflow() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║   SendPost Go SDK - ESP Example Workflow                     ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	// Check if API keys are set
	if e.accountAPIKey == "YOUR_ACCOUNT_API_KEY_HERE" || e.subAccountAPIKey == "YOUR_SUB_ACCOUNT_API_KEY_HERE" {
		fmt.Println("\n⚠️  WARNING: Please set your API keys!")
		fmt.Println("   Set environment variables:")
		fmt.Println("   - SENDPOST_ACCOUNT_API_KEY")
		fmt.Println("   - SENDPOST_SUB_ACCOUNT_API_KEY")
		fmt.Println("   Or modify the constants in main.go")
		fmt.Println()
	}

	// Step 1: List existing sub-accounts (or create new one)
	e.ListSubAccounts()

	// Step 2: Create webhook for event notifications
	e.CreateWebhook()
	e.ListWebhooks()

	// Step 3: Add and verify domain
	e.AddDomain()
	e.ListDomains()

	// Step 4: Manage IPs and IP pools (before sending emails)
	e.ListIPs()
	e.CreateIPPool()
	e.ListIPPools()

	// Step 5: Send emails (using the created IP pool)
	e.SendTransactionalEmail()
	e.SendMarketingEmail()

	// Step 6: Monitor statistics
	e.GetSubAccountStats()
	e.GetAggregateStats()

	// Step 7: Get account-level overview
	e.GetAccountStats()

	// Step 8: Retrieve message details (at the end to give system time to store data)
	e.GetMessageDetails()

	fmt.Println("\n╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Workflow Complete!                                          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}

func main() {
	example := NewESPExample()
	example.RunCompleteWorkflow()
}
