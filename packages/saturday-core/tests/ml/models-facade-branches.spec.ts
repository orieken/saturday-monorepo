
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ModelsFacade } from '../../src/ml/facades/models.facade';

describe('ModelsFacade Branch Coverage', () => {
    let mockSite: any;
    let facade: ModelsFacade;

    beforeEach(() => {
        mockSite = {
        };
        facade = new ModelsFacade(mockSite);
    });

    it('should filter listModels', async () => {
        // Mock listModels internal array is hardcoded in the method (based on previous view).
        // If it was querying external, we'd mock it. But here we test the filtering logic on the hardcoded data.
        // It has 2 models: anomaly (homepage) and visual (gallery).
        
        let models = await facade.listModels();
        expect(models).toHaveLength(2);

        models = await facade.listModels({ type: 'anomaly' });
        expect(models).toHaveLength(1);
        expect(models[0].type).toBe('anomaly');

        models = await facade.listModels({ type: 'missing' });
        expect(models).toHaveLength(0);

        models = await facade.listModels({ label: 'homepage' });
        expect(models).toHaveLength(1);

        models = await facade.listModels({ minAccuracy: 0.94 });
        expect(models).toHaveLength(1); // 0.95 > 0.94

        models = await facade.listModels({ createdAfter: new Date('2025-01-10') });
        expect(models).toHaveLength(1); // 2025-01-15 > 10
    });

    it('should handle getModelInfo nulls', async () => {
        expect(await facade.getModelInfo('homepage_baseline', '1.0')).toBeDefined();
        expect(await facade.getModelInfo('homepage_baseline', '9.9')).toBeNull();
        expect(await facade.getModelInfo('missing')).toBeNull();
    });

    it('should handle getModelPerformance nulls', async () => {
        expect(await facade.getModelPerformance('missing')).toBeNull();
        // Present
        expect(await facade.getModelPerformance('homepage_baseline')).toBeDefined();
    });

    it('should handle load/delete/deploy errors', async () => {
        await expect(facade.loadModel('missing')).rejects.toThrow();
        await expect(facade.deleteModel('missing')).rejects.toThrow();
        await expect(facade.deployModel('missing', '1.0', 'prod' as any)).rejects.toThrow();
        // saveModel doesn't seem to throw on 'missing' because it creates new version
        await facade.saveModel('new', {});
    });

    it('should compare models logic', async () => {
        // Need to spy on getModelPerformance
        const spy = vi.spyOn(facade, 'getModelPerformance');
        spy.mockResolvedValueOnce({ accuracy: 0.9 } as any).mockResolvedValueOnce({ accuracy: 0.8 } as any);
        
        // A > B (>0.05)
        let res = await facade.compareModels('A', 'B');
        expect(res.comparison.recommendation).toContain('significantly better');
        expect(res.comparison.performanceBetter).toBe('A');

        // A slightly better (0.821 vs 0.8) -> Diff 0.021 > 0.02
        spy.mockReset();
        spy.mockResolvedValueOnce({ accuracy: 0.821 } as any).mockResolvedValueOnce({ accuracy: 0.8 } as any);
        res = await facade.compareModels('A', 'B');
        expect(res.comparison.recommendation).toContain('slightly better');

        // Similar
        spy.mockReset();
        spy.mockResolvedValueOnce({ accuracy: 0.8 } as any).mockResolvedValueOnce({ accuracy: 0.8 } as any);
        res = await facade.compareModels('A', 'B');
        expect(res.comparison.recommendation).toContain('similarly');
        
        // Error if not found
        spy.mockReset();
        spy.mockResolvedValueOnce(null).mockResolvedValueOnce(null);
        await expect(facade.compareModels('A', 'B')).rejects.toThrow();
    });
    
    it('should optimize and schedule', async () => {
        // Optimize missing
        // Need to spy getModelInfo because we rely on hardcoded list
        // Or pass 'homepage_baseline' which exists in hardcoded list
        await expect(facade.optimizeModel('missing')).rejects.toThrow();
        
        const opt = await facade.optimizeModel('homepage_baseline');
        expect(opt.newVersion).toBeDefined();

        const sched = await facade.scheduleRetraining('l', { frequency: 'daily', minAccuracyThreshold: 0.9, maxDataAge: '1d', autoApprove: true });
        expect(sched.scheduled).toBe(true);
        expect(sched.nextTraining.getTime()).toBeGreaterThan(Date.now());
    });

    it('should cleanup old models', async () => {
        // Default listModels returns 2 models
        // keepVersions: 1.
        // If query returns multiple versions for same label, it keeps newest.
        // Existing hardcoded: 'homepage_baseline' v1.0, 'product_gallery_standard' v2.1.
        // Each label has 1 version. So nothing deleted.
        let res = await facade.cleanupOldModels({ keepVersions: 1, minAge: '1d' });
        expect(res.modelsDeleted).toBe(0);

        // Mock listModels to return old versions
        vi.spyOn(facade, 'listModels').mockResolvedValue([
             { label: 'A', version: '1.0', createdAt: new Date(0) } as any,
             { label: 'A', version: '2.0', createdAt: new Date() } as any
        ]);
        
        // Should delete 1.0
        // excludeLabels logic
        res = await facade.cleanupOldModels({ keepVersions: 1, minAge: '0d', excludeLabels: ['B'] });
        expect(res.modelsDeleted).toBe(1);
        expect(res.deletedModels[0]).toContain('1.0');
    });

    it('should cover rollback', async () => {
        // Mock getCurrentDeployment
        const res = await facade.rollbackModel('l', 'production');
        expect(res.rolledBack).toBe(true);
        
        // If no rollback version?
        // Mock getCurrentDeployment to return null or no rollbackVersion
        // facade.getCurrentDeployment is private. Mock it using spyOn if protected or via prototype?
        // Or create test subclass.
    });
    
    it('should deploy logic', async () => {
        const res = await facade.deployModel('homepage_baseline', '1.0', 'production');
        expect(res.deployed).toBe(true);
    });
});
