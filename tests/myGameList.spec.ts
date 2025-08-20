import { test, expect } from '@playwright/test';
import { config } from 'dotenv';
import { addGame, loginWith, registerWith, searchGame } from './helpers';
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

        expect(page.getByText("Insufficient password length.")).toBeVisible()

    })
    test("Should not be able to register when user already exists", async ({ page }) => {
        await registerWith(page, "tester1", "tester12@gmail.com", "123456A")
        await page.goto("register")
        // Try to register a second time
        await registerWith(page, "tester", "tester12@gmail.com", "123456A")

        expect(page.getByText("User already exists.")).toBeVisible()
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
        expect(page.getByText("View your gamelist!")).toBeVisible
    })
    test("Should not login with incorrect credentials", async ({ page }) => {
        await loginWith(page, "tester123", "12345")
        expect(page.getByText("Incorrect password or username.")).toBeVisible()
    })
})

test.describe("GameList", () => {
    test("Search results are rendered in view", async ({ page }) => {
        await searchGame(page, "metroid")
        const result = page.getByRole("img")
        const resultCount = result.count()
        expect(page.getByRole("img")).toBeVisible()
        expect.soft(resultCount, "the expected count 17, this may be subject to change").toBe(17)
    })
    test("Search results view is rendered after pressing enter", async ({page}) => {
        await searchGame(page, "metroid")
        await page.keyboard.press("Enter")
        expect(page.getByRole("heading")).toBeVisible()
    })

    test.beforeEach(async ({ page, request }) => {
        await request.post(`${url}/reset`)
        await request.post(`${url}/user/register`, {
            data: {
                email: "tester12@gmail.com",
                password: "123456A",
                username: "tester123"
            }
        })
        await request.post(`${url}/user/login`, {
            data: {
                username: "tester123",
                password: "123456A"
            }
        })
    })
    test("Should add game to list", async ({page}) => {
        await addGame(page, "mario")
        expect(page.getByText("Game successfully added!")).toBeVisible()
    })
    test("Should edit game status after adding it to list", async ({page}) => {
        await addGame(page, "zelda")
        expect.soft(page.getByText("Game successfully added!")).toBeVisible()
        await page.getByText("Playing").click()
        expect.soft(page.getByText("Completed")).toBeVisible()
        await page.getByText("Completed").click()

        expect(page.getByText("Game successfully updated!")).toBeVisible()
        expect(page.getByText("Playing")).not.toBeVisible()
    })
    test("Should render warning when gamelist is empty", async ({page}) => {
        await page.goto("/gamelist/tester123")
        expect(page.getByText("Gamelist is empty.")).toBeVisible()
    })
    test("Should delete game from list", async ({page}) => {
        await addGame(page, "metroid")
        await page.goto("/gamelist/tester123")
        let entries = page.getByText("Delete")
        expect.soft(entries.count()).toBeGreaterThan(0)
        await entries[0].click()
        entries = page.getByText("Delete")
        console.log("entries on page:" + entries)
    })
})