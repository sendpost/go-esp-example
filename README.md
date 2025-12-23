# SendPost Go SDK - ESP Example

This project provides a comprehensive example demonstrating how Email Service Providers (ESPs) can use the SendPost Go SDK to manage email sending operations.

## Overview

The example demonstrates a complete ESP workflow including:

1. **Sub-Account Management** - Create and manage sub-accounts for different clients or use cases
2. **Webhook Setup** - Configure webhooks to receive real-time email event notifications
3. **Domain Management** - Add and verify sending domains
4. **Email Sending** - Send transactional and marketing emails
5. **Message Tracking** - Retrieve message details for tracking and debugging
6. **Statistics & Analytics** - Monitor email performance via sub-account stats, IP stats, and IP pool stats
7. **IP Pool Management** - Create and manage IP pools for better deliverability control

## Prerequisites

- Go 1.18 or higher
- SendPost account with:
  - Account API Key (for account-level operations)
  - Sub-Account API Key (for sub-account-level operations)

## Setup

### 1. Clone or Navigate to the Project

```bash
cd example-sdk-go
```

### 2. Configure API Keys

You can set API keys in two ways:

#### Option A: Environment Variables (Recommended)

```bash
export SENDPOST_ACCOUNT_API_KEY="your_account_api_key_here"
export SENDPOST_SUB_ACCOUNT_API_KEY="your_sub_account_api_key_here"
```

#### Option B: Edit the Source Code

Edit `main.go` and update the constants in the `NewESPExample` function:

```go
accountAPIKey := "your_account_api_key_here"
subAccountAPIKey := "your_sub_account_api_key_here"
```

### 3. Update Configuration Values

Edit `main.go` and update the constants at the top of the file:

- `testFromEmail` - Your verified sender email address
- `testToEmail` - Recipient email address
- `testDomainName` - Your sending domain
- `webhookURL` - Your webhook endpoint URL

## Running the Example

### Install Dependencies

```bash
go mod download
```

### Run the Complete Workflow

```bash
go run main.go
```

This will execute the complete ESP workflow demonstrating all features.

## Project Structure

```
example-sdk-go/
├── go.mod              # Go module definition
├── README.md           # This file
├── .gitignore          # Git ignore file
└── main.go             # Main example program
```

## Workflow Steps

The example demonstrates the following workflow:

### Step 1: Sub-Account Management
- List all sub-accounts
- Create new sub-accounts for different clients or use cases

### Step 2: Webhook Configuration
- Create webhooks to receive email event notifications
- Configure which events to receive (delivered, opened, clicked, bounced, etc.)

### Step 3: Domain Management
- Add sending domains
- View DNS records needed for domain verification
- List all domains

### Step 4: Email Sending
- Send transactional emails (order confirmations, receipts, etc.)
- Send marketing emails (newsletters, promotions, etc.)
- Configure tracking (opens, clicks)
- Add custom headers and fields

### Step 5: Message Tracking
- Retrieve message details by message ID
- View delivery information, IP used, submission time, etc.

### Step 6: Statistics & Analytics
- Get sub-account statistics (processed, delivered, opens, clicks, bounces, etc.)
- Get aggregate statistics
- Get account-level statistics across all sub-accounts

### Step 7: IP and IP Pool Management
- List all dedicated IPs
- Create IP pools for better deliverability control
- View IP pool configurations

## Key Features Demonstrated

### Email Sending
- **Transactional Emails**: Order confirmations, receipts, notifications
- **Marketing Emails**: Newsletters, promotions, campaigns
- **Tracking**: Open tracking, click tracking
- **Customization**: Custom headers, custom fields, groups

### Statistics & Monitoring
- **Sub-Account Stats**: Daily statistics for a specific sub-account
- **Aggregate Stats**: Overall performance metrics
- **Account Stats**: Statistics across all sub-accounts
- **Performance Metrics**: Open rates, click rates, delivery rates

### Infrastructure Management
- **Sub-Accounts**: Organize sending by client, product, or use case
- **Domains**: Add and verify sending domains
- **IPs**: Monitor dedicated IP addresses
- **IP Pools**: Group IPs for better deliverability control

### Event Handling
- **Webhooks**: Receive real-time notifications for email events
- **Event Types**: Processed, delivered, dropped, bounced, opened, clicked, unsubscribed, spam

## API Keys Explained

