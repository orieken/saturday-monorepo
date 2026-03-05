# Saturday Framework - API Testing Analysis & Recommendations

## 🎯 Current State Analysis

### What We Have

**UI Testing (Fully Supported):**
- ✅ `BaseSite` - Site-level organization
- ✅ `BasePage` - Page Object pattern
- ✅ `BaseElement` - Element abstraction
- ✅ `BaseFlow` - Multi-step user journeys
- ✅ `BaseFilter` - Data filtering

**API Testing (Partially Planned):**
- ⚠️ `ServiceGenerationRequest` model exists
- ⚠️ `EndpointDefinition` model exists
- ❌ No `BaseService` class in saturday-core
- ❌ No Service generator implementation
- ❌ No API testing templates

---

## 🔍 API Testing Requirements

### Core Concepts

**Service Pattern (Parallel to Site/Page):**
```
Site (UI)           →  Service (API)
├── Pages           →  ├── Endpoints
├── Elements        →  ├── Request/Response Models
└── Flows           →  └── API Flows (multi-endpoint sequences)
```

### What We Need

#### 1. **BaseService Class** (saturday-core)
```typescript
// packages/saturday-core/src/base/base-service.ts
import { APIRequestContext } from '@playwright/test';

export abstract class BaseService {
  protected baseUrl: string;
  protected context: APIRequestContext;
  protected endpoints: Map<string, Endpoint>;

  constructor(context: APIRequestContext, baseUrl: string) {
    this.context = context;
    this.baseUrl = baseUrl;
    this.endpoints = new Map();
    this.initializeEndpoints();
  }

  protected abstract initializeEndpoints(): void;

  protected registerEndpoint(name: string, method: string, path: string): void {
    this.endpoints.set(name, { method, path });
  }

  async request(endpointName: string, options?: RequestOptions): Promise<APIResponse> {
    const endpoint = this.endpoints.get(endpointName);
    if (!endpoint) {
      throw new Error(`Endpoint ${endpointName} not found`);
    }

    const url = `${this.baseUrl}${endpoint.path}`;
    return this.context.request({
      method: endpoint.method,
      url,
      ...options
    });
  }

  // Convenience methods
  async get(endpointName: string, params?: any): Promise<APIResponse> { }
  async post(endpointName: string, data?: any): Promise<APIResponse> { }
  async put(endpointName: string, data?: any): Promise<APIResponse> { }
  async delete(endpointName: string): Promise<APIResponse> { }
}
```

#### 2. **Service Template** (saturday-mcp)
```typescript
// internal/templates/data/service.tmpl
import { BaseService } from '@orieken/saturday-core';
import { APIRequestContext } from '@playwright/test';

export class {{pascalCase .Name}}Service extends BaseService {
  constructor(context: APIRequestContext, baseUrl: string) {
    super(context, baseUrl);
  }

  protected initializeEndpoints(): void {
    {{range .Endpoints}}
    this.registerEndpoint('{{.Name}}', '{{.Method}}', '{{.Path}}');
    {{end}}
  }

  // Typed methods for each endpoint
  {{range .Endpoints}}
  async {{camelCase .Name}}({{if eq .Method "POST" "PUT" "PATCH"}}data?: any{{end}}): Promise<any> {
    const response = await this.{{lower .Method}}('{{.Name}}'{{if eq .Method "POST" "PUT" "PATCH"}}, data{{end}});
    return response.json();
  }
  {{end}}
}
```

#### 3. **API Flow Pattern**
```typescript
// Similar to UI Flows, but for API sequences
export class CheckoutAPIFlow extends BaseFlow {
  constructor(
    private userService: UserService,
    private cartService: CartService,
    private orderService: OrderService
  ) {
    super();
  }

  async execute(userId: string, items: CartItem[]): Promise<Order> {
    // 1. Get user
    const user = await this.userService.getUser(userId);
    
    // 2. Add items to cart
    for (const item of items) {
      await this.cartService.addItem(userId, item);
    }
    
    // 3. Create order
    const order = await this.orderService.createOrder(userId);
    
    return order;
  }
}
```

---

## 📋 Proposed TODOs

### TODO-012: BaseService Implementation (saturday-core)
**Priority**: HIGH  
**Effort**: 2-3 hours

**Deliverables:**
- [ ] `packages/saturday-core/src/base/base-service.ts`
- [ ] Request/Response type definitions
- [ ] Error handling utilities
- [ ] Authentication helpers (Bearer, API Key, etc.)
- [ ] Request interceptors
- [ ] Response validators
- [ ] Comprehensive tests

**Example Generated Service:**
```typescript
export class UserService extends BaseService {
  async getUser(id: string): Promise<User> {
    const response = await this.get('getUser', { params: { id } });
    return response.json();
  }

  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await this.post('createUser', data);
    return response.json();
  }
}
```

---

### TODO-013: Service Generator (saturday-mcp)
**Priority**: HIGH  
**Effort**: 1-2 hours (follows existing pattern!)

**Deliverables:**
- [ ] `internal/generators/service_generator.go`
- [ ] `internal/templates/data/service.tmpl`
- [ ] Service generator tests
- [ ] MCP tool integration (`generate_service`)
- [ ] JSON schema for service generation

**Pattern**: Identical to Site/Page/Flow generators!

---

### TODO-014: API Testing Examples
**Priority**: MEDIUM  
**Effort**: 2-3 hours

**Deliverables:**
- [ ] Example API service in ye-olde-magic-shop
- [ ] API test scenarios (Cucumber)
- [ ] Mixed UI + API test examples
- [ ] Documentation

