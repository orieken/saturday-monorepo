
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ModelsFacade } from '../../src/ml/facades/models.facade';
import { BaseSite } from '../../src/base/base-site';

describe('ModelsFacade', () => {
    let mockSite: any;
    let facade: ModelsFacade;

    beforeEach(() => {
       mockSite = {};
       facade = new ModelsFacade(mockSite as BaseSite);
    });

    it('should initialize models lazy', () => {
        expect(facade.visualModel).toBeDefined();
        expect(facade.anomalyModel).toBeDefined();
        expect(facade.regressionModel).toBeDefined();
    });

    it('should list models with filters', async () => {
        const all = await facade.listModels();
        expect(all.length).toBeGreaterThan(0);

        const anomalyModels = await facade.listModels({ type: 'anomaly' });
        expect(anomalyModels.every(m => m.type === 'anomaly')).toBe(true);

        const filteredByLabel = await facade.listModels({ label: 'homepage' });
        expect(filteredByLabel).toHaveLength(1);
    });

    it('should get model info', async () => {
        const info = await facade.getModelInfo('homepage_baseline');
        expect(info).not.toBeNull();
        expect(info?.label).toBe('homepage_baseline');
        
        const infoVer = await facade.getModelInfo('product_gallery_standard', '2.1');
        expect(infoVer?.version).toBe('2.1');
    });

    it('should get model performance', async () => {
        const perf = await facade.getModelPerformance('homepage_baseline');
        expect(perf).not.toBeNull();
        expect(perf?.accuracy).toBeGreaterThan(0);
    });

    it('should save, load and delete model', async () => {
        const ver = await facade.saveModel('new_model', {});
        expect(ver).toBe('0.1'); // Assuming base calculation logic

        await expect(facade.loadModel('homepage_baseline')).resolves.toEqual({});
        await expect(facade.loadModel('non_existent')).rejects.toThrow();

        await expect(facade.deleteModel('homepage_baseline')).resolves.toBe(true);
        await expect(facade.deleteModel('non_existent')).rejects.toThrow();
    });

    it('should compare models', async () => {
        const comparison = await facade.compareModels('homepage_baseline', 'product_gallery_standard');
        expect(comparison.comparison.recommendation).toContain('better performance');
    });

    it('should optimize model', async () => {
        const res = await facade.optimizeModel('homepage_baseline');
        expect(res.improvementPercent).toBeDefined();
        expect(res.newVersion).toBeDefined();
    });

    it('should schedule retraining', async () => {
        const res = await facade.scheduleRetraining('model', {
             frequency: 'daily', minAccuracyThreshold: 0.9, maxDataAge: '1d', autoApprove: true 
        });
        expect(res.scheduled).toBe(true);
        expect(res.nextTraining.getTime()).toBeGreaterThan(Date.now());
    });

    it('should deploy and rollback', async () => {
        const deploy = await facade.deployModel('homepage_baseline', '1.0', 'production');
        expect(deploy.deployed).toBe(true);

        const rollback = await facade.rollbackModel('homepage_baseline', 'production');
        expect(rollback.rolledBack).toBe(true);
        expect(rollback.currentVersion).toBe('2.0'); // mocked return
    });

    it('should get usage stats', async () => {
        const stats = await facade.getModelUsageStats('label', { from: new Date(), to: new Date() });
        expect(stats.totalValidations).toBe(1250);
    });

    it('should cleanup old models', async () => {
         const result = await facade.cleanupOldModels({ keepVersions: 1, minAge: '1d' });
         expect(result.modelsDeleted).toBeGreaterThanOrEqual(0);
    });
});
