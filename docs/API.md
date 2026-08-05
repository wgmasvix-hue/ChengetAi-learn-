# ChengetAi API Documentation

## Base URL

```
http://localhost:8000/api/v1
```

## Authentication

All endpoints (except `/auth/login`, `/auth/register`) require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## Response Format

All responses are JSON:

### Success Response
```json
{
    "success": true,
    "data": { ... },
    "message": "Operation successful"
}
```

### Error Response
```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "Error description"
    }
}
```

## HTTP Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Permission denied
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## API Endpoints

### Health & Status

#### Check Health
```
GET /health
```

Returns service health status.

#### Check Readiness
```
GET /ready
```

Returns service readiness status.

### Authentication

#### Register User
```
POST /auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "secure_password",
    "full_name": "User Name",
    "country": "ZW"
}
```

Response:
```json
{
    "success": true,
    "data": {
        "user_id": "uuid",
        "email": "user@example.com",
        "access_token": "jwt_token",
        "refresh_token": "refresh_token"
    }
}
```

#### Login
```
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password"
}
```

#### Logout
```
POST /auth/logout
Authorization: Bearer <token>
```

#### Refresh Token
```
POST /auth/refresh
Content-Type: application/json

{
    "refresh_token": "refresh_token"
}
```

### Users

#### Get User Profile
```
GET /users/{user_id}
Authorization: Bearer <token>
```

#### Update Profile
```
PUT /users/{user_id}
Authorization: Bearer <token>
Content-Type: application/json

{
    "full_name": "New Name",
    "bio": "Bio text",
    "avatar_url": "https://..."
}
```

#### Search Users
```
GET /users/search?q=query&limit=10&offset=0
Authorization: Bearer <token>
```

### Schools

#### Get School
```
GET /schools/{school_id}
Authorization: Bearer <token>
```

#### List Schools
```
GET /schools?limit=20&offset=0
Authorization: Bearer <token>
```

### Wallet

#### Get Wallet Balance
```
GET /wallet/balance
Authorization: Bearer <token>
```

Response:
```json
{
    "success": true,
    "data": {
        "balance": "1250.50",
        "currency": "USD",
        "pending": "100.00",
        "total_earned": "5000.00"
    }
}
```

#### Get Transactions
```
GET /wallet/transactions?limit=20&offset=0
Authorization: Bearer <token>
```

#### Request Withdrawal
```
POST /wallet/withdraw
Authorization: Bearer <token>
Content-Type: application/json

{
    "amount": "100.00",
    "method": "bank_transfer"
}
```

### Search

#### Search Resources
```
GET /search?q=query&type=resource&limit=20&offset=0
Authorization: Bearer <token>
```

Query Parameters:
- `q` - Search query (required)
- `type` - Resource type (resource, quiz, course)
- `limit` - Results per page (default 20, max 100)
- `offset` - Pagination offset (default 0)
- `sort` - Sort by (relevance, date, views)
- `filter` - Filter by metadata

### Analytics

#### Get Learning Stats
```
GET /analytics/learning?period=month
Authorization: Bearer <token>
```

Response:
```json
{
    "success": true,
    "data": {
        "total_views": 150,
        "active_resources": 12,
        "completion_rate": 75,
        "avg_session_duration": "23m"
    }
}
```

#### Get Creator Stats
```
GET /analytics/creator?period=month
Authorization: Bearer <token>
```

Response:
```json
{
    "success": true,
    "data": {
        "total_revenue": "1250.50",
        "resource_earnings": { ... },
        "top_resources": [ ... ]
    }
}
```

## Error Codes

| Code | Message | HTTP Status |
|------|---------|-------------|
| `INVALID_REQUEST` | Request validation failed | 400 |
| `UNAUTHORIZED` | Authentication required | 401 |
| `PERMISSION_DENIED` | Insufficient permissions | 403 |
| `NOT_FOUND` | Resource not found | 404 |
| `CONFLICT` | Resource already exists | 409 |
| `RATE_LIMIT` | Too many requests | 429 |
| `INTERNAL_ERROR` | Internal server error | 500 |

## Rate Limiting

- **Anonymous**: 100 requests/hour
- **Authenticated**: 1000 requests/hour
- **Premium**: Unlimited

Rate limit info is in response headers:
- `X-RateLimit-Limit` - Total requests allowed
- `X-RateLimit-Remaining` - Requests remaining
- `X-RateLimit-Reset` - Unix timestamp of reset time

## Pagination

List endpoints support pagination:

```
GET /resources?limit=20&offset=0
```

Response:
```json
{
    "success": true,
    "data": [ ... ],
    "pagination": {
        "total": 1000,
        "limit": 20,
        "offset": 0,
        "pages": 50
    }
}
```

## Webhooks

To receive events, subscribe to webhooks (future implementation):

```
POST /webhooks/subscribe
Authorization: Bearer <token>
Content-Type: application/json

{
    "url": "https://your-server.com/webhook",
    "events": ["resource.created", "user.joined"]
}
```

Events will be POSTed to your endpoint.

## Code Examples

### JavaScript/Fetch
```javascript
// Login
const response = await fetch('http://localhost:8000/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        email: 'user@example.com',
        password: 'password'
    })
});

const { data } = await response.json();
const token = data.access_token;

// Get balance
const balanceResponse = await fetch('http://localhost:8000/api/v1/wallet/balance', {
    headers: { 'Authorization': `Bearer ${token}` }
});

const { data: balance } = await balanceResponse.json();
console.log(`Balance: ${balance.balance} ${balance.currency}`);
```

### Python/Requests
```python
import requests

BASE_URL = 'http://localhost:8000/api/v1'

# Login
response = requests.post(f'{BASE_URL}/auth/login', json={
    'email': 'user@example.com',
    'password': 'password'
})

token = response.json()['data']['access_token']

# Get balance
headers = {'Authorization': f'Bearer {token}'}
balance = requests.get(f'{BASE_URL}/wallet/balance', headers=headers)
print(f"Balance: {balance.json()['data']['balance']}")
```

### Go/HTTP
```go
package main

import (
    "net/http"
    "fmt"
)

func main() {
    client := &http.Client{}
    
    // Login
    req, _ := http.NewRequest("POST", 
        "http://localhost:8000/api/v1/auth/login", nil)
    // ... set body and headers
    
    resp, _ := client.Do(req)
    defer resp.Body.Close()
    
    // ... parse response
}
```

## Support

For API questions or issues:
- GitHub Issues: https://github.com/wgmasvix-hue/ChengetAi-learn-/issues
- Email: support@chengetai.africa
- Documentation: https://docs.chengetai.africa
