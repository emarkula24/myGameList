import { test, expect } from '@playwright/test';
import axios from "axios"
import { config } from 'dotenv';
config({ path: './.env' })
const url = process.env.VITE_BACKEND_URL
test.beforeEach(async () => {
    console.log(url)
    const response = await axios.post(`http://127.0.0.1:8081/reset`) 
    
})
test("Has content", async ({ page }) => {
    await page.goto('/')
    await expect(page).toHaveTitle(/myGameList/)
})

test.describe("Register", () => {
    test("Should be able to register with approved values", async ({ page }) => {
        await page.goto("register")

        // Fill form
        await page.getByPlaceholder("Enter username").fill("tester");
        await page.getByPlaceholder("Enter email").fill("tester@gmail.com");
        await page.getByPlaceholder("Enter password").fill("aB123456@");
        await page.getByPlaceholder("Confirm password").fill("aB123456@");

        await page.getByRole("button", { name: "register" }).click()

        // expect(page.getByText("Register Success!")).toBeVisible()
        await page.waitForURL('**/login')

        expect(page.url()).toContain('/login')
    })
    test("Should not be able to register with short password", async ({ page }) => {
        await page.goto("register")

        // Fill form
        await page.getByPlaceholder("Enter username").fill("tester");
        await page.getByPlaceholder("Enter email").fill("tester@gmail.com");
        await page.getByPlaceholder("Enter password").fill("1234");
        await page.getByPlaceholder("Confirm password").fill("1234");

        await page.getByRole("button", { name: "register" }).click()

        expect(page.getByText("Insufficient password length."))

    })
    test("Should not be able to register when user already exists", async ({ page }) => {
        await page.goto("register")

        // Fill form
        await page.getByPlaceholder("Enter username").fill("tester");
        await page.getByPlaceholder("Enter email").fill("tester@gmail.com");
        await page.getByPlaceholder("Enter password").fill("aB123456@");
        await page.getByPlaceholder("Confirm password").fill("aB123456@");

        await page.getByRole("button", { name: "register" }).click()

        await page.goto("register")

        // Fill form a second time
        await page.getByPlaceholder("Enter username").fill("tester");
        await page.getByPlaceholder("Enter email").fill("tester@gmail.com");
        await page.getByPlaceholder("Enter password").fill("aB123456@");
        await page.getByPlaceholder("Confirm password").fill("aB123456@");

        await page.getByRole("button", { name: "register" }).click()

        expect(page.getByText("User already exists."))
    })
})