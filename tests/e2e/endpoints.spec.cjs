const { test, expect } = require('@playwright/test');

test('style.css should return a stylesheet', async ({ request }) => {
    const response = await request.get(`${process.env.APP_URL}/style.css`);

    expect(response.status()).toBe(200);

    expect(response.headers()['content-type']).toBe('text/css; charset=utf-8');
});

test('/ should return a html', async ({ request }) => {
    const response = await request.get(`${process.env.APP_URL}/`);

    expect(response.status()).toBe(200);

    expect(response.headers()['content-type']).toBe('texthtml; charset=utf-8');
});

