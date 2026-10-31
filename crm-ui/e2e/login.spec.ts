import { expect, test } from "@playwright/test";

// Uses a real staff account: create one with `go run ./cmd/create-admin` in ads-platform-back.
const mobile = process.env.E2E_STAFF_MOBILE;
const password = process.env.E2E_STAFF_PASSWORD;

test("login reaches the dashboard", async ({ page }) => {
  test.skip(!mobile || !password, "Set E2E_STAFF_MOBILE and E2E_STAFF_PASSWORD");
  await page.goto("/login");
  await page.getByLabel(/موبایل|Mobile/).fill(mobile!);
  await page.getByLabel(/رمز عبور|Password/).fill(password!);
  await page.getByRole("button", { name: /ورود|Sign in/ }).click();
  await expect(page).toHaveURL(/\/dashboard/);
});

test("wrong password stays on the login page", async ({ page }) => {
  test.skip(!mobile, "Set E2E_STAFF_MOBILE");
  await page.goto("/login");
  await page.getByLabel(/موبایل|Mobile/).fill(mobile!);
  await page.getByLabel(/رمز عبور|Password/).fill("definitely-wrong-password");
  await page.getByRole("button", { name: /ورود|Sign in/ }).click();
  await expect(page).toHaveURL(/\/login/);
});
