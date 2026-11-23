import { BaseFlow } from '@orieken/saturday-core';
import type { CesiumFixture } from '@orieken/saturday-playwright-cesium';
import type { ReagentMapPage } from '../ReagentMapPage';

export interface GeoPoint {
  lat: number;
  lng: number;
}

export interface DrawLineParams {
  cesium: CesiumFixture;
  from: GeoPoint;
  to: GeoPoint;
  save?: boolean;
}

/**
 * Choose the line tool, place A, place B, and optionally save (D-6). Each
 * point is placed with a real click on the globe through `cesium.clickAt`.
 */
export class DrawLineFlow extends BaseFlow {
  async execute(params?: Record<string, unknown>): Promise<void> {
    const { cesium, from, to, save } = params as unknown as DrawLineParams;
    const map = this.getPage<ReagentMapPage>('reagentMapPage');
    await map.startLine();
    await cesium.clickAt(from.lat, from.lng);
    await cesium.clickAt(to.lat, to.lng);
    if (save) await map.save();
  }
}
