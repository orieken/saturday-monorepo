import { APIRequestContext, APIResponse } from '@playwright/test';

const MOCK_API_URL = process.env.MOCK_API_URL || 'http://localhost:8001';

export interface GeoPoint {
  lat: number;
  lng: number;
}

export interface ReagentSite extends GeoPoint {
  id: string;
  name: string;
}

export interface Expedition {
  id: string;
  name?: string;
  owner: string;
  from: GeoPoint;
  to: GeoPoint;
}

/** The shop's reagent-map API, as the F1 suite and E1's read-back (AC11) use it. */
export class ShopApi {
  private token?: string;

  constructor(private readonly request: APIRequestContext) {}

  async signIn(email: string, password = 'password123'): Promise<void> {
    const response = await this.request.post(`${MOCK_API_URL}/api/auth/login`, {
      data: { email, password },
    });
    this.token = (await response.json()).token;
  }

  async reagentSites(): Promise<ReagentSite[]> {
    return (await this.request.get(`${MOCK_API_URL}/api/reagent-sites`)).json();
  }

  async visibleExpeditions(): Promise<Expedition[]> {
    const response = await this.request.get(`${MOCK_API_URL}/api/expeditions`, {
      headers: this.authHeaders(),
    });
    return response.json();
  }

  async saveExpedition(from: GeoPoint, to: GeoPoint): Promise<APIResponse> {
    return this.request.post(`${MOCK_API_URL}/api/expeditions`, {
      data: { from, to },
      headers: this.authHeaders(),
    });
  }

  private authHeaders(): Record<string, string> {
    return this.token ? { Authorization: `Bearer ${this.token}` } : {};
  }
}
