export interface CucumberStepRef {
  text: string;
  line: number;
}

export interface CucumberScenarioRef {
  id: string; // kebab-case from scenario name
  line: number; // line number in .feature file
  name: string;
  tags: string[];
  file: string;
  steps: CucumberStepRef[];
}

export interface CucumberFeatureRef {
  id: string; // kebab-case from feature name
  name: string;
  file: string;
  description: string;
  scenarios: CucumberScenarioRef[];
}

export interface CucumberIndex {
  framework: 'cucumber';
  suiteId: string;
  features: CucumberFeatureRef[];
}

export const demoEcommerceCucumberIndex: CucumberIndex = {
  framework: 'cucumber',
  suiteId: 'demo-ecommerce',
  features: [
    {
      id: 'checkout-process-order-confirmation',
      name: 'Checkout Process - Order Confirmation',
      file: 'checkout_process.feature',
      description: 'Happy path checkout flow',
      scenarios: [
        {
          id: 'proceeding-to-shipping-information-as-a-logged-in-user',
          line: 3,
          name: 'Proceeding to shipping information as a logged-in user',
          tags: ['@smoke'],
          file: 'checkout_process.feature',
          steps: [
            { text: 'Given I am logged in and have items in my cart', line: 4 },
            { text: 'When I proceed to checkout', line: 5 },
            { text: 'Then I should be on the shipping information page', line: 6 }
          ]
        },
        {
          id: 'entering-new-shipping-address',
          line: 10,
          name: 'Entering new shipping address',
          tags: [],
          file: 'checkout_process.feature',
          steps: [
            { text: 'Given I am on the shipping information page', line: 11 },
            { text: 'When I enter my address', line: 12 },
            { text: 'Then my shipping info should be saved', line: 13 }
          ]
        }
      ]
    },
    {
      id: 'product-search',
      name: 'Product Search',
      file: 'product_search.feature',
      description: 'Basic product search',
      scenarios: [
        {
          id: 'searching-for-a-product-by-name',
          line: 2,
          name: 'Searching for a product by name',
          tags: ['@regression'],
          file: 'product_search.feature',
          steps: [
            { text: 'Given I am on the home page', line: 3 },
            { text: 'When I search for "keyboard"', line: 4 },
            { text: 'Then I should see search results containing "keyboard"', line: 5 }
          ]
        }
      ]
    }
  ]
};

export type { CucumberFeatureRef as TCucumberFeatureRef, CucumberScenarioRef as TCucumberScenarioRef };
