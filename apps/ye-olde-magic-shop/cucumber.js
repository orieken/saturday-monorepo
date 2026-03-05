console.log('[cucumber.js] Loading configuration...');

const formatOptions = {
  snippetInterface: 'async-await',
};

const timestamp = new Date().toISOString().replace(/[:-]/g, '').replace(/\..+/, '');

// Allow runner to override report paths via environment variables
const htmlReport = process.env.CUCUMBER_HTML_REPORT || `./reports/cucumber/report-${timestamp}.html`;
const jsonReport = process.env.CUCUMBER_JSON_REPORT || `./reports/cucumber/report-${timestamp}.json`;

const formatters = [
  `html:${htmlReport}`,
  `json:${jsonReport}`
];

if (process.env.ENABLE_OTEL !== 'false') {
  console.log('[cucumber.js] Enabling OTel formatter');
  formatters.push('@orieken/saturday-cucumber-otel-formatter:./reports/otel-debug.log');
}

module.exports = {
  default: {
    format: formatters,
    formatOptions,
    requireModule: [
      'ts-node/register',
      'dotenv/config'
    ],
    require: [
      'lib/**/*.ts',
      'features/**/*.ts',
    ],
  },

};
