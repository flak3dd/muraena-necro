const puppeteer = require('puppeteer');
const fs = require('fs');
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

function ask(question) {
    return new Promise(resolve => {
        rl.question(question, resolve);
    });
}

(async () => {
    const targetUrl = await ask('Enter target URL (default: https://ignitioncasino.cc): ') || 'https://ignitioncasino.cc';
    
    console.log(`\n[+] Launching Chromium browser (GUI mode)...`);
    console.log('[i] The browser window will open. Navigate and interact with the site.');
    console.log('[i] Press ENTER in this console when ready to capture analysis.\n');
    
    const browser = await puppeteer.launch({
        headless: false,
        devtools: true,
        ignoreHTTPSErrors: true,
        defaultViewport: null,
        args: [
            '--start-maximized',
            '--no-sandbox',
            '--disable-setuid-sandbox',
            '--disable-web-security',
            '--disable-features=IsolateOrigins,site-per-process',
            '--ignore-certificate-errors',
            '--ignore-certificate-errors-spki-list'
        ]
    });

    const pages = await browser.pages();
    const page = pages[0];
    
    await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36');
    
    console.log(`[+] Navigating to ${targetUrl}...`);
    
    try {
        await page.goto(targetUrl, { 
            waitUntil: 'networkidle2', 
            timeout: 30000 
        });
        console.log('[✓] Page loaded successfully!');
    } catch (err) {
        console.log(`[!] Initial navigation error: ${err.message}`);
        console.log('[i] Browser is still open. You can manually navigate to the site.');
    }
    
    // Wait for user to interact
    await ask('\n[?] Press ENTER when you are ready to capture analysis (after navigating to login page, etc.)...\n');
    
    console.log('\n[+] Capturing comprehensive analysis...\n');
    
    const analysis = await page.evaluate(() => {
        const results = {
            url: window.location.href,
            title: document.title,
            forms: [],
            inputs: [],
            buttons: [],
            links: [],
            scripts: [],
            iframes: [],
            cookies: [],
            localStorage: [],
            sessionStorage: [],
            subdomains: new Set(),
            externalOrigins: new Set(),
            apiEndpoints: new Set(),
            dataAttributes: []
        };

        // Analyze ALL forms in detail
        document.querySelectorAll('form').forEach((form, idx) => {
            const formData = {
                index: idx,
                id: form.id || null,
                name: form.name || null,
                action: form.action || null,
                method: form.method || 'GET',
                classes: Array.from(form.classList),
                enctype: form.enctype || null,
                target: form.target || null,
                autocomplete: form.autocomplete || null,
                inputs: []
            };
            
            form.querySelectorAll('input, select, textarea, button').forEach(field => {
                formData.inputs.push({
                    tag: field.tagName.toLowerCase(),
                    type: field.type || null,
                    name: field.name || null,
                    id: field.id || null,
                    placeholder: field.placeholder || null,
                    value: field.value || null,
                    classes: Array.from(field.classList),
                    required: field.required,
                    autocomplete: field.autocomplete || null,
                    pattern: field.pattern || null,
                    minlength: field.minLength || null,
                    maxlength: field.maxLength || null,
                    disabled: field.disabled,
                    readonly: field.readOnly
                });
            });
            
            results.forms.push(formData);
        });

        // Analyze ALL inputs (including those outside forms)
        document.querySelectorAll('input, select, textarea').forEach(input => {
            results.inputs.push({
                tag: input.tagName.toLowerCase(),
                type: input.type || null,
                name: input.name || null,
                id: input.id || null,
                placeholder: input.placeholder || null,
                value: input.value || null,
                classes: Array.from(input.classList),
                required: input.required,
                autocomplete: input.autocomplete || null,
                pattern: input.pattern || null,
                'data-*': Array.from(input.attributes)
                    .filter(attr => attr.name.startsWith('data-'))
                    .map(attr => ({ name: attr.name, value: attr.value }))
            });
        });

        // Analyze ALL buttons
        document.querySelectorAll('button, input[type="submit"], input[type="button"], [role="button"]').forEach(btn => {
            results.buttons.push({
                tag: btn.tagName.toLowerCase(),
                type: btn.type || null,
                text: btn.textContent?.trim() || btn.value || null,
                id: btn.id || null,
                name: btn.name || null,
                classes: Array.from(btn.classList),
                onclick: btn.onclick ? 'has-onclick' : null,
                href: btn.href || null,
                disabled: btn.disabled,
                'data-*': Array.from(btn.attributes)
                    .filter(attr => attr.name.startsWith('data-'))
                    .map(attr => ({ name: attr.name, value: attr.value }))
            });
        });

        // Analyze important links
        document.querySelectorAll('a').forEach(link => {
            const href = link.href;
            const text = link.textContent?.trim().toLowerCase();
            
            if (text && (
                text.includes('login') || 
                text.includes('sign in') || 
                text.includes('register') || 
                text.includes('sign up') ||
                text.includes('account') ||
                text.includes('profile') ||
                text.includes('cashier') ||
                text.includes('deposit') ||
                text.includes('withdraw') ||
                text.includes('dashboard') ||
                text.includes('lobby')
            )) {
                results.links.push({
                    text: link.textContent?.trim(),
                    href: href,
                    id: link.id || null,
                    classes: Array.from(link.classList),
                    target: link.target || null
                });
            }

            // Collect subdomains and external origins
            try {
                const url = new URL(href);
                if (url.hostname !== window.location.hostname) {
                    const currentDomain = window.location.hostname.split('.').slice(-2).join('.');
                    if (url.hostname.includes(currentDomain)) {
                        results.subdomains.add(url.hostname);
                    } else {
                        results.externalOrigins.add(url.hostname);
                    }
                }
            } catch (e) {}
        });

        // Analyze scripts for API endpoints
        document.querySelectorAll('script').forEach(script => {
            const src = script.src;
            if (src) {
                results.scripts.push({ type: 'external', src: src });
            }
            
            const content = script.textContent || '';
            
            // Find API endpoints
            const apiPatterns = [
                /(https?:)?\/\/[^\s"']+\/api\/[^\s"']*/gi,
                /['"]\/api\/[^'"]+['"]/gi,
                /['"]\/auth\/[^'"]+['"]/gi,
                /['"]\/login[^'"]*['"]/gi,
                /['"]\/session[^'"]+['"]/gi
            ];
            
            apiPatterns.forEach(pattern => {
                const matches = content.match(pattern);
                if (matches) {
                    matches.forEach(match => {
                        const cleaned = match.replace(/['"]/g, '');
                        results.apiEndpoints.add(cleaned);
                    });
                }
            });
        });

        // Analyze iframes
        document.querySelectorAll('iframe').forEach(iframe => {
            results.iframes.push({
                src: iframe.src || null,
                id: iframe.id || null,
                name: iframe.name || null,
                classes: Array.from(iframe.classList)
            });
        });

        // Get cookies
        results.cookies = document.cookie.split(';').map(c => {
            const parts = c.trim().split('=');
            return { name: parts[0], value: parts[1] || '' };
        }).filter(c => c.name);

        // Check localStorage
        try {
            for (let i = 0; i < localStorage.length; i++) {
                const key = localStorage.key(i);
                results.localStorage.push({
                    key: key,
                    value: localStorage.getItem(key)?.substring(0, 100) + '...'
                });
            }
        } catch (e) {}

        // Check sessionStorage
        try {
            for (let i = 0; i < sessionStorage.length; i++) {
                const key = sessionStorage.key(i);
                results.sessionStorage.push({
                    key: key,
                    value: sessionStorage.getItem(key)?.substring(0, 100) + '...'
                });
            }
        } catch (e) {}

        // Collect all data-* attributes
        document.querySelectorAll('[data-api], [data-endpoint], [data-url], [data-action]').forEach(el => {
            results.dataAttributes.push({
                tag: el.tagName.toLowerCase(),
                id: el.id || null,
                classes: Array.from(el.classList),
                attributes: Array.from(el.attributes)
                    .filter(attr => attr.name.startsWith('data-'))
                    .map(attr => ({ name: attr.name, value: attr.value }))
            });
        });

        // Convert Sets to Arrays
        results.subdomains = Array.from(results.subdomains);
        results.externalOrigins = Array.from(results.externalOrigins);
        results.apiEndpoints = Array.from(results.apiEndpoints);

        return results;
    });

    // Get network requests
    const client = await page.target().createCDPSession();
    await client.send('Network.enable');
    
    console.log('[+] Capturing screenshot...');
    const timestamp = Date.now();
    await page.screenshot({ 
        path: `C:\\Users\\j\\phish-deploy\\screenshot-${timestamp}.png`,
        fullPage: true 
    });
    console.log(`[✓] Screenshot saved: screenshot-${timestamp}.png`);

    // Get page HTML
    const html = await page.content();
    fs.writeFileSync(`C:\\Users\\j\\phish-deploy\\page-${timestamp}.html`, html);
    console.log(`[✓] HTML saved: page-${timestamp}.html`);

    // Compile final report
    const finalReport = {
        timestamp: new Date().toISOString(),
        target: analysis.url,
        title: analysis.title,
        
        forms: {
            count: analysis.forms.length,
            details: analysis.forms
        },
        
        inputs: {
            count: analysis.inputs.length,
            byType: analysis.inputs.reduce((acc, inp) => {
                acc[inp.type || 'unknown'] = (acc[inp.type || 'unknown'] || 0) + 1;
                return acc;
            }, {}),
            details: analysis.inputs
        },
        
        buttons: {
            count: analysis.buttons.length,
            details: analysis.buttons
        },
        
        links: {
            count: analysis.links.length,
            details: analysis.links
        },
        
        apiEndpoints: analysis.apiEndpoints,
        subdomains: analysis.subdomains,
        externalOrigins: analysis.externalOrigins,
        
        storage: {
            cookies: analysis.cookies,
            localStorage: analysis.localStorage,
            sessionStorage: analysis.sessionStorage
        },
        
        iframes: analysis.iframes,
        dataAttributes: analysis.dataAttributes,
        
        configRecommendations: {
            phishing_domain: 'ignitioncasino.cc',
            target_domain: new URL(analysis.url).hostname,
            subdomains: analysis.subdomains,
            
            tracking_paths: [
                ...new Set([
                    '/api/login',
                    '/api/auth',
                    '/login',
                    '/signin',
                    ...analysis.apiEndpoints.filter(ep => 
                        ep.includes('login') || 
                        ep.includes('auth') || 
                        ep.includes('session')
                    )
                ])
            ],
            
            cookie_patterns: analysis.cookies.map(c => c.name),
            localStorage_keys: analysis.localStorage.map(ls => ls.key),
            
            important_fields: [
                ...new Set([
                    ...analysis.inputs.filter(i => i.name).map(i => i.name),
                    'email', 'username', 'password', 'token', 'session'
                ])
            ],
            
            form_actions: analysis.forms.map(f => f.action).filter(Boolean)
        }
    };

    // Save comprehensive report
    const reportPath = `C:\\Users\\j\\phish-deploy\\analysis-${timestamp}.json`;
    fs.writeFileSync(reportPath, JSON.stringify(finalReport, null, 2));
    
    console.log(`\n[+] Analysis complete! Report saved to: analysis-${timestamp}.json\n`);
    
    // Print summary
    console.log('========== COMPREHENSIVE ANALYSIS SUMMARY ==========');
    console.log(`URL: ${analysis.url}`);
    console.log(`Title: ${analysis.title}`);
    console.log(`\nForms: ${analysis.forms.length}`);
    analysis.forms.forEach((form, idx) => {
        console.log(`  Form ${idx + 1}: ${form.action || 'no action'} (${form.method})`);
        console.log(`    Inputs: ${form.inputs.length}`);
        form.inputs.forEach(inp => {
            console.log(`      - ${inp.tag} [${inp.type}] name="${inp.name}" id="${inp.id}"`);
        });
    });
    
    console.log(`\nInputs (all): ${analysis.inputs.length}`);
    const inputTypes = analysis.inputs.reduce((acc, inp) => {
        acc[inp.type || 'unknown'] = (acc[inp.type || 'unknown'] || 0) + 1;
        return acc;
    }, {});
    Object.entries(inputTypes).forEach(([type, count]) => {
        console.log(`  ${type}: ${count}`);
    });
    
    console.log(`\nButtons: ${analysis.buttons.length}`);
    console.log(`Auth-related links: ${analysis.links.length}`);
    console.log(`API Endpoints found: ${analysis.apiEndpoints.length}`);
    analysis.apiEndpoints.forEach(api => console.log(`  - ${api}`));
    
    console.log(`\nSubdomains: ${analysis.subdomains.length}`);
    analysis.subdomains.forEach(sub => console.log(`  - ${sub}`));
    
    console.log(`\nCookies: ${analysis.cookies.length}`);
    analysis.cookies.forEach(cookie => console.log(`  - ${cookie.name}`));
    
    console.log(`\nLocalStorage keys: ${analysis.localStorage.length}`);
    analysis.localStorage.forEach(ls => console.log(`  - ${ls.key}`));
    
    console.log('\n===================================================\n');
    
    await ask('[?] Press ENTER to close browser and exit...');
    
    await browser.close();
    rl.close();
    
    console.log('\n[✓] Done!\n');
})().catch(err => {
    console.error('[!] Error:', err);
    rl.close();
    process.exit(1);
});
