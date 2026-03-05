export default {
    resourceAttributes: {
        'service.name': 'playwright-test',
        'service.instance.id': 'dev-instance-1',
        'service.version': '1.0.0',
        'service.runner': 'playwright',
        'service.testrun.name': process.env.TEST_RUN_NAME || 'unknown'
    },
    testAttributes: (test) => {
        return {
            'custom.test.tags': test.tags ? test.tags.map(t => t.tag).join(',') : '',
            'test.browser': process.env.BROWSER || 'chrome',
            'custom.test.annotations': test.annotations ? test.annotations.map(a => `${a.type}:${a.description}`).join(';') : ''
        };
    },
    stepAttributes: (step) => {
        return {
            'custom.step.category': step.category
        };
    }
};
