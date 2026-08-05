# Authentication Service

The Authentication Service handles all user authentication and authorization for the ChengetAi Platform.

## Responsibilities

- User registration and login
- JWT token generation and validation
- OAuth2 integration (Google, Apple)
- Session management
- Password reset and 2FA
- Token refresh and revocation

## Port

Port 8001

## Building & Running

```bash
cd apps/auth-service
make build
make run
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/verify` - Verify JWT token

### OAuth2
- `GET /api/v1/auth/oauth2/google` - Google OAuth2 redirect
- `GET /api/v1/auth/oauth2/apple` - Apple OAuth2 redirect

### Password
- `POST /api/v1/auth/password/reset` - Request password reset
- `POST /api/v1/auth/password/change` - Change password

### 2FA
- `POST /api/v1/auth/2fa/enable` - Enable two-factor authentication
- `POST /api/v1/auth/2fa/verify` - Verify 2FA code

## Next Steps

- [ ] Implement JWT token management
- [ ] Add OAuth2 providers
- [ ] Implement 2FA (TOTP)
- [ ] Add session management
- [ ] Add audit logging