**Example Scenario:**
```gherkin
Feature: Order API
  Scenario: Create order via API
    Given I have a valid auth token
    When I create an order via API with:
      | product | quantity |
      | Wand    | 1        |
    Then the API response status should be 201
    And the order should have a valid ID
    And I can retrieve the order via API
```

---

### TODO-015: Request/Response Models Generator
**Priority**: MEDIUM  
**Effort**: 2-3 hours

**Deliverables:**
- [ ] Generate TypeScript interfaces from JSON schemas
- [ ] Generate request/response types
- [ ] Validation helpers
- [ ] Mock data generators

**Example:**
```typescript
// Generated from OpenAPI/JSON Schema
export interface CreateOrderRequest {
  userId: string;
  items: OrderItem[];
  shippingAddress: Address;
}

export interface CreateOrderResponse {
  orderId: string;
  status: OrderStatus;
  total: number;
  createdAt: string;
}
```

---

### TODO-016: Contract Testing Support
**Priority**: LOW  
**Effort**: 3-4 hours

**Deliverables:**
- [ ] Pact integration
- [ ] Contract validation
- [ ] Schema validation
- [ ] API versioning support

---

## 🎯 Recommended Implementation Order

### Phase 1: Core API Testing (Essential)
1. **TODO-012**: BaseService Implementation
2. **TODO-013**: Service Generator
3. **TODO-014**: API Testing Examples

**Result**: Full API testing capability with code generation

### Phase 2: Enhanced API Testing (Nice to Have)
4. **TODO-015**: Request/Response Models Generator
5. **TODO-016**: Contract Testing Support

---

## 💡 Design Decisions

### 1. **Use Playwright's API Testing**
**Why**: Already integrated, consistent with UI testing

```typescript
import { test, expect } from '@playwright/test';

test('API test', async ({ request }) => {
  const userService = new UserService(request, 'https://api.example.com');
  const user = await userService.getUser('123');
  expect(user.id).toBe('123');
});
```

### 2. **Service Pattern Mirrors Site Pattern**
**Why**: Consistency, reusability, familiar to users

```
Site                Service
├── Pages           ├── Endpoints
├── Flows           ├── API Flows
└── Elements        └── Models
```

### 3. **Support Mixed UI + API Tests**
**Why**: Real-world scenarios often need both

```typescript
test('E2E: Create order via UI, verify via API', async ({ page, request }) => {
  // UI: Add to cart
  const site = new MagicShopSite(page, baseUrl);
  await site.goToPage('products');
  await site.page('products').addToCart('Wand');
  
  // API: Verify order created
  const orderService = new OrderService(request, apiUrl);
  const orders = await orderService.getUserOrders(userId);
  expect(orders).toHaveLength(1);
});
```

### 4. **Generate from OpenAPI/Swagger**
**Future Enhancement**: Import OpenAPI spec, generate services

```bash
saturday-mcp generate-service --from-openapi api-spec.yaml
```

---

## 📊 Comparison: UI vs API Testing

| Aspect | UI Testing | API Testing |
|--------|-----------|-------------|
| **Base Class** | BaseSite | BaseService |
| **Components** | Pages, Elements | Endpoints, Models |
| **Flows** | User journeys | API sequences |
| **Context** | Browser Page | API Request Context |
| **Generator** | ✅ Implemented | ⚠️ Needs TODO-013 |
| **Core Support** | ✅ saturday-core | ⚠️ Needs TODO-012 |

---

## 🎯 Benefits of API Testing in Saturday

### 1. **Faster Test Execution**
- API tests run 10-100x faster than UI tests
- No browser overhead
- Parallel execution

### 2. **Better Test Coverage**
- Test edge cases easily
- Test error conditions
- Test data validation

### 3. **Hybrid Testing**
- Setup via API, verify via UI
- Create data via API, interact via UI
- Best of both worlds

### 4. **Consistent Patterns**
- Same generator pattern
- Same testing approach
- Same framework concepts

---

## 🚀 Quick Start (After Implementation)

### 1. Generate a Service
```
Ask Claude: "Use Saturday MCP to generate an API service for Users with endpoints:
- GET /users/:id (getUser)
- POST /users (createUser)
- PUT /users/:id (updateUser)
- DELETE /users/:id (deleteUser)"
```

### 2. Use in Tests
```typescript
import { test } from '@playwright/test';
import { UserService } from './services/user-service';

test('User CRUD operations', async ({ request }) => {
  const service = new UserService(request, 'https://api.example.com');
  
  // Create
  const user = await service.createUser({ name: 'John', email: 'john@example.com' });
  
  // Read
  const retrieved = await service.getUser(user.id);
  
  // Update
  await service.updateUser(user.id, { name: 'John Doe' });
  
  // Delete
  await service.deleteUser(user.id);
});
```

---

## 📝 Summary

### Current Gap
- ✅ UI testing fully supported
- ❌ API testing not yet supported
- ⚠️ Models exist but no implementation

### Recommendation
**Implement TODO-012 and TODO-013 ASAP**

**Why**:
1. Follows existing patterns (easy to implement)
2. Completes the testing story
3. Enables hybrid UI + API testing
4. High value, low effort

**Effort**: ~4-5 hours total for both TODOs

### Next Steps
1. Review this analysis
2. Approve TODO-012 and TODO-013
3. Implement BaseService in saturday-core
4. Implement Service Generator in saturday-mcp
5. Create examples in ye-olde-magic-shop
6. Update documentation

---

**Would you like me to implement TODO-012 and TODO-013 now?** 🚀
