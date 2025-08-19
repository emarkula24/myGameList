import { test, expect } from '@playwright/test';
import { config } from 'dotenv';
import { loginWith, registerWith } from './helpers';
config({ path: './.env' })
const url = process.env.VITE_BACKEND_URL


test("Has content", async ({ page }) => {
    await page.goto('/')
    await expect(page).toHaveTitle(/myGameList/)
})

test.describe("Register", () => {
    test.beforeEach(async ({ page, request }) => {
        await request.post(`${url}/reset`)
        await page.goto("/register")
    })
    test("Should be able to register with approved values", async ({ page }) => {
        await registerWith(page, "tester", "tester12@gmail.com", "123456A")

        // expect(page.getByText("Register Success!")).toBeVisible()
        await page.waitForURL('**/login')

        expect(page.url()).toContain('/login')
    })
    test("Should not be able to register with short password", async ({ page }) => {
        await registerWith(page, "tester", "tester1@gmail.com", "1234")

        expect(page.getByText("Insufficient password length."))

    })
    test("Should not be able to register when user already exists", async ({ page }) => {
        await registerWith(page, "tester1", "tester12@gmail.com", "123456A")
        await page.goto("register")
        // Try to register a second time
        await registerWith(page, "tester", "tester12@gmail.com", "123456A")

        expect(page.getByText("User already exists."))
    })
})

test.describe("Login", () => {
    test.beforeEach(async ({ page, request }) => {
        await request.post(`${url}/reset`)
        await request.post(`${url}/user/register`, {
            data: {
                email: "tester12@gmail.com",
                password: "123456A",
                username: "tester123"
            }
        })
        await page.goto("/login")
    })
    test("Should login with correct credentials", async ({ page }) => {
        await loginWith(page, "tester123", "123456A")
        await page.waitForURL("/")
        expect(page.getByText("View your gamelist!"))
    })

})