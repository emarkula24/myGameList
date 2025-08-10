import { test, expect } from '@playwright/test';
import { text } from 'stream/consumers';

test("Has content", async ({page}) => {
    await page.goto('http://127.0.0.1:5173')
    await expect(page).toHaveTitle(/myGameList/)
})

test.describe("Register", () => {
    test("Should be able to register with approved values",  async ({ page }) => {
        await page.goto("/register")
        const oldUrl = page.url()
        const textboxes = await page.getByRole("textbox").all()
        await textboxes[0].fill("tester")
        await textboxes[1].fill("tester@gmail.com")
        await textboxes[2].fill("aB123456@")
        await textboxes[3].fill("aB123456@")

        await page.getByRole("button", {name: "register"}).click()

        expect(page.getByText("Register Success!")).toBeVisible()

        await page.waitForTimeout(800)
        const newUrl = page.url()
        expect(newUrl).not.toBe(oldUrl)
    })
})