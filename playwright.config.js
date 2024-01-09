const { devices } = require('@playwright/test');

module.exports = {
  projects: [
    {
      name: 'Chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'Firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'WebKit',
      use: { ...devices['Desktop Safari'] },
    },
  ],

  retries: 2,

  use: {
    baseURL: 'http://localhost:3000',
  },

  timeout: 30000,

  outputDir: 'test-results/',
};
