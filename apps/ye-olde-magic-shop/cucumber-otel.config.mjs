export default {
    resourceAttributes: {
        'service.name': 'cucumber-test',
        'service.instance.id': 'dev-instance-1',
        'service.version': '1.0.0',
        'service.runner': 'cucumber',
        'service.testrun.name': process.env.TEST_RUN_NAME || 'unknown'
    },
    scenarioAttributes: (pickle, gherkinDocument) => {
        return {
            'custom.feature.name': gherkinDocument.feature?.name || 'unknown',
            'custom.scenario.tags': pickle.tags.map(t => t.name).join(',')
        };
    },
    stepAttributes: (pickleStep, gherkinStep) => {
        return {
            'custom.step.keyword': gherkinStep.keyword
        };
    }
};
