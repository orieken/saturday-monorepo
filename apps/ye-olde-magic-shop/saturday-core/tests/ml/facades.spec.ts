
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TrainersFacade } from '../../src/ml/facades/trainers.facade';
import { BaseSite } from '../../src/base/base-site';

describe('TrainersFacade', () => {
  let mockSite: any;
  let facade: TrainersFacade;
  let mockPage: any;

  beforeEach(() => {
     mockPage = {
         setViewportSize: vi.fn(),
     };
    
    mockSite = {
       getPage: vi.fn().mockReturnValue({ visit: vi.fn(), captureForTraining: vi.fn() }),
       getFlow: vi.fn().mockReturnValue({ execute: vi.fn() }),
       pageObject: mockPage
    };
    facade = new TrainersFacade(mockSite as BaseSite);
  });

  it('should initialize trainers lazy', () => {
    expect(facade.screenshotTrainer).toBeDefined();
    expect(facade.elementTrainer).toBeDefined();
    expect(facade.pageTrainer).toBeDefined();
    expect(facade.flowTrainer).toBeDefined();
  });

  it('should train current page', async () => {
    const result = await facade.trainCurrentPage('test-label');
    expect(result.success).toBe(true);
    expect(result.modelLabel).toBe('test-label');
  });

  it('should train multiple pages', async () => {
    const results = await facade.trainMultiplePages([{ pageName: 'home', label: 'home-label' }]);
    expect(results).toHaveLength(1);
    expect(mockSite.getPage).toHaveBeenCalledWith('home');
  });

  it('should train complete user journey', async () => {
      const results = await facade.trainCompleteUserJourney('test-journey', [
          { type: 'page', target: 'home' },
          { type: 'element', target: '.btn' },
          { type: 'flow', target: 'login', params: {} }
      ]);
      expect(results).toHaveLength(3);
  });
  
  it('should train responsive design', async () => {
      const results = await facade.trainResponsiveDesign('responsive', [
          { name: 'mobile', width: 375, height: 667 }
      ]);
      expect(results).toHaveLength(1);
      expect(mockPage.setViewportSize).toHaveBeenCalledWith({ width: 375, height: 667 });
  });

  it('should train across browsers', async () => {
      const results = await facade.trainAcrossBrowsers('label', ['chrome', 'firefox']);
      expect(results).toHaveLength(2);
  });

  it('should get training stats', async () => {
      const stats = await facade.getTrainingStats();
      expect(stats.totalSamples).toBe(0);
  });

  it('should compare model versions', async () => {
      const comparison = await facade.compareModelVersions('label', ['v1', 'v2']);
      expect(comparison).toHaveLength(2);
  });
});
