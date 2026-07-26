import { test, expect } from '@playwright/test';

const BASE_URL = 'https://vivi-37863aed9d29.herokuapp.com';

test.describe('Vivi Demo E2E Tests', () => {

  test.beforeEach(async ({ page }) => {
    // Capture console errors
    const consoleErrors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') consoleErrors.push(msg.text());
    });
    page.on('pageerror', err => consoleErrors.push(err.message));
    
    // Store for later assertion
    (page as any)._consoleErrors = consoleErrors;
  });

  test('1. Page loads and renders main layout', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    
    // Check the main shell loads
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Check topbar
    await expect(page.locator('.topbar')).toBeVisible();
    await expect(page.locator('.vivi-title')).toContainText('Vivi');
    
    // Check sidebar
    await expect(page.locator('.sidebar')).toBeVisible();
    
    // Check the 3 panels exist
    await expect(page.locator('#leads-panel')).toBeVisible();
    await expect(page.locator('#panel-chat')).toBeVisible();
    await expect(page.locator('#details-panel')).toBeVisible();
    
    await page.screenshot({ path: 'test-results/01-initial-load.png', fullPage: true });
  });

  test('2. Chat panel renders with Vivi greeting', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Wait for chat to load with messages
    await page.waitForTimeout(3000); // Wait for polling to fetch conversation
    
    // Check chat panel has the shell
    await expect(page.locator('.chat-top')).toBeVisible();
    await expect(page.locator('.barra-entrada')).toBeVisible();
    await expect(page.locator('#input-mensaje')).toBeVisible();
    await expect(page.locator('#btn-enviar')).toBeVisible();
    
    // Check for messages (Vivi greeting)
    const messages = page.locator('.burbuja');
    const messageCount = await messages.count();
    console.log(`Message count: ${messageCount}`);
    
    // Check for date/time display on messages
    const timestamps = page.locator('.hora');
    const timestampCount = await timestamps.count();
    console.log(`Timestamp count: ${timestampCount}`);
    
    if (timestampCount > 0) {
      const firstTimestamp = await timestamps.first().textContent();
      console.log(`First timestamp text: "${firstTimestamp}"`);
    }
    
    await page.screenshot({ path: 'test-results/02-chat-panel.png', fullPage: true });
  });

  test('3. Send a message and get Vivi response', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Wait for initial messages to load
    await page.waitForTimeout(3000);
    
    // Count initial messages
    const initialMessageCount = await page.locator('.burbuja').count();
    console.log(`Initial messages: ${initialMessageCount}`);
    
    // Check if there's an error alert
    const errorAlert = page.locator('[role="alert"]');
    if (await errorAlert.count() > 0) {
      const alertText = await errorAlert.first().textContent();
      console.log(`Error alert visible: "${alertText}"`);
    }
    
    // Type a message
    const input = page.locator('#input-mensaje');
    await input.fill('Hola, estoy interesada en comprar vivienda');
    
    // Send the message
    await page.locator('#btn-enviar').click();
    
    // Wait for message to appear and response
    await page.waitForTimeout(5000);
    
    // Count messages after sending
    const afterMessageCount = await page.locator('.burbuja').count();
    console.log(`Messages after sending: ${afterMessageCount}`);
    
    // Check if our message appeared
    const userMessages = page.locator('.burbuja.derecha');
    const userMessageCount = await userMessages.count();
    console.log(`User messages: ${userMessageCount}`);
    
    // Check if Vivi responded 
    const viviMessages = page.locator('.burbuja.izquierda');
    const viviMessageCount = await viviMessages.count();
    console.log(`Vivi messages: ${viviMessageCount}`);
    
    await page.screenshot({ path: 'test-results/03-after-message.png', fullPage: true });
  });

  test('4. Leads panel shows leads list', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Wait for leads to load
    await page.waitForTimeout(5000);
    
    // Check leads panel
    const leadsPanel = page.locator('#leads-panel');
    await expect(leadsPanel).toBeVisible();
    
    // Check if lead list has items
    const leadRows = page.locator('.lead-row');
    const leadCount = await leadRows.count();
    console.log(`Lead count: ${leadCount}`);
    
    // Check filter tabs
    const filterTabs = page.locator('.filter');
    const filterCount = await filterTabs.count();
    console.log(`Filter tab count: ${filterCount}`);
    
    // Check search input
    const searchInput = page.locator('#buscar-lead');
    if (await searchInput.count() > 0) {
      console.log('Search input present: YES');
    }
    
    await page.screenshot({ path: 'test-results/04-leads-panel.png', fullPage: true });
  });

  test('5. Details panel shows ficha or empty state', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    await page.waitForTimeout(5000);
    
    const detailsPanel = page.locator('#details-panel');
    const detailsContent = await detailsPanel.innerHTML();
    console.log(`Details panel content length: ${detailsContent.length}`);
    console.log(`Details panel first 200 chars: ${detailsContent.substring(0, 200)}`);
    
    await page.screenshot({ path: 'test-results/05-details-panel.png', fullPage: true });
  });

  test('6. Demo buttons (Avanzar Tiempo + Reiniciar) work', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Check demo buttons
    const avanzarBtn = page.locator('#btn-avanzar-tiempo');
    const reiniciarBtn = page.locator('#btn-reiniciar-demo');
    
    const avanzarVisible = await avanzarBtn.isVisible().catch(() => false);
    const reiniciarVisible = await reiniciarBtn.isVisible().catch(() => false);
    
    console.log(`Avanzar tiempo button visible: ${avanzarVisible}`);
    console.log(`Reiniciar button visible: ${reiniciarVisible}`);
    
    if (avanzarVisible) {
      await avanzarBtn.click();
      await page.waitForTimeout(2000);
      
      const aviso = page.locator('#demo-aviso');
      if (await aviso.count() > 0) {
        const avisoText = await aviso.textContent();
        console.log(`Demo aviso: "${avisoText}"`);
      }
    }
    
    await page.screenshot({ path: 'test-results/06-demo-buttons.png', fullPage: true });
  });

  test('7. Check date/time display in messages', async ({ page }) => {
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Wait for conversation to load
    await page.waitForTimeout(5000);
    
    // Check timestamps
    const timestamps = page.locator('.hora');
    const count = await timestamps.count();
    console.log(`Timestamp elements: ${count}`);
    
    for (let i = 0; i < Math.min(count, 5); i++) {
      const text = await timestamps.nth(i).textContent();
      console.log(`  Timestamp ${i}: "${text}"`);
    }
    
    // Also check if messages have valid 'creado_en' dates - we do this by checking the rendered time
    const allBubbles = page.locator('.burbuja');
    const bubbleCount = await allBubbles.count();
    console.log(`Total chat bubbles: ${bubbleCount}`);
    
    for (let i = 0; i < Math.min(bubbleCount, 5); i++) {
      const bubbleHTML = await allBubbles.nth(i).innerHTML();
      console.log(`  Bubble ${i} HTML: ${bubbleHTML.substring(0, 200)}`);
    }
    
    await page.screenshot({ path: 'test-results/07-timestamps.png', fullPage: true });
  });

  test('8. Console errors check', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });
    page.on('pageerror', err => errors.push(err.message));
    
    await page.goto(BASE_URL, { waitUntil: 'networkidle' });
    await expect(page.locator('.shell')).toBeVisible({ timeout: 15000 });
    
    // Wait and interact
    await page.waitForTimeout(5000);
    
    console.log(`Console errors: ${errors.length}`);
    errors.forEach((e, i) => console.log(`  Error ${i}: ${e}`));
    
    // Network failures check
    const failedRequests: string[] = [];
    page.on('requestfailed', request => {
      failedRequests.push(`${request.method()} ${request.url()} - ${request.failure()?.errorText}`);
    });
    
    await page.waitForTimeout(3000);
    console.log(`Failed network requests: ${failedRequests.length}`);
    failedRequests.forEach((r, i) => console.log(`  Failed ${i}: ${r}`));
  });
});
