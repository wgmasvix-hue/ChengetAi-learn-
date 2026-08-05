# User Service

The User Service manages user profiles, preferences, and settings.

## Responsibilities

- User profile management
- User preferences and settings
- Role and permission management
- User search and discovery
- Profile analytics

## Port

Port 8002

## Building & Running

```bash
cd apps/user-service
make build
make run
```

## API Endpoints

### Users
- `GET /api/v1/users/{id}` - Get user profile
- `PUT /api/v1/users/{id}` - Update user profile
- `DELETE /api/v1/users/{id}` - Delete user account
- `GET /api/v1/users/{id}/preferences` - Get user preferences
- `PUT /api/v1/users/{id}/preferences` - Update preferences

### Profiles
- `GET /api/v1/users/search` - Search users
- `GET /api/v1/users/{id}/activity` - User activity

## Next Steps

- [ ] Implement user profile management
- [ ] Add user preferences
- [ ] Implement role management
- [ ] Add search functionality
- [ ] Add analytics tracking
