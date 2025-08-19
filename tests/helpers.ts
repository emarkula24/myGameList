import { Page } from "@playwright/test"

export const loginWith = async (page: Page, username: string, password: string) => {
    await page.getByPlaceholder("Enter username").fill(username)
    await page.getByPlaceholder("Enter password").fill(password)

    await page.getByRole("button", { name: "Login" }).click()
}

export const registerWith = async (page: Page, username: string, email: string, password: string) => {
    await page.getByPlaceholder("Enter username").fill(username);
    await page.getByPlaceholder("Enter email").fill(email);
    await page.getByPlaceholder("Enter password").fill(password);
    await page.getByPlaceholder("Confirm password").fill(password);

    await page.getByRole("button", { name: "register" }).click()
}