### Account API Key (`X-Account-ApiKey`)
Used for account-level operations:
- Creating and managing sub-accounts
- Managing IPs and IP pools
- Creating webhooks
- Getting account-level statistics
- Retrieving messages

### Sub-Account API Key (`X-SubAccount-ApiKey`)
Used for sub-account-level operations:
- Sending emails
- Managing domains
- Managing suppressions
- Getting sub-account statistics

## Example Output

When you run the example, you'll see output like:

```
╔═══════════════════════════════════════════════════════════════╗
║   SendPost Go SDK - ESP Example Workflow                     ║
╚═══════════════════════════════════════════════════════════════╝

=== Step 1: Listing All Sub-Accounts ===
Retrieving all sub-accounts...
✓ Retrieved 3 sub-account(s)
  - ID: 50441
    Name: API
    API Key: pR0YIuxYSbVwmQi2Y8Qs
    ...

=== Step 2: Creating Webhook ===
Creating webhook...
  URL: https://your-webhook-endpoint.com/webhook
✓ Webhook created successfully!
  ID: 12345
  ...

...
```

## Error Handling

The example includes comprehensive error handling. If an operation fails, you'll see:
- HTTP status code
- Error response details
- Error message for debugging

Common issues:
- **401 Unauthorized**: Invalid or missing API key
- **403 Forbidden**: Resource already exists or insufficient permissions
- **404 Not Found**: Resource ID doesn't exist
- **422 Unprocessable Entity**: Invalid request body or parameters

## Code Structure

The example is organized into a single `ESPExample` struct with methods for each workflow step:

- `NewESPExample()` - Initializes the example with API client
- `ListSubAccounts()` - Lists all sub-accounts
- `CreateSubAccount()` - Creates a new sub-account
- `CreateWebhook()` - Creates a webhook
- `ListWebhooks()` - Lists all webhooks
- `AddDomain()` - Adds a sending domain
- `ListDomains()` - Lists all domains
- `SendTransactionalEmail()` - Sends a transactional email
- `SendMarketingEmail()` - Sends a marketing email
- `GetMessageDetails()` - Retrieves message details
- `GetSubAccountStats()` - Gets sub-account statistics
- `GetAggregateStats()` - Gets aggregate statistics
- `ListIPs()` - Lists all IPs
- `CreateIPPool()` - Creates an IP pool
- `ListIPPools()` - Lists all IP pools
- `GetAccountStats()` - Gets account-level statistics
- `RunCompleteWorkflow()` - Runs the complete workflow

## Authentication Context

The example uses Go's `context.Context` to pass authentication credentials:

```go
// Account-level operations
ctx := context.WithValue(
    context.Background(),
    sendpost.ContextAPIKeys,
    map[string]sendpost.APIKey{
        "accountAuth": {Key: accountAPIKey},
    },
)

// Sub-account-level operations
ctx := context.WithValue(
    context.Background(),
    sendpost.ContextAPIKeys,
    map[string]sendpost.APIKey{
        "subAccountAuth": {Key: subAccountAPIKey},
    },
)
```

## Building the Example

To build a binary:

```bash
go build -o esp-example main.go
```

Then run:

```bash
./esp-example
```

## Testing Individual Steps

You can modify the `main()` function to test individual steps:

```go
func main() {
    example := NewESPExample()
    
    // Test only email sending
    example.SendTransactionalEmail()
    example.SendMarketingEmail()
    
    // Or test only statistics
    example.GetSubAccountStats()
    example.GetAggregateStats()
}
```

## Next Steps

After running the example:

1. **Customize for Your Use Case**: Modify the example to match your specific requirements
2. **Integrate with Your Application**: Use the SDK in your own Go application
3. **Set Up Webhooks**: Configure your webhook endpoint to receive email events
4. **Monitor Statistics**: Set up regular monitoring of your email performance
5. **Optimize Deliverability**: Use IP pools and domain verification to improve deliverability

## Additional Resources

- [SendPost API Documentation](https://docs.sendpost.io)
- [SendPost Go SDK](https://github.com/sendpost/sendpost_go_sdk)
- [SendPost Developer Portal](https://app.sendpost.io)

## Support

For questions or issues:
- Email: hello@sendpost.io
- Website: https://sendpost.io
- Documentation: https://docs.sendpost.io

## License

This example is provided as-is for demonstration purposes.

