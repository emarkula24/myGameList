import { test, expect } from '@playwright/test';
import { text } from 'stream/consumers';

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
})