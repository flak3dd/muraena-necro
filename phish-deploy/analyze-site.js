const puppeteer = require('puppeteer');
const fs = require('fs');

(async () => {
    const targetUrl = 'https://ignitioncasino.cc/';
    
    console.log(`[+] Launching browser...`);
    const browser = await puppeteer.launch({
        headless: true,
        ignoreHTTPSErrors: true,
        args: [
            '--no-sandbox',
            '--disable-setuid-sandbox',
            '--disable-web-security',
            '--disable-features=IsolateOrigins,site-per-process',
            '--ignore-certificate-errors',
            '--ignore-certificate-errors-spki-list'
        ]
    });

    const page = await browser.newPage();
    await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36');
    
    console.log(`[+] Navigating to ${targetUrl}...`);
    await page.goto(targetUrl, { waitUntil: 'networkidle2', timeout: 60000 });
    
    console.log(`[+] Analyzing page structure...`);
    
    const analysis = await page.evaluate(() => {
        const results = {
            url: window.location.href,
            title: document.title,
            forms: [],
            inputs: [],
            buttons: [],
            links: [],
            apis: [],
            cookies: [],
            localStorage: [],
            subdomains: new Set(),
            externalOrigins: new Set()
        };

        // Analyze forms
        document.querySelectorAll('form').forEach((form, idx) => {
            results.forms.push({
                index: idx,
                id: form.id || null,
                name: form.name || null,
                action: form.action || null,
                method: form.method || 'GET',
                classes: Array.from(form.classList),
                inputs: Array.from(form.querySelectorAll('input')).map(input => ({
                    type: input.type,
                    name: input.name,
                    id: input.id,
                    placeholder: input.placeholder,
                    required: input.required,
                    autocomplete: input.autocomplete
                }))
            });
        });

        // Analyze all inputs (including those outside forms)
        document.querySelectorAll('input').forEach(input => {
            results.inputs.push({
                type: input.type,
                name: input.name || null,
                id: input.id || null,
                placeholder: input.placeholder || null,
                value: input.value || null,
                classes: Array.from(input.classList),
                required: input.required,
                autocomplete: input.autocomplete || null
            });
        });

        // Analyze buttons
        document.querySelectorAll('button, input[type="submit"], input[type="button"], a[role="button"]').forEach(btn => {
            results.buttons.push({
                tag: btn.tagName.toLowerCase(),
                type: btn.type || null,
                text: btn.textContent?.trim() || btn.value || null,
                id: btn.id || null,
                classes: Array.from(btn.classList),
                onclick: btn.onclick ? 'has-onclick' : null,
                href: btn.href || null
            });
        });

        // Analyze important links (login, register, etc.)
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
                text.includes('withdraw')
            )) {
                results.links.push({
                    text: link.textContent?.trim(),
                    href: href,
                    id: link.id || null,
                    classes: Array.from(link.classList)
                });
            }

            // Collect subdomains and external origins
            try {
                const url = new URL(href);
                if (url.hostname !== window.location.hostname) {
                    if (url.hostname.includes(window.location.hostname.split('.').slice(-2).join('.'))) {
                        results.subdomains.add(url.hostname);
                    } else {
                        results.externalOrigins.add(url.hostname);
                    }
                }
            } catch (e) {}
        });

        // Check for API endpoints in scripts
        document.querySelectorAll('script').forEach(script => {
            const content = script.textContent || '';
            const apiMatches = content.match(/(https?:)?\/\/[^\s"']+\/api\/[^\s"']*/gi);
            if (apiMatches) {
                apiMatches.forEach(match => results.apis.push(match));
            }
        });

        // Get cookies
        results.cookies = document.cookie.split(';').map(c => c.trim().split('=')[0]);

        // Check localStorage keys
        try {
            for (let i = 0; i < localStorage.length; i++) {
                results.localStorage.push(localStorage.key(i));
            }
        } catch (e) {}

        // Convert Sets to Arrays for JSON serialization
        results.subdomains = Array.from(results.subdomains);
        results.externalOrigins = Array.from(results.externalOrigins);

        return results;
    });

    console.log(`[+] Checking for login modal/page...`);
    
    // Try to find and click login button
    try {
        const loginSelectors = [
            'button[class*="login"]',
            'a[class*="login"]',
            'button:contains("Login")',
            'button:contains("Sign In")',
            '[data-test*="login"]',
            '#login',
            '.login-button',
            'a[href*="login"]'
        ];

        for (const selector of loginSelectors) {
            try {
                await page.click(selector, { timeout: 2000 });
                console.log(`[+] Clicked login button: ${selector}`);
                await page.waitForTimeout(3000);
                break;
            } catch (e) {}
        }
    } catch (e) {}

    // Check for login form after clicking
    const loginFormAnalysis = await page.evaluate(() => {
        const loginResults = {
            loginFields: [],
            authEndpoints: [],
            validationRules: []
        };

        // Look for email/username fields
        const emailSelectors = [
            'input[type="email"]',
            'input[name*="email"]',
            'input[name*="username"]',
            'input[id*="email"]',
            'input[id*="username"]',
            'input[placeholder*="email"]',
            'input[placeholder*="username"]'
        ];

        emailSelectors.forEach(selector => {
            document.querySelectorAll(selector).forEach(input => {
                loginResults.loginFields.push({
                    type: 'email/username',
                    selector: selector,
                    name: input.name,
                    id: input.id,
                    placeholder: input.placeholder,
                    autocomplete: input.autocomplete
                });
            });
        });

        // Look for password fields
        document.querySelectorAll('input[type="password"]').forEach(input => {
            loginResults.loginFields.push({
                type: 'password',
                name: input.name,
                id: input.id,
                placeholder: input.placeholder,
                autocomplete: input.autocomplete
            });
        });

        // Check for patterns in scripts that might indicate auth endpoints
        document.querySelectorAll('script').forEach(script => {
            const content = script.textContent || '';
            
            // Look for common auth endpoint patterns
            const authPatterns = [
                /['"]\/api\/auth[^'"]*['"]/gi,
                /['"]\/api\/login[^'"]*['"]/gi,
                /['"]\/login[^'"]*['"]/gi,
                /['"]\/signin[^'"]*['"]/gi,
                /['"]\/authenticate[^'"]*['"]/gi,
                /['"]\/session[^'"]*['"]/gi
            ];

            authPatterns.forEach(pattern => {
                const matches = content.match(pattern);
                if (matches) {
                    matches.forEach(match => {
                        loginResults.authEndpoints.push(match.replace(/['"]/g, ''));
                    });
                }
            });
        });

        return loginResults;
    });

    // Capture screenshot
    console.log(`[+] Capturing screenshot...`);
    await page.screenshot({ 
        path: 'C:\\Users\\j\\phish-deploy\\ignition-homepage.png',
        fullPage: true 
    });

    // Try to navigate to common auth pages
    const authPages = ['/login', '/signin', '/account/login', '/auth/login'];
    const authAnalysis = {};

    for (const authPath of authPages) {
        try {
            console.log(`[+] Checking ${authPath}...`);
            const response = await page.goto(`https://ignitioncasino.cc${authPath}`, {
                waitUntil: 'networkidle2', 
                timeout: 10000 
            });
            
            if (response.status() === 200) {
                authAnalysis[authPath] = {
                    status: 200,
                    url: page.url(),
                    forms: await page.evaluate(() => {
                        return Array.from(document.querySelectorAll('form')).map(f => ({
                            action: f.action,
                            method: f.method,
                            inputs: Array.from(f.querySelectorAll('input')).map(i => ({
                                type: i.type,
                                name: i.name,
                                id: i.id
                            }))
                        }));
                    })
                };
            }
        } catch (e) {
            authAnalysis[authPath] = { status: 'failed', error: e.message };
        }
    }

    // Compile final analysis
    const finalReport = {
        timestamp: new Date().toISOString(),
        target: targetUrl,
        analysis: analysis,
        loginForm: loginFormAnalysis,
        authPages: authAnalysis,
        recommendations: {
            phishing_domain: 'ignitioncasino.cc',
            target_domain: 'www.ignitioncasino.eu',
            subdomains: analysis.subdomains,
            tracking_paths: [
                '/api/login',
                '/api/auth',
                '/login',
                '/signin',
                '/api/session'
            ],
            cookie_patterns: analysis.cookies,
            important_fields: [
                ...new Set([
                    ...loginFormAnalysis.loginFields.map(f => f.name).filter(Boolean),
                    'email', 'username', 'password', 'token', 'session'
                ])
            ]
        }
    };

    // Save report
    const reportPath = 'C:\\Users\\j\\phish-deploy\\ignition-analysis.json';
    fs.writeFileSync(reportPath, JSON.stringify(finalReport, null, 2));
    console.log(`\n[+] Analysis complete! Report saved to: ${reportPath}`);
    
    // Print summary
    console.log(`\n========== ANALYSIS SUMMARY ==========`);
    console.log(`URL: ${analysis.url}`);
    console.log(`Title: ${analysis.title}`);
    console.log(`Forms found: ${analysis.forms.length}`);
    console.log(`Inputs found: ${analysis.inputs.length}`);
    console.log(`Buttons found: ${analysis.buttons.length}`);
    console.log(`Auth links: ${analysis.links.length}`);
    console.log(`Subdomains: ${analysis.subdomains.join(', ') || 'none'}`);
    console.log(`Login fields: ${loginFormAnalysis.loginFields.length}`);
    console.log(`=====================================\n`);

    await browser.close();
})().catch(err => {
    console.error('[!] Error:', err);
    process.exit(1);
});
